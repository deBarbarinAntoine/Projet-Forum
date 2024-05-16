CREATE DATABASE IF NOT EXISTS serverTemplate;

USE serverTemplate;

CREATE TABLE IF NOT EXISTS sessions (
    token CHAR(43) PRIMARY KEY,
    data BLOB NOT NULL,
    expiry TIMESTAMP(6) NOT NULL
);

CREATE INDEX sessions_expiry_idx ON sessions (expiry);

CREATE USER 'sT_Manager'@'localhost'
    IDENTIFIED BY '$er!/3rT3mpI4t3';
GRANT SELECT, INSERT, UPDATE, DELETE
    ON serverTemplate.sessions
    TO 'sT_Manager'@'localhost';