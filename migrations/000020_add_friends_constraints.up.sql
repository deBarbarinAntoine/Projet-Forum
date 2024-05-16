ALTER TABLE friends
    ADD CONSTRAINT fk_friends_Id_users_1 FOREIGN KEY(Id_users_1) REFERENCES users(Id_users),
    ADD CONSTRAINT fk_friends_Id_users_2 FOREIGN KEY(Id_users_2) REFERENCES users(Id_users);