CREATE TABLE IF NOT EXISTS tags_users(
                           Id_users INT unsigned,
                           Id_tags INT unsigned NOT NULL,
                           Created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
                           Updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                           PRIMARY KEY(Id_users, Id_tags)
)ENGINE = INNODB;