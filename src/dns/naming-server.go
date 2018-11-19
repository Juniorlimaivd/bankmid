package dns

import (
	"fmt"
	"log"

	"../common"
)

// NamingService handles information about all the registered services
type NamingService struct {
	services map[string]*common.Service
}

func (dns *NamingService) addService(service *common.Service) {
	log.Printf("Adding service %s in IP: %s and Port %d", service.Name, service.IP, service.Port)
	dns.services[service.Name] = service
}

func (dns *NamingService) getService(name string) *common.Service {
	log.Printf("Geting service of name %s", name)
	return dns.services[name]
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
				pkt := ns.marshaller.Marshall(s)
				ns.srh.send(pkt)
			}

		}

		ns.srh.close()

	}
}
