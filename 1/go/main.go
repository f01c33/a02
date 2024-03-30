package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var (
	DB   *sql.DB
	stmt *sql.Stmt
)

type tb1 struct {
	ID   int       `json:"id"`
	Text string    `json:"texto"`
	Dt   time.Time `json:"dt"`
}

func (t tb1) save() (err error) {
	if stmt == nil {
		if stmt, err = DB.Prepare("INSERT INTO tb02 (col_texto, col_dt) VALUES ($1, $2)"); err != nil {
			return fmt.Errorf("unable to prepare statement: %v", err)
		}
	}
	result, err := stmt.Exec(t.Text, t.Dt)
	if err != nil {
		return fmt.Errorf("unable to run statement: %v", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("unable to get affected rows: %v", err)
	}
	if affected == 0 {
		return fmt.Errorf("row didn't get inserted")
	}
	return nil
}

func endpointTb01(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("wrong method, use POST"))
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("unable to read body %v", err)))
		return
	}
	var t tb1
	json.Unmarshal(body, &t)
	if t.Text == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid text"))
		return
	}
	t.Dt = time.Now()
	err = t.save()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("unable to save tb01: %v", err)))
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func main() {
	db, err := sql.Open("pgx", "postgres://postgres:tb01@localhost:5432/tb01")
	if err != nil {
		panic(fmt.Errorf("unable to connect to database: %v", err))

	}
	defer db.Close()
	DB = db
	http.HandleFunc("/tb01", endpointTb01)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
