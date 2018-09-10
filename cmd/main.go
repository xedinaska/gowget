package main

import (
	"github.com/xedinaska/gowget/downloader"
	"github.com/xedinaska/gowget/drawer"
	"log"
	"os"
	"sync"
	"time"
)

//progressDelay - delay that used to show download progress
const progressDelay = 1 * time.Second

func main() {
	urls := os.Args[1:]

	if len(urls) <= 0 {
		log.Fatalf("[FATAL] at lease one URL should be provided")
	}

	log.Printf("[INFO] download has been started.. %d file(s) in progress", len(urls))

	loader := downloader.Downloader{
		WaitGroup: new(sync.WaitGroup),
		Folder:    "./etc/downloads",
	}

	//start download progress && show possible errors
	failed := loader.Start(urls)
	if len(failed) > 0 {
		for file, msg := range failed {
			log.Printf("file `%s`: %s", file, msg)
		}
	}

	//if all downloads are broken - show fatal error
	if len(failed) == len(urls) {
		log.Fatalf("[FATAL] all downloads are broken")
	}

	//use drawer to print results as table
	tbl := drawer.Table{
		Header: loader.FileNames(),
	}

	tbl.DrawHeader()

	//show progress each `progressDelay` second (by default delay is 1sec)
	t := time.NewTicker(progressDelay)
	defer t.Stop()

	go func() {
		for {
			select {
			case <-t.C:
				tbl.DrawRow(loader.Progress())
			}
		}
	}()

	loader.Wait()
	tbl.DrawRow(loader.Progress())

	log.Printf("[INFO] ..complete")
}
