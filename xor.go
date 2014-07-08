package main

import (
	"flag"
	//	"fmt"
	"io"
    "log"
	"os"
)

var outputFlag = flag.String("output", "", "Output")

//var stringFlag = flag.String("string", "", "String")

func init() {
	flag.StringVar(outputFlag, "o", "-", "Output")
	//	flag.StringVar(stringFlag, "i", "", "String")
}

func main() {

	flag.Parse()
	var filesFlag []string = flag.Args()
	//fmt.Printf("stringflag: %v", *stringFlag)

	if len(filesFlag) < 2 {
		log.Fatal("Need at least two files to XOR")
	}

	// setup file inputs
	fds := make([]*os.File, len(filesFlag))
	for i, filename := range filesFlag {
		fd, err := os.Open(filename)
		if err != nil {
			log.Fatal(err)
		}
		//fmt.Printf("%v arg %v", i, filename)
		fds[i] = fd
		// close fi on exit and check for its returned error
		defer func() {
			if err := fd.Close(); err != nil {
				log.Fatal(err)
			}
		}()
	}

	// setup output
	var fo *os.File
	if *outputFlag == "-" {
		fo = os.Stdout
	} else {
		fd, err := os.Create(*outputFlag)
		if err != nil {
			log.Fatal(err)
		}
		fo = fd
	}
	// close fo on exit and check for its returned error
	defer func() {
		if err := fo.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	// make a buffer to keep chunks that are read

	bufs := make([][]byte, len(filesFlag))
	for i, _ := range filesFlag {
		bufs[i] = make([]byte, 1024)
	}
	bytesread := make([]int, len(filesFlag))

	for {
		for i, _ := range filesFlag {
			n, err := fds[i].Read(bufs[i])
			if err != nil && err != io.EOF {
				log.Fatal(err)
			}
			bytesread[i] = n
		}

		nout := smallest(bytesread)

		bufo := make([]byte, 1024)
		bufo = xorBytes(bufs)

		if nout == 0 {
			break
		}

		// write a chunk
		if _, err := fo.Write(bufo[:nout]); err != nil {
			log.Fatal(err)
		}
	}
}

func xorChannels(in1 <-chan int, in2 <-chan int, out chan int) {

	for {
		sum := 0
		select {
		case sum = <-in1:
			sum += <-in2
		case sum = <-in2:
			sum += <-in1
		}
		out <- sum
	}
}

func smallest(i []int) int {
	j := i[0]
	for _, m := range i {
		if m < j {
			j = m
		}
	}
	return j
}

func xorBytes(b [][]byte) []byte {

	b_len := len(b[0])
	for _, m := range b {
		if len(m) != b_len {
			log.Fatal("length mismatch!")
		}
	}
	br := make([]byte, b_len)
	for i := range b[0] {
		br[i] = 0
		for _, m := range b {
			br[i] = br[i] ^ m[i]
		}
	}
	return br
}
