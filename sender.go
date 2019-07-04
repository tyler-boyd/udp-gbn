package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"time"
)

const (
	// WindowSize int
	WindowSize int32 = 10
	// PacketSize int
	PacketSize int32 = 500
)

func inWindow(windowStart int32, seqNum int32) bool {
	return (seqNum >= windowStart && seqNum < windowStart+WindowSize) ||
		(seqNum < windowStart && windowStart+WindowSize >= 32 && seqNum < (windowStart+WindowSize)%32)
}

func send() {
	host := flag.Arg(0)
	sendPort := flag.Arg(1)
	recvPort := flag.Arg(2)
	fileName := flag.Arg(3)
	windowStart := int32(0)
	fmt.Println("Starting sender")

	// Get ready to send
	conn, err := net.Dial("udp", host+":"+sendPort)
	if err != nil {
		fmt.Printf("Send UDP dial error: %v\n", err)
		return
	}
	defer conn.Close()

	// Read file
	contents, err := ioutil.ReadFile(fileName)
	n := int32(len(contents))
	nPackets := int32((n + PacketSize - 1) / PacketSize)

	// Prepare packets
	packets := make([]Packet, nPackets)
	for i := int32(0); i < nPackets; i++ {
		packets[i].Type = Data
		packets[i].Seqnum = int32(i) % 32
		packets[i].Data = contents[i*PacketSize : (i+1)*PacketSize]
		if (i+1)*PacketSize > int32(len(contents)) {
			packets[i].Data = contents[i*PacketSize:]
		}
	}

	// Ack receiver
	ackChan := make(chan int32)
	eotChan := make(chan int32)
	go func() {
		// Spawn new thread and listen for ACKs
		pc, err := net.ListenPacket("udp", host+":"+recvPort)
		if err != nil {
			fmt.Printf("Some error: %v\n", err)
			return
		}
		defer pc.Close()

		buffer := make([]byte, 512)

		for {
			n, _, err := pc.ReadFrom(buffer)
			if err != nil {
				return
			}

			p, err := NewPacket(buffer[:n])
			if err != nil {
				fmt.Printf("Some error: %v\n", err)
				return
			}
			// fmt.Println("[AckRec] Got packet: ", p)
			if p.Type == Ack && inWindow(windowStart, p.Seqnum) {
				ackChan <- p.Seqnum
			}
			if p.Type == Eot {
				eotChan <- 0
			}
		}
	}()

	// Wait for ack receiver to start
	time.Sleep(time.Millisecond * 500)

	// Send initial packets
	for i := int32(0); i < WindowSize && i < nPackets; i++ {
		fmt.Println("[SEND] Sending ", packets[i].Seqnum)
		conn.Write(packets[i].GetBytes())
	}

	// Wait for acks/timeouts
outer:
	for {
		select {
		case ack := <-ackChan:
			fmt.Println("[SEND] Got ack", ack)
			toSend := (1 + ack - windowStart) % 32
			if toSend < 0 {
				toSend = -toSend
			}
			windowStart += toSend
			for i := WindowSize - toSend; i < WindowSize; i++ {
				if i+windowStart >= nPackets || i+windowStart == 3 {
					break
				}
				p := &packets[i+windowStart]
				fmt.Println("[SEND] Sending ", p.Seqnum)
				conn.Write(p.GetBytes())
			}
			if windowStart >= nPackets {
				fmt.Println("Done sending file")
				break outer
			}
		case <-time.After(time.Millisecond * 1000):
			fmt.Println("Timeout")
			for i := windowStart; i < WindowSize+windowStart && i < nPackets; i++ {
				fmt.Println("Re-sending packet ", packets[i].Seqnum)
				conn.Write(packets[i].GetBytes())
			}
		}
	}

	// Send Eot
	p := &Packet{
		Type:   Eot,
		Seqnum: windowStart,
	}
	conn.Write(p.GetBytes())
	// Wait for Eot
	<-eotChan
}
