ALTER TABLE tags
    ADD CONSTRAINT fk_tags_Id_author FOREIGN KEY(Id_author) REFERENCES users(Id_users) ON DELETE SET NULL;