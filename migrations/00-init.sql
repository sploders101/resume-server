PRAGMA foreign_keys = true;

CREATE TABLE migrations (
    key TEXT NOT NULL PRIMARY KEY,
    value INTEGER
);

CREATE TABLE resumes (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    created_at INTEGER NOT NULL DEFAULT unixepoch(),
    value TEXT
);

CREATE TABLE cover_letters (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    created_at INTEGER NOT NULL DEFAULT unixepoch(),
    value TEXT
);

CREATE TABLE distributables (
    id BYTES NOT NULL PRIMARY KEY,
    created_at INTEGER NOT NULL DEFAULT unixepoch(),
    resume_revision INTEGER REFERENCES resumes(id),
    cover_letter INTEGER REFERENCES cover_letters(id)
);

INSERT INTO migrations (key, value) VALUES ('version', 1);
