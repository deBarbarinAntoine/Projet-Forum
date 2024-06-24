ALTER TABLE threads_tags
    ADD CONSTRAINT fk_threads_tags_Id_tags FOREIGN KEY(Id_tags) REFERENCES tags(Id_tags) ON DELETE CASCADE,
    ADD CONSTRAINT fk_threads_tags_Id_threads FOREIGN KEY(Id_threads) REFERENCES threads(Id_threads) ON DELETE CASCADE;