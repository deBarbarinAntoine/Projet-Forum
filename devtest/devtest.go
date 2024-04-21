package main

import (
	"Projet-Forum/internal/db"
	"Projet-Forum/internal/models"
	"Projet-Forum/server"
	"database/sql"
	"log"
	"time"
)

func main() {
	go server.Run()

	time.Sleep(time.Second * 5)

	users, err := db.GetAllUsers()
	if err != nil {
		log.Println(err)
	}
	log.Printf("Found %d users:\n%+v\n", len(users), users)

	user, err := db.GetUserById(2)
	if err != nil {
		log.Println(err)
	}
	log.Printf("Found user:\n%+v\n", user)

	user, err = db.GetUserByLogin("thorgdar@gmail.com")
	if err != nil {
		log.Println(err)
	}
	log.Printf("Found user:\n%+v\n", user)

	newUser := new(models.User)
	newUser.Username = "toto"
	newUser.Email = "toto@gmail.com"
	newUser.HashedPwd = sql.NullString{
		String: "asldEDkvniDh@K#oesjh!G3289tigSAna9e$H?uduJFDXjfpoCa&wjWkoPOtewery=",
		Valid:  true,
	}
	newUser.Salt = sql.NullString{
		String: "sjh!G3289tigSAna9e$H?uduJFDXjfpoCa==",
		Valid:  true,
	}
	newUser.AvatarPath = sql.NullString{
		String: "/img/avatars/totoAvatar.png",
		Valid:  true,
	}
	newUser.BirthDate = sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}
	newUser.Bio = sql.NullString{
		String: "",
		Valid:  false,
	}
	newUser.Signature = sql.NullString{
		String: "",
		Valid:  false,
	}
	exists, err2 := db.IsLogin(newUser.Username)
	if exists {
		user, err = db.GetUserByLogin("toto")
		if err != nil {
			log.Println(err)
		}
		err = db.DeleteUser(user)
		if err != nil {
			log.Println(err)
		}
	}
	if err2 != nil {
		log.Println(err2)
	}
	err = db.CreateUser(*newUser)
	if err != nil {
		log.Println(err)
	}

	users, err = db.GetAllUsers()
	if err != nil {
		log.Println(err)
	}
	log.Printf("Found %d users:\n%+v\n", len(users), users)

	user, err = db.GetUserByLogin("toto")
	if err != nil {
		log.Println(err)
	}
	log.Printf("Found user:\n%+v\n", user)

	updatedFields := make(map[string]any)
	updatedFields[models.UserFields.Status] = "active"
	err = db.UpdateUser(user, updatedFields)
	if err != nil {
		log.Println(err)
	}

	users, err = db.GetAllUsers()
	if err != nil {
		log.Println(err)
	}
	log.Printf("Found %d users:\n%+v\n", len(users), users)
}
