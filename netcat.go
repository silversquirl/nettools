package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"syscall"
)

func errln(args ...interface{}) {
	fmt.Fprintln(os.Stderr, args...)
}

func unboxSysOpError(err error) syscall.Errno {
	if e2, ok := err.(*net.OpError); ok {
		if e3, ok := e2.Err.(*os.SyscallError); ok {
			if errno, ok := e3.Err.(syscall.Errno); ok {
				return errno
			}
		}
	}
	return 0
}

func chanRead(r io.Reader, ch chan []byte) {
	defer close(ch)
	for {
		buf := make([]byte, 1024)
		n, err := r.Read(buf)
		ch <- buf[:n]
		if err == io.EOF {
			break
		} else if e := unboxSysOpError(err); e == syscall.ECONNRESET || e == syscall.WSAECONNRESET {
			// The remote has closed the connection
			break
		} else if err != nil {
			errln("Error reading:", err)
			break
		}
	}
}

func connectToStdio(conn io.ReadWriter) {
	netCh := make(chan []byte)
	go chanRead(conn, netCh)
	ttyCh := make(chan []byte)
	go chanRead(os.Stdin, ttyCh)

	var buf []byte
	for {
		select {
		case buf = <-netCh:
			os.Stdout.Write(buf)
		case buf = <-ttyCh:
			conn.Write(buf)
		}

		// 0-length buffer means EOF or error
		if len(buf) == 0 {
			break
		}
	}
}

func acceptConnection(lis net.Listener) int {
	conn, err := lis.Accept()
	if err != nil {
		errln("Error accepting connection:", err)
		return 1
	}
	defer conn.Close()
	connectToStdio(conn)
	return 0
}

func Main() int {
	listenMode := flag.Bool("l", false, "Listen for connections")
	keepListening := flag.Bool("k", false, "Keep the socket open after connections are closed. Ignored without -l.")
	flag.Parse()
	if flag.NArg() == 0 {
		errln("Please specify an address to connect to")
		return 1
	}
	addr := flag.Arg(0)

	if *listenMode {
		lis, err := net.Listen("tcp", addr)
		if err != nil {
			errln("Error opening socket:", err)
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
			errln("Error connecting:", err)
			return 1
		}
		defer conn.Close()
		connectToStdio(conn)
	}
	return 0
}

func main() {
	os.Exit(Main())
}
