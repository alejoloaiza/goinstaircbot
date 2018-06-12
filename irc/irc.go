package irc

import (
	"bufio"
	"fmt"
	"goinstaircbot/config"
	"goinstaircbot/extra"
	"goinstaircbot/instagram"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

var (
	Context     string
	Connection  net.Conn
	FromIRCChan chan string
	ToIRCChan   chan string
)

func StartIRCprocess() {
	FromIRCChan = make(chan string)
	ToIRCChan = make(chan string)
MainCycle:
	for {
		Connection, err := net.Dial("tcp", config.Localconfig.IRCServerPort)

		if err != nil {
			log.Println(err)
			time.Sleep(2000 * time.Millisecond)
			continue MainCycle
		}

		fmt.Fprintln(Connection, "NICK "+config.Localconfig.IRCNick)
		fmt.Fprintln(Connection, "USER "+config.Localconfig.IRCUser)
		fmt.Fprintln(Connection, "JOIN "+config.Localconfig.IRCChannels)
		go RoutineWriter(Connection)
		MyReader := bufio.NewReader(Connection)
	ReaderCycle:
		for {

			message, err := MyReader.ReadString('\n')
			if err != nil {
				log.Println(time.Now().Format(time.Stamp) + ">>>" + err.Error())
				if io.EOF == err {
					Connection.Close()
					log.Println("server closed connection")
				}
				time.Sleep(2000 * time.Millisecond)
				break ReaderCycle
			}

			log.Print(time.Now().Format(time.Stamp) + ">>" + message)

			text := strings.Split(message, " ")
			//log.Println("Number of objects in text: "+ strconv.Itoa(len(text)))
			var respond bool = false
			var response string
			if len(text) >= 4 && text[1] == "PRIVMSG" {
				respond = true
				var repeat bool = true
				var respondTo string
				if text[2][0:1] == "#" {
					// logic to respond the same thing to a channel / repeater BOT
					respondTo = text[2]
					Context = respondTo
				} else {
					userto := strings.Split(text[0], "!")
					respondTo = userto[0][1:]
					Context = respondTo
					// logic to respond the same thing to a user / repeater BOT
				}
				// If its a command BOT will execute the command given
				if text[3] == ":!cmd" {
					repeat = false
					commandresponse := ProcessCommand(text[4:])
					response = "PRIVMSG " + respondTo + " :" + commandresponse

				}
				// If is not a command BOT will repeat the same thing
				if repeat == true {
					response = "PRIVMSG " + respondTo + " " + strings.Join(text[3:], " ")

				}
			}
			if len(text) == 2 && text[0] == "PING" {
				response = "PONG " + text[1]
				respond = true
			}
			// This checks if the received text requires response or not, and respond according to the above logic

			if respond == true {
				fmt.Fprintln(Connection, response)
				log.Println(time.Now().Format(time.Stamp) + "<<" + response)
			}

		}
		time.Sleep(2000 * time.Millisecond)
	}

}

func ProcessCommand(command []string) string {

	var bodyString string
	/*
		var UserToFollow string = ""
		if strings.TrimSpace(command[0]) == "stop" {
			OutChan <- "stop"
			bodyString = "Command received... processing"
		}
		if len(command) >= 3 && strings.TrimSpace(command[0]) == "set" {
			if extra.RemoveEnds(command[1]) == "proxy" {
				config.Localconfig.UseProxy, _ = strconv.ParseBool(extra.RemoveEnds(command[2]))
				bodyString = "Proxy set to " + strconv.FormatBool(config.Localconfig.UseProxy)
			}
			if extra.RemoveEnds(command[1]) == "logout" {
				instagram.InstaLogout()
				bodyString = "Logged out"

			}
		}
	*/
	if len(command) >= 3 && strings.TrimSpace(command[0]) == "init" {
		var arg2 int
		var err error
		arg1 := extra.RemoveEnds(command[1])
		if extra.IsInteger(extra.RemoveEnds(command[2])) {
			arg2, err = strconv.Atoi(extra.RemoveEnds(command[2]))
		}

		if err != nil {
			return ""
		}
		switch arg1 {
		case "follow":
			go ExecuteFollowProcess(arg2)
		case "message":
			go ExecuteMessageProcess(arg2)
		case "chatbot":
			go ExecuteChatbotProcess()
		}
		bodyString = "Command received... processing"
	}

	return bodyString
}

func ExecuteFollowProcess(Limit int) {
	instagram.Login(FromIRCChan, ToIRCChan)
	instagram.SyncFollowingDBfromApp()
	instagram.StartFollowingWithMediaLikes(Limit)
	defer instagram.InstaLogout()
}

func ExecuteMessageProcess(Limit int) {
	instagram.Login(FromIRCChan, ToIRCChan)
	instagram.SyncFollowingDBfromApp()
	instagram.StartSendingNewMessages(Limit)
	defer instagram.InstaLogout()
}
func ExecuteChatbotProcess() {
	instagram.Login(FromIRCChan, ToIRCChan)
	instagram.StartChatbot()
	defer instagram.InstaLogout()
}

func RoutineWriter(Response net.Conn) {
	for {
		var err error
		select {
		case msg := <-ToIRCChan:
			if Context != "" {
				_, err = fmt.Fprintln(Response, "PRIVMSG "+Context+" :"+msg)
			} else {
				_, err = fmt.Fprintln(Response, "PRIVMSG "+config.Localconfig.IRCChannels+" :"+msg)
			}
		}
		if err != nil {
			log.Println("Error during Write in RoutineWriter")
			break
		}
	}
}
