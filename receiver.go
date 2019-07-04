package main

import (
	"flag"
	"fmt"
	"net"
	"os"
)

func recv() {
	host := flag.Arg(0)
	recvPort := flag.Arg(1)
	sendPort := flag.Arg(2)
	fileName := flag.Arg(3)
	type PacketStatus struct {
		p    *Packet
		good bool
	}
	window := make([]PacketStatus, WindowSize)
	windowStart := int32(0)
	f, err := os.Create(fileName)
	if err != nil {
		fmt.Printf("Error opening output file: %v\n", err)
		return
	}
	defer f.Close()
	fmt.Println("Starting receiver on ", host+":"+recvPort)
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
		fmt.Println("[RECV] Got packet: ", p)
		if p.Type == Eot {
			conn, err := net.Dial("udp", host+":"+sendPort)
			defer conn.Close()
			if err != nil {
				fmt.Printf("Dial err: %v\n", err)
				return
			}
			conn.Write(p.GetBytes())
			fmt.Println("Done receiving file")
			return
		}
		if p.Type == Data && inWindow(windowStart, p.Seqnum) {
			offset := (p.Seqnum + WindowSize - windowStart) % WindowSize
			window[offset].good = true
			window[offset].p = p
			if offset == 0 {
				// Move window
				for window[0].good {
					f.Write(window[0].p.Data)
					window = append(window[1:], PacketStatus{}) // Shift left
					windowStart++
				}
				p := &Packet{
					Type:   Ack,
					Seqnum: windowStart - 1,
				}
				conn, err := net.Dial("udp", host+":"+sendPort)
				defer conn.Close()
				if err != nil {
					fmt.Printf("Dial err: %v\n", err)
					return
				}
				conn.Write(p.GetBytes())
			}
			fmt.Println("Window:", windowStart)
			for i := int32(0); i < WindowSize; i++ {
				if window[i].good {
					fmt.Printf("%d\t", window[i].p.Seqnum)
				} else {
					fmt.Print("-\t")
				}
			}
			fmt.Println("")
		}
	}
}
