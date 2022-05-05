package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
)

func receiver(filename *string, conn *net.UDPConn) int {
	var expected uint16 = 0
	var pkt *Packet
	// recieve
	for {
		// TODO: receive DATA and send ACK if exepcted seqno arrives
		// NOTE: Don't forget to write the data
		// NOTE: You'll need the addr returned from recv in order to
		// send back to the sender.

		rcv, addr, ok := recv(conn, 0)

		//write to file
		f, err := os.OpenFile(*filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			panic(err)
		}

		defer f.Close()

		if _, err = f.WriteString(string(rcv.dat)); err != nil {
			panic(err)
		}

		pkt = make_ack_pkt(rcv.hdr.ackno)

		if rcv.hdr.seqno == expected && ok {
			//test print
			fmt.Println("seqno: " + strconv.Itoa(int(rcv.hdr.seqno)))

			//send ack
			send(pkt, conn, addr)
		}

		// TODO: break out of infinte loop after FINACK
		//if not corrupted and seqno == expected
		if ok && rcv.hdr.seqno == expected {
			pkt := make_finack_pkt(expected)
			send(pkt, conn, addr)
			break
		}
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
