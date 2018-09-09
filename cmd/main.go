package main

import (
	"git.oxagile.com/mmgfusppointment/gowget/downloader"
	"git.oxagile.com/mmgfusppointment/gowget/drawer"
	"log"
	"os"
	"sync"
	"time"
)

const progressDelay = 1 * time.Second

func main() {
	urls := os.Args[1:]

	if len(urls) <= 0 {
		log.Fatalf("[FATAL] at lease one URL should be provided")
	}

	log.Printf("[INFO] download has been started.. %d file(s) in progress", len(urls))

	loader := downloader.Downloader{
		WaitGroup: new(sync.WaitGroup),
		Folder: "./etc/downloads",
	}
	loader.Start(urls)

	t := time.NewTicker(progressDelay)
	defer t.Stop()

	tbl := drawer.Table{
		Header: loader.GetFileNames(),
	}

	tbl.DrawHeader()

	go func() {
		for {
			select {
			case <-t.C:
				tbl.DrawRow(loader.GetProgress())
			}
		}
	}()

	loader.Wait()
	tbl.DrawRow(loader.GetProgress())

	log.Printf("[INFO] ..complete")
}
