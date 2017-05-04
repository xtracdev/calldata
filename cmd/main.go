package main

import (
	"compress/gzip"
	"fmt"
	"github.com/xtracdev/calldata"
	"github.com/xtracdev/calldata/apitimings"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("usage: " + os.Args[0] + " <file path>")
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err.Error())
	}
	defer file.Close()

	var reader io.Reader = file

	gzipReader, err := gzip.NewReader(file)
	if err == nil {
		fmt.Println("assuming gzip encoding")
		reader = gzipReader
		defer gzipReader.Close()
	} else {
		fmt.Println("Error creating gzip reader... read as uncompressed", err.Error())
	}

	sr, err := sumoreader.NewSumoReader(reader)
	if err != nil {
		log.Fatal(err)
	}

	var count = 1
	for sr.Scan() {
		line := sr.Text()
		if strings.Contains(line, "{") {
			//fmt.Println(sr.Text())
			//fmt.Println("...create api timings record")
			at, err := apitimings.NewAPITimingRec(line)
			if err != nil {
				fmt.Println("not a timing record", err.Error())
				continue
			}
			//fmt.Println("create call record")
			cr, err := at.CallRecord()
			if err != nil {
				fmt.Println("Not a call record", err.Error())
				continue
			}
			fmt.Println(cr)
			calls, err := at.ServiceCalls()

			if err != nil {
				//fmt.Println(err.Error())
				continue
			}

			for _, c := range calls {
				fmt.Println(c)
			}
		}
		count++
	}

	if err := sr.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("read", sr.Lines())
}
