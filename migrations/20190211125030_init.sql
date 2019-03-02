-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE people (
    person_id SERIAL PRIMARY KEY,

    login_id VARCHAR (8) NOT NULL UNIQUE,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    presentation_title TEXT,
    date_time TIMESTAMP UNIQUE,
    presenter BOOLEAN NOT NULL,
    CONSTRAINT presenter_not_null CHECK (
        (presenter = true AND date_time IS NOT NULL AND presentation_title IS NOT NULL)
        OR
        (presenter = false AND date_time IS NULL AND presentation_title IS NULL)
    )
);

CREATE TABLE questions (
    question_id SERIAL PRIMARY KEY,

    question_text TEXT NOT NULL UNIQUE,
    question_type VARCHAR (2) NOT NULL,
    CONSTRAINT type_verifier CHECK (
        (question_type = 'mc' OR question_type = 'la')
    )
);

CREATE TABLE answers (
    answer_id SERIAL PRIMARY KEY,

    person_id INTEGER NOT NULL REFERENCES people(person_id),
    presenter_id INTEGER NOT NULL REFERENCES people(person_id),
    question_id INTEGER NOT NULL REFERENCES questions(question_id),
    answer_text TEXT,

    UNIQUE(person_id, presenter_id, question_id)
);

CREATE TABLE mc_options (
    options TEXT NOT NULL UNIQUE
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE answers;
DROP TABLE questions;
DROP TABLE people;
DROP TABLE mc_options;