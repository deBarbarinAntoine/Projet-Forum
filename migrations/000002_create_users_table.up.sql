CREATE TABLE IF NOT EXISTS users(
                      Id_users INT unsigned not null auto_increment unique,
                      Username VARCHAR(25) not null,
                      Email VARCHAR(35) not null,
                      Password CHAR(128),
                      Salt CHAR(88),
                      Avatar_path VARCHAR(125),
                      Role VARCHAR(20) NOT NULL DEFAULT 'normal',
                      Birth_date DATE,
                      Created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
                      Updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
                      Visited_at DATETIME,
                      Bio VARCHAR(255) unicode,
                      Signature VARCHAR(255) unicode,
                      Status VARCHAR(20) NOT NULL DEFAULT 'to-confirm',
                      PRIMARY KEY(Id_users),
                      UNIQUE(Username),
                      UNIQUE(Email)
)ENGINE = INNODB;