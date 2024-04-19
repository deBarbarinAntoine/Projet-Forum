USE forum;

INSERT INTO users(username, email, password, salt, avatar_path, role, status) VALUES ('Thorgan', 'thorgdar@gmail.com', '91d90921ac94f38d23fdf5394143a76d7f613b714ba8e8220383f4e0fa3c82f80f0b3aee3fcf5f0586e172678ec91843ac8e5c5dc2ea30904782743f96341290', 'qF7W1W5Mc9PQeCyPM57NJGLuvdzpieSluBlBKDEFYLOqPenU41L698RcKMilXZDgifk2wXJYqhp2sv8fZRuQMw==', '/img/avatar/myavatar.png', 'admin', 'active');

INSERT INTO categories(Name, Id_author)
VALUES ('Coding', '1');

INSERT INTO threads(Title, Description, Is_public, Status, Id_author, Id_categories)
VALUES ('How to handle a MySQL database in web development', 'This is a thread to discuss about MySQL database\'s integration into a backend code.', true, 'active', 1, 2);

INSERT INTO posts(Content, Id_authors, Id_threads)
VALUES ('Hey guys, I\'m struggling to create and link my database with my backend in Go, have you any advice on this subject?', 1, 1);

SELECT
    u.Username AS Username,
    c.Name AS 'Category Name',
    t.Title AS Title,
    p.Content AS Post
FROM threads AS t
         INNER JOIN forum.categories c on t.Id_categories = c.Id_categories
         INNER JOIN forum.users u on t.Id_author = u.Id_users
         INNER JOIN forum.posts p on t.Id_threads = p.Id_threads
ORDER BY p.Updated_at;