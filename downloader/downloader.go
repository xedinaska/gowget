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

//Downloader is used to download remote file to a local drive. WaitGroup & Folder should be initialized before usage.
//Folder property - downloads folder
type Downloader struct {
	*sync.WaitGroup
	Folder string
	items  []*item
}

//item represents file that should be downloaded.
//It contains pointer to file, URL, item size & count of already downloaded bytes (by io.Reader interface implementation)
type item struct {
	io.Reader
	url        string
	name       string
	file       *os.File
	size       float64
	downloaded float64
}

//Read is used to calculate download progress
func (i *item) Read(p []byte) (int, error) {
	n, err := i.Reader.Read(p)
	i.downloaded += float64(n)
	return n, err
}

//Start accepts slice of URLs that should be downloaded and returns map of failed downloads (filename=>message format)
func (d *Downloader) Start(urls []string) map[string]string {
	failed := make(map[string]string, 0)
	for _, u := range urls {
		filename := path.Base(u)

		f, err := os.Create(fmt.Sprintf("%s/%s", d.Folder, filename))
		if err != nil {
			failed[filename] = fmt.Sprintf("unable to create file: `%s`", err.Error())
			continue
		}

		h, err := http.Head(u)
		if err != nil {
			f.Close()
			failed[filename] = fmt.Sprintf("unable to retrieve HEAD: `%s`", err.Error())
			continue
		}

		defer h.Body.Close()

		size, err := strconv.Atoi(h.Header.Get("Content-Length"))
		if err != nil {
			f.Close()
			failed[filename] = fmt.Sprintf("unable to get file length: `%s`", err.Error())
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

	return failed
}

//FileNames returns name of downloaded files
func (d *Downloader) FileNames() []string {
	names := make([]string, len(d.items))
	for c, i := range d.items {
		names[c] = i.name
	}
	return names
}

//Progress returns download progress for all files (percentage)
func (d *Downloader) Progress() []string {
	progress := make([]string, len(d.items))
	for c, i := range d.items {
		progress[c] = fmt.Sprintf("%.2f%%", 100*(i.downloaded/i.size))
	}
	return progress
}

//save is used to get file content & copy it to file on drive. Should be called as goroutine
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
