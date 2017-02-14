package main

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func (bot *Object) GuildMemberAdd(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
	// Check to make sure the object is not nil.
	if bot.System != nil {
		// Make sure the Autorole is not empty.
		if bot.System.Autorole != "" {
			// Default: Set the auto role channel to #general (Main Channel)
			var autochan string = m.GuildID
			// Check if Autorole channel is empty.
			if bot.System.Channels.Autorole != "" {
				autochan = bot.System.Channels.Autorole
			}
			// Grab the guild information from State.
			g, err := s.State.Guild(m.GuildID)
			if err != nil {
				fmt.Println(err)
			}
			// Iterate through the guild roles.
			for _, r := range g.Roles {
				// Make sure the role exists!
				if r.Name == bot.System.Autorole {
					// Add the role to their object.
					m.Roles = append(m.Roles, r.ID)
					// Make it official with discord!
					err = s.GuildMemberEdit(m.GuildID, m.User.ID, m.Roles)
					if err == nil {
						s.ChannelMessageSend(autochan, "I have given <@"+m.User.ID+"> the role "+bot.System.Autorole)
					} else {
						s.ChannelMessageSend(autochan, "I don't have permissions to autorole "+m.User.Username)
					}
				}
			}
		}

		// Check to see if greeting is not empty.
		if bot.System.Greeting != "" {
			// Default: Set the greet channel to #general (Main Channel)
			var greetchan string = m.GuildID
			// Check for the {pm} key.
			if strings.Contains(bot.System.Greeting, "{pm}") {
				// Open a Private Channel with the user.
				k, err := s.UserChannelCreate(m.User.ID)
				if err != nil {
					fmt.Println(err)
					return
				}
				// Set greet channel to Private Channel ID.
				greetchan = k.ID
			}
			// Remove the {keys}
			bot.System.Greeting = strings.Replace(bot.System.Greeting, "{pm}", "", -1)
			bot.System.Greeting = strings.Replace(bot.System.Greeting, "{user}", "<@"+m.User.ID+">", -1)
			bot.System.Greeting = strings.Replace(bot.System.Greeting, "{/user}", m.User.Username, -1)
			_, err = s.ChannelMessageSend(greetchan, bot.System.Greeting)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

func (bot *Object) GuildMemberRemove(s *discordgo.Session, m *discordgo.GuildMemberRemove) {
	// Set Bye Message to #general (Main Channel)
	var byechan string = m.GuildID
	// Make sure object is not nil.
	if bot.System != nil {
		// Make sure the Bye Message is not empty.
		if bot.System.ByeMsg != "" {
			// Check for the {pm} Key.
			if strings.Contains(bot.System.ByeMsg, "{pm}") {
				// Open a Private Channel with the user.
				k, err := s.UserChannelCreate(m.User.ID)
				if err != nil {
					fmt.Println(err)
				}
				// Set bye channel to Private Channel ID.
				byechan = k.ID
			}
			// Remove the {keys}
			bot.System.ByeMsg = strings.Replace(bot.System.ByeMsg, "{pm}", "", -1)
			bot.System.ByeMsg = strings.Replace(bot.System.ByeMsg, "{user}", "<@"+m.User.ID+">", -1)
			bot.System.ByeMsg = strings.Replace(bot.System.ByeMsg, "{/user}", m.User.Username, -1)
			_, err := s.ChannelMessageSend(byechan, bot.System.ByeMsg)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

func (bot *Object) Initiate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if bot.System == nil {
		bot.Register(s, m)
	} else {
		bot.Task(s, m)
	}
	fmt.Println(m.Author.ID, m.Content)
}
