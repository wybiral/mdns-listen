package main

import (
	"encoding/json"
	"log"
	"net"
	"os"

	"github.com/wybiral/mdns-listen/packet"
)

const bufferSize = 8 * 1024

func main() {
	// Resolve UDP address
	addr, err := net.ResolveUDPAddr("udp", "224.0.0.251:5353")
	if err != nil {
		log.Fatal(err)
	}
	// Create multicast listener
	c, err := net.ListenMulticastUDP("udp", nil, addr)
	if err != nil {
		log.Fatal(err)
	}
	buffer := make([]byte, bufferSize)
	enc := json.NewEncoder(os.Stdout)
	for {
		// Read packet
		n, from, err := c.ReadFrom(buffer)
		if err != nil {
			log.Fatal(err)
		}
		// Parse packet
		p, err := packet.Parse(buffer[:n], from)
		if err != nil {
			log.Fatal(err)
		}
		// Encode packet
		err = enc.Encode(p)
		if err != nil {
			log.Fatal(err)
		}
	}
}
