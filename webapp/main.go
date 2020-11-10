package main

import (
	"encoding/json"
	"strconv"
)
import "fmt"
import "github.com/gorilla/mux"
import _ "github.com/gorilla/mux"
import "io/ioutil"
import "log"
import "net/http"
import _ "io/ioutil"

type task struct {
	ID      int    `json:ID`
	Name    string `json:Name`
	Content string `json:Content`
}

type allTasks []task

var tasks = allTasks{
	{
		ID:      1,
		Name:    "Task One",
		Content: "Some content",
	},
}

func getTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func getTaskById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskId, err := strconv.Atoi(vars["idTask"])
	if err != nil {
		fmt.Fprintf(w, "Invalid ID")
	}

	for _, task := range tasks {
		if task.ID == taskId {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(task)
		}
	}
}
func deleteTaskById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskId, err := strconv.Atoi(vars["idTask"])
	if err != nil {
		fmt.Fprintf(w, "Invalid ID")
	}
	for i, task := range tasks {
		if task.ID == taskId {
			tasks = append(tasks[:i], tasks[i+1:]...)
			fmt.Fprintf(w, "The task with ID %v has been remove succesfully", taskId)
		}
	}

}

func updateTaskById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskId, err := strconv.Atoi(vars["idTask"])
	var updatedTask task
	if err != nil {
		fmt.Fprintf(w, "Invalid ID")
	}
	rqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Please Enter Valid Data")
	}
	json.Unmarshal(rqBody, &updatedTask)

	for i, task := range tasks {
		if task.ID == taskId {
			tasks = append(tasks[:i], tasks[i+1:]...)
			updatedTask.ID = taskId
			tasks = append(tasks, updatedTask)

			fmt.Fprintf(w, "The task with ID %v has been updated successfully", taskId)
		}
	}
}

func createTask(w http.ResponseWriter, r *http.Request) {
	var newTask task
	reqBody, err := ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Fprintf(w, "Insert a Valid Task")
	}
	json.Unmarshal(reqBody, &newTask)

	newTask.ID = len(tasks) + 1
	tasks = append(tasks, newTask)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTask)

}

func indexRoute(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to my API")

}

func main() {

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", indexRoute).Methods("GET")
	router.HandleFunc("/task", getTasks).Methods("GET")
	router.HandleFunc("/task", createTask).Methods("POST")
	router.HandleFunc("/task/{idTask}", getTaskById).Methods("GET")
	router.HandleFunc("/task/{idTask}", deleteTaskById).Methods("DELETE")
	router.HandleFunc("/task/{idTask}", updateTaskById).Methods("PUT")

	log.Fatal(http.ListenAndServe(":3000", router))

}
