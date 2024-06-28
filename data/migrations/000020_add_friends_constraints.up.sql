ALTER TABLE friends
    ADD CONSTRAINT fk_friends_Id_users_from FOREIGN KEY(Id_users_from) REFERENCES users(Id_users) ON DELETE CASCADE,
    ADD CONSTRAINT fk_friends_Id_users_to FOREIGN KEY(Id_users_to) REFERENCES users(Id_users) ON DELETE CASCADE;