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
	From string
	To   string
}

type skinparam struct {
	Param string
	Value string
}
