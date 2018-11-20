package client

import (
	"log"

	"../common"
)

// Requestor ...
type Requestor struct {
	crh        *ClientRequestHandler
	marshaller *common.Marshaller
	usuario    string
	senha      string
}

func newRequestor(user string, password string) *Requestor {
	crh := new(ClientRequestHandler)
	return &Requestor{
		crh:        crh,
		marshaller: new(common.Marshaller),
		usuario:    user,
		senha:      password,
	}
}

func (r *Requestor) getServiceInfo(name string) (string, int, string) {
	crh := newClientRequestHandler("localhost", 5555)
	crh.connect()

	requestInfo := common.RequestInfo{Name: name, Usuario: r.usuario, Senha: r.senha}

	data := r.marshaller.Marshall(requestInfo)

	consultPkt := common.ConsultPkt{ConsultType: "consult", Data: data}

	pkt := r.marshaller.Marshall(consultPkt)

	crh.send(pkt)

	retData := crh.receive()

	returnPkt := new(common.ConsultReturnPkt)

	r.marshaller.Unmarshall(retData, returnPkt)

	return returnPkt.Servico.IP, int(returnPkt.Servico.Port), returnPkt.Key
}

func (r *Requestor) invoke(request common.RequestPkt) *common.ReturnPkt {

	host, port, key := r.getServiceInfo(request.MethodName)
	log.Printf("Service %s on %s,%d", request.MethodName, host, port)
	log.Printf("Key to encrypt: %s", key)
	r.crh = newClientRequestHandler(host, port)
	r.crh.connect()

	marshContent := r.marshaller.Marshall(request)

	r.crh.send(marshContent)

	marshRet := r.crh.receive()
	var resPkt common.ReturnPkt
	r.marshaller.Unmarshall(marshRet, &resPkt)

	return &resPkt
}
