package client

import (
	"log"

	"../common"
)

// Requestor ...
type Requestor struct {
	crh        *ClientRequestHandler
	marshaller *common.Marshaller
}

func newRequestor(host string, port int) *Requestor {
	crh := new(ClientRequestHandler)
	return &Requestor{
		crh:        crh,
		marshaller: new(common.Marshaller),
	}
}

func (r *Requestor) getServiceInfo(name string) (string, int) {
	crh := newClientRequestHandler("localhost", 5555)
	crh.connect()

	requestInfo := common.RequestInfo{Name: name}

	data := r.marshaller.Marshall(requestInfo)

	consultPkt := common.ConsultPkt{ConsultType: "consult", Data: data}

	pkt := r.marshaller.Marshall(consultPkt)

	crh.send(pkt)

	retData := crh.receive()

	returnPkt := new(common.Service)

	r.marshaller.Unmarshall(retData, returnPkt)

	return returnPkt.IP, int(returnPkt.Port)
}

func (r *Requestor) invoke(request common.RequestPkt) *common.ReturnPkt {

	host, port := r.getServiceInfo(request.MethodName)
	log.Printf("Service %s on %s,%d", request.MethodName, host, port)
	r.crh = newClientRequestHandler(host, port)
	r.crh.connect()

	marshContent := r.marshaller.Marshall(request)

	r.crh.send(marshContent)

	marshRet := r.crh.receive()
	var resPkt common.ReturnPkt
	r.marshaller.Unmarshall(marshRet, &resPkt)

	return &resPkt
}
