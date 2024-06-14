CREATE TABLE IF NOT EXISTS posts_users(
                            Id_users INTEGER UNSIGNED,
                            Id_posts INTEGER UNSIGNED,
                            Created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
                            Updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                            Emoji CHAR(1) NOT NULL,
                            Version INTEGER NOT NULL DEFAULT 1,
                            PRIMARY KEY(Id_users, Id_posts)
)ENGINE = INNODB;