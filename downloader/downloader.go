package downloader

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"sync"
)

type Downloader struct {
	*sync.WaitGroup
	Folder string
	items []*item
}

func (d *Downloader) Start(urls []string) {
	for _, u := range urls {
		filename := path.Base(u)

		f, err := os.Create(fmt.Sprintf("%s/%s", d.Folder, filename))
		if err != nil {
			log.Printf("[ERROR] unable to create file: `%s`", err.Error())
			continue
		}

		h, err := http.Head(u)
		if err != nil {
			f.Close()
			log.Printf("[ERROR] failed to fetch URL head: `%s`", err.Error())
			continue
		}

		defer h.Body.Close()

		size, err := strconv.Atoi(h.Header.Get("Content-Length"))
		if err != nil {
			f.Close()
			log.Printf("[ERROR] unable to get file length: `%s`", err.Error())
			continue
		}

		i := &item{
			url:  u,
			name: filename,
			size: float64(size),
			file: f,
		}
		d.items = append(d.items, i)
		d.Add(1)
	}

	for _, i := range d.items {
		go d.save(i)
	}
}

func (d *Downloader) GetFileNames() []string {
	names := make([]string, len(d.items))
	for c, i := range d.items {
		names[c] = i.name
	}
	return names
}

func (d *Downloader) GetProgress() []string {
	progress := make([]string, len(d.items))
	for c, i := range d.items {
		progress[c] = fmt.Sprintf("%.2f%%", 100*(i.downloaded/i.size))
	}
	return progress
}

func (d *Downloader) save(i *item) {
	defer i.file.Close()
	defer d.Done()

	resp, err := http.Get(i.url)
	if err != nil {
		log.Printf("[ERROR] unable to make GET call: `%s`", err.Error())
		return
	}

	defer resp.Body.Close()

	i.Reader = resp.Body
	if _, err = io.Copy(i.file, i); err != nil {
		log.Printf("[ERROR] unable to copy file content: `%s`", err.Error())
		return
	}
}
