ALTER TABLE categories
    ADD CONSTRAINT fk_categories_Id_author FOREIGN KEY(Id_author) REFERENCES users(Id_users) ON DELETE SET NULL,
    ADD CONSTRAINT fk_categories_Id_parent_categories FOREIGN KEY(Id_parent_categories) REFERENCES categories(Id_categories);