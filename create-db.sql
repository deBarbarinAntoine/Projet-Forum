DROP DATABASE IF EXISTS forum;

CREATE DATABASE IF NOT EXISTS forum;

USE forum;

DROP TABLE IF EXISTS users;
CREATE TABLE users(
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

DROP TABLE IF EXISTS categories;
CREATE TABLE categories(
                           Id_categories INT unsigned auto_increment,
                           Name VARCHAR(50) NOT NULL,
                           Created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
                           Updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                           Id_author INT unsigned,
                           Id_parent_categories INT unsigned,
                           PRIMARY KEY(Id_categories),
                           UNIQUE(Name)
)ENGINE = INNODB;

DROP TABLE IF EXISTS tags;
CREATE TABLE tags(
                     Id_tags INT unsigned auto_increment,
                     Name VARCHAR(50) NOT NULL,
                     Created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
                     Updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                     Id_author INT unsigned,
                     PRIMARY KEY(Id_tags),
                     UNIQUE(Name)
)ENGINE = INNODB;

DROP TABLE IF EXISTS threads;
CREATE TABLE threads(
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

DROP TABLE IF EXISTS posts;
CREATE TABLE posts(
                      Id_posts INT unsigned auto_increment,
                      Content VARCHAR(1020) unicode NOT NULL,
                      Created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
                      Updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                      Id_authors INT unsigned,
                      Id_parent_posts INT unsigned,
                      Id_threads INT unsigned NOT NULL,
                      PRIMARY KEY(Id_posts)
)ENGINE = INNODB;

DROP TABLE IF EXISTS threads_users;
CREATE TABLE threads_users(
                       Id_users INT unsigned,
                       Id_threads INT unsigned NOT NULL,
                       Created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
                       Updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                       PRIMARY KEY(Id_users, Id_threads)
)ENGINE = INNODB;

DROP TABLE IF EXISTS tags_users;
CREATE TABLE tags_users(
                         Id_users INT unsigned,
                         Id_tags INT unsigned NOT NULL,
                         Created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
                         Updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                         PRIMARY KEY(Id_users, Id_tags)
)ENGINE = INNODB;

DROP TABLE IF EXISTS posts_users;
CREATE TABLE posts_users(
                      Id_users INT unsigned,
                      Id_posts INT unsigned NOT NULL,
                      Created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
                      Updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                      Emoji CHAR(1) unicode NOT NULL,
                      PRIMARY KEY(Id_users, Id_posts)
)ENGINE = INNODB;

DROP TABLE IF EXISTS friends;
CREATE TABLE friends(
                         Id_users_1 INT unsigned NOT NULL,
                         Id_users_2 INT unsigned NOT NULL,
                         Created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
                         Updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                         Status VARCHAR(20) NOT NULL,
                         PRIMARY KEY(Id_users_1, Id_users_2)
)ENGINE = INNODB;

DROP TABLE IF EXISTS posts_tags;
CREATE TABLE posts_tags(
                     Id_tags INT unsigned NOT NULL,
                     Id_posts INT unsigned NOT NULL,
                     PRIMARY KEY(Id_tags, Id_posts)
)ENGINE = INNODB;

ALTER TABLE categories
    ADD CONSTRAINT FOREIGN KEY(Id_author) REFERENCES users(Id_users) ON DELETE SET NULL,
    ADD CONSTRAINT FOREIGN KEY(Id_parent_categories) REFERENCES categories(Id_categories);

ALTER TABLE tags
    ADD CONSTRAINT FOREIGN KEY(Id_author) REFERENCES users(Id_users) ON DELETE SET NULL;

ALTER TABLE threads
    ADD CONSTRAINT FOREIGN KEY(Id_author) REFERENCES users(Id_users) ON DELETE SET NULL,
    ADD CONSTRAINT FOREIGN KEY(Id_categories) REFERENCES categories(Id_categories);

ALTER TABLE posts
    ADD CONSTRAINT FOREIGN KEY(Id_authors) REFERENCES users(Id_users) ON DELETE SET NULL,
    ADD CONSTRAINT FOREIGN KEY(Id_parent_posts) REFERENCES posts(Id_posts),
    ADD CONSTRAINT FOREIGN KEY(Id_threads) REFERENCES threads(Id_threads);

ALTER TABLE threads_users
    ADD CONSTRAINT FOREIGN KEY(Id_users) REFERENCES users(Id_users) ON DELETE CASCADE,
    ADD CONSTRAINT FOREIGN KEY(Id_threads) REFERENCES threads(Id_threads);

ALTER TABLE tags_users
    ADD CONSTRAINT FOREIGN KEY(Id_users) REFERENCES users(Id_users) ON DELETE CASCADE,
    ADD CONSTRAINT FOREIGN KEY(Id_tags) REFERENCES tags(Id_tags);

ALTER TABLE posts_users
    ADD CONSTRAINT FOREIGN KEY(Id_users) REFERENCES users(Id_users) ON DELETE CASCADE,
    ADD CONSTRAINT FOREIGN KEY(Id_posts) REFERENCES posts(Id_posts);

ALTER TABLE friends
    ADD CONSTRAINT FOREIGN KEY(Id_users_1) REFERENCES users(Id_users),
    ADD CONSTRAINT FOREIGN KEY(Id_users_2) REFERENCES users(Id_users);

ALTER TABLE posts_tags
    ADD CONSTRAINT FOREIGN KEY(Id_tags) REFERENCES tags(Id_tags),
    ADD CONSTRAINT FOREIGN KEY(Id_posts) REFERENCES posts(Id_posts);