package main

import (
	"flag"
	"log"
	"os"
	"strconv"

	"./src/client"
	"./src/dns"
	"./src/server"
)

func createServer(localPort int, dnsAddr string, dnsPort int) {
	accs := make(map[string]*Account)
	initialBalance := 1000.0
	accsNumber := []int{1, 2, 3, 4, 5, 6, 11, 435, 43232, 5}
	for i := 0; i < len(accsNumber); i++ {
		accID := "ACC" + strconv.Itoa(accsNumber[i]+1)
		accs[accID] = &Account{Balance: initialBalance}
	}
	accManager := AccountsManager{Accs: accs}

	invoker := server.NewInvoker(&accManager, localPort, dnsAddr, dnsPort)
	invoker.Invoke()
}

func createClient(dnsAddr string, dnsPort int) {
	proxy := client.NewProxy("ACC4", "pudim", dnsAddr, dnsPort)
	balance := proxy.Withdraw("ACC4", 50.0)
	log.Printf("Balance: %s", balance)
}

func createDNS(port int) {
	dnsServer := dns.NamingServer{}
	dnsServer.Start(port)
}

func main() {
	var err error

	mwType := flag.String(
		"type",
		"",
		"Describes the middleware type to be initialized\n* Available options\n- client\n- server\n- dns")
	port := flag.Int("port", -1, "")
	dnsAddr := flag.String("dnsAddr", "", "")
	dnsPort := flag.Int("dnsPort", -1, "")
	flag.Parse()

	if *mwType == "" {
		*mwType = os.Getenv("MW_TYPE")
		if *mwType == "" {
			log.Fatalln("Type is required")
		}
	}

	if *port == -1 {
		*port, err = strconv.Atoi(os.Getenv("PORT"))
		if err != nil {
			*port = 5000
		}
	}

	if *dnsAddr == "" {
		*dnsAddr = os.Getenv("DNS_ADDR")
		if *dnsAddr == "" {
			*dnsAddr = "localhost"
		}
	}

	if *dnsPort == -1 {
		*dnsPort, err = strconv.Atoi(os.Getenv("DNS_PORT"))
		if err != nil {
			*dnsPort = 80
		}
	}

	switch *mwType {
	case "server":
		createServer(*port, *dnsAddr, *dnsPort)
	case "client":
		createClient(*dnsAddr, *dnsPort)
	case "dns":
		createDNS(*port)
	}
}
