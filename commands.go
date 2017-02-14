package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// GetPageContents - Get page content based on URL.
// url: Valid url of image.
func GetPageContents(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return body, err
}

// CommandHandler - Handle the commands and Auto Response System
// bot: Main Object with all your settings.
// s: The Current Session between the bot and discord
// m: The Message Object sent back from Discord.
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

// MasterCommands - Commands Available to people with Manage Server Permissions.
// bot: Main Object with all your settings.
// s: The Current Session between the bot and discord
// m: The Message Object sent back from Discord.
func (bot *Object) MasterCommands(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Set the Prefix variable.
	var prefix string = bot.System.Prefix

	// Set the autorole command.
	if strings.Contains(m.Content, prefix+"autorole ") {
		role := strings.Split(m.Content, prefix+"autorole ")[1]
		// Validate the role submitted.
		if bot.GetRoleID(s, role) == "" {
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
					_, err = s.ChannelMessageSend(m.ChannelID, fmt.Sprint(err))
					if err != nil {
						fmt.Println(err)
					}
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
			_, err = s.ChannelMessageSend(m.ChannelID, fmt.Sprint(err))
			if err != nil {
				fmt.Println(err)
			}
			return
		}
		// Write entire database to file.
		err = ioutil.WriteFile("config.json", js, 0777)
		if err != nil {
			fmt.Println(err)
			_, err = s.ChannelMessageSend(m.ChannelID, fmt.Sprint(err))
			if err != nil {
				fmt.Println(err)
			}
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
				_, err = s.ChannelMessageSend(m.ChannelID, fmt.Sprint(err))
				if err != nil {
					fmt.Println(err)
				}
				return
			}
		}
		io, err := ioutil.ReadFile("autoresponse.json")
		if err != nil {
			fmt.Println(err)
			_, err = s.ChannelMessageSend(m.ChannelID, fmt.Sprint(err))
			if err != nil {
				fmt.Println(err)
			}
			return
		}
		json.Unmarshal(io, &ars)
		ars[trigger] = response
		eo, err := json.MarshalIndent(ars, "", "  ")
		if err != nil {
			fmt.Println(err)
			_, err = s.ChannelMessageSend(m.ChannelID, fmt.Sprint(err))
			if err != nil {
				fmt.Println(err)
			}
			return
		}
		err = ioutil.WriteFile("autoresponse.json", eo, 0777)
		if err != nil {
			fmt.Println(err)
			_, err = s.ChannelMessageSend(m.ChannelID, fmt.Sprint(err))
			if err != nil {
				fmt.Println(err)
			}
			return
		}
		_, err = s.ChannelMessageSend(m.ChannelID, "Added `"+trigger+"` with the response:\n```"+response+"```")
		if err != nil {
			fmt.Println(err)
		}
	}

	// Command to delete a rule from the A.R.S - delauto triggername or delauto &triggername
	if strings.HasPrefix(m.Content, prefix+"delauto ") {
		trigger := strings.Replace(m.Content, prefix+"delauto ", "", -1)
		if trigger != "" {
			var ars map[string]string
			io, err := ioutil.ReadFile("autoresponse.json")
			if err != nil {
				fmt.Println(err)
				_, err = s.ChannelMessageSend(m.ChannelID, fmt.Sprint(err))
				if err != nil {
					fmt.Println(err)
				}
				return
			}
			json.Unmarshal(io, &ars)
			if _, ok := ars[trigger]; ok {
				delete(ars, trigger)
				js, err := json.MarshalIndent(ars, "", "  ")
				if err != nil {
					fmt.Println(err)
					_, err = s.ChannelMessageSend(m.ChannelID, fmt.Sprint(err))
					if err != nil {
						fmt.Println(err)
					}
					return
				}
				err = ioutil.WriteFile("autoresponse.json", js, 0777)
				if err != nil {
					fmt.Println(err)
					_, err = s.ChannelMessageSend(m.ChannelID, fmt.Sprint(err))
					if err != nil {
						fmt.Println(err)
					}
				}
			}
		}
	}

	// Change Bot Avatar.
	if strings.HasPrefix(m.Content, prefix+"avatar") {
		var img []byte
		if strings.Replace(m.Content, prefix+"avatar", "", -1) == "" {
			img, err = ioutil.ReadFile("avatar.png")
			if err != nil {
				fmt.Println(err)
				_, err = s.ChannelMessageSend(m.ChannelID, fmt.Sprint(err))
				if err != nil {
					fmt.Println(err)
				}
				return
			}
		} else {
			// Collect the avatar from a link.
			url := strings.Replace(m.Content, prefix+"avatar ", "", -1)
			img, err = GetPageContents(url)
			if err != nil {
				fmt.Println(err)
				_, err = s.ChannelMessageSend(m.ChannelID, fmt.Sprint(err))
				if err != nil {
					fmt.Println(err)
				}
				return
			}
		}
		base64 := base64.StdEncoding.EncodeToString(img)
		avatar := fmt.Sprintf("data:image/png;base64,%s", string(base64))
		_, err = s.UserUpdate("", "", s.State.User.Username, avatar, "")
		if err != nil {
			fmt.Println(err)
			_, err = s.ChannelMessageSend(m.ChannelID, fmt.Sprint(err))
			if err != nil {
				fmt.Println(err)
			}
			return
		}
		_, err = s.ChannelMessageSend(m.ChannelID, "Successfully changed my avatar.")
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

// Parse - Auto Response Keys into Data and than send Response (if any)
// bot: Main Object with all your settings.
// s: The Current Session between the bot and discord
// m: The Message Object sent back from Discord.
// Auto Response System is Licensed by the MIT and used by many other bots on Discord.
// Originally attempted with Paradox Bot (4/11/2016)
// Pefected with Echo 2.0!
func (bot *Object) Parse(s *discordgo.Session, m *discordgo.MessageCreate, trigger, response string) {
	var ChannelID string
	ChannelID = m.ChannelID
	// The bot cannot trigger the A.R.S (Would be really bad)
	if m.Author.ID == s.State.User.ID {
		return
	}

	// If the user is a bot than we're going to ignore it.
	if m.Author.Bot == true {
		return
	}

	// Show channel name (with mention) #general, #lobby etc..
	if strings.Contains(response, "{chan}") {
		response = strings.Replace(response, "{chan}", "<#"+m.ChannelID+">", -1)
	}

	// Show channel topic.
	if strings.Contains(response, "{topic}") {
		ch, err := s.State.Channel(m.ChannelID)
		if err != nil {
			fmt.Println(err)
			return
		}
		response = strings.Replace(response, "{topic}", ch.Topic, -1)
	}

	// List all the roles in a guild.
	if strings.Contains(response, "{listroles}") {
		var glob string
		g, err := s.State.Guild(bot.Guild)
		if err != nil {
			fmt.Println(err)
			return
		}
		if g.Roles != nil {
			for _, r := range g.Roles {
				glob = glob + r.Name + "\n"
			}
			if glob != "" {
				response = strings.Replace(response, "{listroles}", glob, -1)
			} else {
				response = strings.Replace(response, "{listroles}", "[Empty]", -1)
			}
		} else {
			response = strings.Replace(response, "{listroles}", "[Empty]", -1)
		}
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
