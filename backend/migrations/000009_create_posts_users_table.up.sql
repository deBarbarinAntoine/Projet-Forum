CREATE TABLE IF NOT EXISTS posts_users(
                            Id_users INT unsigned,
                            Id_posts INT unsigned NOT NULL,
                            Created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
                            Updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                            Emoji CHAR(1) unicode NOT NULL,
                            PRIMARY KEY(Id_users, Id_posts)
)ENGINE = INNODB;