# GoInstaIRCBot is a IRC Bot that you can control and interacts with Instagram

`Okay, who doesn't love IRC bots? ...... hummm, noone? Okay I do! I grew up within IRC World`

## Features

* **AutoFollowing: It will automatically follow new people based on biography preferences and name preferences.**
* **AutoMessaging: It will automatically send direct messages to contacts, you can configure what to send.**
* **Chatbot automation: If someone responds to you message, you can have a full conversation based on Google DialogFlow chatbot framework.**


## Package installation 
`go get -u -v github.com/alejoloaiza/goinstaircbot`

## Package dependencies
* **For Instagram interaction.**
- `go get -u -v github.com/ahmdrz/goinsta`
* **For Google Dialog Flow Chatbot.**
- `go get -u -v cloud.google.com/go/dialogflow/apiv2`
- `go get -u -v google.golang.org/genproto/googleapis/cloud/dialogflow/v2`
* **For ORM with PostgreSQL Database.**
- `go get -u -v github.com/jinzhu/gorm`

## Config file
You must provide a config file fullpath as unique argument to the program. Here(https://github.com/alejoloaiza/goinstaircbot/blob/master/CONFIG.md) you can find a guide on how to set up your config file.

## Commands

Once your bot is configured an joins a channel, you can talk to him on channel or on private message, this are the commands:

 **1) !cmd init follow `<number of new followers>`**
- `<number of new followers>` give a limit number of followers you want to add. Try not to add more than 150 contacts by day.
* **Once the bot recieves this command, it will start doing this steps:**
* **a) Check the contacts you are following and the media of those contacts.**
* **b) Check if that media has likes, and see if those likes belong to people with names listed withing the config parameter PreferredNames.**
* **c) Check if that contact has any of the words present in the config parameter BiographyPreference.**
* **d) It will start following the contact that matches with those preferences.**

 **2) !cmd init message `<number of messages>`**
* **- `<number of messages>` give a number of new messages you want to send.**
* **Once the bot recieves this command, it will start doing this steps:**
* **a) Check all the contacts that already have been messaged before (Direct message inbox).**
* **b) Send the defined message on config variable OpeningLine to the following contacts that were not messaged before.**

 **3) !cmd init chatbot**
* **Once the bot recieves this command, it will start doing this steps:**
* **a) Every once a while will check if someone has responded any of his messages.**
* **b) If someone responded, it will respond again based on a chatbot flow defined in Google Dialog Flow and mapped to this appliation with the config variables DialogFlowProjectID and DialogFlowLangCode.**
