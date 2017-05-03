package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

// AddStatusCommand - Adds a status messages to Frostbyte.
func (bot *Object) AddStatusCommand(s *discordgo.Session, m *discordgo.MessageCreate, prefix string) {
	if strings.Contains(m.Content, prefix+"addstatus ") {
		status := strings.Replace(m.Content, prefix+"addstatus ", "", -1)
		err = bot.AddStatus(status)
		if err == nil {
			js, err := json.MarshalIndent(bot, "", "  ")
			if err != nil {
				fmt.Println(err)
			} else {
				ioutil.WriteFile("config.json", js, 0777)
				_, err = s.ChannelMessageSend(m.ChannelID, "Status has been added to collection.")
				if err != nil {
					fmt.Println(err)
				}
			}
		} else {
			_, err = s.ChannelMessageSend(m.ChannelID, fmt.Sprint(err))
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

// DelStatusCommand - allows you to delete existing status messages in frostbyte.
func (bot *Object) DelStatusCommand(s *discordgo.Session, m *discordgo.MessageCreate, prefix string) {
	if strings.Contains(m.Content, prefix+"delstatus ") {
		status := strings.Replace(m.Content, prefix+"delstatus ", "", -1)
		err = bot.RemoveStatus(status)
		if err == nil {
			js, err := json.MarshalIndent(bot, "", "  ")
			if err == nil {
				ioutil.WriteFile("config.json", js, 0777)
				_, err = s.ChannelMessageSend(m.ChannelID, "Status has been removed from collection.")
				if err != nil {
					fmt.Println(err)
				}
			} else {
				fmt.Println(err)
			}
		} else {
			_, err = s.ChannelMessageSend(m.ChannelID, fmt.Sprint(err))
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

// ViewStatusCommand - Views the existing status messages.
func (bot *Object) ViewStatusCommand(s *discordgo.Session, m *discordgo.MessageCreate, prefix string) {
	if m.Content == prefix+"viewstatus" {
		var glob string
		for _, st := range bot.System.Status {
			glob = glob + st + "\n"
		}
		if glob == "" {
			_, err = s.ChannelMessageSend(m.ChannelID, "You don't have any status messages set.")
			if err != nil {
				fmt.Println(err)
			}
		} else {
			_, err = s.ChannelMessageSend(m.ChannelID, "```"+glob+"```")
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

// GreetCommand - greet new people with a message.
func (bot *Object) GreetCommand(s *discordgo.Session, m *discordgo.MessageCreate, prefix string) {
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
}

// AutoRoleCommand - give new people roles.
func (bot *Object) AutoRoleCommand(s *discordgo.Session, m *discordgo.MessageCreate, prefix string) {
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
}

// ChangeAvatar - Changes bot avatar
func (bot *Object) ChangeAvatar(s *discordgo.Session, m *discordgo.MessageCreate, prefix string) {

	// Change Bot Avatar.
	if strings.HasPrefix(m.Content, prefix+"avatar") {
		var img []byte
		if strings.Replace(m.Content, prefix+"avatar", "", -1) == "" {
			img, err = ioutil.ReadFile("avatar.jpg")
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

// AddARS - adds rule to your A.R.S.
func (bot *Object) AddARS(s *discordgo.Session, m *discordgo.MessageCreate, prefix string) {
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
}

// DeleteARS - removes rule from your A.R.S.
func (bot *Object) DeleteARS(s *discordgo.Session, m *discordgo.MessageCreate, prefix string) {
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
					return
				}
				s.ChannelMessageSend(m.ChannelID, "I've removed the rule `"+trigger+"` from your A.R.S")
			} else {
				s.ChannelMessageSend(m.ChannelID, "The trigger doesn't exist in your A.R.S Rules.")
			}
		}
	}
}

// InfoCommand - Displays an Embed with some information.
func (bot *Object) InfoCommand(s *discordgo.Session, m *discordgo.MessageCreate, prefix string) {
	if m.Content == prefix+"info" {
		t := time.Unix(0, start)
		elapsed := time.Since(t)
		upt := fmt.Sprintf("%s", elapsed)
		upti := strings.Split(upt, ".")
		uptime := upti[0]
		stats, err := s.State.Guild(bot.Guild)
		if err != nil {
			fmt.Println(err)
			s.ChannelMessageSend(m.ChannelID, fmt.Sprint(err))
		}
		channelcount := strconv.Itoa(len(stats.Channels))
		membercount := strconv.Itoa(len(stats.Members))
		rolecount := strconv.Itoa(len(stats.Roles))
		var dgo []*discordgo.MessageEmbedField
		field1 := &discordgo.MessageEmbedField{
			Name:   "Current Members",
			Value:  membercount,
			Inline: true,
		}
		field2 := &discordgo.MessageEmbedField{
			Name:   "Channel Count",
			Value:  channelcount,
			Inline: true,
		}
		field3 := &discordgo.MessageEmbedField{
			Name:   "Role Count",
			Value:  rolecount,
			Inline: true,
		}
		field4 := &discordgo.MessageEmbedField{
			Name:   "Uptime",
			Value:  uptime + "s",
			Inline: true,
		}
		field5 := &discordgo.MessageEmbedField{
			Name:   "Github",
			Value:  "[Open Source](https://github.com/proxikal/Frostbyte)",
			Inline: true,
		}

		dgo = append(dgo, field1)
		dgo = append(dgo, field2)
		dgo = append(dgo, field3)
		dgo = append(dgo, field4)
		dgo = append(dgo, field5)
		co1 := strings.TrimSpace("4286f4")
		color, _ := strconv.ParseInt(co1, 16, 0)
		obj := &discordgo.MessageEmbed{
			URL:         "https://discord.gg/9PRs6xH",
			Type:        "rich",
			Title:       "Frostbyte v0.0.1 (Silver)",
			Description: "",
			Color:       int(color),
			Footer: &discordgo.MessageEmbedFooter{
				Text:    "Need more help? Join our Server!",
				IconURL: "https://xtclabs.net/img/byte-icon.png",
			},
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: "https://xtclabs.net/img/byte-icon.png",
			},
			Fields: dgo,
		}
		_, err = s.ChannelMessageSendEmbed(m.ChannelID, obj)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

// ViewARS - Displays an Embed with some information.
func (bot *Object) ViewARS(s *discordgo.Session, m *discordgo.MessageCreate, prefix string) {
	if m.Content == prefix+"viewauto" {
		var ars map[string]string
		io, err := ioutil.ReadFile("autoresponse.json")
		if err != nil {
			fmt.Println(err)
			s.ChannelMessageSend(m.ChannelID, fmt.Sprint(err))
		} else {
			err = json.Unmarshal(io, &ars)
			if err != nil {
				fmt.Println(err)
				s.ChannelMessageSend(m.ChannelID, fmt.Sprint(err))
			}
			var glob string
			for t := range ars {
				glob = glob + "Trigger: " + t + "\n"
			}
			if len(glob) < 1999 {
				_, err = s.ChannelMessageSend(m.ChannelID, "```ruby\n"+glob+"```")
				if err != nil {
					fmt.Println(err)
				}
			} else {
				_, err = s.ChannelMessageSend(m.ChannelID, "A.R.S Exceeds 2000 characters (discords message limit) You will need to manually view your A.R.S")
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	}
}

// InspectARS - Displays an Embed with some information.
func (bot *Object) InspectARS(s *discordgo.Session, m *discordgo.MessageCreate, prefix string) {
	if strings.HasPrefix(m.Content, prefix+"inspect ") {
		var ars map[string]string
		trigger := strings.Replace(m.Content, prefix+"inspect ", "", -1)

		io, err := ioutil.ReadFile("autoresponse.json")
		if err != nil {
			fmt.Println(err)
			s.ChannelMessageSend(m.ChannelID, fmt.Sprint(err))
		} else {
			err = json.Unmarshal(io, &ars)
			if err != nil {
				fmt.Println(err)
				s.ChannelMessageSend(m.ChannelID, fmt.Sprint(err))
			}
			var glob string
			for t, r := range ars {
				if t == trigger {
					glob = r
				}
			}
			if len(glob) < 1999 {
				_, err = s.ChannelMessageSend(m.ChannelID, "Response for `"+trigger+"` ```ruby\n"+glob+"```")
				if err != nil {
					fmt.Println(err)
				}
			} else {
				_, err = s.ChannelMessageSend(m.ChannelID, "A.R.S Exceeds 2000 characters (discords message limit) You will need to manually view your A.R.S")
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	}
}
