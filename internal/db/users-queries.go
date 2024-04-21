package db

// get all Users from Users table (no params, selecting *).
var getAllUsersQuery = `
SELECT
    * 
FROM 
    users 
ORDER BY 
    Visited_at;`

// get User matching a specific id (param: id).
var getUserByIdQuery = `
SELECT
    * 
FROM 
    users 
WHERE 
    Id_users = ?;`

// get count of Users matching a specific set of Credentials (params: username or email, password, salt).
var getUserByLoginQuery = `
SELECT 
    * 
FROM 
    users 
WHERE 
    Username = ? OR Email = ?;`

// update user visited_at field matching an id (param: id).
var updateUserVisitedAtQuery = `
UPDATE 
    users 
SET 
    Visited_at = CURRENT_TIMESTAMP 
WHERE 
    Id_users = ?
LIMIT 1;`

// update user matching an id (params: fields and values in SQL format [%s], id).
var updateUserQuery = `
UPDATE 
    users 
SET 
    %s
WHERE 
    Id_users = ?
LIMIT 1;`

// create user (params: username, email, password, salt, avatarPath, BirthDate, Bio, Signature).
var createUserQuery = `
INSERT INTO
	users (
	       Username, 
	       Email, 
	       Password, 
	       Salt, 
	       Avatar_path, 
	       Birth_date, 
	       Bio, 
	       Signature
	       )
VALUES 
    (?, ?, ?, ?, ?, ?, ?, ?);`

// delete user matching an id (param: id).
var deleteUserQuery = `
DELETE FROM
           users
WHERE
    Id_users = ?
LIMIT 1;`

// check if a username or email exists (param: username or email).
var isLoginQuery = `
SELECT
	COUNT(*)
FROM
    users
WHERE
    Username = ? OR Email = ?;`
