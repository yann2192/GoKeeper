package gokeeper

import (
	"bytes"
	"io"
	"os"
)

func readAll(f *os.File) ([]byte, error) {
	var buffer bytes.Buffer
	tmp := make([]byte, 2048)
	var totalsize int = 0
	for {
		n, err := f.Read(tmp)
		if err == io.EOF {
			totalsize += n
			buffer.Write(tmp)
			break
		} else if err != nil {
			return nil, err
		}
		buffer.Write(tmp)
		totalsize += n
	}
	return buffer.Bytes()[:totalsize], nil
}
