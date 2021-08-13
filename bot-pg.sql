CREATE TABLE IF NOT EXISTS servers (
    id TEXT PRIMARY KEY NOT NULL,
    staff_role TEXT,
    suggestion_channel TEXT,
    report_channel TEXT,
    auto_approve INT,
    auto_reject INT,
    approve_emoji TEXT,
    reject_emoji TEXT,
    delete_approved BOOLEAN,
    delete_rejected BOOLEAN,
    submit_blacklist TEXT,
    premium BOOLEAN NOT NULL
)

CREATE TABLE IF NOT EXISTS suggestions (
    id TEXT PRIMARY KEY NOT NULL, -- s_abc45
    context TEXT NOT NULL,
    author TEXT NOT NULL,
    avatar TEXT NOT NULL,
    channel TEXT NOT NULL,
    message TEXT NOT NULL,
    status INT NOT NULL,
    timestamp TIMESTAMP NOT NULL
)

CREATE TABLE IF NOT EXISTS suggestions (
    id TEXT PRIMARY KEY NOT NULL, -- r_abc45
    context TEXT NOT NULL,
    author TEXT NOT NULL,
    avatar TEXT NOT NULL,
    channel TEXT NOT NULL,
    message TEXT NOT NULL,
    status INT NOT NULL,
    timestamp TIMESTAMP NOT NULL
)

CREATE TABLE IF NOT EXISTS polls (
    id TEXT PRIMARY KEY NOT NULL, -- p_abc45
    context TEXT NOT NULL,
    channel TEXT NOT NULL,
    message TEXT NOT NULL,
    timestamp TIMESTAMP NOT NULL
)