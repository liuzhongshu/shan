package main

import (
	"fmt"
	"os"

	"github.com/liuzhongshu/shan/commands"
)

func main() {
	err := commands.RootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
