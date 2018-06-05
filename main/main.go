package main

import (
	"fmt"
	"goinstaircbot/config"
	"goinstaircbot/db"
	"goinstaircbot/instagram"
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
	//insta, err := goinsta.Import("~/.goinsta")
	fmt.Println("New")

	// export your configuration
	// after exporting you can use Import function instead of New function.
	fmt.Println("Sync")
	//db.DBInsertPostgres_Following("muamuamua")
	instagram.Login()
	//instagram.SyncFollowingDBfromApp()
	instagram.LoadFollowingFromDB()
	instagram.StartFollowingWithMediaLikes(100)
}
