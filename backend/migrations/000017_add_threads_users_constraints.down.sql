ALTER TABLE threads_users
    DROP FOREIGN KEY fk_threads_users_Id_users,
    DROP FOREIGN KEY fk_threads_users_Id_threads;