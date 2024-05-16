ALTER TABLE posts_tags
    ADD CONSTRAINT fk_posts_tags_Id_tags FOREIGN KEY(Id_tags) REFERENCES tags(Id_tags),
    ADD CONSTRAINT fk_posts_tags_Id_posts FOREIGN KEY(Id_posts) REFERENCES posts(Id_posts);