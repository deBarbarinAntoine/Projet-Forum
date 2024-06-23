CREATE TABLE IF NOT EXISTS categories(
                           Id_categories INTEGER UNSIGNED PRIMARY KEY AUTO_INCREMENT,
                           Name VARCHAR(70) UNIQUE NOT NULL,
                           Created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
                           Updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                           Id_author INTEGER UNSIGNED,
                           Id_parent_categories INTEGER UNSIGNED,
                           Version INTEGER NOT NULL DEFAULT 1
)ENGINE = INNODB;