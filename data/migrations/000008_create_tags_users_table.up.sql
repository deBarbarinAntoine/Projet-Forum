CREATE TABLE IF NOT EXISTS tags_users(
                           Id_users INTEGER UNSIGNED,
                           Id_tags INTEGER UNSIGNED,
                           Created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
                           Updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                           Version INTEGER NOT NULL DEFAULT 1,
                           PRIMARY KEY(Id_users, Id_tags)
)ENGINE = INNODB;