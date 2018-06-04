package instagram

import (
	"fmt"
	"goinstabot/config"
	"goinstabot/db"

	goinsta "gopkg.in/ahmdrz/goinsta.v2"
)

var Insta *goinsta.Instagram

func Login() {
	Insta = goinsta.New(config.Localconfig.InstaUser, config.Localconfig.InstaPass)
	if err := Insta.Login(); err != nil {
		fmt.Println(err)
		return
	}
}
func SyncFollowingDBfromApp() {
	var followingusersDB []string
	var followingusersApp []string
	followingusersDB = getAllFollowing_FromDB()
	followingusersApp = getAllFollowing_FromInstagram()

	if len(followingusersApp) > 0 {
	DBvsApp:
		for _, dbu := range followingusersDB {
			var found bool = false
			for _, dbapp := range followingusersApp {
				if dbu == dbapp {
					found = true
				}
			}
			if !found {
				dberr := db.DBDeletePostgres_Following(dbu)
				if dberr != nil {
					continue DBvsApp
				}
				// This means the user was there before in DB but was deleted from App (unfollowed manually)
				// TODO: Logic to add this user to blacklist and Block it in instagram

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
				// This means the user was added and DB doesnt have it , so we should add it there
				dberr := db.DBInsertPostgres_Following(dbapp)
				if dberr != nil {
					continue AppvsDB
				}
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
	users := Insta.Account.Following()
	for users.Next() {
		for _, user := range users.Users {
			followingusers = append(followingusers, user.Username)
		}
	}
	return followingusers
}
