package dns

import (
	"fmt"
	"log"

	"../common"
)

// NamingService handles information about all the registered service and userss
type NamingService struct {
	services map[string]*common.Service
	users    map[string]*common.User
}

func (dns *NamingService) addService(service *common.Service) {
	log.Printf("Adding service %s in IP: %s and Port %d", service.Name, service.IP, service.Port)
	dns.services[service.Name] = service
}

func (dns *NamingService) getService(name string) *common.Service {
	log.Printf("Getting service of name %s", name)
	return dns.services[name]
}

func (dns *NamingService) getKey(user string, password string) string {
	log.Printf("Getting user %s key", user)

	if dns.users[user] == nil {
		return ""
	}

	if password == dns.users[user].Password {
		log.Printf("Key: %s", dns.users[user].Key)
		return dns.users[user].Key
	}
	return ""

}

// NamingServer perfoms all the operations
type NamingServer struct {
	dns        *NamingService
	srh        *ServerRequestHandler
	marshaller *common.Marshaller
}

// Start ...
func (ns *NamingServer) Start() {
	ns.dns = new(NamingService)
	ns.dns.services = make(map[string]*common.Service)
	ns.dns.users = make(map[string]*common.User)
	ns.dns.users["ACC4"] = &common.User{Username: "ACC4",
		Password: "pudim",
		Key:      "6368616e676520746869732070617373776f726420746f206120736563726574"}

	for {
		ns.srh = newServerRequestHandler(5555)
		data := ns.srh.receive()
		pkt := new(common.ConsultPkt)

		ns.marshaller.Unmarshall(data, pkt)
		fmt.Printf("packet type: ")
		fmt.Println(pkt.ConsultType)
		switch pkt.ConsultType {

		case "register":
			{
				s := new(common.Service)
				ns.marshaller.Unmarshall(pkt.Data, s)
				ns.dns.addService(s)
			}

		case "consult":
			{
				requestInfo := new(common.RequestInfo)
				ns.marshaller.Unmarshall(pkt.Data, requestInfo)

				s := ns.dns.getService(requestInfo.Name)
				key := ns.dns.getKey(requestInfo.Username, requestInfo.Password)
				returnPkt := new(common.ConsultReturnPkt)

				if key != "" {
					returnPkt.ServiceInfo = s
					returnPkt.Key = key
				}

				pkt := ns.marshaller.Marshall(returnPkt)
				ns.srh.send(pkt)
			}
		case "consultname":
			{
				requestInfo := new(common.RequestInfo)
				ns.marshaller.Unmarshall(pkt.Data, requestInfo)

				returnPkt := new(common.ConsultReturnPkt)

				if ns.dns.users[requestInfo.Name] != nil {
					key := ns.dns.users[requestInfo.Name].Key
					returnPkt.Key = key
				}

				pkt := ns.marshaller.Marshall(returnPkt)
				ns.srh.send(pkt)
			}

		}

		ns.srh.close()

	}
}
