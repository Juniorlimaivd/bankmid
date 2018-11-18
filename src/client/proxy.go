package client

import (
	"reflect"

	"../common"
)

// Proxy ...
type Proxy struct {
	host      string
	port      int
	requestor *Requestor
}

func NewProxy(host string, port int) *Proxy {
	return &Proxy{
		host:      host,
		port:      port,
		requestor: newRequestor(host, port),
	}
}

func (p *Proxy) GetBalance(accountID string) float64 {
	reqPkt := common.NewRequestPkt("GetBalance", accountID)
	retPkt := p.requestor.invoke(reqPkt)
	return reflect.ValueOf(retPkt.ReturnValue).Float()
}

func (p *Proxy) Withdraw(accountID string, amount float64) string {
	reqPkt := common.NewRequestPkt("Withdraw", accountID, amount)
	retPkt := p.requestor.invoke(reqPkt)
	return reflect.ValueOf(retPkt.ReturnValue).String()
}

func (p *Proxy) Deposit(accountID string, amount float64) string {
	reqPkt := common.NewRequestPkt("Deposit", accountID, amount)
	retPkt := p.requestor.invoke(reqPkt)
	return reflect.ValueOf(retPkt.ReturnValue).String()
}

func (p *Proxy) Transfer(payerID string, payeeID string, amount float64) string {
	reqPkt := common.NewRequestPkt("Transfer", payerID, payeeID, amount)
	retPkt := p.requestor.invoke(reqPkt)
	return reflect.ValueOf(retPkt.ReturnValue).String()
}
