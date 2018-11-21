package server

import (
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"reflect"
	"time"

	"../common"
)

// MethodInfo ...
type MethodInfo struct {
	Method    reflect.Method
	ArgsType  []reflect.Type
	ReplyType reflect.Type
}

// Invoker handles the directing from ServerRequestHandler to the correct remote method
type Invoker struct {
	dnsAddr    string
	dnsPort    int
	localPort  int
	srh        *ServerRequestHandler
	marshaller *common.Marshaller
	methods    map[string]*MethodInfo
	object     interface{}
}

// NewInvoker creates a new invoker
func NewInvoker(object interface{}, localPort int, dnsAddr string, dnsPort int) *Invoker {

	inv := Invoker{
		dnsAddr:    dnsAddr,
		dnsPort:    dnsPort,
		localPort:  localPort,
		srh:        new(ServerRequestHandler),
		marshaller: new(common.Marshaller),
		object:     object}

	inv.registerMethods()
	inv.srh, _ = newServerRequestHandler(localPort)
	return &inv
}

func (i *Invoker) registerMethodInDNS(name string) {
	log.Printf("Attempt to connect to DNS (%s:%d)", i.dnsAddr, i.dnsPort)
	dnsSrh := newClientRequestHandler(i.dnsAddr, i.dnsPort)
	err := dnsSrh.connect()
	if err != nil {
		log.Fatalf("It was not possible to connect to %s:%d", i.dnsAddr, i.dnsPort)
	}
	localIP := GetOutboundIP()
	service := common.Service{Name: name, IP: localIP, Port: int32(i.localPort), AccessLevel: 1}
	log.Printf("Registering IP: %s | Port: %d", localIP, i.localPort)
	data := i.marshaller.Marshall(service)

	consult := common.ConsultPkt{ConsultType: "register", Data: data}

	pkt := i.marshaller.Marshall(consult)

	err = dnsSrh.send(pkt)
	if err != nil {
		log.Printf("Error sending package %s", err)
	}
	ret := dnsSrh.receive()
	returnPkt := new(common.RegisterResultPkt)
	err = i.marshaller.Unmarshall(ret, &returnPkt)
	if err != nil {
		log.Println("Error trying to register on DNS")
		return
	}
	log.Printf("%s:%d was successfully registered", returnPkt.IP, returnPkt.Port)
}

func (i *Invoker) registerMethods() {
	methods := make(map[string]*MethodInfo)

	objectType := reflect.TypeOf(i.object)

	for j := 0; j < objectType.NumMethod(); j++ {

		method := objectType.Method(j)
		methodName := method.Name
		methodType := method.Type

		argsType := []reflect.Type{}

		for k := 0; k < methodType.NumIn(); k++ {
			argsType = append(argsType, methodType.In(k))
		}

		if methodType.NumOut() != 1 {
			log.Printf("invoker.registerMethods: method %q has %d output parameters; needs exactly one\n", methodName, methodType.NumOut())
			continue
		}
		returnType := methodType.Out(0)
		i.registerMethodInDNS(methodName)
		methods[methodName] = &MethodInfo{Method: method, ArgsType: argsType, ReplyType: returnType}
	}
	i.methods = methods
}

func (i *Invoker) handleOperation(method reflect.Method, args []reflect.Value) (interface{}, error) {
	in := []reflect.Value{reflect.ValueOf(i.object)}
	for _, arg := range args {
		in = append(in, arg)
	}

	resultv := method.Func.Call(in)
	return resultv[0].Interface(), nil
}

func (i *Invoker) parseMethod(methodName string, argsI []interface{}) (reflect.Method, []reflect.Value, error) {
	methodInf := i.methods[methodName]
	args := []reflect.Value{}
	for _, i := range argsI {
		args = append(args, reflect.ValueOf(i))
	}

	if methodInf == nil {
		log.Fatalln("invoker.handleOperation: unknown request")
	}

	if len(args)+1 != len(methodInf.ArgsType) {
		log.Fatalf("invoker.handleOperation: request has %d parameters; needs exactly %d\n", len(args), len(methodInf.ArgsType))
	}
	return methodInf.Method, args, nil
}

func (i *Invoker) handleRequestPkt(requestPkt *common.RequestPkt) common.ReturnPkt {
	log.Printf("Received %s request from \"%s\"", requestPkt.MethodName, i.srh.remoteAddr)

	method, args, _ := i.parseMethod(requestPkt.MethodName, requestPkt.Args)

	ret, err := i.handleOperation(method, args)

	return common.ReturnPkt{MethodName: requestPkt.MethodName, ReturnValue: ret, Err: err}
}

func (i *Invoker) separatePacket(data []byte) *common.Request {
	req := new(common.Request)
	i.marshaller.Unmarshall(data, &req)
	return req
}

func (i *Invoker) getUserKey(request *common.Request) string {
	dnsSrh := newClientRequestHandler(i.dnsAddr, i.dnsPort)
	dnsSrh.connect()

	requestInfo := common.RequestInfo{Name: request.Username, Username: "", Password: ""}

	data := i.marshaller.Marshall(requestInfo)

	consultPkt := common.ConsultPkt{ConsultType: "consultname", Data: data}

	pkt := i.marshaller.Marshall(consultPkt)

	dnsSrh.send(pkt)

	ret := dnsSrh.receive()

	returnPkt := new(common.ConsultReturnPkt)

	i.marshaller.Unmarshall(ret, returnPkt)

	return returnPkt.Key

}

func (i *Invoker) createEmptyPacket() *common.ReturnPkt {
	return &common.ReturnPkt{MethodName: "", ReturnValue: nil, Err: errors.New("Invalid requisition")}
}

func (i *Invoker) sendEmptyResponse(srh *ServerRequestHandler, keyData []byte) {
	retPkt := i.createEmptyPacket()
	pkt := i.marshaller.Marshall(retPkt)
	encryptedContent := common.Encrypt(keyData, pkt)
	srh.send(encryptedContent)
}

func (i *Invoker) handleConnection(srh *ServerRequestHandler) {

	for {
		data, err := srh.receive()

		switch {
		case err == io.EOF:
			log.Printf("close this connection.\n   ---")
			i.srh.accept()
			continue
		case err != nil:
			log.Printf("\nError reading command. Got: \n", err)
			continue
		}

		req := i.separatePacket(data)

		key := i.getUserKey(req)

		if key == "" {
			log.Printf("No key found for user.")
			i.sendEmptyResponse(srh, []byte(""))
			i.srh.connection.Close()
			break
		}

		keyData, _ := hex.DecodeString(key)

		data = common.Decrypt(keyData, req.Data)

		if len(data) == 0 {
			log.Printf("Error decrypting message. Wrong key or data.")
			i.sendEmptyResponse(srh, []byte(""))
			i.srh.connection.Close()
			break
		}

		request := new(common.RequestPkt)

		if err == nil {
			err = i.marshaller.Unmarshall(data, &request)
		}

		switch {
		case err == io.EOF:
			log.Printf("close this connection.\n   ---")
			i.srh.accept()
			continue
		case err != nil:
			log.Printf("\nError reading command. Got: %err", err)
			continue
		}

		go func() {
			start := time.Now()
			returnPkt := i.handleRequestPkt(request)
			pkt := i.marshaller.Marshall(returnPkt)

			encryptedContent := common.Encrypt(keyData, pkt)

			if len(encryptedContent) == 0 {
				fmt.Printf("Failed encrypting message.")
				i.sendEmptyResponse(srh, []byte(""))
				i.srh.connection.Close()
				return
			}

			srh.send(encryptedContent)
			end := time.Now()
			log.Printf("%s - %.2f us", returnPkt, float64(end.Sub(start).Nanoseconds()/1000.))

		}()
	}

}

// Invoke invokes the invoker
func (i *Invoker) Invoke() {
	for {
		srh := i.srh
		srh.accept()
		go i.handleConnection(srh)
	}
}
