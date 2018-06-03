package main

import (
	"fmt"
	"goinstabot/config"
	"goinstabot/db"
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
	//insta, err := goinsta.Import("~/.goinsta")
	fmt.Println("New")

	// export your configuration
	// after exporting you can use Import function instead of New function.
	fmt.Println("Sync")

}
