CREATE DATABASE if not exists forum;

USE forum;

DROP TABLE IF EXISTS users;
CREATE TABLE users(
                      Id_users INT unsigned not null auto_increment unique,
                      Username VARCHAR(25) not null,
                      Email VARCHAR(35) not null,
                      Password CHAR(128),
                      Salt CHAR(88),
                      Avatar_path VARCHAR(125),
                      Role VARCHAR(20),
                      Birth_date DATETIME,
                      Created_at DATETIME,
                      Updated_at DATETIME,
                      Visited_at DATETIME,
                      Bio VARCHAR(255) unicode,
                      Signature VARCHAR(255) unicode,
                      Status VARCHAR(20),
                      PRIMARY KEY(Id_users),
                      UNIQUE(Username),
                      UNIQUE(Email)
);

DROP TABLE IF EXISTS categories;
CREATE TABLE categories(
                           Id_categories INT unsigned not null auto_increment unique,
                           Name VARCHAR(50) NOT NULL,
                           Created_at DATETIME,
                           Updated_at DATETIME,
                           Id_author INT unsigned NOT NULL,
                           Id_parent_categories INT unsigned NOT NULL,
                           PRIMARY KEY(Id_categories),
                           UNIQUE(Name),
                           FOREIGN KEY(Id_author) REFERENCES users(Id_users),
                           FOREIGN KEY(Id_parent_categories) REFERENCES categories(Id_categories)
);

DROP TABLE IF EXISTS tags;
CREATE TABLE tags(
                     Id_tags INT unsigned not null auto_increment unique,
                     Name VARCHAR(50) unique NOT NULL,
                     Created_at DATETIME,
                     Updated_at DATETIME,
                     Id_author INT unsigned NOT NULL,
                     PRIMARY KEY(Id_tags),
                     UNIQUE(Name),
                     FOREIGN KEY(Id_author) REFERENCES users(Id_users)
);

DROP TABLE IF EXISTS threads;
CREATE TABLE threads(
                        Id_threads INT unsigned not null auto_increment unique,
                        Title VARCHAR(62) unicode NOT NULL,
                        Description VARCHAR(255) unicode,
                        Is_public boolean NOT NULL,
                        Created_at DATETIME,
                        Updated_at DATETIME,
                        Status VARCHAR(20) NOT NULL,
                        Id_author INT unsigned NOT NULL,
                        Id_categories INT unsigned NOT NULL,
                        PRIMARY KEY(Id_threads),
                        UNIQUE(Title),
                        FOREIGN KEY(Id_author) REFERENCES users(Id_users),
                        FOREIGN KEY(Id_categories) REFERENCES categories(Id_categories)
);

DROP TABLE IF EXISTS posts;
CREATE TABLE posts(
                      Id_posts INT unsigned not null auto_increment unique,
                      Content VARCHAR(1020) unicode NOT NULL,
                      Created_at DATETIME,
                      Updated_at DATETIME,
                      Id_authors INT unsigned NOT NULL,
                      Id_parent_posts INT unsigned NOT NULL,
                      Id_threads INT unsigned NOT NULL,
                      PRIMARY KEY(Id_parent_posts),
                      FOREIGN KEY(Id_authors) REFERENCES users(Id_users),
                      FOREIGN KEY(Id_parent_posts) REFERENCES posts(Id_posts),
                      FOREIGN KEY(Id_threads) REFERENCES threads(Id_threads)
);

DROP TABLE IF EXISTS follow;
CREATE TABLE follow(
                       Id_users INT unsigned NOT NULL,
                       Id_threads INT unsigned NOT NULL,
                       Created_at DATETIME,
                       Updated_at DATETIME,
                       PRIMARY KEY(Id_users, Id_threads),
                       FOREIGN KEY(Id_users) REFERENCES users(Id_users),
                       FOREIGN KEY(Id_threads) REFERENCES threads(Id_threads)
);

DROP TABLE IF EXISTS favorite;
CREATE TABLE favorite(
                         Id_users INT unsigned NOT NULL,
                         Id_tags INT unsigned NOT NULL,
                         Created_at DATETIME,
                         Updated_at DATETIME,
                         PRIMARY KEY(Id_users, Id_tags),
                         FOREIGN KEY(Id_users) REFERENCES users(Id_users),
                         FOREIGN KEY(Id_tags) REFERENCES tags(Id_tags)
);

DROP TABLE IF EXISTS react;
CREATE TABLE react(
                      Id_users INT unsigned NOT NULL,
                      Id_posts INT unsigned NOT NULL,
                      Created_at DATETIME,
                      Updated_at DATETIME,
                      Emoji CHAR(1) unicode NOT NULL,
                      PRIMARY KEY(Id_users, Id_posts),
                      FOREIGN KEY(Id_users) REFERENCES users(Id_users),
                      FOREIGN KEY(Id_posts) REFERENCES posts(Id_posts)
);

DROP TABLE IF EXISTS befriend;
CREATE TABLE befriend(
                         Id_users_1 INT unsigned NOT NULL,
                         Id_users_2 INT unsigned NOT NULL,
                         Created_at DATETIME,
                         Updated_at DATETIME,
                         Status VARCHAR(20) NOT NULL,
                         PRIMARY KEY(Id_users_1, Id_users_2),
                         FOREIGN KEY(Id_users_1) REFERENCES users(Id_users),
                         FOREIGN KEY(Id_users_2) REFERENCES users(Id_users)
);

DROP TABLE IF EXISTS have;
CREATE TABLE have(
                     Id_tags INT unsigned NOT NULL,
                     Id_posts INT unsigned NOT NULL,
                     PRIMARY KEY(Id_tags, Id_posts),
                     FOREIGN KEY(Id_tags) REFERENCES tags(Id_tags),
                     FOREIGN KEY(Id_posts) REFERENCES posts(Id_posts)
);
