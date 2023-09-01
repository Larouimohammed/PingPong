package main

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"time"
	"os"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

func main() {

	PingPong()
	//fmt.Println(reader.Reader())
}
func PingPong() {
	//icmp ipv4 listen
	host := os.Args[1]
	ange := os.Args[2]
	int, _ := strconv.Atoi(ange)

	packetconn, err := icmp.ListenPacket("ip4:1", "")
	if err != nil {
		log.Fatal(err)
	}
	defer packetconn.Close()
	//send icmp message echo
	msg := &icmp.Message{
		Type: ipv4.ICMPTypeEchoReply,
		Code: 0,
		Body: &icmp.Echo{
			ID:   os.Getpid() & 0xffff,
			Seq:  1,
			Data: []byte("hello"),
		},
	}
	for i := 0; i < int; i++ {

		wb, err := msg.Marshal(nil)

		//fmt.Println(string(wb))
		if err != nil {
			log.Fatal(err)
		}

		start := time.Now()
		if _, err := packetconn.WriteTo(wb, &net.IPAddr{IP: net.ParseIP(host)}); err != nil {
			log.Fatal(err)
		}

		//receive icmp messgae reply
		rb := make([]byte, 1500)
		err = packetconn.SetReadDeadline(time.Now().Add(5 * time.Second))
		n, peer, err := packetconn.ReadFrom(rb)
		duration := time.Since(start)
		//fmt.Println(string(rb))
		if err != nil {
			log.Fatal(err)
		}

		rm, err := icmp.ParseMessage(1, rb[:n])
		if err != nil {
			log.Fatal(err)
		}

		switch rm.Type {
		case ipv4.ICMPTypeEchoReply:

			fmt.Printf("received %v Bytes from %v seq=%v in time=%v MS\n", rm.Body.Len(1), peer, rm.Body.(*icmp.Echo).Seq, duration.Microseconds())

		default:
			fmt.Printf("Failed: %+v\n", rm)
		}
		msg.Body.(*icmp.Echo).Seq = msg.Body.(*icmp.Echo).Seq + 1

	}
}
