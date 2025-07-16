package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
)

func appendStringToBytes(b []byte, str string) []byte {
	/**
	[]byte(str)... -> Tells the append function to add each element of the __str__ variable individually
		    	  	  rather that appending the entire slice (as the JavaScript spread operator does)
	*/
	return append(b, []byte(str)...)
}

func writeGzipBuffer(s string) (string, error) {
	var buffer bytes.Buffer
	writer := gzip.NewWriter(&buffer)
	writeOp, err := writer.Write([]byte(s))
	if err != nil {
		return "", err
	}
	err = writer.Close()
	if err != nil {
		return "", err
	}
	fmt.Printf("Wrote %d bytes of gzip compressed data", writeOp)

	return buffer.String(), nil
}
