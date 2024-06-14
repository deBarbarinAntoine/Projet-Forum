ALTER TABLE posts
    DROP FOREIGN KEY fk_posts_Id_authors,
    DROP FOREIGN KEY fk_posts_Id_parent_posts,
    DROP FOREIGN KEY fk_posts_Id_threads;