package instagram

import (
	"fmt"
	"goinstaircbot/config"
	"goinstaircbot/db"
	"log"
	"strings"
	"time"

	goinsta "github.com/ahmdrz/goinsta"
)

var (
	Insta          *goinsta.Instagram
	Following      = make(map[string]int)
	Blocked        = make(map[string]int)
	PreferredNames = make(map[string]int)
	FollowingList  []string
)

func Login() {
	if Insta == nil {
		Insta = goinsta.New(config.Localconfig.InstaUser, config.Localconfig.InstaPass)
		if err := Insta.Login(); err != nil {
			fmt.Println(err)
			return
		}
	}

}
func InstaLogout() {
	if Insta != nil {
		Insta.Logout()
		Insta = nil
	}
}
func LoadMappings() {
	PreferredNames = make(map[string]int)
	for _, pref := range config.Localconfig.PreferredNames {
		PreferredNames[pref] = 1
	}
}
func LoadFollowingFromDB() {
	followingusersDB := getAllFollowing_FromDB()
	FollowingList = followingusersDB
	for _, fuser := range followingusersDB {
		Following[fuser] = 1
	}
}
func LoadBlockedFromDB() {
	blockedDB := getAllBlocked_FromDB()
	for _, buser := range blockedDB {
		Blocked[buser] = 1
	}
}
func StartFollowingWithMediaLikes(Limit int) {
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			_, ok = r.(error)
			if !ok {
				fmt.Errorf("pkg: %v", r)
			}
		}
	}()
	var FollowCount int = 0
	for _, myUser := range FollowingList {

		user, err := Insta.Profiles.ByName(myUser)
		log.Println("Checking user: " + myUser)

		if err != nil {
			log.Println(err)
		}
		media := user.Feed()
	MediaLoop:
		for media.Next() {

			fmt.Printf("Printing %d items\n", len(media.Items))
			for _, item := range media.Items {
				item.SyncLikers()
				for _, liker := range item.Likers {
					if FollowCount >= Limit {
						break MediaLoop
					}
					time.Sleep(1 * time.Second)
					fmt.Printf("Checking liker: %s \n", liker.Username)
					fullname := strings.Split(liker.FullName, " ")
					firstname := strings.ToLower(fullname[0])
					if PreferredNames[firstname] == 1 && Blocked[liker.Username] != 1 && Following[liker.Username] != 1 {
						time.Sleep(1 * time.Second)
						profile, err := Insta.Profiles.ByID(liker.ID)
						if err != nil {
							continue
						}
						biography := strings.ToLower(profile.Biography)
					PreferenceLoop:
						for _, pref := range config.Localconfig.BiographyPreference {
							if strings.Contains(biography, pref) {
								profile.Follow()
								FollowCount++
								Following[profile.Username] = 1
								fmt.Printf("Following >>> %s\n", liker.Username)
								break PreferenceLoop
							}
						}
					}
				}
			}

		}
	}
}

func SyncMappings(followingList []string, blockedList []string) {
	Following = make(map[string]int)
	for _, follow := range followingList {
		Following[follow] = 1
	}
	Blocked = make(map[string]int)
	for _, block := range blockedList {
		Blocked[block] = 1
	}
}
func SyncFollowingDBfromApp() {
	var followingusersDB []string
	var followingusersApp []string
	followingusersDB = getAllFollowing_FromDB()
	followingusersApp = getAllFollowing_FromInstagram()
	FollowingList = followingusersApp
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
				dberr = db.DBInsertPostgres_Blocked(dbu)
				if dberr != nil {
					continue DBvsApp
				}

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
	SyncMappings(followingusersApp, db.DBSelectPostgres_Blocked())

}

func getAllFollowing_FromDB() []string {
	var followingusers []string
	followingusers = db.DBSelectPostgres_Following()
	return followingusers
}
func getAllBlocked_FromDB() []string {
	var blockedusers []string
	blockedusers = db.DBSelectPostgres_Blocked()
	return blockedusers
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
