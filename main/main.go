package main

import (
	"fmt"
	"goinstabot/config"
	"goinstabot/db"
	"os"

	goinsta "gopkg.in/ahmdrz/goinsta.v2"
)

var insta *goinsta.Instagram

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
	insta = goinsta.New(config.Localconfig.InstaUser, config.Localconfig.InstaPass)

	// also you can use New function from gopkg.in/ahmdrz/goinsta.v2/utils
	fmt.Println("Login")

	if err := insta.Login(); err != nil {
		fmt.Println(err)
		return
	}

	// export your configuration
	// after exporting you can use Import function instead of New function.
	fmt.Println("Sync")

}
func UpdateFollowingDBFrom() {
	var followingusersDB []string
	var followingusersApp []string
	followingusersDB = getAllFollowing_FromDB()
	followingusersApp = getAllFollowing_FromInstagram()

DBvsApp:
	for _, dbu := range followingusersDB {
		var found bool = false
		for _, dbapp := range followingusersApp {
			if dbu == dbapp {
				found = true
			}
		}
		if !found {
			// This means the user was there before in DB but was deleted from App (unfollowed manually)
			// TODO: Logic to add this user to blacklist and Blockit in instagram
			continue DBvsApp
		}
	}

AppvsDB:
	for _, dbapp := range followingusersApp {
		var found bool = false
		for _, dbu := range followingusersDB {
			if dbu == dbapp {
				found = true
			}
		}
		if !found {
			// This means the user was added and DB doesnt have it , so we should add it
			dberr := db.DBInsertPostgres_Following(dbapp)
			if dberr != nil {
				continue AppvsDB
			}
		}
	}
}
func getAllFollowing_FromDB() []string {
	var followingusers []string
	followingusers = db.DBSelectPostgres_Following()
	return followingusers
}
func getAllFollowing_FromInstagram() []string {
	var followingusers []string
	users := insta.Account.Following()
	for users.Next() {

		fmt.Println("Next:", users.NextID)

		for _, user := range users.Users {
			followingusers = append(followingusers, user.Username)

		}
	}
	return followingusers
}
