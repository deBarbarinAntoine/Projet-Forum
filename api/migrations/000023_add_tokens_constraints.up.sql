ALTER TABLE tokens
    ADD CONSTRAINT fk_tokens_Id_users FOREIGN KEY(Id_users) REFERENCES users(Id_users) ON DELETE CASCADE;