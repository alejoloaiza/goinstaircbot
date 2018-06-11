package instagram

import (
	"fmt"
	"goinstaircbot/chatbot"
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
	FromIRCChan    chan string
	ToIRCChan      chan string
	InboxUsers     = make(map[string]int)
)

func Login(fromirc chan string, toirc chan string) {
	FromIRCChan = fromirc
	ToIRCChan = toirc
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
func getInbox_FromInstagram() []string {
	var inboxusers []string
	err := Insta.Inbox.Sync()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	for _, conversation := range Insta.Inbox.Conversations {
		for _, iuser := range conversation.Users {
			if iuser.Username != config.Localconfig.InstaUser {
				inboxusers = append(inboxusers, iuser.Username)
			}
		}
	}
	for Insta.Inbox.Next() {
		for _, conversation := range Insta.Inbox.Conversations {
			for _, iuser := range conversation.Users {
				if iuser.Username != config.Localconfig.InstaUser {
					inboxusers = append(inboxusers, iuser.Username)
				}
			}
		}
	}
	return inboxusers
}
func LoadBlockedFromDB() {
	blockedDB := getAllBlocked_FromDB()
	for _, buser := range blockedDB {
		Blocked[buser] = 1
	}
}
func StartFollowingWithMediaLikes(Limit int) {
	var Rejected = make(map[string]int)

	defer func() {
		if r := recover(); r != nil {
			var ok bool
			_, ok = r.(error)
			if !ok {
				sendMessage(ToIRCChan, fmt.Sprintf("Recover from error: %v", r))
			}
		}
	}()
	var FollowCount int = 0
	for _, myUser := range FollowingList {

		user, err := Insta.Profiles.ByName(myUser)

		sendMessage(ToIRCChan, fmt.Sprintf("Checking user >>> %s ", myUser))

		if err != nil {
			log.Println(err)
		}
		media := user.Feed()
	MediaLoop:
		for media.Next() {
			for _, item := range media.Items {
				item.SyncLikers()
				for _, liker := range item.Likers {
					match := false
					if FollowCount >= Limit {
						sendMessage(ToIRCChan, fmt.Sprintf("Finished with #%v ", FollowCount))
						break MediaLoop
					}
					time.Sleep(1 * time.Second)

					fullname := strings.Split(liker.FullName, " ")
					firstname := strings.ToLower(fullname[0])
					if Rejected[liker.Username] != 1 && PreferredNames[firstname] == 1 && Blocked[liker.Username] != 1 && Following[liker.Username] != 1 {
						time.Sleep(1 * time.Second)
						profile, err := Insta.Profiles.ByID(liker.ID)
						if err != nil {
							continue
						}
						biography := strings.ToLower(profile.Biography)
						sendMessage(ToIRCChan, fmt.Sprintf("Checking liker >>> %s ", liker.Username))

					PreferenceLoop:
						for _, pref := range config.Localconfig.BiographyPreference {
							if strings.Contains(biography, pref) {
								profile.Follow()
								FollowCount++
								Following[profile.Username] = 1
								match = true
								sendMessage(ToIRCChan, fmt.Sprintf("Following #%v>>> %s ", FollowCount, liker.Username))
								time.Sleep(10 * time.Second)
								break PreferenceLoop
							}
						}
						if !match {
							Rejected[profile.Username] = 1
						}
					}
				}
			}

		}
	}
}
func StartSendingNewMessages(Limit int) {
	inboxusers := getInbox_FromInstagram()
	for _, iuser := range inboxusers {
		InboxUsers[iuser] = 1
	}
	sendMessage(ToIRCChan, "Inbox fully loaded")
	for _, myUser := range FollowingList {
		if InboxUsers[myUser] != 1 {
			// TODO LOGIC TO SEND RANDOM MESSAGE
			newmsguser, err := Insta.Profiles.ByName(myUser)
			if err != nil {
				fmt.Println(err.Error())
			}
			text := config.Localconfig.OpeningLine
			Insta.Inbox.New(newmsguser, text)
			InboxUsers[myUser] = 1
			sendMessage(ToIRCChan, fmt.Sprintf("Message sent to: %s ", myUser))
			time.Sleep(10 * time.Minute)

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
	sendMessage(ToIRCChan, "Started Sync of Following")
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
	sendMessage(ToIRCChan, "Finished Sync of Following")

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
func sendMessage(toirc chan string, message string) {
	log.Printf(message)
	toirc <- message
}
func StartChatbot() {
	//TODO
}
func GetResponseFromChatbot(text string, username string) string {
	var projectID = config.Localconfig.DialogFlowProjectID
	var langCode = config.Localconfig.DialogFlowLangCode
	response, err := chatbot.DetectIntentText(projectID, username, text, langCode)
	if err != nil {
		return err.Error()
	}
	return response
}
