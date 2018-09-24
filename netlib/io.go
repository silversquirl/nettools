package netlib

import (
	"fmt"
	"io"
	"os"
)

type multiReadWriter struct {
	r io.Reader
	w io.Writer
}

func (rw multiReadWriter) Read(p []byte) (n int, err error) {
	return rw.r.Read(p)
}

func (rw multiReadWriter) Write(p []byte) (n int, err error) {
	return rw.w.Write(p)
}

func MultiReadWriter(r io.Reader, w io.Writer) io.ReadWriter {
	return multiReadWriter{r, w}
}

var Stdio = MultiReadWriter(os.Stdin, os.Stdout)

func Errln(args ...interface{}) {
	fmt.Fprintln(os.Stderr, args...)
}

func chanRead(r io.Reader, ch chan []byte) {
	defer close(ch)
	for {
		buf := make([]byte, 1024)
		n, err := r.Read(buf)
		ch <- buf[:n]
		if err == io.EOF {
			break
		} else if ConnResetError(err) {
			// The remote has closed the connection
			break
		} else if err != nil {
			Errln("Error reading:", err)
			break
		}
	}
}

func ConnectBidirectional(rw1, rw2 io.ReadWriter) {
	ch1 := make(chan []byte)
	go chanRead(rw1, ch1)
	ch2 := make(chan []byte)
	go chanRead(rw2, ch2)

	var buf []byte
	for {
		select {
		case buf = <-ch1:
			rw2.Write(buf)
		case buf = <-ch2:
			rw1.Write(buf)
		}

		// 0-length buffer means EOF or error
		if len(buf) == 0 {
			break
		}
	}
}
