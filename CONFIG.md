# Config file guide:

Okay basically there is a config.json that you need to pass as argument to the program. Or to be more accurate, you need to pass the full path of the config.json file like this.

`
./main ~/mybot/config.json
`

Now what shoud be there inside this file? lets take a look:
```
{ 
"InstaUser":"myusername",
"InstaPass":"mypassword",
"PreferredNames":["name 1", "name 2"],
"BiographyPreference": ["preference 1","preference 2"],
"BlacklistUsers":["blacklisted 1","blacklisted 2"],
"IRCNick":"goinstaircbot",
"IRCChannels":"#goinstaircbot",
"IRCUser":"ircbot ircbot ircbot ircbot ircbot ircbot",
"IRCServerPort":"irc.freenode.org:6667",
"DBHost":"localhost",
"DBPort":"5432",
"DBUser":"mydbuser",
"DBPass":"mydbpass",
"DBName":"postgres",   
"DialogFlowProjectID":"chatbot-dialogflow",
"DialogFlowLangCode":"en",
"OpeningLine":"hello!!! how are u? I really like that picture."
}
```
Now here I will explain each parameter an why and when is used.
* **InstaUser: The user id of your instagram account, its used during Login.**
* **InstaPass: The password of your instagram account, its used during Login.**
* **PreferredNames: This list of names are considered during Auto following, Bot will only add people with this names.**
* **BiographyPreference: This list of biography words are considered during Auto following, Bot will only add people with any of this words present on his instagram biography.**
* **BlacklistUsers: This are users that you dont want to follow, because of some reason.**
* **IRCNick: This is the nickname of the bot on the IRC server.**
* **IRCChannels: This is the channel the bot will join once connected.**
* **IRCUser: This the IRC user of the bot, leave it like it is in the sample if you dont know what to put.**
* **IRCServerPort: This is the server and port combination that the bot will connect to.**
* **DBHost: This is the hostname of the PostgreSQL server.**
* **DBPort: The port of your PostgreSQL Server.**
* **DBUser: The user to connect to you PostgreSQL DB.**
* **DBPass: The password of the user to connect to your PostgreSQL DB.**
* **DBName: The name of your DB, where your tables will be created.**
* **DialogFlowProjectID: This project ID is created on Google Cloud Platform. Check this [link](https://cloud.google.com/dialogflow-enterprise/docs/quickstart-client-libraries) if you dont have idea of what Im talking.**
* **DialogFlowLangCode: This is the language of the Bot, en for English, es for Spanish, and so on.**
* **OpeningLine: Its the first sentence our bot is going to send to his followers.**

# Database
**Why a DB is required?** Okay basically the bot has some sort of memory stored in the DB.
**Which memory?** For example if you manually unfollow someone, Bot will detect this and will not follow this contact again in future.

**Important:** For database config you just need to give a user that has enough privileges to create tables on a DB, you dont need to create anything, the ORM will take care fo that.

# Google Dialog Flow
You should have a environment variable **GOOGLE_APPLICATION_CREDENTIALS** pointing to the json file of your (GCP) Google Cloud Platform configuration. For more info about this, check this [link](https://cloud.google.com/docs/authentication/getting-started).

If this is not set, the Bot will still work but only the Chatbot functionality will not work.