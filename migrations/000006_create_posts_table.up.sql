CREATE TABLE IF NOT EXISTS posts(
                      Id_posts INT unsigned auto_increment,
                      Content VARCHAR(1020) unicode NOT NULL,
                      Created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
                      Updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                      Id_authors INT unsigned,
                      Id_parent_posts INT unsigned,
                      Id_threads INT unsigned NOT NULL,
                      PRIMARY KEY(Id_posts)
)ENGINE = INNODB;