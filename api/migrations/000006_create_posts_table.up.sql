CREATE TABLE IF NOT EXISTS posts(
                      Id_posts INTEGER UNSIGNED PRIMARY KEY AUTO_INCREMENT,
                      Content VARCHAR(1020) NOT NULL,
                      Created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
                      Updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                      Id_author INTEGER UNSIGNED,
                      Id_parent_posts INTEGER UNSIGNED,
                      Id_threads INTEGER UNSIGNED,
                      Version INTEGER NOT NULL DEFAULT 1
)ENGINE = INNODB;