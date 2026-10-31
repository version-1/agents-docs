-- Rebuild all referencing tables so historical revisions may share a fingerprint.
CREATE TABLE questions_v2 (
    id INTEGER PRIMARY KEY,
    repository_id INTEGER NOT NULL REFERENCES repositories(id),
    content TEXT NOT NULL,
    fingerprint TEXT NOT NULL,
    replaces INTEGER REFERENCES questions_v2(id),
    status TEXT NOT NULL DEFAULT 'active' CHECK(status IN ('active','deleted','superseded')),
    created_at TEXT NOT NULL
);
INSERT INTO questions_v2 SELECT * FROM questions ORDER BY id;
CREATE TABLE progress_v2 (
    question_id INTEGER PRIMARY KEY REFERENCES questions_v2(id),
    presentations INTEGER NOT NULL DEFAULT 0 CHECK(presentations >= 0),
    streak INTEGER NOT NULL DEFAULT 0 CHECK(streak BETWEEN 0 AND 3),
    due_at TEXT,
    mastered INTEGER NOT NULL DEFAULT 0 CHECK(mastered IN (0,1)),
    generation INTEGER NOT NULL DEFAULT 0
);
INSERT INTO progress_v2 SELECT * FROM progress;
CREATE TABLE attempts_v2 (
    id TEXT PRIMARY KEY,
    question_id INTEGER NOT NULL REFERENCES questions_v2(id),
    presented_at TEXT NOT NULL,
    ordinal INTEGER NOT NULL,
    generation INTEGER NOT NULL,
    snapshot TEXT NOT NULL,
    choices TEXT NOT NULL,
    answered_at TEXT,
    selected TEXT,
    correct INTEGER CHECK(correct IN (0,1)),
    feedback TEXT,
    CHECK ((answered_at IS NULL AND selected IS NULL AND correct IS NULL AND feedback IS NULL)
        OR (answered_at IS NOT NULL AND selected IS NOT NULL AND correct IS NOT NULL AND feedback IS NOT NULL))
);
INSERT INTO attempts_v2 SELECT * FROM attempts;
DROP TABLE attempts;
DROP TABLE progress;
DROP TABLE questions;
ALTER TABLE questions_v2 RENAME TO questions;
ALTER TABLE progress_v2 RENAME TO progress;
ALTER TABLE attempts_v2 RENAME TO attempts;
CREATE INDEX questions_repository ON questions(repository_id, status);
CREATE INDEX questions_fingerprint ON questions(repository_id, fingerprint, status);
CREATE INDEX attempts_question ON attempts(question_id, presented_at);
PRAGMA user_version = 2;
