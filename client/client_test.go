package client

import (
	"testing"
)

func TestNames(t *testing.T) {
	qty := 5
	names, err := GetNames(qty)
	if err != nil {
		t.Error(err)
	}

	if len(names) != qty {
		t.Errorf("Expected [%d] received [%d]\n", qty, len(names))
	}
}

func TestJokes(t *testing.T) {
	joke, err := GetJoke("a", "b", []string{"nerd"})
	if err != ErrorsInvalidCategory {
		t.Error("invalid category error should have occurred")
	}

	joke, err = GetJoke("y", "z", []string{"nerdy"})
	if err != nil {
		t.Error(err)
	}

	if joke.Type != "success" {
		t.Errorf("expected [%s] got [%s]", "success", joke.Type)
	}
}
