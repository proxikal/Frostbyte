package main

type Object struct {
	Token  string
	Guild  string
	System *System
}

type System struct {
	Prefix   string
	Autorole string
	Greeting string
	ByeMsg   string
	Channels *Channels
	Messages []*Messages
}

type Channels struct {
	Autorole string
	Greeting string
	ByeMsg   string
}

type Messages struct {
	ID        string
	Author    string
	Channel   string
	Timestamp int64
}
