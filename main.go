package main

import (
	"bad_bot/common"
	"flag"
	"fmt"
)

var (
	Dir   string
	Token string
)

func init() {
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.StringVar(&Dir, "d", "", "Script directory")
	flag.Parse()
}

func main() {
	common.Start(Token, Dir)
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	<-make(chan struct{})
}
