package main

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// ParseServer - Auto Response Parse Server keys
func (bot *Object) ParseServer(s *discordgo.Session, m *discordgo.MessageCreate, trigger, response string) string {
	// Show channel name (with mention) #general, #lobby etc..
	if strings.Contains(response, "{chan}") {
		response = strings.Replace(response, "{chan}", "<#"+m.ChannelID+">", -1)
	}
	// Show channel topic.
	if strings.Contains(response, "{topic}") {
		ch, err := s.State.Channel(m.ChannelID)
		if err != nil {
			fmt.Println(err)
			return ""
		}
		response = strings.Replace(response, "{topic}", ch.Topic, -1)
	}

	// List all the roles in a guild.
	if strings.Contains(response, "{listroles}") {
		var glob string
		g, err := s.State.Guild(bot.Guild)
		if err != nil {
			fmt.Println(err)
			return ""
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
	return response
}

// ParseDirection - Auto Response Parse direction (Last Stage) (pm, redirect or basic channel)
func (bot *Object) ParseDirection(s *discordgo.Session, m *discordgo.MessageCreate, trigger, response string) {
	var ChannelID string
	ChannelID = m.ChannelID
	if strings.Contains(response, "{redirect:") {
		ch := strings.Split(response, "{redirect:")[1]
		ch = strings.Split(ch, "}")[0]
		g, err := s.State.Guild(bot.Guild)
		if err != nil {
			fmt.Println(err)
		}
		for _, c := range g.Channels {
			if c.ID == ch {
				ChannelID = c.ID
			}
		}
		response = strings.Replace(response, "{redirect:"+ch+"}", "", -1)
	}

	if strings.Contains(response, "{pm}") {
		k, err := s.UserChannelCreate(m.Author.ID)
		if err != nil {
			fmt.Println(err)
			return
		}
		ChannelID = k.ID
		response = strings.Replace(response, "{pm}", "", -1)
	}

	_, err := s.ChannelMessageSend(ChannelID, response)
	if err != nil {
		fmt.Println(err)
	}
}
