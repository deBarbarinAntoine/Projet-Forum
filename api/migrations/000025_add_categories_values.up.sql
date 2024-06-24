INSERT INTO categories (Id_categories, Name, Id_author, Id_parent_categories)
VALUES (1, 'Web Dev', 3, NULL),
       (2, 'Game Dev', 4, NULL),
       (3, 'Hardware', 5, NULL),
       (4, 'Server', 3, NULL),
       (5, 'RESTful APIs', 3, 1),
       (6, 'Frontend tricks', 4, 1),
       (7, 'Best smartphone buys', 5, 3),
       (8, 'Razer for ever', 4, 3),
       (9, 'About languages', 5, 1);