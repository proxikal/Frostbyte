# Frostbyte ![Build Status](https://api.travis-ci.org/proxikal/Frostbyte.svg?branch=master)
> Open Source Discord Bot written in Golang (DiscordGo library by bwmarrin)  
Comes with a few commands and a light weight A.R.S  
  
## Config Frostbyte
When you clone this branch you will see `config.json` with two entries.  
`Token` -> Discord Bot Token
`Guild` -> Your Guild ID
  
Once you setup and run Frostbyte for the first time to initate the databse  
You can set your `.greet` or `.autorole`  
Once you set one of these systems your `config.json` file will change.
  
You will have the option to set AutoRole, Greeting, ByeMsg and even the channels right from the config.json.  

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
{pm}    - Pms the user
{user}  - Mentions the user
{/user} - Says the users name.
```
More commands and keys coming soon!
  
## Want to contribute?
> Make a pull request! If it **passes** I will merge!  
  

