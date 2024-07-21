package models

type Task struct {
	Id      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

type TaskID struct {
	Id int `json:"id"`
}

type Tasks struct {
	Tasks []Task `json:"tasks"`
}
type ErrorResponse struct {
	Error string `json:"error"`
}
