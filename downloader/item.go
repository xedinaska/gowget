package downloader

import (
	"io"
	"os"
)

type item struct {
	io.Reader
	url        string
	name       string
	file       *os.File
	size       float64
	downloaded float64
}

func (i *item) Read(p []byte) (int, error) {
	n, err := i.Reader.Read(p)
	i.downloaded += float64(n)
	return n, err
}
