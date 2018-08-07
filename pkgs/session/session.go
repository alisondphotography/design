package session

import (
	"math/rand"
	"database/sql"
	"net/http"
	"strconv"
)

func Set(db *sql.DB, w http.ResponseWriter, accountID string) (int, error) {
	id := rand.Intn(10000000000000)

	query := "INSERT INTO session(id, account_id) VALUES(?,?);"
	stmt, err := db.Prepare(query)
	if err != nil {
		return id, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id, accountID)

	http.SetCookie(w, &http.Cookie{Name: "sessionID", Value: strconv.Itoa(id), Path: "/"})

	return id, err
}

func Get(r *http.Request) (string, error) {
	cookie, err := r.Cookie("sessionID")
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

func Queries() []string {
	// todo: can id's be unique forever?
	// todo: is session a sqlite keyword?
	return []string{
		"CREATE TABLE session(" +
			"id INTEGER NOT NULL UNIQUE," +
			"account_id INTEGER NOT NULL REFERENCES account(id)" +
		");",
	}
}
