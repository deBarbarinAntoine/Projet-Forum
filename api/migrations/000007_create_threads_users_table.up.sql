CREATE TABLE IF NOT EXISTS threads_users(
                              Id_users INTEGER UNSIGNED,
                              Id_threads INTEGER UNSIGNED,
                              Created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
                              Updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                              Version INTEGER NOT NULL DEFAULT 1,
                              PRIMARY KEY(Id_users, Id_threads)
)ENGINE = INNODB;