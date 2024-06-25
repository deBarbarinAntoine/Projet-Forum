ALTER TABLE posts_users
    ADD CONSTRAINT fk_posts_users_Id_users FOREIGN KEY(Id_users) REFERENCES users(Id_users) ON DELETE SET NULL,
    ADD CONSTRAINT fk_posts_users_Id_posts FOREIGN KEY(Id_posts) REFERENCES posts(Id_posts) ON DELETE CASCADE;