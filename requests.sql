USE forum;

-- Query to fetch all threads with the author's username, the parent category name and the current category name, the thread's title and the posts' contents.
SELECT
    pc.Name AS 'Parent Category',
    c.Name AS 'Category Name',
    u.Username AS 'Thread Author',
    t.Title AS 'Thread Title',
    pu.Username AS 'Post Author',
    p.Content AS Post
FROM threads AS t
         INNER JOIN forum.categories c on t.Id_categories = c.Id_categories
         INNER JOIN forum.users u on t.Id_author = u.Id_users
         INNER JOIN forum.posts p on t.Id_threads = p.Id_threads
         INNER JOIN forum.users pu on pu.Id_users = p.Id_authors
         LEFT OUTER JOIN forum.categories pc on c.Id_parent_categories = pc.Id_categories
ORDER BY p.Updated_at;