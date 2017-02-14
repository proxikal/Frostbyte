package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func (bot *Object) CommandHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if bot.System == nil {
		return
	}
	var prefix string = bot.System.Prefix

	// Check user for Manage Server permission.
	if IsManager(s, bot.Guild, m.Author.ID) == true {
		// Redirect to Master Commands.
		bot.MasterCommands(s, m)
	}
	// Execute Auto Response System.
	if strings.Contains(m.Content, prefix+"auto ") == false && strings.Contains(m.Content, prefix+"delauto ") == false {
		var ars map[string]string
		if _, err := os.Stat("autoresponse.json"); err == nil {
			io, err := ioutil.ReadFile("autoresponse.json")
			if err == nil {
				json.Unmarshal(io, &ars)
				for t, r := range ars {
					if strings.Contains(t, "&") {
						// Using the Contains system.
						if strings.Contains(m.Content, t) {
							bot.Parse(s, m, t, r)
						}
					} else {
						// Just a basic trigger.
						if m.Content == t {
							bot.Parse(s, m, t, r)
						}
					}
				}
			}
		}
	}
}

func (bot *Object) MasterCommands(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Set the Prefix variable.
	var prefix string = bot.System.Prefix

	// Set the autorole command.
	if strings.Contains(m.Content, prefix+"autorole ") {
		role := strings.Split(m.Content, prefix+"autorole ")[1]
		// Validate the role submitted.
		if GetRoleID(s, bot.Guild, role) == "" {
			_, err = s.ChannelMessageSend(m.ChannelID, "The role `"+role+"` doesn't exist.")
			if err != nil {
				fmt.Println(err)
			}
		} else {
			// Make sure the Object isn't empty.
			if bot.System == nil {
				fmt.Println("Database Error")
				return
			}
			// Set Autorole to the suggested role.
			bot.System.Autorole = role
			// Marshal (Pretty) for saving to a file.
			js, err := json.MarshalIndent(bot, "", "  ")
			if err == nil {
				// Write entire database to file.
				err = ioutil.WriteFile("config.json", js, 0777)
				if err != nil {
					fmt.Println(err)
				}
				_, err = s.ChannelMessageSend(m.ChannelID, "You have set the auto role to `"+role+"`")
				if err != nil {
					fmt.Println(err)
				}
			} else {
				fmt.Println(err)
			}
		}
	}

	// Set the greet message.
	if strings.Contains(m.Content, prefix+"greet ") {
		// Make sure the Object isn't empty.
		if bot.System == nil {
			fmt.Println("Database Error")
			return
		}
		msg := strings.Split(m.Content, prefix+"greet ")[1]
		// Set the Greet Message
		bot.System.Greeting = msg
		// Marshal (Pretty) for saving to a file.
		js, err := json.MarshalIndent(bot, "", "  ")
		if err != nil {
			fmt.Println(err)
			return
		}
		// Write entire database to file.
		err = ioutil.WriteFile("config.json", js, 0777)
		if err != nil {
			fmt.Println(err)
		} else {
			_, err = s.ChannelMessageSend(m.ChannelID, "You have set the greeting to `"+msg+"`")
			if err != nil {
				fmt.Println(err)
			}
		}
	}

	// add auto to your autoresponse.
	if strings.Contains(m.Content, prefix+"auto ") {
		var ars map[string]string
		data := strings.Replace(m.Content, prefix+"auto ", "", -1)
		trigger := strings.Split(data, "={init}")[0]
		response := strings.Split(data, "={init}")[1]

		if _, err := os.Stat("autoresponse.json"); err != nil {
			// File doesn't exist, make a blank template.
			err = ioutil.WriteFile("autoresponse.json", []byte("{}"), 0777)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
		io, err := ioutil.ReadFile("autoresponse.json")
		if err != nil {
			fmt.Println(err)
			return
		}
		json.Unmarshal(io, &ars)
		ars[trigger] = response
		eo, err := json.MarshalIndent(ars, "", "  ")
		if err != nil {
			fmt.Println(err)
			return
		}
		err = ioutil.WriteFile("autoresponse.json", eo, 0777)
		if err != nil {
			fmt.Println(err)
			return
		}
		_, err = s.ChannelMessageSend(m.ChannelID, "Added `"+trigger+"` with the response:\n```"+response+"```")
		if err != nil {
			fmt.Println(err)
		}
	}

	if strings.HasPrefix(m.Content, prefix+"delauto ") {
		trigger := strings.Replace(m.Content, prefix+"delauto ", "", -1)
		if trigger != "" {
			var ars map[string]string
			io, err := ioutil.ReadFile("autoresponse.json")
			if err != nil {
				fmt.Println(err)
				return
			}
			json.Unmarshal(io, &ars)
			if _, ok := ars[trigger]; ok {
				delete(ars, trigger)
				js, err := json.MarshalIndent(ars, "", "  ")
				if err != nil {
					fmt.Println(err)
					return
				}
				err = ioutil.WriteFile("autoresponse.json", js, 0777)
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	}
}

func (bot *Object) Parse(s *discordgo.Session, m *discordgo.MessageCreate, trigger, response string) {
	var ChannelID string = m.ChannelID

	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.Contains(response, "{user}") {
		response = strings.Replace(response, "{user}", "<@"+m.Author.ID+">", -1)
	}

	if strings.Contains(response, "{/user}") {
		response = strings.Replace(response, "{/user}", m.Author.Username, -1)
	}

	if strings.Contains(response, "{pm}") {
		k, err := s.UserChannelCreate(m.Author.ID)
		if err != nil {
			fmt.Println(err)
			return
		}
		ChannelID = k.ID
	}

	_, err := s.ChannelMessageSend(ChannelID, response)
	if err != nil {
		fmt.Println(err)
	}
}
