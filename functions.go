package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/bwmarrin/discordgo"
)

func IsManager(s *discordgo.Session, GuildID string, AuthorID string) bool {
	// Check the user permissions of the guild.
	perms, err := s.State.UserChannelPermissions(AuthorID, GuildID)
	if err == nil {
		if (perms & discordgo.PermissionManageServer) > 0 {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

func (bot *Object) Save() {
	for {
		<-time.After(5 * time.Minute)
		js, err := json.MarshalIndent(bot, "", "  ")
		if err == nil {
			ioutil.WriteFile("config.json", js, 0777)
		}
	}
}

func GetRoleID(s *discordgo.Session, GuildID string, role string) string {
	var id string
	r, err := s.State.Guild(GuildID)
	if err == nil {
		for _, v := range r.Roles {
			if v.Name == role {
				id = v.ID
			}
		}
	}
	return id
}

func MemberHasRole(s *discordgo.Session, GuildID string, AuthorID string, role string) bool {
	therole := GetRoleID(s, GuildID, role)
	z, err := s.State.Member(GuildID, AuthorID)
	if err != nil {
		z, err = s.GuildMember(GuildID, AuthorID)
		if err != nil {
			fmt.Println("Error ->", err)
			return false
		}
	}
	for r := range z.Roles {
		if therole == z.Roles[r] {
			return true
		}
	}
	return false
}

func (bot *Object) Register(s *discordgo.Session, m *discordgo.MessageCreate) {
	// check and make sure the server already exists in my collection.
	if bot.System != nil {
		return
	}
	c, err := s.State.Channel(m.ChannelID)
	if err != nil {
		fmt.Println(err)
		return
	}

	bot.Guild = c.GuildID
	chn := &Channels{
		Autorole: "",
		Greeting: "",
		ByeMsg:   "",
	}

	// Create a new Info pointer.
	info := &System{
		Prefix:   ".",
		Autorole: "",
		Greeting: "",
		ByeMsg:   "",
		Channels: chn,
	}
	// Add our Info object to the bot map.
	bot.System = info
}

func (bot *Object) Task(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Don't track the bots messages.
	if m.Author.ID == s.State.User.ID {
		return
	}

	if bot.System == nil {
		return
	}
	// Create a new pointer.
	msg := &Messages{
		ID:        m.ID,
		Author:    m.Author.ID,
		Channel:   m.ChannelID,
		Timestamp: time.Now().Unix(),
	}
	// Add this Message to our Info object.
	bot.System.Messages = append(bot.System.Messages, msg)
}
