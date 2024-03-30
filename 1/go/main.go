package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var (
	DB   *sql.DB
	stmt *sql.Stmt
)

type tb1 struct {
	ID   int       `json:"id"`
	Text string    `json:"col_texto"`
	Dt   time.Time `json:"col_dt"`
}

func (t tb1) save() (err error) {
	if stmt == nil {
		stmt, err = DB.Prepare("INSERT INTO tb01 (id, col_texto, col_dt) VALUES ($1, $2, $3)")
	}
	if err != nil {
		return fmt.Errorf("unable to prepare statement: %v", err)
	}
	result, err := stmt.Exec(t.ID, t.Text, t.Dt)
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

// func endpointTb01(w http.ResponseWriter, r *http.Request) {

// }

func main() {
	db, err := sql.Open("pgx", "postgres://postgres:tb01@localhost:5432/tb01")
	if err != nil {
		panic(fmt.Errorf("unable to connect to database: %v", err))

	}
	defer db.Close()
	DB = db
	err = tb1{
		ID:   3,
		Text: "test2",
		Dt:   time.Now(),
	}.save()
	if err != nil {
		panic(err)
	}
	// http.HandleFunc("/tb01", endpointTb01)
	// log.Fatal(http.ListenAndServe(":8080", nil))
}
