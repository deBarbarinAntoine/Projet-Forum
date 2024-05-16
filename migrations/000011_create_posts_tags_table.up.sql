CREATE TABLE IF NOT EXISTS posts_tags(
                           Id_tags INT unsigned NOT NULL,
                           Id_posts INT unsigned NOT NULL,
                           PRIMARY KEY(Id_tags, Id_posts)
)ENGINE = INNODB;