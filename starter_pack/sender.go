package main

import (
	"fmt"
	"io/ioutil"
	"net"
)

func sender(filename *string, conn *net.UDPConn) int {
	var seqno uint16 = 0
	data, err := ioutil.ReadFile(*filename)
	if err != nil {
		fmt.Println("Error, cannot read file.")
		return 2
	}

	var pkt *Packet

	for start := 0; start < len(data); start += 255 {
		end := start + 255
		if end > len(data) {
			end = len(data)
		}

		pkt = make_data_pkt(data[start:end], seqno+255)

		// TODO: send DATA and get ACK
		n, _ := send(pkt, conn, nil)

		rcv, _, ok := recv(conn, n)
		if ok && isACK(rcv, rcv.hdr.ackno) {
			fmt.Println("success...")
		}
	}
	// TODO: send FIN and get FINACK
	pkt = make_fin_pkt(seqno)
	n, _ := send(pkt, conn, nil)

	rcv, _, ok := recv(conn, n)
	if !(ok && isACK(rcv, seqno)) {
		return 3
	}

	// TODO: return 0 for success, 3 for failure
	return 0
}

func make_data_pkt(data []byte, seqno uint16) *Packet {
	var pkt *Packet = &Packet{}

	pkt.hdr.flag = DATA
	pkt.hdr.seqno = seqno
	pkt.hdr.len = uint8(len(data))
	pkt.dat = data

	return pkt
}

func make_fin_pkt(seqno uint16) *Packet {
	var pkt *Packet = &Packet{}

	pkt.hdr.seqno = seqno
	pkt.hdr.flag = FIN

	return pkt
}

func isACK(pkt *Packet, expected uint16) bool {
	// TODO: return true if ACK (including FINACK) and ackno is what is expected
	isAck := false

	if (pkt.hdr.flag == ACK || pkt.hdr.flag == FINACK) && pkt.hdr.ackno == expected {
		isAck = true
	}

	return isAck
}
