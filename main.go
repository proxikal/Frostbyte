package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/bwmarrin/discordgo"
)

var bot Object
var err error

func main() {

	// Load the config.json file.
	io, err := ioutil.ReadFile("config.json")
	if err != nil {
		fmt.Println(err)
	}
	// Load json into struct: Object
	json.Unmarshal(io, &bot)

	// Login to discord. You can use a token or email, password arguments.
	dg, err := discordgo.New("Bot " + bot.Token)
	if err != nil {
		fmt.Println(err)
		return
	}
	// Register new server
	dg.AddHandler(bot.Initiate)
	// Frostbyte Command handler.
	dg.AddHandler(bot.CommandHandler)
	// Greet Message & Autorole.
	dg.AddHandler(bot.GuildMemberAdd)
	// Bye Message.
	dg.AddHandler(bot.GuildMemberRemove)

	// Save the database every x minutes (Default is 5m)
	go bot.Save()
	dg.Open()

	// Simple way to keep program running until any key press.
	var input string
	fmt.Scanln(&input)
	return
}
