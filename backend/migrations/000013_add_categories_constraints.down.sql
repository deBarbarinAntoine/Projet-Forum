ALTER TABLE categories
    DROP FOREIGN KEY fk_categories_Id_author,
    DROP FOREIGN KEY fk_categories_Id_parent_categories;