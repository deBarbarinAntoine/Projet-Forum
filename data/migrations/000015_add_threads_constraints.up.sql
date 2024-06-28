ALTER TABLE threads
    ADD CONSTRAINT fk_threads_Id_author FOREIGN KEY(Id_author) REFERENCES users(Id_users) ON DELETE SET NULL,
    ADD CONSTRAINT fk_threads_Id_categories FOREIGN KEY(Id_categories) REFERENCES categories(Id_categories);