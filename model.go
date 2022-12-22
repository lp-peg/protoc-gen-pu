package main

// PlantUML class.
type class struct {
	Name    string
	Members []member
}

type member struct {
	Name       string
	Type       string
	IsRepeated bool
}

type reference struct {
	From class
	To   class
}