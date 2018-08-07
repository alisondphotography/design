package main

import (
	"net/http"
	"html/template"
	"log"
	"passionorbust.org/ross/website/pkgs/account"
	"passionorbust.org/ross/website/pkgs/farm"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"passionorbust.org/ross/website/pkgs/session"
)

// we're in a war, so really treat it like it, get this shit up and don't be flashy
// [x] register account
// [x] register farm (invite only? or anyone. probably anyone...) / create landing page
// [ ] post updates since fb doesn't show you all of them
// - agree to meeting my standards as a permaculture farm
// - i make clear to consumers that they need to decide for themselves as to whether or not they can trust this farm
// - vouched for by my network, their network, get educated and go and visit farm
// i help them sell
// - what do customers want? are they right?
// - what is right?
// - standard sales
// - subscription
// - pickup
// - go to grocery store... (only support grocery stores that are completely on board with my vision?)
// encourage others to farm / garden for themselves and as a business

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/account/register", account.RegisterHandler)
	http.HandleFunc("/farm/register", farm.RegisterHandler)
	http.HandleFunc("/farm/profile", farm.ProfileHandler)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

var indexTemplate = template.Must(template.ParseFiles("template/root.html", "template/index.html"))
func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	renderTemplate(w, indexTemplate, nil)
}

func renderTemplate(w http.ResponseWriter, t *template.Template, data interface{}) {
	err := t.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func init() {
	os.Remove("pob.db")
	var queries []string
	queries = append(queries, account.Queries()...)
	queries = append(queries, session.Queries()...)
	queries = append(queries, farm.Queries()...)


	db, err := sql.Open("sqlite3", "pob.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	for _, query := range queries {
		_, err = db.Exec(query)
		if err != nil {
			log.Fatalln(err, query)
		}
	}
}