package client

import (
	"encoding/hex"
	"log"

	"../common"
)

var DNSKEY = "6368612e676520746869732070617373776f726420746f206120736563726e21"

// Requestor ...
type Requestor struct {
	dnsAddr    string
	dnsPort    int
	crh        *ClientRequestHandler
	marshaller *common.Marshaller
	username   string
	password   string
}

func newRequestor(user string, password string, dnsAddr string, dnsPort int) *Requestor {
	crh := new(ClientRequestHandler)
	return &Requestor{
		dnsAddr:    dnsAddr,
		dnsPort:    dnsPort,
		crh:        crh,
		marshaller: new(common.Marshaller),
		username:   user,
		password:   password,
	}
}

func (r *Requestor) getServiceInfo(name string) (string, int, string) {
	crh := newClientRequestHandler(r.dnsAddr, r.dnsPort)
	crh.connect()

	requestInfo := common.RequestInfo{Name: name, Username: r.username, Password: r.password}

	data := r.marshaller.Marshall(requestInfo)

	consultPkt := common.ConsultPkt{ConsultType: "consult", Data: data}

	pkt := r.marshaller.Marshall(consultPkt)

	keyData, _ := hex.DecodeString(DNSKEY)

	pkt = common.Encrypt(keyData, pkt)

	crh.send(pkt)

	retData := crh.receive()

	retData = common.Decrypt(keyData, retData)

	returnPkt := new(common.ConsultReturnPkt)

	r.marshaller.Unmarshall(retData, returnPkt)

	if returnPkt.ServiceInfo == nil {
		return "", 0, returnPkt.Key
	}

	return returnPkt.ServiceInfo.IP, int(returnPkt.ServiceInfo.Port), returnPkt.Key
}

func (r *Requestor) invoke(request common.RequestPkt) *common.ReturnPkt {

	host, port, key := r.getServiceInfo(request.MethodName)

	log.Printf("Service %s on %s,%d", request.MethodName, host, port)
	log.Printf("Key to encrypt: %s", key)

	if key == "" {
		log.Printf("Autentication failed. Invalid credentials.")
		return new(common.ReturnPkt)
	} else if port < 1000 {
		log.Printf("Service not found.")
		return new(common.ReturnPkt)
	}

	r.crh = newClientRequestHandler(host, port)
	r.crh.connect()

	marshContent := r.marshaller.Marshall(request)
	keyData, _ := hex.DecodeString(key)

	encryptedContent := common.Encrypt(keyData, marshContent)

	if len(encryptedContent) == 0 {
		log.Printf("Failed Encrypting.")
		return new(common.ReturnPkt)
	}

	content := common.Request{Username: r.username, Data: encryptedContent}

	packet := r.marshaller.Marshall(content)

	r.crh.send(packet)

	marshRet := r.crh.receive()

	if len(marshRet) == 0 {
		log.Printf("Invalid received packet. Verify your requisition.")
		return new(common.ReturnPkt)
	}

	decrypted := common.Decrypt(keyData, marshRet)

	if len(decrypted) == 0 {
		log.Printf("Failed Decrypting.")
		return new(common.ReturnPkt)
	}

	var resPkt common.ReturnPkt
	r.marshaller.Unmarshall(decrypted, &resPkt)

	return &resPkt
}
