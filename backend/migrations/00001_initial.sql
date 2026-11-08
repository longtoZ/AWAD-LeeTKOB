-- +goose Up
CREATE TABLE users (
    id uuid PRIMARY KEY,
    email text NOT NULL UNIQUE,
    password_hash text NOT NULL,
    role text NOT NULL CHECK (role IN ('student', 'ta', 'admin'))
);
CREATE TABLE sessions (token text PRIMARY KEY, data bytea NOT NULL, expiry timestamptz NOT NULL);
CREATE INDEX sessions_expiry_idx ON sessions (expiry);
CREATE TABLE labs (id uuid PRIMARY KEY, title text NOT NULL);
CREATE TABLE enrollments (lab_id uuid NOT NULL REFERENCES labs(id), student_id uuid NOT NULL REFERENCES users(id), PRIMARY KEY (lab_id, student_id));
CREATE TABLE submissions (id uuid PRIMARY KEY, lab_id uuid NOT NULL REFERENCES labs(id), student_id uuid NOT NULL REFERENCES users(id), source text NOT NULL, created_at timestamptz NOT NULL);
CREATE TABLE jobs (id uuid PRIMARY KEY, submission_id uuid NOT NULL UNIQUE REFERENCES submissions(id), status text NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'processing', 'completed', 'failed')), created_at timestamptz NOT NULL DEFAULT now());
CREATE TABLE idempotency_keys (student_id uuid NOT NULL REFERENCES users(id), lab_id uuid NOT NULL REFERENCES labs(id), key text NOT NULL, fingerprint bytea NOT NULL, response json NOT NULL, created_at timestamptz NOT NULL DEFAULT now(), PRIMARY KEY (student_id, lab_id, key));

-- +goose Down
DROP TABLE idempotency_keys, jobs, submissions, enrollments, labs, sessions, users;
