USE forum;

INSERT INTO users (Username, Email, Password, Salt, Role, Birth_date, Bio, Signature)
VALUES ('admin', 'admin@example.com', 'hashed_password', 'random_salt', 'admin', '1990-01-01', 'Welcome to the forum!', 'Happy modding!');

INSERT INTO users (Username, Email, Password, Salt, Role, Birth_date, Bio, Signature)
VALUES ('user1', 'user1@example.com', 'hashed_password', 'random_salt', 'user', '2000-12-31', 'I love cats!', 'Meow!');

INSERT INTO users (Username, Email, Password, Salt, Role, Birth_date, Bio)
VALUES ('user2', 'user2@example.com', 'hashed_password', 'random_salt', 'user', '1985-05-15', NULL);

INSERT INTO users(username, email, password, salt, avatar_path, role, status) VALUES ('Thorgan', 'thorgdar@gmail.com', '91d90921ac94f38d23fdf5394143a76d7f613b714ba8e8220383f4e0fa3c82f80f0b3aee3fcf5f0586e172678ec91843ac8e5c5dc2ea30904782743f96341290', 'qF7W1W5Mc9PQeCyPM57NJGLuvdzpieSluBlBKDEFYLOqPenU41L698RcKMilXZDgifk2wXJYqhp2sv8fZRuQMw==', '/img/avatar/myavatar.png', 'admin', 'active');

INSERT INTO categories (Name, Id_author)
VALUES ('General', 1);

INSERT INTO categories (Name, Id_author)
VALUES ('Tech', 4);

INSERT INTO categories (Name, Id_author, Id_parent_categories)
VALUES ('Gaming', 1, 1);

INSERT INTO categories(Name, Id_author, Id_parent_categories)
VALUES ('Coding', 4, 2);


INSERT INTO tags (Name, Id_author)
VALUES ('PHP', 1);

INSERT INTO tags (Name, Id_author)
VALUES ('JavaScript', 1);

INSERT INTO tags (Name, Id_author)
VALUES ('MMORPG', 2);

INSERT INTO tags (Name, Id_author)
VALUES ('FPS', 2);


INSERT INTO threads (Title, Description, Is_public, Id_author, Id_categories)
VALUES ('Welcome message!', 'A warm welcome to all new users!', 1, 1, 1);

INSERT INTO threads (Title, Description, Is_public, Id_author, Id_categories)
VALUES ('Learning PHP', 'Need help getting started with PHP?', 1, 4, 4);

INSERT INTO threads (Title, Description, Is_public, Id_author, Id_categories)
VALUES ('Best MMORPGs 2024', 'Share your favorite MMORPGs!', 2, 3, 3);

INSERT INTO threads(Title, Description, Is_public, Status, Id_author, Id_categories)
VALUES ('How to handle a MySQL database in web development', 'This is a thread to discuss about MySQL database\'s integration into a backend code.', true, 'active', 1, 4);


INSERT INTO posts (Content, Id_authors, Id_parent_posts, Id_threads)
VALUES ('Glad to be here!', 2, NULL, 2);

INSERT INTO posts (Content, Id_authors, Id_parent_posts, Id_threads)
VALUES ('I recommend checking out official documentation!', 1, NULL, 2);

INSERT INTO posts (Content, Id_authors, Id_parent_posts, Id_threads)
VALUES ('World of Warcraft is a classic!', 3, NULL, 3);

INSERT INTO posts(Content, Id_authors, Id_threads)
VALUES ('Hey guys, I\'m struggling to create and link my database with my backend in Go, have you any advice on this subject?', 2, 4);

INSERT INTO posts(Content, Id_authors, Id_threads)
VALUES ('Just look for an online tutorial if you don\'t listen during classes.', 4, 4);