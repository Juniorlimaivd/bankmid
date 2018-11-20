package client

import (
	"reflect"

	"../common"
)

// Proxy ...
type Proxy struct {
	requestor *Requestor
}

func NewProxy(user string, password string) *Proxy {
	return &Proxy{
		requestor: newRequestor(user, password),
	}
}

func (p *Proxy) GetBalance(accountID string) float64 {
	reqPkt := common.NewRequestPkt("GetBalance", accountID)
	retPkt := p.requestor.invoke(reqPkt)
	if retPkt.ReturnValue == nil {
		return 0.0
	}
	return reflect.ValueOf(retPkt.ReturnValue).Float()
}

func (p *Proxy) Withdraw(accountID string, amount float64) string {
	reqPkt := common.NewRequestPkt("Withdraw", accountID, amount)
	retPkt := p.requestor.invoke(reqPkt)
	if retPkt.ReturnValue == nil {
		return "Operation failed"
	}
	return reflect.ValueOf(retPkt.ReturnValue).String()
}

func (p *Proxy) Deposit(accountID string, amount float64) string {
	reqPkt := common.NewRequestPkt("Deposit", accountID, amount)
	retPkt := p.requestor.invoke(reqPkt)
	if retPkt.ReturnValue == nil {
		return "Operation failed"
	}
	return reflect.ValueOf(retPkt.ReturnValue).String()
}

func (p *Proxy) Transfer(payerID string, payeeID string, amount float64) string {
	reqPkt := common.NewRequestPkt("Transfer", payerID, payeeID, amount)
	retPkt := p.requestor.invoke(reqPkt)
	if retPkt.ReturnValue == nil {
		return "Operation failed"
	}
	return reflect.ValueOf(retPkt.ReturnValue).String()
}
