ALTER TABLE posts_users
    DROP FOREIGN KEY fk_posts_users_Id_users,
    DROP FOREIGN KEY fk_posts_users_Id_posts;