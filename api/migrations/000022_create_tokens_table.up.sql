CREATE TABLE IF NOT EXISTS tokens (
    Hash CHAR(60) UNIQUE NOT NULL,
    Id_users INTEGER UNSIGNED,
    Expiry TIMESTAMP NOT NULL,
    Scope VARCHAR(80) NOT NULL
)