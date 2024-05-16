CREATE TABLE IF NOT EXISTS threads_users(
                              Id_users INT unsigned,
                              Id_threads INT unsigned NOT NULL,
                              Created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
                              Updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                              PRIMARY KEY(Id_users, Id_threads)
)ENGINE = INNODB;