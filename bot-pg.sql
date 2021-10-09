CREATE TABLE IF NOT EXISTS servers (
    -- Required
    id TEXT PRIMARY KEY NOT NULL,
    premium BOOLEAN NOT NULL,

    -- Extra
    staff_role TEXT,
    sug_channel TEXT,
    rep_channel TEXT,
    auto_approve INT,
    auto_reject INT,
    approve_emoji TEXT,
    reject_emoji TEXT,
    del_approved BOOLEAN,
    del_rejected BOOLEAN,

    -- Future update
    blacklist TEXT
);