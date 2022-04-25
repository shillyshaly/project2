package main

import (
	"fmt"
	"io/ioutil"
	"net"
)

func receiver(filename *string, conn *net.UDPConn) int {
	var expected uint16 = 0
	var pkt *Packet
	// recieve
	for {
		// TODO: receive DATA and send ACK if exepcted seqno arrives
		// NOTE: Don't forget to write the data - DONE
		// NOTE: You'll need the addr returned from recv in order to
		// send back to the sender.
		rcv, addr, ok := recv(conn, 0)

		err := ioutil.WriteFile(*filename, rcv.dat, 0)
		if err != nil {
			fmt.Println("Error, cannot read file.")
			return 2
		}

		if rcv.hdr.seqno == expected {
			pkt = make_ack_pkt(rcv.hdr.ackno)
		}
		send(pkt, conn, addr)

		// TODO: break out of infinte loop after FINACK
	}

	return 0
}

func make_ack_pkt(ackno uint16) *Packet {
	var pkt *Packet = &Packet{}

	pkt.hdr.ackno = ackno
	pkt.hdr.flag = ACK

	return pkt
}

func make_finack_pkt(ackno uint16) *Packet {
	var pkt *Packet = &Packet{}

	pkt.hdr.flag = FINACK
	pkt.hdr.ackno = ackno

	return pkt
}
