package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"scheduler/models"
	"strconv"

	_ "modernc.org/sqlite"
)

const (
	tableName string = "scheduler"
	limit     int    = 50
)

type TaskDB struct {
	DB *sql.DB
}

func OpenDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("error opening database %q: %v", path, err)
	}
	return db, nil
}

func (d *TaskDB) Close() {
	d.DB.Close()
}

func (d *TaskDB) CreateDbObject(expressions ...string) error {
	for _, expression := range expressions {
		_, err := d.DB.Exec(fmt.Sprintf(expression, tableName))
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *TaskDB) Init(fileName string) error {
	const (
		createDbExpression string = `
		CREATE TABLE IF NOT EXISTS %[1]s
		(id INTEGER PRIMARY KEY AUTOINCREMENT,
		date VARCHAR(8) NOT NULL DEFAULT "",
		title VARCHAR(256) NOT NULL DEFAULT "",
		comment TEXT NOT NULL DEFAULT "",
		repeat VARCHAR(128) NOT NULL DEFAULT "");`
		createIndexExpression = `CREATE INDEX IF NOT EXISTS %[1]s_idx_date ON %[1]s(date);`
	)
	appPath, err := os.Getwd()
	if err != nil {
		return err
	}

	dbPath := filepath.Join(appPath, fileName)
	_, err = os.Stat(dbPath)
	if err != nil {
		log.Printf("Creating database at %s", dbPath)
		d.DB, err = OpenDB(dbPath)
		if err != nil {
			return err
		}
		err := d.CreateDbObject(createDbExpression, createIndexExpression)
		if err != nil {
			return err
		}
	} else {
		log.Printf("Opening database at %s", dbPath)
		d.DB, err = OpenDB(dbPath)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *TaskDB) InsertTask(task models.Task) (int, error) {
	const insertTaskExpression = `INSERT INTO %s (date, title, comment, repeat) VALUES ($1, $2, $3, $4);`

	res, err := d.DB.Exec(fmt.Sprintf(insertTaskExpression, tableName), task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (d *TaskDB) SelectTasks() ([]models.Task, error) {
	const selectTasksExpression = `SELECT id,date,title,comment,repeat FROM %s ORDER BY date LIMIT %d;`

	tasks := make([]models.Task, 0)
	res, err := d.DB.Query(fmt.Sprintf(selectTasksExpression, tableName, limit))
	if err != nil {
		log.Printf("error query selecting tasks: %v\n", err)
		return nil, err
	}
	defer res.Close()
	for res.Next() {
		var task models.Task
		err := res.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err == sql.ErrNoRows {
			return tasks, nil
		} else if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	if res.Err() != nil {
		log.Printf("error scan rows tasks: %v\n", res.Err())
		return nil, res.Err()
	}
	return tasks, nil
}

func (d *TaskDB) SelectTask(id string) (models.Task, error) {
	const selectTaskExpression = `SELECT id,date,title,comment,repeat FROM %s WHERE id = $1;`
	var task models.Task

	if _, err := strconv.Atoi(id); err != nil {
		return task, fmt.Errorf("wrong task id ")
	}
	res := d.DB.QueryRow(fmt.Sprintf(selectTaskExpression, tableName), id)
	err := res.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return task, err
	}
	return task, nil
}

func (d *TaskDB) UpdateTask(task models.Task) error {
	const updateTaskExpression = `UPDATE %s SET date = $1, title = $2, comment = $3, repeat = $4 WHERE id = $5`

	res, err := d.DB.Exec(fmt.Sprintf(updateTaskExpression, tableName), task.Date, task.Title, task.Comment, task.Repeat, task.Id)
	if err != nil {
		return err
	}
	result, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if result == 1 {
		return nil
	}
	return fmt.Errorf("update failed, expected 1 row affected, got %d", result)
}

func (d *TaskDB) DeleteTask(id string) error {
	const deleteTaskExpression = `DELETE FROM %s WHERE id = $1;`

	res, err := d.DB.Exec(fmt.Sprintf(deleteTaskExpression, tableName), id)
	if err != nil {
		return err
	}
	result, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if result == 1 {
		return nil
	}
	return fmt.Errorf("delete failed, expected 1 row affected, got %d", result)
}
