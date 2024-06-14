CREATE TABLE IF NOT EXISTS friends(
                        Id_users_1 INT unsigned NOT NULL,
                        Id_users_2 INT unsigned NOT NULL,
                        Created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
                        Updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                        Status VARCHAR(20) NOT NULL,
                        PRIMARY KEY(Id_users_1, Id_users_2)
)ENGINE = INNODB;