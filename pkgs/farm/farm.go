package farm

import (
	"html/template"
	"net/http"
	"database/sql"
	"passionorbust.org/ross/website/pkgs/session"
	"log"
	"passionorbust.org/ross/website/pkgs/account"
)

var tmplPath = "pkgs/farm/template"

type RegisterForm struct {
	Handle string
	HandleError string
	Location string
	LocationError string
}

var errEmpty = "Must not be empty"

func (f *RegisterForm) IsHandleValid() bool {
	if len(f.Handle) < 1 {
		f.HandleError = errEmpty
		return false
	}
	return true
}

func (f *RegisterForm) IsLocationValid() bool {
	if len(f.Location) < 1 {
		f.LocationError = errEmpty
		return false
	}
	return true
}

func (f *RegisterForm) IsValid() bool {
	isValid := true

	if !f.IsHandleValid() {
		isValid = false
	}

	if !f.IsLocationValid() {
		isValid = false
	}

	return isValid
}

func (f *RegisterForm) Submit(db *sql.DB, accountID int64) (int64, error) {
	// todo: consider if i should be using session id instead of account id
	query := "INSERT INTO farm(account_id, handle, location) VALUES(?,?,?);"
	stmt, err := db.Prepare(query)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(accountID, f.Handle, f.Location)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	return id, err
}

var registerTemplate = template.Must(template.ParseFiles(
	"template/root.html",
	tmplPath+"/register.html",
))

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var f RegisterForm
	if r.Method == "POST" {
		f.Handle = r.FormValue("handle")
		f.Location = r.FormValue("location")
		if f.IsValid() {
			db, err := sql.Open("sqlite3", "pob.db")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			sessionID, err := session.Get(r)
			if err != nil {
				// todo: could be no cookie set, or invalid sessionID
				log.Fatal("do something here")
			}
			accountID, err := account.IDBySessionID(db, sessionID)
			if err != nil {
				// todo: could be error or no accountID
				log.Fatal("do something here")
			}
			_, err = f.Submit(db, accountID)
			if err == nil {
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


func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	var profileTemplate = template.Must(template.ParseFiles("template/root.html", tmplPath + "/profile.html"))
	renderTemplate(w, profileTemplate, nil)
}

// todo: repeated in multiple packages
func renderTemplate(w http.ResponseWriter, t *template.Template, data interface{}) {
	err := t.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func Queries() []string {
	return []string{
		"CREATE TABLE farm(" +
			"id INTEGER NOT NULL PRIMARY KEY," +
			"account_id INTEGER NOT NULL REFERENCES account(id)," +
			"handle TEXT NOT NULL UNIQUE," +
			"location TEXT NOT NULL" +
		");",
	}
}