// CMPT 315 (Winter 2019)
// Assignment 2
// Author: Jamie Rajewski
//
// This file implements the backend API that communicates with
// a PostGreSQL database
package main

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

var port int

type Handler struct {
	*Database
}

type keyString string

type customEncoder interface {
	Encode(v interface{}) error
}

type ErrorMessage struct {
	XMLName xml.Name `json:"-" xml:"errorMessage"`
	Code    int      `json:"code" xml:"code"`
	Status  string   `json:"status" xml:"status"`
	Details string   `json:"details" xml:"details"`
}

// ADD LOGGER MESSAGE TYPE

func (h *Handler) handleGetPresenter(w http.ResponseWriter, r *http.Request) {
	encoder := r.Context().Value(keyString("Encoder")).(customEncoder)

	data := mux.Vars(r)
	presenterIdText, ok := data["id"]
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(ErrorMessage{Code: http.StatusInternalServerError,
			Status:  http.StatusText(http.StatusInternalServerError),
			Details: "Could not get variable from URI"})
		return
	}

	presenterId, err := strconv.Atoi(presenterIdText)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		encoder.Encode(ErrorMessage{Code: http.StatusBadRequest,
			Status:  http.StatusText(http.StatusBadRequest),
			Details: "Could not convert presenterId to integer"})
		return
	}

	presenter, err := h.GetOnePresenter(presenterId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		encoder.Encode(ErrorMessage{Code: http.StatusNotFound,
			Status:  http.StatusText(http.StatusNotFound),
			Details: "Could not find the presenter in the database"})
		return
	}
	// Attach the reviewer's user ID
	userID := r.Context().Value(keyString("PersonId")).(int)
	presenter.UserID = &userID

	encoder.Encode(presenter)
}

func (h *Handler) handleGetPresenters(w http.ResponseWriter, r *http.Request) {
	encoder := r.Context().Value(keyString("Encoder")).(customEncoder)

	presenters, err := h.GetPresenters()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(ErrorMessage{Code: http.StatusInternalServerError,
			Status:  http.StatusText(http.StatusInternalServerError),
			Details: "Could not retrieve presenters from database"})
		return
	}

	encoder.Encode(presenters)
}

func (h *Handler) handleGetPresentationTitle(w http.ResponseWriter, r *http.Request) {
	encoder := r.Context().Value(keyString("Encoder")).(customEncoder)

	data := mux.Vars(r)
	presenterIdText, ok := data["id"]
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(ErrorMessage{Code: http.StatusInternalServerError,
			Status:  http.StatusText(http.StatusInternalServerError),
			Details: "Could not get variable from URI"})
		return
	}

	presenterId, err := strconv.Atoi(presenterIdText)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		encoder.Encode(ErrorMessage{Code: http.StatusBadRequest,
			Status:  http.StatusText(http.StatusBadRequest),
			Details: "Could not convert presenterId to integer"})
		return
	}

	title, err := h.GetPresentationTitle(presenterId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		encoder.Encode(ErrorMessage{Code: http.StatusNotFound,
			Status:  http.StatusText(http.StatusNotFound),
			Details: "Could not find the presentation title for the presenter in the database"})
		return
	}

	encoder.Encode(title)
}

func (h *Handler) handleGetPresentationTitles(w http.ResponseWriter, r *http.Request) {
	encoder := r.Context().Value(keyString("Encoder")).(customEncoder)

	titles, err := h.GetPresentationTitles()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(ErrorMessage{Code: http.StatusInternalServerError,
			Status:  http.StatusText(http.StatusInternalServerError),
			Details: "Could not retrieve presentation titles from database"})
		return
	}

	encoder.Encode(titles)
}

func (h *Handler) handleGetQuestion(w http.ResponseWriter, r *http.Request) {
	encoder := r.Context().Value(keyString("Encoder")).(customEncoder)

	data := mux.Vars(r)
	questionIdText, ok := data["id"]
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(ErrorMessage{Code: http.StatusInternalServerError,
			Status:  http.StatusText(http.StatusInternalServerError),
			Details: "Could not get question from URI"})
		return
	}

	questionId, err := strconv.Atoi(questionIdText)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		encoder.Encode(ErrorMessage{Code: http.StatusBadRequest,
			Status:  http.StatusText(http.StatusBadRequest),
			Details: "Could not convert questionId to integer"})
		return
	}

	question, err := h.GetQuestion(questionId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		encoder.Encode(ErrorMessage{Code: http.StatusNotFound,
			Status:  http.StatusText(http.StatusNotFound),
			Details: "Could not find the question in the database"})
		return
	}

	encoder.Encode(question)
}

func (h *Handler) handleGetQuestions(w http.ResponseWriter, r *http.Request) {
	encoder := r.Context().Value(keyString("Encoder")).(customEncoder)

	questions, err := h.GetQuestions()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(ErrorMessage{Code: http.StatusInternalServerError,
			Status:  http.StatusText(http.StatusInternalServerError),
			Details: "Could not retrieve questions from the database"})
		return
	}

	encoder.Encode(questions)
}

func (h *Handler) handleGetPriorResponses(w http.ResponseWriter, r *http.Request) {
	encoder := r.Context().Value(keyString("Encoder")).(customEncoder)

	uncastPersonId := r.Context().Value(keyString("PersonId"))
	personId, ok := uncastPersonId.(int)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(ErrorMessage{Code: http.StatusInternalServerError,
			Status:  http.StatusText(http.StatusInternalServerError),
			Details: "Failed to obtain PersonId from context"})
		return
	}

	data := mux.Vars(r)
	presenterIdText, ok := data["id"]
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(ErrorMessage{Code: http.StatusInternalServerError,
			Status:  http.StatusText(http.StatusInternalServerError),
			Details: "Could not get personId from URI"})
		return
	}

	presenterId, err := strconv.Atoi(presenterIdText)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		encoder.Encode(ErrorMessage{Code: http.StatusBadRequest,
			Status:  http.StatusText(http.StatusBadRequest),
			Details: "Could not convert personId to integer"})
		return
	}

	var response AnswerWrapper
	response, err = h.GetPriorResponses(personId, presenterId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		encoder.Encode(ErrorMessage{Code: http.StatusNotFound,
			Status:  http.StatusText(http.StatusNotFound),
			Details: "Could not retrieve prior responses from the database"})
		return
	}

	encoder.Encode(response)
}

func (h *Handler) handlePostResponse(w http.ResponseWriter, r *http.Request) {
	encoder := r.Context().Value(keyString("Encoder")).(customEncoder)

	uncastPersonId := r.Context().Value(keyString("PersonId"))
	personId, ok := uncastPersonId.(int)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(ErrorMessage{Code: http.StatusInternalServerError,
			Status:  http.StatusText(http.StatusInternalServerError),
			Details: "Failed to obtain PersonId from context"})
		return
	}
	response := Answer{}
	defer r.Body.Close()

	// Assume that content will be sent as JSON from the client for simplicity
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		encoder.Encode(ErrorMessage{Code: http.StatusBadRequest,
			Status:  http.StatusText(http.StatusBadRequest),
			Details: "Failed to decode; Request body must be in JSON format (Is content-type application/json?)"})
		return
	}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&response)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		encoder.Encode(ErrorMessage{Code: http.StatusBadRequest,
			Status:  http.StatusText(http.StatusBadRequest),
			Details: "Could not decode body of request prior to posting"})
		return
	}

	err = h.VerifyMC(response.AnswerText, response.QuestionId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		encoder.Encode(ErrorMessage{Code: http.StatusBadRequest,
			Status:  http.StatusText(http.StatusBadRequest),
			Details: "Invalid multiple choice selection"})
		return
	}

	response.PersonId = personId
	err = h.PostResponse(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(ErrorMessage{Code: http.StatusInternalServerError,
			Status:  http.StatusText(http.StatusInternalServerError),
			Details: "Could not post response to database"})
		return
	}
}

func (h *Handler) handleUpdateResponse(w http.ResponseWriter, r *http.Request) {
	encoder := r.Context().Value(keyString("Encoder")).(customEncoder)

	uncastPersonId := r.Context().Value(keyString("PersonId"))
	personId, ok := uncastPersonId.(int)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(ErrorMessage{Code: http.StatusInternalServerError,
			Status:  http.StatusText(http.StatusInternalServerError),
			Details: "Failed to obtain PersonId from context"})
		return
	}
	response := Answer{}
	defer r.Body.Close()

	// Assume that content will be sent as JSON from the client for simplicity
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		encoder.Encode(ErrorMessage{Code: http.StatusBadRequest,
			Status:  http.StatusText(http.StatusBadRequest),
			Details: "Failed to decode; Request body must be in JSON format (Is content-type application/json?)"})
		return
	}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&response)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		encoder.Encode(ErrorMessage{Code: http.StatusBadRequest,
			Status:  http.StatusText(http.StatusBadRequest),
			Details: "Could not decode body of request prior to posting"})
		return
	}

	// This will check if it is a MC and verify; if not, do nothing and continue
	err = h.VerifyMC(response.AnswerText, response.QuestionId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		encoder.Encode(ErrorMessage{Code: http.StatusBadRequest,
			Status:  http.StatusText(http.StatusBadRequest),
			Details: "Invalid multiple choice selection"})
		return
	}

	response.PersonId = personId

	var result AnswerWrapper
	result, err = h.UpdateResponse(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(ErrorMessage{Code: http.StatusInternalServerError,
			Status:  http.StatusText(http.StatusInternalServerError),
			Details: "Failed to update response"})
		return
	}
	encoder.Encode(result)
}

func (h *Handler) handleDeleteResponse(w http.ResponseWriter, r *http.Request) {
	encoder := r.Context().Value(keyString("Encoder")).(customEncoder)

	uncastPersonId := r.Context().Value(keyString("PersonId"))
	personId, ok := uncastPersonId.(int)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(ErrorMessage{Code: http.StatusInternalServerError,
			Status:  http.StatusText(http.StatusInternalServerError),
			Details: "Failed to obtain PersonId from context"})
		return
	}

	data := mux.Vars(r)
	answerIdText, ok := data["id"]
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(ErrorMessage{Code: http.StatusInternalServerError,
			Status:  http.StatusText(http.StatusInternalServerError),
			Details: "Could not get answerId from URI"})
		return
	}
	answerId, err := strconv.Atoi(answerIdText)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		encoder.Encode(ErrorMessage{Code: http.StatusBadRequest,
			Status:  http.StatusText(http.StatusBadRequest),
			Details: "Could not convert answerId to integer"})
		return
	}

	var answerPersonId int
	answerPersonId, err = h.GetSerialFromAnswer(answerId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(ErrorMessage{Code: http.StatusInternalServerError,
			Status:  http.StatusText(http.StatusInternalServerError),
			Details: "Failed to obtain PersonId from AnswerId"})
		return
	}

	if personId != answerPersonId {
		w.WriteHeader(http.StatusForbidden)
		encoder.Encode(ErrorMessage{Code: http.StatusForbidden,
			Status:  http.StatusText(http.StatusForbidden),
			Details: "You do not have permission to delete a response for this user"})
		return
	}

	err = h.DeleteResponse(answerId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(ErrorMessage{Code: http.StatusInternalServerError,
			Status:  http.StatusText(http.StatusInternalServerError),
			Details: "Failed to delete response"})
		return
	}
}

func (h *Handler) addEncoderMw(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		encodeType := r.Header.Get("Accept")
		var encoder interface{}
		switch encodeType {
		case "application/xml":
			encoder = xml.NewEncoder(w)
		// If anything other than xml, default to json
		default:
			encoder = json.NewEncoder(w)
		}
		var key keyString = "Encoder"
		ctx := context.WithValue(r.Context(), key, encoder)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *Handler) authenticateMw(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Check for errors once here; in the actual handlers, skip error-checking
		// stage since it was already done in the middleware here
		uncastEncoder := r.Context().Value(keyString("Encoder"))
		encoder, ok := uncastEncoder.(customEncoder)
		if !ok {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			token := strings.TrimPrefix(authHeader, "Bearer ")
			if len(token) != len(authHeader) {
				_, err := h.Authenticate(token)
				if err != nil {
					w.WriteHeader(http.StatusUnauthorized)
					w.Header().Add("WWW-Authenticate", "Bearer realm=\"Access to survey\"")
					encoder.Encode(ErrorMessage{Code: http.StatusUnauthorized,
						Status:  http.StatusText(http.StatusUnauthorized),
						Details: "This resource requires a valid bearer token"})
					return
				}
				next.ServeHTTP(w, r)
			} else {
				w.WriteHeader(http.StatusUnauthorized)
				w.Header().Add("WWW-Authenticate", "Bearer realm=\"Access to survey\"")
				encoder.Encode(ErrorMessage{Code: http.StatusUnauthorized,
					Status:  http.StatusText(http.StatusUnauthorized),
					Details: "This resource requires a valid bearer token"})
				return
			}
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			w.Header().Add("WWW-Authenticate", "Bearer realm=\"Access to survey\"")
			encoder.Encode(ErrorMessage{Code: http.StatusUnauthorized,
				Status:  http.StatusText(http.StatusUnauthorized),
				Details: "This resource requires a valid bearer token"})
			return
		}
	})
}

func (h *Handler) attachUserIdMw(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Because the MWs are used in order, the error-handling has already been
		// done previously
		encoder := r.Context().Value(keyString("Encoder")).(customEncoder)

		authHeader := r.Header.Get("Authorization")
		// FIX THIS TO NOT GO OUT OF RANGE by using TrimPrefix()
		token := strings.Split(authHeader, " ")[1]

		personId, err := h.GetSerialFromLogin(token)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			encoder.Encode(ErrorMessage{Code: http.StatusInternalServerError,
				Status:  http.StatusText(http.StatusInternalServerError),
				Details: "Failed to retrieve personId"})
			return
		}
		var key keyString = "PersonId"
		ctx := context.WithValue(r.Context(), key, personId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// SOURCED FROM NICK BOERS' main.go IN LAB 4, CMPT315
// Use the init to parse command line arguments
func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `usage: %s [-p port]

Options:
`, path.Base(os.Args[0]))
		flag.PrintDefaults()
	}

	flag.IntVar(&port, "p", 8080, "port")

	flag.Parse()
}

func main() {
	log.SetOutput(os.Stdout)

	connect := "dbname=assign user=postgres host=localhost port=5432 sslmode=disable"
	db, err := OpenDatabase(connect)
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}
	defer db.Close()

	handlers := Handler{
		db,
	}

	// High-level router
	router := mux.NewRouter()

	// Mid-level (API) router
	apiRouter := router.PathPrefix("/api/v1").Subrouter()

	// Low-level routers
	peopleRouter := apiRouter.PathPrefix("/people").Subrouter()
	questionsRouter := apiRouter.PathPrefix("/questions").Subrouter()
	answersRouter := apiRouter.PathPrefix("/answers").Subrouter()

	peopleRouter.HandleFunc("/presenters/{id:[0-9]+}", handlers.handleGetPresenter).Methods("GET")
	peopleRouter.HandleFunc("/presenters", handlers.handleGetPresenters).Methods("GET")
	peopleRouter.HandleFunc("/presentationTitles/{id:[0-9]+}", handlers.handleGetPresentationTitle).Methods("GET")
	peopleRouter.HandleFunc("/presentationTitles", handlers.handleGetPresentationTitles).Methods("GET")

	questionsRouter.HandleFunc("/{id:[0-9]+}", handlers.handleGetQuestion).Methods("GET")
	questionsRouter.HandleFunc("/", handlers.handleGetQuestions).Methods("GET")

	answersRouter.HandleFunc("/{id:[0-9]+}", handlers.handleGetPriorResponses).Methods("GET")
	answersRouter.HandleFunc("/", handlers.handlePostResponse).Methods("POST")
	answersRouter.HandleFunc("/", handlers.handleUpdateResponse).Methods("PUT")
	answersRouter.HandleFunc("/{id:[0-9]+}", handlers.handleDeleteResponse).Methods("DELETE")

	router.PathPrefix("/").Handler(http.FileServer(http.Dir("dist")))

	router.Use(handlers.addEncoderMw)
	apiRouter.Use(handlers.authenticateMw)
	apiRouter.Use(handlers.attachUserIdMw)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), router))
}
