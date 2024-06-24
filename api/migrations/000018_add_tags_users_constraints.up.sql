ALTER TABLE tags_users
    ADD CONSTRAINT fk_tags_users_Id_users FOREIGN KEY(Id_users) REFERENCES users(Id_users) ON DELETE CASCADE,
    ADD CONSTRAINT fk_tags_users_Id_tags FOREIGN KEY(Id_tags) REFERENCES tags(Id_tags) ON DELETE CASCADE;