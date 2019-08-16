package service

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/dfense/appchallenge/client"
	log "github.com/sirupsen/logrus"
)

var (
	ErrInvalidEndpoint = errors.New("invalid endpoint")
	ErrDequeueFailed   = errors.New("dequeue failed on namecache")
)

// errorResponse - struct to use on http json body return error
type errorResponse struct {
	Code        int    `json:"code"`        // http.Status standard error code
	Description string `json:"description"` // err.Error() verbose description of what went wrong
}

// NewServices - defines all the REST API endpoints, URLs, functions, and starts
// the http server. It returns a reference so it can be gracefully stopped by the caller
//
// --- extrememly simple http rest server. ---
// TODO defer to a more elegant mux/router http library for adding more routes,
// such as gorilla or httprouter (lightweight/powerful)
func NewHttpService(address string) *http.Server {

	// TODO store address, and verify if a server has already been started
	// on that address before... return error before trying to start
	mux := http.NewServeMux()
	mux.HandleFunc("/", nameJoke)
	srv := &http.Server{
		Addr:    address,
		Handler: mux,
	}
	go func() {
		// throw fatal error immediately on startup. no need to go further
		// since the only reason this program exists
		log.Infof("Starting HTTP server on address: %s", address)
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe(): %s", err)
		}
	}()
	return srv
}

// nameJokes - GET handler for service. If no url match or improper VERB return error
// the ONLY endpoint this server will handle !!
func nameJoke(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		if r.URL.Path != "/" { // handle the root url only
			WriteJsonErrorResponse(w, http.StatusBadRequest, ErrInvalidEndpoint.Error())
			return
		}

		person, err := GetNameCache().Dequeue()
		if err != nil {
			WriteJsonErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// Serve the resource.
		joke, err := client.GetJoke(person.First, person.Last, []string{"nerdy"})
		if err != nil {
			log.Errorf("Error on GetJoke %s", err)
			WriteJsonErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		WriteJsonResponse(w, "GET", joke)

	case http.MethodPost:
		WriteJsonErrorResponse(w, http.StatusBadRequest, "POST not allowed on this endpoint")
		return
	case http.MethodPut:
		WriteJsonErrorResponse(w, http.StatusBadRequest, "PUT not allowed on this endpoint")
		return
	case http.MethodDelete:
		WriteJsonErrorResponse(w, http.StatusBadRequest, "DELETE not allowed on this endpoint")
		return
	default:
		WriteJsonErrorResponse(w, http.StatusBadRequest, "HTTP Method not recognized. Only GET is allowed")
		return
	}
}

// Write a formatted error response using the ResponseWriter
func WriteJsonErrorResponse(w http.ResponseWriter, errCode int, errDescription string) {

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(errCode)
	data := errorResponse{
		Code:        errCode,
		Description: errDescription,
	}

	// provide server side log of error
	log.Error(errDescription)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Error(err)
	}
}

// WriteJsonResponse - write a formatted success response to ResponseWriter.
// requestType is either "GET", "PUT", "POST", or "DELETE".
func WriteJsonResponse(w http.ResponseWriter, requestType string, data interface{}) {

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if requestType == "POST" {
		w.WriteHeader(http.StatusCreated)
		return
	} else {
		w.WriteHeader(http.StatusOK)
	}
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Error(err)
	}
}
