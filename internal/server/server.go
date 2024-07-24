package server

import (
	"log"
	"net/http"
	"scheduler/internal/database"
	"strconv"
)

var db database.TaskDB

func RunServer(port int, dbfile string) error {
	db = database.TaskDB{}
	err := db.Init(dbfile)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	log.Println("Running server on port", port)
	http.Handle("/", http.FileServer(http.Dir("./web")))
	http.HandleFunc("/api/nextdate", NextDateHandler)
	http.HandleFunc("/api/task", TaskHandler)
	http.HandleFunc("/api/tasks", TasksHandler)
	http.HandleFunc("/api/task/done", TaskDoneHandler)
	err = http.ListenAndServe(":"+strconv.Itoa(port), nil)
	if err != nil {
		return err
	}
	return nil
}
