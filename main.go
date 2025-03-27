package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

var tasksFile = "tasks.json"
var serverPassword string
var sessionID string // Identificador global de sesión

type Task struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

func generateSessionID() string {
	// Crear un generador de números aleatorios con ChaCha8
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Generar 16 bytes aleatorios
	randomBytes := make([]byte, 16)
	_, err := r.Read(randomBytes)
	if err != nil {
		panic(err) // Manejar el error adecuadamente
	}

	// Convertir los bytes aleatorios a una representación hexadecimal
	return hex.EncodeToString(randomBytes)
}

func isAuthenticated(c *gin.Context) bool {
	session := sessions.Default(c)
	authenticated := session.Get("authenticated")
	savedSessionID := session.Get("session_id")

	return authenticated != nil && authenticated.(bool) && savedSessionID == sessionID
}

func loadTasks() ([]Task, error) {
	var tasks []Task
	file, err := os.ReadFile(tasksFile)
	if err != nil {
		if os.IsNotExist(err) {
			// Si el archivo no existe, crearlo
			_, err := os.Create(tasksFile) // Crear el archivo vacío
			if err != nil {
				return nil, err
			}
			return tasks, nil // Retornar lista vacía
		}
		return nil, err
	}

	// Si el archivo está vacío, retorna una lista vacía
	if len(file) == 0 {
		return tasks, nil
	}

	err = json.Unmarshal(file, &tasks)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func saveTasks(tasks []Task) error {
	file, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(tasksFile, file, 0644)
}

func homePage(c *gin.Context) {
	if isAuthenticated(c) {
		c.Redirect(http.StatusFound, "/tasks")
	} else {
		c.Redirect(http.StatusFound, "/login")
	}
}

func loginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

func loginPost(c *gin.Context) {
	password := c.PostForm("password")

	if password == serverPassword {
		session := sessions.Default(c)
		session.Set("authenticated", true)
		session.Set("session_id", sessionID)
		session.Save()
		c.Redirect(http.StatusFound, "/tasks")
	} else {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{"error": "Incorrect password"})
	}
}

func logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.Redirect(http.StatusFound, "/login")
}

func getTasks(c *gin.Context) {
	if !isAuthenticated(c) {
		c.Redirect(http.StatusFound, "/login")
		return
	}
	tasks, err := loadTasks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Si no hay tareas, simplemente devolver una lista vacía
	if tasks == nil {
		tasks = []Task{}
	}
	c.JSON(http.StatusOK, tasks)
}

func createTask(c *gin.Context) {
	if !isAuthenticated(c) {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	var newTask Task
	if err := c.ShouldBindJSON(&newTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tasks, err := loadTasks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Validar que no exista una tarea con el mismo título
	for _, t := range tasks {
		if t.Title == newTask.Title {
			c.JSON(http.StatusBadRequest, gin.H{"error": "A task with that title already exists"})
			return
		} else if newTask.Title == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "The title cannot be empty"})
			return

		} else if newTask.Description == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "The description cannot be empty"})
			return
		} else if len(newTask.Title) > 50 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "The title cannot be longer than 50 characters"})
			return
		}
	}

	// Calcular un ID único basado en el máximo existente
	maxID := 0
	for _, t := range tasks {
		if t.ID > maxID {
			maxID = t.ID
		}
	}
	newTask.ID = maxID + 1
	newTask.Completed = false
	tasks = append(tasks, newTask)

	if err := saveTasks(tasks); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Task created"})
}

// Actualiza una tarea por ID
func updateTask(c *gin.Context) {
	if !isAuthenticated(c) {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	id := c.Param("id")
	taskID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	var updatedTask Task
	if err := c.ShouldBindJSON(&updatedTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tasks, err := loadTasks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Validar que no exista otra tarea con el mismo título
	for _, t := range tasks {
		if t.Title == updatedTask.Title && t.ID != taskID {
			c.JSON(http.StatusBadRequest, gin.H{"error": "A task with that title already exists"})
			return
		} else if updatedTask.Title == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "The title cannot be empty"})
			return
		} else if updatedTask.Description == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "The description cannot be empty"})
			return
		} else if len(updatedTask.Title) > 50 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "The title cannot be longer than 50 characters"})
			return
		}
	}

	// Buscar la tarea por ID y actualizar solo el título y la descripción
	found := false
	for i, task := range tasks {
		if task.ID == taskID {
			tasks[i].Title = updatedTask.Title
			tasks[i].Description = updatedTask.Description
			found = true
			break
		}
	}

	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	if err := saveTasks(tasks); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task updated"})
}

func deleteTask(c *gin.Context) {
	if !isAuthenticated(c) {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	tasks, err := loadTasks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	id := c.Param("id")
	taskID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	found := false
	for i, task := range tasks {
		if task.ID == taskID {
			tasks = append(tasks[:i], tasks[i+1:]...)
			found = true
			break
		}
	}

	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	if err := saveTasks(tasks); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task deleted"})
}

func tasksPage(c *gin.Context) {
	if !isAuthenticated(c) {
		c.Redirect(http.StatusFound, "/login")
		return
	}
	c.HTML(http.StatusOK, "tasks.html", nil)
}

func main() {
	fmt.Print("Enter a password to start the server: ")
	fmt.Scanln(&serverPassword)

	sessionID = generateSessionID()

	store := cookie.NewStore([]byte("secret"))

	r := gin.Default()
	r.Use(sessions.Sessions("my_session", store))

	r.LoadHTMLGlob("templates/*")

	r.GET("/", homePage)
	r.GET("/login", loginPage)
	r.POST("/login", loginPost)
	r.GET("/logout", logout)

	r.GET("/tasks", tasksPage)
	r.GET("/api/tasks", getTasks)
	r.POST("/api/tasks", createTask)
	r.PUT("/api/tasks/:id", updateTask)
	r.DELETE("/api/tasks/:id", deleteTask)

	r.Static("/static", "./static")

	if err := r.Run(":8080"); err != nil {
		fmt.Println(err)
	}
}
