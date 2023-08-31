
CREATE TABLE IF NOT EXISTS identities(
    username TEXT PRIMARY KEY UNIQUE,
    password TEXT
);

CREATE TABLE IF NOT EXISTS pieces(
    rid SERIAL PRIMARY KEY UNIQUE,
    owner TEXT,
    meta TEXT,
    content BYTEA,
    salt BYTEA,
    iv BYTEA,
);
