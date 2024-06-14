CREATE TABLE IF NOT EXISTS categories(
                           Id_categories INT unsigned auto_increment,
                           Name VARCHAR(50) NOT NULL,
                           Created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
                           Updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                           Id_author INT unsigned,
                           Id_parent_categories INT unsigned,
                           PRIMARY KEY(Id_categories),
                           UNIQUE(Name)
)ENGINE = INNODB;