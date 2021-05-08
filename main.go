package main

import (
	logging "github.com/ipfs/go-log"
	"go-blog/cmd"
)

var log = logging.Logger("main")

func main() {
	logging.SetAllLoggers(logging.LevelDebug)
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
