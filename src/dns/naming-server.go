package dns

import (
	"encoding/hex"
	"log"

	"../common"
)

var DNSKEY = "6368612e676520746869732070617373776f726420746f206120736563726e21"

// NamingService handles information about all the registered service and userss
type NamingService struct {
	services        map[string][]*common.Service
	lastServerIndex map[string]int
	users           map[string]*common.User
}

func (dns *NamingService) addService(service *common.Service) {
	log.Printf("Adding service %s in IP: %s and Port %d", service.Name, service.IP, service.Port)
	dns.services[service.Name] = append(dns.services[service.Name], service)
	dns.lastServerIndex[service.Name] = len(dns.services[service.Name]) - 1
	log.Println("Current service ", dns.lastServerIndex[service.Name])
}

func (dns *NamingService) getService(name string) *common.Service {
	log.Printf("Getting service of name %s", name)
	currServerIndex := dns.lastServerIndex[name]
	log.Println("Current service ", currServerIndex)
	dns.lastServerIndex[name]++
	if dns.lastServerIndex[name] == len(dns.services[name]) {
		dns.lastServerIndex[name] = 0
	}

	return dns.services[name][currServerIndex]
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
func (ns *NamingServer) Start(port int) {
	ns.dns = new(NamingService)
	ns.dns.services = make(map[string][]*common.Service)
	ns.dns.lastServerIndex = make(map[string]int)
	ns.dns.users = make(map[string]*common.User)
	ns.dns.users["ACC4"] = &common.User{Username: "ACC4",
		Password:    "pudim",
		Key:         "6368616e676520746869732070617373776f726420746f206120736563726574",
		AccessLevel: 1}

	for {
		var err error
		ns.srh, _ = newServerRequestHandler(port)

		data := ns.srh.receive()

		keyData, _ := hex.DecodeString(DNSKEY)

		data = common.Decrypt(keyData, data)

		pkt := new(common.ConsultPkt)

		err = ns.marshaller.Unmarshall(data, pkt)
		if err != nil {
			log.Printf("It was not possible to unmarshall packet from %s:%d", ns.srh.remoteIP, ns.srh.remotePort)
			continue
		}

		log.Printf("packet type: %s", pkt.ConsultType)
		// fmt.Printf(pkt.ConsultType)
		switch pkt.ConsultType {

		case "register":
			{
				s := new(common.Service)
				ns.marshaller.Unmarshall(pkt.Data, s)
				ns.dns.addService(s)

				returnPkt := &common.RegisterResultPkt{
					IP:   ns.srh.remoteIP,
					Port: s.Port}
				pkt := ns.marshaller.Marshall(returnPkt)
				ns.srh.send(pkt)
			}

		case "consult":
			{
				requestInfo := new(common.RequestInfo)
				ns.marshaller.Unmarshall(pkt.Data, requestInfo)

				s := ns.dns.getService(requestInfo.Name)
				key := ns.dns.getKey(requestInfo.Username, requestInfo.Password)
				returnPkt := new(common.ConsultReturnPkt)

				if key != "" && ns.dns.users[requestInfo.Username] != nil && s != nil {
					if s.AccessLevel <= ns.dns.users[requestInfo.Username].AccessLevel {
						returnPkt.ServiceInfo = s
						returnPkt.Key = key
					}
				}

				pkt := ns.marshaller.Marshall(returnPkt)

				keyData, _ := hex.DecodeString(DNSKEY)

				pkt = common.Encrypt(keyData, pkt)

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

				keyData, _ := hex.DecodeString(DNSKEY)

				pkt = common.Encrypt(keyData, pkt)

				ns.srh.send(pkt)
			}

		}

		ns.srh.close()

	}
}
