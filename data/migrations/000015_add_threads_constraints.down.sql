ALTER TABLE threads
    DROP FOREIGN KEY fk_threads_Id_author,
    DROP FOREIGN KEY fk_threads_Id_categories;