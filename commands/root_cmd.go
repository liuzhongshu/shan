package commands

import (
	"github.com/spf13/cobra"
	"sync"
)
const (
	CORES = 8
)
var RootCmd = &cobra.Command{
	Use:   "shan",
	Short: "shan is a fast directory/file tool",
}

var global = struct {
	cocurrency    int
	verbose   	  bool
	filesChan     chan string
	directoryChan chan string
	sizeChan      chan int64
	doneChan      chan bool
	dirWaitGroup  sync.WaitGroup
}{}

func init() {
	RootCmd.AddCommand(duCmd)

	RootCmd.PersistentFlags().IntVarP(&global.cocurrency,"cocurrency", "c", 4, "cocurrency number")
	RootCmd.PersistentFlags().BoolVarP(&global.verbose,"verbose", "", false, "verbose mode")
}