package appchallenge

import "encoding/json"

// Name - struct to hold result from the name requesting service @https://uinames.com
type Person struct {
	First  string `json:"name"`
	Last   string `json:"surname"`
	Gender string `json:"gender"`
	Region string `json:"region"`
}

// JokeResult - struct to hold the result from
// the chuck norris joke service @http://api.icndb.com/jokes/random
type JokeResult struct {
	Type   string          `json:"type"`
	Result json.RawMessage `json:"value"` // <- keep generic to handle different result bodies
}

type ErrorMessage struct {
	GenericMessage string `json:"value"`
}

// JokeValue - result of the joke detail inside JokeResult above
type JokeValue struct {
	Id         int      `json:"id"`
	Joke       string   `json:"joke"`
	Categories []string `json:"categories"`
}
