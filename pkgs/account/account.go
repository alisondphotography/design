package account

import (
	"html/template"
	"net/http"
	"database/sql"
		"strconv"
	"passionorbust.org/ross/website/pkgs/session"
)

var tmplPath = "pkgs/account/template"

var registerTemplate = template.Must(template.ParseFiles(
	"template/root.html",
	tmplPath+"/register.html",
))

type RegisterForm struct {
	Username string
	UsernameError string
	Password string
	PasswordError string
}

func (f *RegisterForm) IsValid() bool {
	isValid := true

	if !f.UsernameIsValid() {
		isValid = false
	}

	if !f.PasswordIsValid() {
		isValid = false
	}

	return isValid
}

func (f *RegisterForm) UsernameIsValid() bool {
	if len(f.Username) < 1 {
		f.UsernameError = errEmpty
		return false
	}
	return true
}

func (f *RegisterForm) PasswordIsValid() bool {
	if len(f.Password) < 1 {
		f.PasswordError = errEmpty
		return false
	}
	return true
}

func (f *RegisterForm) Submit(db *sql.DB) (int64, error) {
	query := "INSERT INTO account(username, password) VALUES(?,?);"
	stmt, err := db.Prepare(query)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(f.Username, f.Password)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, err
}

// todo: repeating in other packages
var errEmpty = "Must not be empty"

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var f RegisterForm
	if r.Method == "POST" {
		f.Username = r.FormValue("username")
		f.Password = r.FormValue("password")

		if f.IsValid() {
			db, err := sql.Open("sqlite3", "pob.db")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer db.Close()
			id, err := f.Submit(db)
			if err == nil {
				// todo: notice type switch to int
				_, err = session.Set(db, w, strconv.Itoa(int(id)))
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				w.Write([]byte("success"))
				return
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}
	renderTemplate(w, registerTemplate, f)
}

func renderTemplate(w http.ResponseWriter, t *template.Template, data interface{}) {
	err := t.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func Queries() []string{
	return []string{
		"CREATE TABLE account(" +
			"id INTEGER NOT NULL PRIMARY KEY," +
			"username TEXT NOT NULL UNIQUE," +
			"password TEXT NOT NULL" +
		");",
	}
}

func IDBySessionID(db *sql.DB, sessionID string) (int64, error) {
	query := "SELECT account_id FROM session WHERE id = ?;"
	stmt, err := db.Prepare(query)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(sessionID)
	var id int64
	err = row.Scan(&id)
	return id, err
}