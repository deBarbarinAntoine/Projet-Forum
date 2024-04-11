package utils

import (
	"Projet-Forum/internal/models"
	"encoding/json"
	"errors"
	"log/slog"
	"os"
	"reflect"
	"regexp"
	"time"
)

// jsonFile is the models.User's JSON file full path.
var jsonFile = Path + "data/users.json"

// TempUsers is the models.TempUser's array for newly registered models.User
// before they confirm their email address.
var TempUsers []models.TempUser

// LostUsers is the models.TempUser's array for models.User
// that forgot their password, until they access the link sent to their
// email address to change their password.
var LostUsers []models.TempUser

// retrieveUsers
// retrieves all models.User present in jsonFile and stores them in a slice of models.User.
// It returns the slice of models.User and an error.
func retrieveUsers() ([]models.User, error) {
	var users []models.User

	data, err := os.ReadFile(jsonFile)

	if len(data) == 0 {
		return nil, nil
	}

	err = json.Unmarshal(data, &users)
	if err != nil {
		return nil, err
	}

	return users, nil
}

// CheckUser
// checks if the models.User 's username and email are still available in jsonFile and TempUsers.
func CheckUser(user models.User) bool {
	users, err := retrieveUsers()
	if err != nil {
		Logger.Error(GetCurrentFuncName(), slog.Any("output", err))
	}
	for _, singleUser := range users {
		if user.Username == singleUser.Username || user.Email == singleUser.Email {
			return false
		}
	}
	for _, tempUser := range TempUsers {
		if user.Username == tempUser.User.Username || user.Email == tempUser.User.Email {
			return false
		}
	}
	return true
}

// CheckEmail checks the mail's format.
func CheckEmail(email string) bool {
	reg := regexp.MustCompile("^[\\w&#$.%+-]+@[\\w&#$.%+-]+\\.[a-z]{2,6}?$")
	return reg.MatchString(email)
}

// EmailExists return whether the mail address exists in the user's list.
func EmailExists(email string) (bool, models.User) {
	users, err := retrieveUsers()
	if err != nil {
		Logger.Error(GetCurrentFuncName(), slog.Any("output", err))
	}
	for _, singleUser := range users {
		if email == singleUser.Email {
			return true, singleUser
		}
	}
	return false, models.User{}
}

// CheckPasswd
// checks if the password's format is according to the rules.
func CheckPasswd(passwd string) bool {

	// Matches any password containing at least one digit, one lowercase,
	// one uppercase, one symbol and 8 characters in total.
	//regex := regexp.MustCompile(`^(?=.*\d)(?=.*[a-z])(?=.*[A-Z])(?=.*([^\w\s]|_)).{8,}$`) // Alas not supported by the regexp library
	digit := regexp.MustCompile(`\d+`)
	lower := regexp.MustCompile(`[a-z]+`)
	upper := regexp.MustCompile(`[A-Z]+`)
	symbol := regexp.MustCompile(`([^\w\s]|_)+`)
	minLen := regexp.MustCompile(`.{8,}`)
	return digit.MatchString(passwd) && lower.MatchString(passwd) && upper.MatchString(passwd) && symbol.MatchString(passwd) && minLen.MatchString(passwd)
}

// changeUsers
// overwrites jsonFile with `users` in json format.
func changeUsers(users []models.User) {
	data, errJSON := json.MarshalIndent(users, "", "\t")
	if errJSON != nil {
		Logger.Error(GetCurrentFuncName()+" JSON MarshalIndent error!", slog.Any("output", errJSON))
		return
	}
	errWrite := os.WriteFile(jsonFile, data, 0666)
	if errWrite != nil {
		Logger.Error(GetCurrentFuncName()+" WriteFile error!", slog.Any("output", errWrite))
	}
}

// GetIdNewUser
// returns first unused id in jsonFile.
func GetIdNewUser() int {
	users, err := retrieveUsers()
	if err != nil {
		Logger.Error(GetCurrentFuncName(), slog.Any("output", err))
	}
	var id int
	var idFound bool
	for id = 1; !idFound; id++ {
		idFound = true
		for _, user := range users {
			if user.Id == id {
				idFound = false
			}
		}
	}
	id--
	return id
}

// CreateUser
// adds the models.User `newUser` to jsonFile.
func CreateUser(newUser models.User) {
	users, err := retrieveUsers()
	if err != nil {
		Logger.Error(GetCurrentFuncName(), slog.Any("output", err))
	}
	users = append(users, newUser)
	changeUsers(users)
}

// removeUser
// remove the models.User which models.User.Id is sent in argument from jsonFile.
func removeUser(id int) {
	users, err := retrieveUsers()
	if err != nil {
		Logger.Error(GetCurrentFuncName(), slog.Any("output", err))
	}
	for i, user := range users {
		if user.Id == id {
			users = append(users[:i], users[i+1:]...)
		}
	}
	changeUsers(users)
}

// SelectUser
// returns the models.User which models.User.Username matches the `username` argument.
func SelectUser(username string) (models.User, bool) {
	var user models.User
	users, err := retrieveUsers()
	if err != nil {
		Logger.Error(GetCurrentFuncName(), slog.Any("output", err))
	}
	var ok bool
	for _, singleUser := range users {
		if singleUser.Username == username {
			ok = true
			user = singleUser
		}
	}
	return user, ok
}

// UpdateUser
// modifies the models.User in jsonFile that matches
// `updatedUser`'s Id with `updatedUser`'s content.
func UpdateUser(updatedUser models.User) {
	users, err := retrieveUsers()
	if err != nil {
		Logger.Error(GetCurrentFuncName(), slog.Any("output", err))
	}
	for i, user := range users {
		if user.Id == updatedUser.Id {
			users[i] = updatedUser
		}
	}
	changeUsers(users)
}

// deleteTempUser
// removes a specific models.TempUser from TempUsers.
func deleteTempUser(temp models.TempUser) {
	for i, user := range TempUsers {
		if reflect.DeepEqual(user, temp) {
			TempUsers = append(TempUsers[:i], TempUsers[i+1:]...)
		}
	}
}

// deleteLostUser
// removes a specific models.TempUser from LostUsers.
func deleteLostUser(temp models.TempUser) {
	for i, user := range LostUsers {
		if reflect.DeepEqual(user, temp) {
			LostUsers = append(LostUsers[:i], LostUsers[i+1:]...)
		}
	}
}

// PushTempUser
// creates a new user from a models.TempUser
// which Id matches the `id` param.
func PushTempUser(id string) {
	for _, temp := range TempUsers {
		if temp.ConfirmID == id {
			temp.User.Id = GetIdNewUser()
			temp.User.CreationTime = time.Now()
			temp.User.Avatar = "profile-avatar-059.jpg"
			CreateUser(temp.User)
			deleteTempUser(temp)
		}
	}
}

// UpdateLostUser
// updates the models.User's Hash and Salt with the `lost`
// models.TempUser's Hash and Salt sent in the param.
func UpdateLostUser(lost models.TempUser) {
	user, ok := SelectUser(lost.User.Username)
	if !ok {
		Logger.Error(GetCurrentFuncName(), slog.Any("output", errors.New("user not found")))
	}
	user.Salt = lost.User.Salt
	user.HashedPwd = lost.User.HashedPwd
	UpdateUser(user)
	deleteLostUser(lost)
}

// ManageTempUsers
// is a goroutine that periodically removes old models.TempUser from TempUsers and LostUsers.
func ManageTempUsers() {
	time.Sleep(time.Second * 10)
	duration := SetDailyTimer(2)
	for {
		Logger.Info(GetCurrentFuncName(), slog.String("goroutine", "ManageTempUsers"))
		for _, user := range TempUsers {
			if time.Since(user.CreationTime) > time.Hour*12 {
				Logger.Info("TempUser cleared automatically", slog.Any("user", user))
				deleteTempUser(user)
			}
		}
		for _, user := range LostUsers {
			if time.Since(user.CreationTime) > time.Hour*12 {
				Logger.Info("LostUser cleared automatically", slog.Any("user", user))
				deleteTempUser(user)
			}
		}
		time.Sleep(duration)
		duration = time.Hour * 12
	}
}
