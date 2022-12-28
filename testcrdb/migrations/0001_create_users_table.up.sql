CREATE TABLE users
(
    id   INT PRIMARY KEY NOT NULL DEFAULT unique_rowid(),
    name VARCHAR(255)
);