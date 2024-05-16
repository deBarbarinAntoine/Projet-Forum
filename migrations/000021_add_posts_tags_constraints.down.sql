ALTER TABLE posts_tags
    DROP FOREIGN KEY fk_posts_tags_Id_tags,
    DROP FOREIGN KEY fk_posts_tags_Id_posts;