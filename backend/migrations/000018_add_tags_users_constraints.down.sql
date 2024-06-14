ALTER TABLE tags_users
    DROP FOREIGN KEY fk_tags_users_Id_users,
    DROP FOREIGN KEY fk_tags_users_Id_tags;