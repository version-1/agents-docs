CREATE TABLE repositories (
    id INTEGER PRIMARY KEY,
    identity TEXT NOT NULL UNIQUE
);
CREATE TABLE questions (
    id INTEGER PRIMARY KEY,
    repository_id INTEGER NOT NULL REFERENCES repositories(id),
    content TEXT NOT NULL,
    fingerprint TEXT NOT NULL,
    replaces INTEGER REFERENCES questions(id),
    status TEXT NOT NULL DEFAULT 'active' CHECK(status IN ('active','deleted','superseded')),
    created_at TEXT NOT NULL,
    UNIQUE(repository_id, fingerprint)
);
CREATE TABLE progress (
    question_id INTEGER PRIMARY KEY REFERENCES questions(id),
    presentations INTEGER NOT NULL DEFAULT 0 CHECK(presentations >= 0),
    streak INTEGER NOT NULL DEFAULT 0 CHECK(streak BETWEEN 0 AND 3),
    due_at TEXT,
    mastered INTEGER NOT NULL DEFAULT 0 CHECK(mastered IN (0,1)),
    generation INTEGER NOT NULL DEFAULT 0
);
CREATE TABLE attempts (
    id TEXT PRIMARY KEY,
    question_id INTEGER NOT NULL REFERENCES questions(id),
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
CREATE INDEX questions_repository ON questions(repository_id, status);
CREATE INDEX attempts_question ON attempts(question_id, presented_at);
PRAGMA user_version = 1;
