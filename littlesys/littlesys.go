package main

import (
	"flag"
	"log"
	"net"
	"os"
)

var (
	//The IP and port of this node format "ip:port"
	ipPort string
	//Hard coded cluster, this nodes ipPort must occur in the cluster
	pl = []string{"localhost:19000", "localhost:19001", "localhost:19002"}
	//Log output for this node
	logger *log.Logger
	//A collection of pointers to peers locations
	peers map[string]*net.UDPAddr
)

func setFlags() {
	flag.StringVar(&ipPort, "ipPort", "", "-ipPort is the ip and port of this node; format ip:port")
	flag.Parse()
}

//check if the arguments meet the specification if not exit
func checkArgs() {
	//check if the the ip port provided is in the cluster
	found := false
	for _, ipp := range pl {
		if ipp == ipPort {
			found = true
		}
	}
	if !found {
		logger.Fatalf("-ipPort no within hard coded cluster; cluster = %s\n", pl)
	}
}

func main() {
	Init()
}

func Init() {
	logger = log.New(os.Stdout, "[Initalizing] ", log.Lshortfile)
	setFlags()
	checkArgs()
	peers = setupPeers(pl)
}

//setupPeers takes a list of ip:ports and returns a map of their UDP
//addresses. getPeers fails if the address cannot be resloved.
func setupPeers(peerList []string) map[string]*net.UDPAddr {
	P := make(map[string]*net.UDPAddr, len(peerList))
	for i := range peerList {
		if peerList[i] == ipPort {
			continue
		}
		peer := getAddr(peerList[i])
		P[peer.String()] = peer
	}
	return P
}

//get Adder resolves an ip:port to a udp address. It fails if the
//address cannot be resloved.
func getAddr(address string) *net.UDPAddr {
	rAddr, errR := net.ResolveUDPAddr("udp", address)
	if errR != nil {
		logger.Fatal(errR)
	}
	return rAddr
}
