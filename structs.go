package main

// Object - The main structure for the system
// Token: Discord Bot Token
// Guild: Your Guild ID
// System: Collection of Settings.
type Object struct {
	Token  string
	Guild  string
	System *System
}

// System - The main storage for preferences.
// Prefix: Command Prefix.
// Autorole: Valid Discord Role to give members on join.
// Greeting: The welcome message
// Channels: Collection of Channels (Default is #general)
// Messages: Collection of messages.
type System struct {
	Prefix   string
	Autorole string
	Greeting string
	ByeMsg   string
	Channels *Channels
	Messages []*Messages
}

// Channels - The channels storage for preferences
// Autorole: Channel to message for auto role.
// Greetings: Channel for Greet Message (if {pm} is not found)
// ByeMsg: Channel for Bye Message (if {pm} is not found)
type Channels struct {
	Autorole string
	Greeting string
	ByeMsg   string
}

// Messages - The main collection of stored messages.
// ID: Message ID
// Author: Author ID
// Channel: Channel ID
// Timestamp: Message sent timestamp.
type Messages struct {
	ID        string
	Author    string
	Channel   string
	Timestamp int64
}
