package utils

import (
	"Projet-Forum/internal/models"
	"crypto/rand"
	"encoding/base64"
	"log/slog"
	"net/http"
	"time"
)

// SessionsData is an in-memory models.Session data storage.
var SessionsData = make(map[string]models.Session)

// sessionTime is the constant handling the session's maximum opened time without interaction.
const sessionTime time.Duration = time.Hour * 2

// retrieveSessions
//
//	@Description: fetches all sessions present in SessionsData.
//	@return []models.Session
func retrieveSessions() []models.Session {
	var sessions []models.Session
	for _, session := range SessionsData {
		sessions = append(sessions, session)
	}
	return sessions
}

// GetSession
//
//	@Description: fetches the models.Session and sessionId from the cookie present
//	in the *http.Request.
//	@param r
//	@return models.Session
//	@return string
func GetSession(r *http.Request) (models.Session, string) {
	sessionID, err := r.Cookie("updatedCookie")
	if err != nil {
		Logger.Error(GetCurrentFuncName(), slog.Any("output", err))
		return models.Session{}, ""
	}
	return SessionsData[sessionID.Value], sessionID.Value
}

// newConnectionID
//
//	@Description: gets the first unused ConnectionID for the new models.Session.
//	@return int
func newConnectionID() int {
	sessions := retrieveSessions()
	var id int
	var idFound bool
	for id = 1; !idFound; id++ {
		idFound = true
		for _, session := range sessions {
			if session.ConnectionID == id {
				idFound = false
			}
		}
	}
	id--
	return id
}

// OpenSession
//
//	@Description: creates a new models.Session for the user which username matches
//	the username param and writes the cookie that corresponds to its
//	models.Session in the *http.ResponseWriter and *http.Request (for further
//	access).
//	@param w
//	@param username
//	@param r
func OpenSession(w *http.ResponseWriter, username string, r *http.Request) {

	// Generate and set Session ID cookie
	sessionID := generateSessionID()
	// Generate expiration time for the cookie
	expirationTime := time.Now().Add(time.Hour * 2)

	newCookie := &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		Expires:  expirationTime,
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}

	http.SetCookie(*w, newCookie)
	r.AddCookie(newCookie)

	user, _ := SelectUser(username)

	// Update the last connection time in `users.json`.
	user.LastConnection = time.Now()
	UpdateUser(user)

	// Create Session data in memory
	SessionsData[sessionID] = models.Session{
		UserID:         user.Id,
		ConnectionID:   newConnectionID(),
		Username:       username,
		IpAddress:      GetIP(r),
		ExpirationTime: expirationTime,
	}

	Logger.Info("Login", slog.Any("user", SessionsData[sessionID]))
}

// CheckSession checks if there is a cookie in the request
// and if yes, it checks if the corresponding models.Session
// is still valid and returns true if all verifications are ok.
func CheckSession(r *http.Request) bool {
	// Extract session ID from cookie
	cookie, err := r.Cookie("session_id")
	if err != nil || !validateSessionID(cookie.Value) {
		return false
	}
	// Retrieve user data from session
	session, ok := SessionsData[cookie.Value]
	if !ok {
		return false
	}
	// Verify user IP address
	if session.IpAddress != GetIP(r) {
		return false
	}
	// Verify expiration time
	if session.ExpirationTime.Before(time.Now()) {
		Logger.Info("Logout", slog.Any("user", SessionsData[cookie.Value]))
		// deleting previous entry in the SessionsData map
		delete(SessionsData, cookie.Value)
		return false
	}
	return true
}

// RefreshSession
//
//	@Description: refreshes the cookie and the sessionId found in the
//	*http.Request and adds an "updatedCookie" in the *http.Request for further
//	access.
//	@param w
//	@param r
//	@return error
func RefreshSession(w *http.ResponseWriter, r *http.Request) error {
	// generating new sessionID and new expiration time
	newSessionID := generateSessionID()
	newExpirationTime := time.Now().Add(sessionTime)

	var newCookie = &http.Cookie{
		Name:     "session_id",
		Value:    newSessionID,
		HttpOnly: true,
		Secure:   false, // TODO: change when switching to HTTPS in the future.
		Path:     "/",
		Expires:  newExpirationTime,
		SameSite: http.SameSiteStrictMode,
	}

	// setting the new cookie
	http.SetCookie(*w, newCookie)

	// retrieving the in-memory current session data
	cookie, err := r.Cookie("session_id")
	currentSessionData := SessionsData[cookie.Value]

	// updating the sessionID and expirationTime
	currentSessionData.ExpirationTime = newExpirationTime

	// deleting previous entry in the SessionsData map
	delete(SessionsData, cookie.Value)

	// setting the new entry in the SessionsData map
	SessionsData[newSessionID] = currentSessionData

	// adding the new cookie to the request to access it from the targeted handler with the Name "updatedCookie"
	newCookie.Name = "updatedCookie"
	r.AddCookie(newCookie)

	if err != nil {
		return err
	}
	return nil
}

// Logout
//
//	@Description: sets the cookie as expired and clears the SessionsData.
//	@param w
//	@param r
func Logout(w *http.ResponseWriter, r *http.Request) {
	var newCookie = &http.Cookie{
		Name:     "session_id",
		Value:    "",
		HttpOnly: true,
		Secure:   false, // TODO: change when switching to HTTPS in the future.
		Path:     "/",
		MaxAge:   -1,
		SameSite: http.SameSiteStrictMode,
	}

	// setting the new cookie
	http.SetCookie(*w, newCookie)

	// retrieving the in-memory current session data
	cookie, _ := r.Cookie("updatedCookie")

	Logger.Info("Logout", slog.Any("user", SessionsData[cookie.Value]))

	// deleting previous entry in the SessionsData map
	delete(SessionsData, cookie.Value)
}

// generateSessionID
//
//	@Description: generates a new random sessionId.
//	@return string
func generateSessionID() string {
	b := make([]byte, 64)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

// validateSessionID
//
//	@Description: checks if the sessionID has the required length.
//	@param sessionID
//	@return bool
func validateSessionID(sessionID string) bool {
	_, ok := SessionsData[sessionID]
	return len(sessionID) == 88 && ok
}

// isExpired
//
//	@Description: checks if the models.Session is expired.
//	@param session
//	@return bool
func isExpired(session models.Session) bool {
	return session.ExpirationTime.Before(time.Now())
}

// cleanSessions
//
//	@Description: clears all expired models.Session from SessionsData.
func cleanSessions() {
	for sessionID, session := range SessionsData {
		if isExpired(session) {
			Logger.Info("Session cleared automatically", slog.Any("user", session))
			delete(SessionsData, sessionID)
		}
	}
}

// MonitorSessions
//
//	@Description: a simple function that clears all expired models.Session
//	periodically.
//	It is meant to be run as a goroutine.
func MonitorSessions() {
	for {
		time.Sleep(time.Hour)
		Logger.Info(GetCurrentFuncName(), slog.String("goroutine", "MonitorSessions"))
		cleanSessions()
	}
}
