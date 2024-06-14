CREATE TABLE IF NOT EXISTS users(
                      Id_users INTEGER UNSIGNED PRIMARY KEY AUTO_INCREMENT,
                      Username VARCHAR(30) UNIQUE NOT NULL,
                      Email VARCHAR(35) UNIQUE NOT NULL,
                      Hashed_password CHAR(60) NOT NULL,
                      Avatar_path VARCHAR(255) NOT NULL DEFAULT '',
                      Role VARCHAR(20) NOT NULL DEFAULT 'normal',
                      Birth_date DATE,
                      Created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
                      Updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
                      Visited_at DATETIME DEFAULT CURRENT_TIMESTAMP,
                      Bio VARCHAR(255) NOT NULL DEFAULT '',
                      Signature VARCHAR(255) NOT NULL DEFAULT '',
                      Status VARCHAR(20) NOT NULL DEFAULT 'to-confirm',
                      Version INTEGER NOT NULL DEFAULT 1
)ENGINE = INNODB;