package commands

import (
	"github.com/spf13/cobra"
	"os"
	"io/ioutil"
	"log"
	"path/filepath"
	"time"
	"fmt"
)


var duCmd = &cobra.Command{
	Use:   "du",
	Short: "Summarize disk usage of directories",
	Args: cobra.MinimumNArgs(1),
	Run:   duExecute,
}


func duExecute(cmd *cobra.Command, args []string) {
	tstart := time.Now()
	global.filesChan = make(chan string, 256)
	global.directoryChan = make(chan string, 1024)
	global.sizeChan = make(chan int64, 1024)
	global.doneChan = make(chan bool)

	go resultHandler()
	for i := 0; i < global.cocurrency; i++ {
		go processDirectories()
	}

	for _, target := range args {
		fileinfo, err := os.Stat(target)
		if err == nil {
			if fileinfo.IsDir() {
				global.dirWaitGroup.Add(1)
				global.directoryChan <- target
			}
		}
	}

	global.dirWaitGroup.Wait()
	close(global.directoryChan)
	close(global.sizeChan)
	<-global.doneChan
	fmt.Println("time used:", time.Now().Sub(tstart))
}


func processDirectories() {
	for dir := range global.directoryChan {
		processDirectory(dir)
	}
}

func processDirectory(dir string) {
	defer global.dirWaitGroup.Done()
	if (global.verbose) {
		log.Println(dir)
	}

	files,err := ioutil.ReadDir(dir)
	if err!= nil {
		log.Println(err)
	} else {
		for _,f := range files {
			if (f.IsDir()) {
				global.dirWaitGroup.Add(1)
				select {
				case global.directoryChan <- filepath.Join(dir, f.Name()):
				default:
					processDirectory(filepath.Join(dir, f.Name()))
				}
			} else {
				global.sizeChan <- f.Size()
			}
		}
	}
}

func resultHandler() {
	var totalSize int64 = 0
	for size := range global.sizeChan {
		totalSize += size
	}
	log.Println("total size:", totalSize)
	global.doneChan <- true
}
