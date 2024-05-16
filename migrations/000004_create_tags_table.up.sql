CREATE TABLE IF NOT EXISTS tags(
                     Id_tags INT unsigned auto_increment,
                     Name VARCHAR(50) NOT NULL,
                     Created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
                     Updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                     Id_author INT unsigned,
                     PRIMARY KEY(Id_tags),
                     UNIQUE(Name)
)ENGINE = INNODB;