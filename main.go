package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

var bot Object
var start int64
var err error

func main() {

	start = time.Now().UnixNano()
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

	// Handle status messages and the interval.
	go bot.StatusHandler(dg, "2m")

	// Save the database every x minutes (Default is 5m)
	go bot.Save()
	dg.Open()
	go bot.Intro(dg)
	// Simple way to keep program running until any key press.
	var input string
	fmt.Scanln(&input)
	return
}

// CommandHandler - Handle the commands and Auto Response System
// bot: Main Object with all your settings.
// s: The Current Session between the bot and discord
// m: The Message Object sent back from Discord.
func (bot *Object) CommandHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if bot.System == nil {
		return
	}
	var prefix = bot.System.Prefix

	// Check user for Manage Server permission.
	if IsManager(s, bot.Guild, m.Author.ID) == true {
		bot.AddARS(s, m, prefix)            // Listen for .auto command
		bot.DeleteARS(s, m, prefix)         // Listen for .delauto command
		bot.AutoRoleCommand(s, m, prefix)   // listen for .autorole command
		bot.GreetCommand(s, m, prefix)      // Listen for .greet command
		bot.ChangeAvatar(s, m, prefix)      // Listen for .avatar command
		bot.InfoCommand(s, m, prefix)       // Listen for .info command
		bot.ViewARS(s, m, prefix)           // Listen for .viewauto command
		bot.InspectARS(s, m, prefix)        // Listen for .inspect command
		bot.AddStatusCommand(s, m, prefix)  // Listen for .addstatus command
		bot.DelStatusCommand(s, m, prefix)  // Listen for .delstatus command
		bot.ViewStatusCommand(s, m, prefix) // Listen for .viewstatus command
	}
	// Execute Auto Response System.
	if m.Author.ID != s.State.User.ID && m.Author.Bot == false {
		bot.Listen(s, m, prefix)
	}
}

// Listen - Listens for A.R.S Commands.
// bot: Main Object with all your settings.
// s: The Current Session between the bot and discord
// m: The Message Object sent back from Discord.
func (bot *Object) Listen(s *discordgo.Session, m *discordgo.MessageCreate, prefix string) {
	if strings.Contains(m.Content, prefix+"auto ") == false && strings.Contains(m.Content, prefix+"delauto ") == false {
		var ars map[string]string
		if _, err := os.Stat("autoresponse.json"); err != nil {
			return
		}
		io, err := ioutil.ReadFile("autoresponse.json")
		if err != nil {
			fmt.Println(err)
			return
		}
		json.Unmarshal(io, &ars)
		for t, r := range ars {
			if strings.Contains(t, "&") {
				// Using the Contains system.
				if strings.Contains(m.Content, t) {
					content := bot.ParseServer(s, m, t, r)
					bot.ParseDirection(s, m, t, content)
				}
			} else {
				// Just a basic trigger.
				if m.Content == t {
					content := bot.ParseServer(s, m, t, r)
					bot.ParseDirection(s, m, t, content)
				}
			}
		}
	}
}

// StatusHandler - (if set) iterates through your status messages
// bot: Main Object with all your settings.
// s: The Current Session between the bot and discord
// m: The Message Object sent back from Discord.
func (bot *Object) StatusHandler(s *discordgo.Session, duration string) {
	if strings.Contains(duration, "m") == false {
		duration = strings.Replace(duration, "s", "m", -1)
	}
	p, err := time.ParseDuration(duration)
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		if bot.System != nil {
			if len(bot.System.Status) > 0 {
				r := Random(0, len(bot.System.Status))
				data := bot.System.Status[r]
				err = s.UpdateStatus(0, data)
				if err != nil {
					fmt.Println(err)
				}
			}
		}
		<-time.After(p)
	}
}
