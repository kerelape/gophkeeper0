
CREATE TABLE IF NOT EXISTS identities(
    username TEXT PRIMARY KEY UNIQUE,
    password TEXT
);

CREATE TABLE IF NOT EXISTS resources(
    id SERIAL PRIMARY UNIQUE,
    resource INTEGER,
    type INTEGER,
    owner TEXT,
    meta TEXT,
);

CREATE TABLE IF NOT EXISTS pieces(
    id SERIAL PRIMARY KEY UNIQUE,
    content BYTEA,
    salt BYTEA,
    iv BYTEA,
);

CREATE TABLE IF NOT EXISTS blobs(
    id SERIAL PRIMARY KEY UNIQUE,
    path TEXT,
    salt BYTEA,
    iv BYTEA,
);
