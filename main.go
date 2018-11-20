package main

import (
	"flag"
	"log"
	"strconv"

	"./src/client"
	"./src/dns"
	"./src/server"
)

func createServer() {
	accs := make(map[string]*Account)
	initialBalance := 1000.0
	accsNumber := []int{1, 2, 3, 4, 5, 6, 11, 435, 43232, 5}
	for i := 0; i < len(accsNumber); i++ {
		accID := "ACC" + strconv.Itoa(accsNumber[i]+1)
		accs[accID] = &Account{Balance: initialBalance}
	}
	accManager := AccountsManager{Accs: accs}

	invoker := server.NewInvoker(&accManager)
	invoker.Invoke()
}

func createClient() {
	proxy := client.NewProxy("ACC4", "pudim")
	balance := proxy.Withdraw("ACC4", 50.0)
	log.Printf("Balance: %s", balance)
}

func createDNS() {
	dnsServer := dns.NamingServer{}
	dnsServer.Start()
}

func main() {
	mwType := flag.String(
		"type",
		"",
		"Describes the middleware type to be initialized\n* Available options\n- client\n- server\n- dns")
	flag.Parse()

	switch *mwType {
	case "server":
		createServer()
	case "client":
		createClient()
	case "dns":
		createDNS()

	}
}
