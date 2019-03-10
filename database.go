// CMPT 315 (Winter 2019)
// Assignment 2
// Author: Jamie Rajewski
//
// This file implements the necessary queries and structures
// required by the server.

package main

import (
	"database/sql"
	"encoding/xml"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type Database struct {
	*sqlx.DB
}

type PersonWrapper struct {
	People *[]Person `json:"people,omitempty" xml:"people,omitempty"`
	Person *Person   `json:"person,omitempty" xml:"person,omitempty"`
	// Only used when getting the form data back for a presenter from this user
	UserID *int `json:"userID,omitempty" xml:"userID,omitempty"`
}

type QuestionWrapper struct {
	Questions *[]Question `json:"questions,omitempty" xml:"questions,omitempty"`
	Question  *Question   `json:"question,omitempty" xml:"question,omitempty"`
}

type AnswerWrapper struct {
	Answers *[]Answer `json:"answers,omitempty" xml:"answers,omitempty"`
	Answer  *Answer   `json:"answer,omitempty" xml:"answer,omitempty"`
}

type PresentationTitleWrapper struct {
	Titles *[]string `json:"titles,omitempty" xml:"titles,omitempty"`
	Title  *string   `json:"title,omitempty" xml:"title,omitempty"`
}

type Person struct {
	XMLName           xml.Name       `json:"-" xml:"person"`
	PersonId          int            `db:"person_id" json:"personId" xml:"personId"`
	LoginId           string         `db:"login_id" json:"loginId,omitempty" xml:"loginId,omitempty"`
	FirstName         string         `db:"first_name" json:"firstName" xml:"firstName"`
	LastName          string         `db:"last_name" json:"lastName" xml:"lastName"`
	PresentationTitle sql.NullString `db:"presentation_title" json:"presentationTitle" xml:"presentationTitle"`
	DateTime          pq.NullTime    `db:"date_time" json:"dateTime" xml:"dateTime"`
	Presenter         bool           `db:"presenter" json:"presenter" xml:"presenter"`
}

type Question struct {
	XMLName      xml.Name `json:"-" xml:"question"`
	QuestionId   int      `db:"question_id" json:"questionId" xml:"questionId"`
	QuestionText string   `db:"question_text" json:"questionText" xml:"questionText"`
	QuestionType string   `db:"question_type" json:"questionType" xml:"questionType"`
}

type Answer struct {
	XMLName     xml.Name `json:"-" xml:"answer"`
	AnswerId    int      `db:"answer_id" json:"answerId,omitempty" xml:"answerId,omitempty"`
	PersonId    int      `db:"person_id" json:"personId,omitempty" xml:"personId,omitempty"`
	PresenterId int      `db:"presenter_id" json:"presenterId,omitempty" xml:"presenterId,omitempty"`
	QuestionId  int      `db:"question_id" json:"questionId,omitempty" xml:"questionId,omitempty"`
	AnswerText  string   `db:"answer_text" json:"answerText,omitempty" xml:"answerText,omitempty"`
}

// SOURCED FROM NICK BOERS' database.go IN LAB 4, CMPT315
// OpenDatabase attempts to open the database specified by 'connect'
// and returns a handle to it
func OpenDatabase(connect string) (*Database, error) {
	db := Database{}
	var err error
	db.DB, err = sqlx.Connect("postgres", connect)
	if err != nil {
		return nil, err
	}

	return &db, nil
}

// GetOnePresenter gets the associated 'person' data of presenter with 'presenterId'
func (db *Database) GetOnePresenter(presenterId int) (PersonWrapper, error) {
	q := `SELECT person_id, first_name, last_name, presentation_title, date_time, presenter 
			FROM people
			WHERE person_id = $1
			AND presenter = true`
	person := Person{}

	if err := db.Get(&person, q, presenterId); err != nil {
		return PersonWrapper{}, err
	}

	return PersonWrapper{Person: &person}, nil
}

// GetPresenters gets the associated 'person' data of all presenters
func (db *Database) GetPresenters() (PersonWrapper, error) {
	q := `SELECT person_id, first_name, last_name, presentation_title, date_time, presenter
			FROM people
			WHERE presenter = true`
	presenters := []Person{}

	if err := db.Select(&presenters, q); err != nil {
		return PersonWrapper{}, err
	}
	return PersonWrapper{People: &presenters}, nil
}

// GetPresentationTitle gets the assoicated presentation title of a presenter
func (db *Database) GetPresentationTitle(presenterId int) (PresentationTitleWrapper, error) {
	q := `SELECT presentation_title FROM people
			WHERE person_id = $1
			AND presenter = true`
	var title string

	if err := db.Get(&title, q, presenterId); err != nil {
		return PresentationTitleWrapper{}, err
	}

	return PresentationTitleWrapper{Title: &title}, nil
}

// GetPresentationTitles gets all presentation titles
func (db *Database) GetPresentationTitles() (PresentationTitleWrapper, error) {
	q := `SELECT presentation_title FROM people
			WHERE presenter = true`
	var titles []string

	if err := db.Select(&titles, q); err != nil {
		return PresentationTitleWrapper{}, err
	}

	return PresentationTitleWrapper{Titles: &titles}, nil
}

// GetQuestion gets a specific question
func (db *Database) GetQuestion(questionId int) (QuestionWrapper, error) {
	q := `SELECT * FROM questions
			WHERE question_id = $1`
	question := Question{}

	if err := db.Get(&question, q, questionId); err != nil {
		return QuestionWrapper{}, err
	}

	return QuestionWrapper{Question: &question}, nil
}

// GetQuestions gets all questions
func (db *Database) GetQuestions() (QuestionWrapper, error) {
	q := `SELECT * FROM questions`
	questions := []Question{}

	if err := db.Select(&questions, q); err != nil {
		return QuestionWrapper{}, err
	}

	return QuestionWrapper{Questions: &questions}, nil
}

// GetPriorResponses gets the prior responses for a reviewer
func (db *Database) GetPriorResponses(personId, presenterId int) (AnswerWrapper, error) {
	q := `SELECT * FROM answers
			WHERE person_id = $1
			AND presenter_id = $2`
	responses := []Answer{}

	if err := db.Select(&responses, q, personId, presenterId); err != nil {
		return AnswerWrapper{}, err
	}

	return AnswerWrapper{Answers: &responses}, nil
}

// PostResponse posts a response for a question by a user on behalf of a presenter
func (db *Database) PostResponse(response Answer) error {
	q := `INSERT INTO answers VALUES
			(
				DEFAULT,
				$1,
				$2,
				$3,
				$4
			)`

	_, err := db.Exec(q, response.PersonId, response.PresenterId,
		response.QuestionId, response.AnswerText)
	if err != nil {
		return err
	}

	return nil
}

// UpdateResponse changes a response for a specific question
func (db *Database) UpdateResponse(response Answer) (AnswerWrapper, error) {
	q := `UPDATE answers
			SET answer_text = $1
			WHERE person_id = $2
			AND presenter_id = $3
			AND question_id = $4`

	if _, err := db.Exec(q, response.AnswerText, response.PersonId,
		response.PresenterId, response.QuestionId); err != nil {
		return AnswerWrapper{}, err
	}

	// Now verify that the record was updated by pulling the new version from the db and return
	q = `SELECT * FROM answers
			WHERE person_id = $1
			AND presenter_id = $2
			AND question_id = $3`
	newResponse := Answer{}

	if err := db.Get(&newResponse, q, response.PersonId,
		response.PresenterId, response.QuestionId); err != nil {
		return AnswerWrapper{}, err
	}

	return AnswerWrapper{Answer: &newResponse}, nil
}

// DeleteResponse deletes the indicated response
func (db *Database) DeleteResponse(answerId int) error {
	q := `DELETE FROM answers
			WHERE answer_id = $1`

	if _, err := db.Exec(q, answerId); err != nil {
		return err
	}

	return nil
}

// GetSerialFromLogin retrieves the associated serial id for the given login id
func (db *Database) GetSerialFromLogin(loginId string) (int, error) {
	q := `SELECT person_id FROM people
			WHERE login_id = $1`
	var personId int

	if err := db.Get(&personId, q, loginId); err != nil {
		return -1, err
	}
	return personId, nil
}

func (db *Database) VerifyMC(mc string, qid int) error {
	// First check if it's even a MC question
	q := `SELECT question_type FROM questions
			WHERE question_id = $1`
	var qType string

	if err := db.Get(&qType, q, qid); err != nil {
		return err
	}

	// If not, end here
	if qType == "la" {
		return nil
	}

	q = `SELECT options FROM mc_options
			WHERE options = $1`
	var result string

	if err := db.Get(&result, q, mc); err != nil {
		return err
	}
	return nil
}

func (db *Database) GetSerialFromAnswer(answerId int) (int, error) {
	q := `SELECT person_id FROM answers
			WHERE answer_id = $1`
	var personId int

	if err := db.Get(&personId, q, answerId); err != nil {
		return -1, err
	}
	return personId, nil
}

// Authenticate checks the validity of the login ID
func (db *Database) Authenticate(loginId string) (Person, error) {
	q := `SELECT * FROM people
			WHERE login_id = $1`
	person := Person{}

	if err := db.Get(&person, q, loginId); err != nil {
		return person, err
	}

	return person, nil
}
