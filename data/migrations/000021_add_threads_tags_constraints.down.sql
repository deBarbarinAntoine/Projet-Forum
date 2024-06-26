ALTER TABLE threads_tags
    DROP FOREIGN KEY fk_threads_tags_Id_tags,
    DROP FOREIGN KEY fk_threads_tags_Id_threads;