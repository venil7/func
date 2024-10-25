package main

type Item struct {
	Id          int
	Deleted     bool
	Type        string
	By          string
	Time        int
	Text        string
	Dead        bool
	Parent      int
	Poll        int
	Kids        []int
	Url         string
	Score       int
	Title       string
	Parts       []int
	Descendants int
}
