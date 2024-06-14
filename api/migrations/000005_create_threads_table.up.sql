CREATE TABLE IF NOT EXISTS threads(
                        Id_threads INTEGER UNSIGNED PRIMARY KEY AUTO_INCREMENT,
                        Title VARCHAR(125) UNIQUE NOT NULL,
                        Description VARCHAR(1020) NOT NULL DEFAULT '',
                        Is_public BOOLEAN NOT NULL DEFAULT true,
                        Created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
                        Updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                        Status VARCHAR(20) NOT NULL DEFAULT 'active',
                        Id_author INTEGER UNSIGNED,
                        Id_categories INTEGER UNSIGNED,
                        Version INTEGER NOT NULL DEFAULT 1
)ENGINE = INNODB;