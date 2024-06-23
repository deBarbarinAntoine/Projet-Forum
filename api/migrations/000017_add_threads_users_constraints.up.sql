ALTER TABLE threads_users
    ADD CONSTRAINT fk_threads_users_Id_users FOREIGN KEY(Id_users) REFERENCES users(Id_users) ON DELETE CASCADE,
    ADD CONSTRAINT fk_threads_users_Id_threads FOREIGN KEY(Id_threads) REFERENCES threads(Id_threads) ON DELETE CASCADE;