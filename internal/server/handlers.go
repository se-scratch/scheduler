package server

import (
	"encoding/json"
	"log"
	"net/http"
	"scheduler/internal/helper"
	"scheduler/models"
	"time"
)

func NextDateHandler(w http.ResponseWriter, req *http.Request) {
	params := req.URL.Query()
	nowParam := params.Get("now")
	dayParam := params.Get("date")
	repeatParam := params.Get("repeat")

	now, err := time.Parse(helper.DateLayout, nowParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	nextDate, err := helper.NextDate(now, dayParam, repeatParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_, err = w.Write([]byte(nextDate))
	if err != nil {
		log.Printf("error writing response: %v", err)
	}
}

func TaskHandler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		task, err := helper.DecodePostTask(req.Body)
		if err != nil {
			log.Println("decode post task error:", err)
			SendResponse(w, http.StatusBadRequest, err)
			return
		}
		id, err := db.InsertTask(task)
		if err != nil {
			SendResponse(w, http.StatusInternalServerError, err)
		}
		SendResponse(w, http.StatusOK, models.TaskID{Id: id})
		return
	case http.MethodGet:
		id := req.URL.Query().Get("id")
		task, err := db.SelectTask(id)
		if err != nil {
			SendResponse(w, http.StatusBadRequest, err)
			return
		}
		SendResponse(w, http.StatusOK, task)
		return

	case http.MethodPut:
		task, err := helper.DecodePostTask(req.Body)
		if err != nil {
			log.Println(err)
			SendResponse(w, http.StatusBadRequest, err)
			return
		}

		_, err = db.SelectTask(task.Id)
		if err != nil {
			SendResponse(w, http.StatusBadRequest, err)
			return
		}

		err = db.UpdateTask(task)
		if err != nil {
			SendResponse(w, http.StatusInternalServerError, err)
			return
		}
		SendResponse(w, http.StatusOK, map[string]string{})
		return

	case http.MethodDelete:
		id := req.URL.Query().Get("id")
		_, err := db.SelectTask(id)
		if err != nil {
			SendResponse(w, http.StatusBadRequest, err)
			return
		}
		err = db.DeleteTask(id)
		if err != nil {
			SendResponse(w, http.StatusInternalServerError, err)
			return
		}
		SendResponse(w, http.StatusOK, map[string]string{})
		return
	}
}

func TasksHandler(w http.ResponseWriter, req *http.Request) {
	var tasks []models.Task
	tasks, err := db.SelectTasks()
	if err != nil {
		log.Println(err)
		SendResponse(w, http.StatusInternalServerError, err)
		return
	}
	SendResponse(w, http.StatusOK, models.Tasks{Tasks: tasks})
	return
}

func TaskDoneHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		id := req.URL.Query().Get("id")
		task, err := db.SelectTask(id)
		if err != nil {
			SendResponse(w, http.StatusBadRequest, err)
			return
		}
		if task.Repeat == "" {
			err = db.DeleteTask(id)
			if err != nil {
				SendResponse(w, http.StatusInternalServerError, err)
				return
			}
			SendResponse(w, http.StatusOK, map[string]string{})
			return
		}
		now := time.Now()
		task.Date, err = helper.NextDate(now, task.Date, task.Repeat)
		if err != nil {
			SendResponse(w, http.StatusInternalServerError, err)
			return
		}
		err = db.UpdateTask(task)
		if err != nil {
			SendResponse(w, http.StatusInternalServerError, err)
			return
		}
		SendResponse(w, http.StatusOK, map[string]string{})
		return
	}
}

func SendResponse(w http.ResponseWriter, statusCode int, response interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if _, ok := response.(error); ok {
		if err := json.NewEncoder(w).Encode(models.ErrorResponse{Error: response.(error).Error()}); err != nil {
			log.Println("SendResponse() error:", err)
		}
		return
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println("SendResponse() error", err)
	}
	return
}
