![Frostbyte](https://xtclabs.net/img/byte-logo.png)  
![Build Status](https://api.travis-ci.org/proxikal/Frostbyte.svg?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/proxikal/Frostbyte)](https://goreportcard.com/report/github.com/proxikal/Frostbyte) [![Discord Server](https://img.shields.io/badge/discord-xTech%20Labs-blue.svg)](https://discord.gg/9PRs6xH)
> Open Source Discord Bot written in Golang (DiscordGo library by bwmarrin)  
Comes with a few commands and a light weight A.R.S  
  
| Option | Information |
|:--: | :--: |
| [Discord Developers](https://discordapp.com/developers/applications/me) | Register a bot account with Discord! |
| [Discord Go](https://github.com/bwmarrin/discordgo) | DiscordGo Library by: bwmarrin |
| [Discord Go (Go Docs)](https://godoc.org/github.com/bwmarrin/discordgo) | Godocs collection for DiscordGo |
  
## Config Frostbyte
When you clone this branch you will see `config.json` with two entries.  
`Token` -> Discord Bot Token  
`Guild` -> Your Guild ID
  
Once you setup and run Frostbyte for the first time to initate the databse  
You can set your `.greet` or `.autorole`  
Once you set one of these systems your `config.json` file will change.
  
You will have the option to set AutoRole, Greeting, ByeMsg and even the channels right from the config.json!  

```
{
  "Token": "Discord-Token",
  "Guild": "Your Guild ID",
  "System": {
    "Prefix": "!",                        // Bot Prefix
    "Autorole": "Member",                 // Autorole System
    "Greeting": "Testing stuff {/user}!", // Greet Message!
    "ByeMsg": "",                         // Bye Message
    "Channels": {
      "Autorole": "",                     // Channel for Autorole
      "Greeting": "",                     // Channel for greet
      "ByeMsg": ""                        // Channel for bye.
    },
    "Messages": []                        // List of messages in the collection.
  }
}
```
    
**Commands**:
```
.auto trigger={init}Response
.delauto trigger
.autorole role name
.greet Greet Message
```
  
**A.R.S Keys**
```
{pm}           - Pms the user
{user}         - Mentions the user
{/user}        - Says the users name.
{chan}         - Mentions current channel
{listroles}    - Lists all server roles.
{topic}        - Shows current channel topic.
{redirect}     - {redirect:Channel-ID} Redirect msg to another channel.
```
More commands and keys coming soon!
  
## Want to contribute?
> Make a pull request to `develop` If it **passes** I will merge!
  
### Code Specifications
  
```
1. gofmt -s your code!
2. golint your code!
3. English commenting only!
4. No ineffectual assignments!
5. No suspicious constructs!
```
You can run your branch through [Go Report Card](https://goreportcard.com) Which will check for all cases above  
We need to maintain an `A` or `A+` Standard.  
  
### Master branch
> Stable build of Frostbyte available for use!  
Develop branch will be merged to master every few days.  
We will work on a wikipedia explaining features and usage soon!
  
