package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("secret_key") // Ключ для подписи JWT

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Task struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Done   bool   `json:"done"`
	UserID int    `json:"user_id"`
}

var users = map[int]User{}
var tasks = map[int]Task{}
var userCounter = 1
var taskCounter = 1

// Структуры для создания JWT токена
type Claims struct {
	UserID int `json:"user_id"`
	jwt.StandardClaims
}

// Эндпоинт регистрации пользователя
func register(w http.ResponseWriter, r *http.Request) {
	var newUser User
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	newUser.ID = userCounter
	users[userCounter] = newUser
	userCounter++

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newUser)
}

// Эндпоинт логина и генерации JWT токена
func login(w http.ResponseWriter, r *http.Request) {
	var credentials User
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	for _, user := range users {
		if user.Username == credentials.Username && user.Password == credentials.Password {
			// Генерация JWT токена
			expirationTime := time.Now().Add(24 * time.Hour)
			claims := &Claims{
				UserID: user.ID,
				StandardClaims: jwt.StandardClaims{
					ExpiresAt: expirationTime.Unix(),
				},
			}
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			tokenString, err := token.SignedString(jwtKey)
			if err != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"token": tokenString,
			})
			return
		}
	}

	http.Error(w, "Invalid credentials", http.StatusUnauthorized)
}

// Эндпоинт для выхода и аннулирования JWT токена
func logout(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User logged out"))
}

// Эндпоинт для удаления пользователя и аннулирования токена
func deleteUser(w http.ResponseWriter, r *http.Request) {
	tokenString := r.Header.Get("Authorization")
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
		return
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	delete(users, claims.UserID)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User deleted"))
}

// Эндпоинт для получения одной задачи
func getTask(w http.ResponseWriter, r *http.Request) {
	// Получаем ID задачи из URL (например, /tasks/1)
	id := strings.TrimPrefix(r.URL.Path, "/tasks/")
	taskID := 0
	_, err := fmt.Sscanf(id, "%d", &taskID)
	if err != nil || taskID == 0 {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	task, exists := tasks[taskID]
	if !exists {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

// Эндпоинт для получения всех задач
func getAllTasks(w http.ResponseWriter, r *http.Request) {
	var taskList []Task
	for _, task := range tasks {
		taskList = append(taskList, task)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(taskList)
}

// Эндпоинт для создания новой задачи
func createTask(w http.ResponseWriter, r *http.Request) {
	var newTask Task
	if err := json.NewDecoder(r.Body).Decode(&newTask); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	newTask.ID = taskCounter
	tasks[taskCounter] = newTask
	taskCounter++

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTask)
}

// Эндпоинт для загрузки задач в большом объеме
func bulkUpload(w http.ResponseWriter, r *http.Request) {
	var newTasks []Task
	if err := json.NewDecoder(r.Body).Decode(&newTasks); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	for _, task := range newTasks {
		task.ID = taskCounter
		tasks[taskCounter] = task
		taskCounter++
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTasks)
}

// Эндпоинт для обновления задачи
func updateTask(w http.ResponseWriter, r *http.Request) {
	taskID := 1 // Пример, извлеките реальный ID задачи из URL

	var updatedTask Task
	if err := json.NewDecoder(r.Body).Decode(&updatedTask); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	task, exists := tasks[taskID]
	if !exists {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	task.Title = updatedTask.Title
	task.Done = updatedTask.Done
	tasks[taskID] = task

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

// Эндпоинт для удаления задачи
func deleteTask(w http.ResponseWriter, r *http.Request) {
	taskID := 1 // Пример, извлеките реальный ID задачи из URL

	_, exists := tasks[taskID]
	if !exists {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	delete(tasks, taskID)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Task with ID %d has been deleted", taskID)))
}

func main() {
	// Обработчик для статичных файлов (index.html)
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("./static"))))

	http.HandleFunc("/user/register", register)
	http.HandleFunc("/user/login", login)
	http.HandleFunc("/user/logout", logout)
	http.HandleFunc("/user/delete", deleteUser)

	// Эндпоинты для задач
	http.HandleFunc("/tasks", getAllTasks)           // Получить все задачи
	http.HandleFunc("/tasks/create", createTask)     // Создать задачу
	http.HandleFunc("/tasks/bulkupload", bulkUpload) // Загрузить задачи
	http.HandleFunc("/tasks/update/", updateTask)    // Обновить задачу
	http.HandleFunc("/tasks/delete/", deleteTask)    // Удалить задачу
	http.HandleFunc("/tasks/", getTask)              // Получить одну задачу по ID

	log.Println("Server is running on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
