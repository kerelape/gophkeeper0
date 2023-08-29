
CREATE TABLE IF NOT EXISTS identities(
    username TEXT PRIMARY KEY UNIQUE,
    password TEXT
);
