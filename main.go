package main

import (
	"flag"
	"fmt"
)

func main() {
	mode := flag.String("mode", "recv", "Set to 'send' or 'recv' to send or receive, respectively.")
	flag.Parse()
	if *mode == "recv" {
		recv()
	} else if *mode == "send" {
		send()
	} else {
		fmt.Println("Unrecognized mode: ", *mode)
	}
}
