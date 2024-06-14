CREATE TABLE IF NOT EXISTS threads(
                        Id_threads INT unsigned auto_increment,
                        Title VARCHAR(62) unicode NOT NULL,
                        Description VARCHAR(255) unicode,
                        Is_public boolean NOT NULL,
                        Created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
                        Updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                        Status VARCHAR(20) NOT NULL DEFAULT 'active',
                        Id_author INT unsigned,
                        Id_categories INT unsigned NOT NULL,
                        PRIMARY KEY(Id_threads),
                        UNIQUE(Title)
)ENGINE = INNODB;