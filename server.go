package main

import (
	"fmt"
	"log"
	"time"
	"strings"
	"net/http"
	"html/template"
    "database/sql"
    "github.com/lib/pq"
)

type Moment struct {
	ID string
	Body  string
	Tags []string
	Timestamp string
}

var db *sql.DB

// this is called before main
func init() {
    connStr := "user=postgres dbname=timeline password=postgres sslmode=disable"
	var err error
    db, err = sql.Open("postgres", connStr)
    if err != nil {
        fmt.Println("Error connecting to the database:", err)
        return
    }
	
}

// TODO: read from db
// TODO: display with template
func viewHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, blog_post, tags, date FROM moments ORDER BY date DESC")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()


	var items []Moment;
	// Iterate through the result set.
	for rows.Next() {
		var id string
		var blog_post string
		var tags []string
		var timestamp time.Time
		err := rows.Scan(&id, &blog_post, pq.Array(&tags), &timestamp)
		if err != nil {
			log.Fatal(err)
		}
		formattedDate := timestamp.Format("2006/Jan/02 15:04")
		// formattedDate := timestamp.Format("01/Oct 04:40pm")
		items = append(items, Moment{id, blog_post, tags, formattedDate})
	}

	// render template
	t, _ := template.ParseFiles("front.html")
	t.Execute(w, items)

	// Check for errors from iterating over rows.
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}

func newHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("new.html")
	t.Execute(w, nil)
}

func removeWhitespaceFromSlice(slice []string) []string {
    trimmedSlice := make([]string, len(slice))
    for i, str := range slice {
        trimmedSlice[i] = strings.TrimSpace(str)
    }
    return trimmedSlice
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
		body := r.FormValue("body")
		tags := strings.Split(r.FormValue("tags"), ",")
		trimmedTags := removeWhitespaceFromSlice(tags)

		// TODO: save file to db
		// save logic can be moved to that
		// m := &Moment{ID: id.String(), Body: body}
		// err := m.save()
		// if err != nil {
		// 	http.Error(w, err.Error(), http.StatusInternalServerError)
		// 	return
		// }

		// save to db
		query := `INSERT INTO moments (blog_post, tags) VALUES ($1, $2)`
		_, err := db.Query(query, body, pq.Array(trimmedTags))
		if err != nil {
			log.Fatal(err)
		}
		http.Redirect(w, r, "/", http.StatusFound)
}

func viewTagHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	tag := strings.Split(path, "/")[2]

	rows, err := db.Query("SELECT id, blog_post, tags, date FROM moments WHERE $1 = ANY(tags) ORDER BY date DESC", tag)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()


	var items []Moment;
	// Iterate through the result set.
	for rows.Next() {
		var id string
		var blog_post string
		var tags []string
		var timestamp time.Time
		err := rows.Scan(&id, &blog_post, pq.Array(&tags), &timestamp)
		if err != nil {
			log.Fatal(err)
		}
		formattedDate := timestamp.Format("01/Oct 3:4pm")
		items = append(items, Moment{id, blog_post, tags, formattedDate})
	}

	// render template
	t, _ := template.ParseFiles("front.html")
	t.Execute(w, items)

	// Check for errors from iterating over rows.
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	http.HandleFunc("/", viewHandler)
	http.HandleFunc("/new", newHandler)
	http.HandleFunc("/save", saveHandler)
	http.HandleFunc("/tag/", viewTagHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
