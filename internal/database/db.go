package db

import (
	"database/sql"
	"fmt"
	"log"
	_ "modernc.org/sqlite"
	"os"
	"path/filepath"
	"scheduler/models"
)

const (
	tableName          string = "scheduler"
	createDbExpression string = `
		CREATE TABLE IF NOT EXISTS %[1]s
		(id INTEGER PRIMARY KEY AUTOINCREMENT,
		date CHAR(8) NOT NULL DEFAULT "",
		title VARCHAR(128) NOT NULL DEFAULT "",
		comment TEXT NOT NULL DEFAULT "",
		repeat VARCHAR(128) NOT NULL DEFAULT "");`
	createIndexExpression = `CREATE INDEX IF NOT EXISTS %[1]s_idx_date ON %[1]s(date);`
	insertTaskExpression  = `INSERT INTO %s (date, title, comment, repeat) VALUES ($1, $2, $3, $4);`
)

func CreateDbObject(fileName string, expressions ...string) error {
	db, err := sql.Open("sqlite", fileName)
	if err != nil {
		return err
	}
	defer db.Close()
	for _, expression := range expressions {
		//fmt.Println(fmt.Sprintf(expression, tableName))
		_, err = db.Exec(fmt.Sprintf(expression, tableName))
		if err != nil {
			return err
		}
	}
	return nil
}

func InitDb(fileName string) error {
	appPath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	dbFile := filepath.Join(filepath.Dir(appPath), fileName)
	_, err = os.Stat(dbFile)
	if err != nil {
		log.Printf("Creating database at %s", fileName)
		err := CreateDbObject(dbFile, createDbExpression, createIndexExpression)
		if err != nil {
			return err
		}
	}
	return nil
}

func OpenDB(fileName string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", fileName)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func InsertTask(db *sql.DB, task models.Task) (int, error) {
	fmt.Println(task)
	res, err := db.Exec(`INSERT INTO scheduler (date, title, comment, repeat) VALUES (?,?,?,?)`, task.Date, task.Title, task.Comment, task.Repeat)
	//res, err := db.Exec(fmt.Sprintf(insertTaskExpression, tableName), task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}
