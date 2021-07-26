package handler

import (
	"database/sql"
	"io"
	"log"
	"net/http"

	"github.com/catinello/base62"
	"github.com/gorilla/mux"
	"github.com/rishabh96b/minifyurl/app/model"
)

func CreateURLEntry(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	url, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("ERROR", err.Error())
		respondError(w, http.StatusBadRequest, err.Error())
	}
	//Generate the base62 equivalent of the total rows+1 count.
	count, ok := GetAllURLCount(db)
	if !ok {
		log.Println("ERROR", "GetAllURLCount returned non OK response.")
		respondError(w, http.StatusInternalServerError, "Cannot fulfill request")
	}
	shortRef := base62.Encode(count + 1)
	_, err = db.Exec("insert into url(name, shortref) values(?,?)", string(url), shortRef)
	if err != nil {
		log.Println("ERROR", err.Error())
		respondError(w, http.StatusInternalServerError, "Cannot fulfill request")
	}
	w.WriteHeader(http.StatusAccepted)
}

// RedirectToOriginalURL will redirect to the original url from the short ref.
// TODO: modify this func to redirect to url instead of passing a response in json.
func RedirectToOriginalURL(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var urlName string
	value, ok := vars["ref"]
	if !ok {
		respondError(w, http.StatusBadRequest, "")
	}
	row, err := db.Query("select name from url where shortref=?", value)
	if err != nil {
		log.Println("ERROR", err.Error())
		respondError(w, http.StatusInternalServerError, "")
	}
	defer row.Close()
	for row.Next() {
		err := row.Scan(&urlName)
		if err != nil {
			log.Println("ERROR", err.Error())
			respondError(w, http.StatusInternalServerError, "")
		}
	}
	respondJSON(w, http.StatusOK, map[string]string{"url": urlName})
}

func GetAllURLs(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	urls := []model.URL{}
	var (
		id       int
		name     string
		shortRef string
	)

	rows, err := db.Query("select * from url")
	if err != nil {
		log.Println("Querying DB unsuccesfull.")
		respondError(w, http.StatusInternalServerError, "Could not complete the request")
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&id, &name, &shortRef)
		if err != nil {
			log.Println("ERROR", err.Error())
			respondError(w, http.StatusInternalServerError, "Could not complete the request")
		}
		urls = append(urls,
			model.URL{
				Id:       id,
				Name:     name,
				ShortRef: shortRef})
	}
	if err = rows.Err(); err != nil {
		log.Println("ERROR", err.Error())
		respondError(w, http.StatusInternalServerError, "Could not complete the request")
	}
	respondJSON(w, http.StatusOK, urls)
}

func GetURLByID(db *sql.DB, id int) (string, bool) {
	var urlName string
	row, err := db.Query("select name from url where id=?", id)
	if err != nil {
		log.Println("Querying DB unsuccesfull.")
		return "", false
	}
	defer row.Close()
	for row.Next() {
		if err := row.Scan(&urlName); err != nil {
			log.Println("ERROR", err.Error())
		}
	}
	return urlName, true
}

func GetURLByRef(db *sql.DB, ref string) (string, bool) {
	var urlName string
	row, err := db.Query("select name from url where shortref=?", ref)
	if err != nil {
		log.Println("Querying DB unsuccesfull.")
		return "", false
	}
	defer row.Close()
	for row.Next() {
		if err := row.Scan(&urlName); err != nil {
			log.Println("ERROR", err.Error())
		}
	}
	return urlName, true
}

// GetAllURLCount returns the total entries in the url table
func GetAllURLCount(db *sql.DB) (int, bool) {
	var count int
	row, err := db.Query("select count(*) from url")
	if err != nil {
		log.Println("Querying DB unsuccesfull.")
		return 0, false
	}
	defer row.Close()
	for row.Next() {
		if err := row.Scan(&count); err != nil {
			log.Println("ERROR", err.Error())
		}
	}
	return count, true
}
