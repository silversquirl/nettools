package main

import (
	"flag"
	"net"
	"os"

	"github.com/vktec/nettools/netlib"
)

func acceptConnection(lis net.Listener) int {
	conn, err := lis.Accept()
	if err != nil {
		netlib.Errln("Error accepting connection:", err)
		return 1
	}
	defer conn.Close()
	netlib.ConnectBidirectional(conn, netlib.Stdio)
	return 0
}

func Main() int {
	listenMode := flag.Bool("l", false, "Listen for connections")
	keepListening := flag.Bool("k", false, "Keep the socket open after connections are closed. Ignored without -l.")
	flag.Parse()
	if flag.NArg() < 1 {
		netlib.Errln("Please specify an address")
		return 1
	}
	addr := flag.Arg(0)

	if *listenMode {
		lis, err := net.Listen("tcp", addr)
		if err != nil {
			netlib.Errln("Error opening socket:", err)
			return 1
		}
		defer lis.Close()

		if *keepListening {
			for {
				if ret := acceptConnection(lis); ret != 0 {
					return ret
				}
			}
		} else {
			return acceptConnection(lis)
		}
	} else {
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			netlib.Errln("Error connecting:", err)
			return 1
		}
		defer conn.Close()
		netlib.ConnectBidirectional(conn, netlib.Stdio)
	}
	return 0
}

func main() {
	os.Exit(Main())
}
