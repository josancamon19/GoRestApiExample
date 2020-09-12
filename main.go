package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type task struct {
	ID      int    `json:ID`
	Name    string `json:Name`
	Content string `json:Content`
}

type allTasks []task

var tasks = allTasks{
	task{
		ID:      1,
		Name:    "Task 1",
		Content: "Content",
	},
}

func indexRoute(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to rest api 3")
}

func getTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func getTaskBydId(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	println(vars["id"])
	taskId, err := strconv.Atoi(vars["id"])
	if err != nil {
		fmt.Fprintf(w, "Send a valid task id")
	}
	for _, task := range tasks {
		if task.ID == taskId {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(task)
		}
	}
	w.WriteHeader(http.StatusNotFound)
}

func deleteTaskById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskId, err := strconv.Atoi(vars["id"])
	if err != nil {
		fmt.Fprintf(w, "Send a valid task id")
	}
	for i, task := range tasks {
		if task.ID == taskId {
			tasks = append(tasks[:i], tasks[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			fmt.Fprintf(w, "Task with id %v has been deleted", taskId)
		}
	}
	w.WriteHeader(http.StatusNotFound)
}

func createTask(w http.ResponseWriter, r *http.Request) {
	var newTask task
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Insert a valid task")
	}
	json.Unmarshal(reqBody, &newTask)
	newTask.ID = len(tasks) + 1
	tasks = append(tasks, newTask)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTask)
}
func updateTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskId, err := strconv.Atoi(vars["id"])
	if err != nil {
		fmt.Fprintf(w, "Send a valid task id")
	}

	var updatedTask task
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Send valid task data")
	}

	json.Unmarshal(reqBody, &updatedTask)
	for i, task := range tasks {
		if task.ID == taskId {
			updatedTask.ID = taskId
			tasks = append(tasks[:i], append(tasks[i+1:], updatedTask)...)
			fmt.Fprintf(w, "Task with id %v has been updated", taskId)
		}
	}
	w.WriteHeader(http.StatusNotFound)
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", indexRoute)
	router.HandleFunc("/tasks", getTasks).Methods("GET")
	router.HandleFunc("/tasks", createTask).Methods("POST")
	router.HandleFunc("/tasks/{id}", getTaskBydId).Methods("GET")
	router.HandleFunc("/tasks/{id}", deleteTaskById).Methods("DELETE")
	router.HandleFunc("/tasks/{id}", updateTask).Methods("PUT4")
	log.Fatal(http.ListenAndServe(":8000", router))
}
