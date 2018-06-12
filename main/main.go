package main

import (
	"fmt"
	"goinstaircbot/config"
	"goinstaircbot/db"
	"goinstaircbot/instagram"
	"goinstaircbot/irc"
	"os"
)

func main() {
	arg := "../config/config.json"
	if len(os.Args) > 1 {
		arg = os.Args[1]
	}
	fmt.Println(arg)
	_ = config.GetConfig(arg)
	db.DBConnectPostgres()
	defer db.DBClosePostgress()

	instagram.LoadMappings()
	instagram.LoadFollowingFromDB()
	instagram.LoadBlockedFromDB()
	irc.StartIRCprocess()

}
