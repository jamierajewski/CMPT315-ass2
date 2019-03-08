-- +goose Up
-- SQL in this section is executed when the migration is applied.
INSERT INTO people (person_id, login_id, first_name, last_name, presentation_title, date_time, presenter)
VALUES (
    DEFAULT,
    'YvutCBLA',
    'Jamie',
    'Rajewski',
    'React',
    '2019-06-22 19:10:25-07',
    true
);

INSERT INTO people (person_id, login_id, first_name, last_name, presenter)
VALUES (
    DEFAULT,
    'LchajfPu',
    'Nick',
    'Boers',
    false
);

INSERT INTO people (person_id, login_id, first_name, last_name, presentation_title, date_time, presenter)
VALUES (
    DEFAULT,
    'REwPZokD',
    'Lukas',
    'Jenks',
    'Flask and Python',
    '2019-06-22 17:10:25-07',
    true
);

INSERT INTO questions (question_id, question_text, question_type)
VALUES (
    DEFAULT,
    'Preparedness: the presenter was adequately prepared.',
    'mc'
);

INSERT INTO questions (question_id, question_text, question_type)
VALUES (
    DEFAULT,
    'Organization: the presentation material was arranged logically.',
    'mc'
);

INSERT INTO questions (question_id, question_text, question_type)
VALUES (
    DEFAULT,
    'Correctness: the presented facts were correct (to the best of your knowledge).',
    'mc'
);

INSERT INTO questions (question_id, question_text, question_type)
VALUES (
    DEFAULT,
    'Visualization: the visual material included appropriate content/fonts/graphics.',
    'mc'
);

INSERT INTO questions (question_id, question_text, question_type)
VALUES (
    DEFAULT,
    'General introduction: the presentation clearly introduced the broad area containing the topic.',
    'mc'
);

INSERT INTO questions (question_id, question_text, question_type)
VALUES (
    DEFAULT,
    'Motivation: the presentation clearly motivated the specific topic in the context of the broad area.',
    'mc'
);

INSERT INTO questions (question_id, question_text, question_type)
VALUES (
    DEFAULT,
    'Introduction: the presentation clearly introduced the specific topic.',
    'mc'
);

INSERT INTO questions (question_id, question_text, question_type)
VALUES (
    DEFAULT,
    'Tutorial/demonstration: the tutorial/demonstration improved your understanding of the specific topic.',
    'mc'
);

INSERT INTO questions (question_id, question_text, question_type)
VALUES (
    DEFAULT,
    'Multiple-choice questions: at least three multiple-choice questions assessed your understanding of the presented content.',
    'mc'
);

INSERT INTO questions (question_id, question_text, question_type)
VALUES (
    DEFAULT,
    'Answers: the presenters answers to questions were satisfying.',
    'mc'
);

INSERT INTO questions (question_id, question_text, question_type)
VALUES (
    DEFAULT,
    'Provide any comments for the presenter.',
    'la'
);

INSERT INTO questions (question_id, question_text, question_type)
VALUES (
    DEFAULT,
    'Provide any comments for your instructor.',
    'la'
);

INSERT INTO mc_options (options)
VALUES
    ('1'),
    ('2'),
    ('3'),
    ('4'),
    ('5'),
    ('0');

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.

DELETE FROM answers;
DELETE FROM questions;
DELETE FROM people;
DELETE FROM mc_options;