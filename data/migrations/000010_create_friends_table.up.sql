CREATE TABLE IF NOT EXISTS friends(
                        Id_users_from INTEGER UNSIGNED,
                        Id_users_to INTEGER UNSIGNED,
                        Created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
                        Updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                        Status VARCHAR(20) NOT NULL DEFAULT 'pending',
                        PRIMARY KEY(Id_users_from, Id_users_to)
)ENGINE = INNODB;