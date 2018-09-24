package main

import (
	"flag"
	"net"
	"os"

	"github.com/vktec/nettools/netlib"
)

func acceptConnectionPair(l1, l2 net.Listener) int {
	c1, err := l1.Accept()
	if err != nil {
		netlib.Errln("Error accepting connection:", err)
		return 1
	}
	defer c1.Close()

	c2, err := l2.Accept()
	if err != nil {
		netlib.Errln("Error accepting connection:", err)
		return 1
	}
	defer c2.Close()

	netlib.ConnectBidirectional(c1, c2)
	return 0
}

func Main() int {
	listenMode := flag.Bool("l", false, "Listen for connections")
	keepListening := flag.Bool("k", false, "Keep the socket open after connections are closed. Ignored without -l.")
	flag.Parse()
	if flag.NArg() < 2 {
		netlib.Errln("Please specify addresses")
		return 1
	}
	a1 := flag.Arg(0)
	a2 := flag.Arg(1)

	if *listenMode {
		l1, err := net.Listen("tcp", a1)
		if err != nil {
			netlib.Errln("Error opening socket:", err)
			return 1
		}
		defer l1.Close()

		l2, err := net.Listen("tcp", a2)
		if err != nil {
			netlib.Errln("Error opening socket:", err)
			return 1
		}
		defer l2.Close()

		if *keepListening {
			for {
				if ret := acceptConnectionPair(l1, l2); ret != 0 {
					return ret
				}
			}
		} else {
			return acceptConnectionPair(l1, l2)
		}
	} else {
		c1, err := net.Dial("tcp", a1)
		if err != nil {
			netlib.Errln("Error connecting:", err)
			return 1
		}
		defer c1.Close()

		c2, err := net.Dial("tcp", a2)
		if err != nil {
			netlib.Errln("Error connecting:", err)
			return 1
		}
		defer c2.Close()

		netlib.ConnectBidirectional(c1, c2)
	}
	return 0
}

func main() {
	os.Exit(Main())
}
