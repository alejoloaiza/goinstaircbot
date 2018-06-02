package main

import (
	"fmt"
	"goinstabot/config"
	"os"

	goinsta "gopkg.in/ahmdrz/goinsta.v2"
)

var insta *goinsta.Instagram

func main() {
	arg := "../config/config.json"
	if len(os.Args) > 1 {
		arg = os.Args[1]
	}
	_ = config.GetConfig(arg)
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

	err := insta.Inbox.Sync()
	if err != nil {
		panic(err)
	}
	i := 1
	fmt.Printf("Page %d has %d conversations\n", i, len(insta.Inbox.Conversations))

	for insta.Inbox.Next() {
		i++
		fmt.Printf("Page %d has %d conversations\n", i, len(insta.Inbox.Conversations))
	}
	//inst.Inbox.Reset()

}
func getAllFollowing() {

	users := insta.Account.Following()
	var cycle int32 = 0
	for users.Next() {

		fmt.Println("Next:", users.NextID)
		for _, user := range users.Users {
			cycle++
			fmt.Printf("   - %v - %s\n", cycle, user.Username)
		}
	}
}
