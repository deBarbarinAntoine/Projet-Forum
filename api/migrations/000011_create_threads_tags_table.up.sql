CREATE TABLE IF NOT EXISTS threads_tags(
                           Id_threads INTEGER UNSIGNED,
                           Id_tags INTEGER UNSIGNED,
                           PRIMARY KEY(Id_threads, Id_tags)
)ENGINE = INNODB;