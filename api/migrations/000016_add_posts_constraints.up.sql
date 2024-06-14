ALTER TABLE posts
    ADD CONSTRAINT fk_posts_Id_author FOREIGN KEY(Id_author) REFERENCES users(Id_users) ON DELETE SET NULL,
    ADD CONSTRAINT fk_posts_Id_parent_posts FOREIGN KEY(Id_parent_posts) REFERENCES posts(Id_posts),
    ADD CONSTRAINT fk_posts_Id_threads FOREIGN KEY(Id_threads) REFERENCES threads(Id_threads);