package main

import (
	"bytes"
	"encoding/gob"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	t "time"
)

const (
	PASS = iota
	OPEN
)

const (
	MAXTIME = 3
)

var (
	//The IP and port of this node format "ip:port"
	ipPort string
	//connection to send and receive on
	conn *net.UDPConn
	//Hard coded cluster, this nodes ipPort must occur in the cluster
	pl = []string{"localhost:19000", "localhost:19001", "localhost:19002"}
	//Log output for this node
	logger *log.Logger
	//A collection of pointers to peers locations
	peers []*net.UDPAddr
	//Ball is a ball being thrown between peers
	ball bool
	//Open is true if the peer is ready for a pass
	open bool
	//catches is the number of times this peer knows the ball has been
	//caught
	catches int
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
		logger.Fatalf("-ipPort of %s no within hard coded cluster; cluster = %s\n", ipPort, pl)
	}
}

func main() {
	Init()
	readMsg, writeMsg := make(chan Msg, 1), make(chan Msg, 1)
	//Start async reading and writing functions
	go Send(conn, writeMsg)
	go Listen(conn, readMsg)
	openTimer, passTimer := make(<-chan t.Time), make(<-chan t.Time)
	//start random openTimer
	openTimer = t.After(t.Duration((rand.Int() % MAXTIME)) * t.Second)
	//hardcode who starts with the ball
	if ipPort == "localhost:19000" {
		ball = true
		passTimer = t.After(t.Duration((rand.Int()%MAXTIME)*5) * t.Second)
	}

	var message Msg
	for true {
		select {
		case m := <-readMsg:
			switch m.Type {
			case OPEN:
				if ball == true {
					message.Type = PASS
					message.Catches = catches
					message.addr = m.addr
					ball = false
					openTimer = t.After(t.Duration((rand.Int() % MAXTIME)) * t.Second)
					passTimer = nil
					logger.Printf("Okay %s here comes the ball\n", m.addr)
					writeMsg <- message
				}
			case PASS:
				if open == true {
					catches = m.Catches + 1
					ball = true
					open = false
					passTimer = t.After(t.Duration((rand.Int()%MAXTIME)*5) * t.Second)
					logger.Printf("Thanks for the pass %s I caught the ball for the %d time\n", m.addr.String(), catches)
				}
			}
		case <-openTimer:
			open = true
			message.Type = OPEN
			message.Catches = 0
			logger.Printf("Hey everyone Im open pass the ball\n")
			for _, peer := range peers {
				message.addr = peer
				writeMsg <- message
			}
		case <-passTimer:
			message.Type = PASS
			message.Catches = catches
			message.addr = peers[rand.Int()%(len(peers)-1)]
			logger.Printf("I'm done with this ball %s catch it", message.addr.String())
			writeMsg <- message
		}
	}

}

//Listen is an async function for reading messages from the network.
//Incomming messages which are correctly read, and decoded are passed
//to the main loop via a buffered channel.
func Listen(conn *net.UDPConn, readMsg chan Msg) {
	for true {
		var message Msg
		var network bytes.Buffer
		buf := make([]byte, 1024)
		dec := gob.NewDecoder(&network)
		n, addr, err := conn.ReadFromUDP(buf)
		if err != nil {
			logger.Printf("Unable to Read From UDP: Error %s\n", err)
			continue
		}
		//transfer from byte array to buffer because the canonical
		//io.Copy(&network,conn) requires size checks.
		network.Write(buf[0:n])
		err = dec.Decode(&message)
		if err != nil {
			logger.Printf("Cannot decode: Error %s\n", err)
			continue
		}
		message.addr = addr
		//logger.Println(message.String())
		readMsg <- message
	}
}

//Send is an async function for writing message responses to the
//network. Send waits on a buffered channel, and casts response
//strings to arrays of bytes.
func Send(conn *net.UDPConn, writeMsg chan Msg) {
	for true {
		var network bytes.Buffer
		enc := gob.NewEncoder(&network)
		message := <-writeMsg
		err := enc.Encode(message)
		if err != nil {
			logger.Printf(err.Error())
		}
		conn.WriteToUDP(network.Bytes(), message.addr)
		//network.Reset()
	}
}

func Init() {
	logger = log.New(os.Stdout, "[Initalizing]", log.Lshortfile)
	setFlags()
	checkArgs()
	peers = setupPeers(pl)
	conn = getConnection(ipPort)
	logger = log.New(os.Stdout, "["+ipPort+"]", log.Lshortfile)
}

//setupPeers takes a list of ip:ports and returns a map of their UDP
//addresses. getPeers fails if the address cannot be resloved.
func setupPeers(peerList []string) []*net.UDPAddr {
	P := make([]*net.UDPAddr, 0)
	for i := range peerList {
		//do not set up yourself as a peer
		if peerList[i] == ipPort {
			continue
		}
		peer := getAddr(peerList[i])
		P = append(P, peer)
	}
	logger.Println(P)
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

//get connection returns a udp conn which the server listens on.
func getConnection(ip string) *net.UDPConn {
	lAddr, err := net.ResolveUDPAddr("udp", ip)
	if err != nil {
		logger.Fatal(err)
	}
	l, err := net.ListenUDP("udp", lAddr)
	if err != nil {
		logger.Fatal(err)
	}
	return l
}

type Msg struct {
	Catches int
	Type    int
	addr    *net.UDPAddr
}

func (m Msg) String() string {
	return fmt.Sprintf("Catches :%d, Type %d, Address %s", m.Catches, m.Type, m.addr.String())
}
