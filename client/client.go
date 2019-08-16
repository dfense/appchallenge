// client - package to contain all http client calls to external services
package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/dfense/appchallenge"
)

const (
	jokeService = "http://api.icndb.com/jokes/random"
	nameService = "https://uinames.com/api/?amount=%d"
)

// list of valid categories
// TODO query from icndb to get live list
var (

	// create a map of eligable category names
	empty           = struct{}{}
	categoriesTypes = map[string]struct{}{
		"nerdy":    empty,
		"explicit": empty,
	}
	ErrorsInvalidCategory = errors.New("invalid category specified: see http://www.icndb.com/api")

	// amount of time to wait on HTTP GET to services
	timeoutSeconds time.Duration = 5 * time.Second
)

// GetJoke - call external web rest service, and get a random chuck norris joke.
// specify parameters - fname and/or lname to substitute "chuck norris"
//         parameter  - []categories to choose from 1+ joke context/categories
func GetJoke(fname, lname string, categories []string) (*appchallenge.JokeResult, error) {

	for _, cat := range categories {
		if _, ok := categoriesTypes[cat]; !ok {
			return nil, ErrorsInvalidCategory
		}
	}

	// result data struct
	data := &appchallenge.JokeResult{}

	// make HTTP GET for joke. good for https and mutual trust certs
	tr := &http.Transport{
		IdleConnTimeout: timeoutSeconds,
	}

	// create the http request
	req, err := http.NewRequest("GET", jokeService, nil)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	// fill in querystring params
	q := req.URL.Query()
	q.Add("firstName", fname)
	q.Add("lastName", lname)
	// flatten out comma delimited categories strings
	q.Add("limitTo", "["+strings.Join(categories, ",")+"]")
	req.URL.RawQuery = q.Encode()

	// make request
	client := &http.Client{Transport: tr}
	resp, err := client.Get(req.URL.String())
	if err != nil {
		log.Errorf("Error: %s", err)
		return nil, err
	}

	// get first pass json info
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(data)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	// grab secondary json joke object or error string
	// TODO thought i would do more based on what the data.Type was, but
	// not real useful if i just return the raw message back to client.
	// seems sufficient, so can remove unless a reason to be smarter...
	switch data.Type {
	case "success":
		var jokeDetail appchallenge.JokeValue

		if err := json.Unmarshal([]byte(data.Result), &jokeDetail); err != nil {
			log.Error(err)
			return nil, err
		}

	default:
		var e string
		// fmt.Printf("Value %s\n", string([]byte(data.Result)))
		if err := json.Unmarshal([]byte(data.Result), &e); err != nil {
			log.Error(err)
			return nil, err
		}
	}

	defer resp.Body.Close()
	return data, err
}

// GetNames - convenience method to retrieve a list of people from external service
// params
// qty - number of records requested for service
// return
// []appchallenge.Person - list of people retrieved
// error - any errors that occured, and will have empty Person slice
func GetNames(qty int) ([]appchallenge.Person, error) {

	// make an slice of person
	var people []appchallenge.Person

	client := http.Client{
		Timeout: timeoutSeconds,
	}
	// create the http request
	q := fmt.Sprintf(nameService, qty)
	resp, err := client.Get(q)
	if err != nil {
		return people, err
	}

	// get first pass json info
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&people)
	if err != nil {
		log.Error(err)
	}
	return people, err
}
