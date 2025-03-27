package main

import (
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	_ "modernc.org/sqlite"
)

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

func loadTasks(db *sql.DB) ([]Task, error) {
	// Consultar las tareas
	rows, err := db.Query("SELECT id, title, description, completed FROM tasks")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Completed); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}
func saveTasks(db *sql.DB, tasks []Task) error {
	// Comenzar una transacción
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insertar o actualizar las tareas en la base de datos
	for _, task := range tasks {
		// Verificar si la tarea ya existe
		var exists bool
		err := tx.QueryRow("SELECT EXISTS(SELECT 1 FROM tasks WHERE id = ?)", task.ID).Scan(&exists)
		if err != nil {
			return err
		}

		if exists {
			// Si la tarea existe, actualizarla
			_, err = tx.Exec("UPDATE tasks SET title = ?, description = ?, completed = ? WHERE id = ?",
				task.Title, task.Description, task.Completed, task.ID)
			if err != nil {
				return err
			}
		} else {
			// Si la tarea no existe, insertarla
			_, err = tx.Exec("INSERT INTO tasks (id, title, description, completed) VALUES (?, ?, ?, ?)",
				task.ID, task.Title, task.Description, task.Completed)
			if err != nil {
				return err
			}
		}
	}

	// Confirmar la transacción
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
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

func getTasks(db *sql.DB, c *gin.Context) {
	if !isAuthenticated(c) {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	// Obtener las tareas desde la base de datos
	tasks, err := loadTasks(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Si no hay tareas, devolver una lista vacía
	if tasks == nil {
		tasks = []Task{}
	}
	c.JSON(http.StatusOK, tasks)
}

func createTask(db *sql.DB, c *gin.Context) {
	if !isAuthenticated(c) {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	var newTask Task
	// Intentar enlazar el JSON recibido con la estructura Task
	if err := c.ShouldBindJSON(&newTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Cargar las tareas existentes desde la base de datos
	tasks, err := loadTasks(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Validar que no exista una tarea con el mismo título
	if newTask.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "The title cannot be empty"})
		return
	} else if newTask.Description == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "The description cannot be empty"})
		return
	} else if len(newTask.Title) > 50 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "The title cannot be longer than 50 characters"})
		return
	} else if len(newTask.Description) > 300 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "The description cannot be longer than 300 characters"})
		return
	}

	// Verificar si ya existe una tarea con el mismo título
	for _, t := range tasks {
		if t.Title == newTask.Title {
			c.JSON(http.StatusBadRequest, gin.H{"error": "A task with that title already exists"})
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

	// Guardar las tareas actualizadas
	if err := saveTasks(db, tasks); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Responder con un mensaje de éxito
	c.JSON(http.StatusCreated, gin.H{"message": "Task created"})
}

// Actualiza una tarea por ID
func updateTask(db *sql.DB, c *gin.Context) {
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

	tasks, err := loadTasks(db)
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
		} else if len(updatedTask.Description) > 300 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "The description cannot be longer than 300 characters"})
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

	if err := saveTasks(db, tasks); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task updated"})
}

func deleteTask(db *sql.DB, c *gin.Context) {
	if !isAuthenticated(c) {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	// Obtener el ID de la tarea desde la URL
	id := c.Param("id")
	taskID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	// Cargar las tareas desde la base de datos
	tasks, err := loadTasks(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load tasks"})
		return
	}

	// Buscar la tarea por ID y eliminarla de la lista
	var found bool
	for i, t := range tasks {
		if t.ID == taskID {
			// Eliminar la tarea de la lista
			tasks = append(tasks[:i], tasks[i+1:]...)
			found = true
			break
		}
	}

	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	// Eliminar la tarea de la base de datos
	_, err = db.Exec("DELETE FROM tasks WHERE id = ?", taskID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Enviar una respuesta de éxito
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

	var err error
	db, err := sql.Open("sqlite", "./tasks.db")
	if err != nil {
		log.Fatal("Can't open database", err)
	} else if err := db.Ping(); err != nil {
		log.Fatal("Can't connect to database", err)
	}

	// Crear la tabla "tasks" si no existe
	createTableSQL := `CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT,
		description TEXT,
		completed BOOLEAN
		);`
	if _, err := db.Exec(createTableSQL); err != nil {
		log.Fatal("Can't create table", err)
	}

	defer db.Close()

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
	r.GET("/api/tasks", func(ctx *gin.Context) {
		getTasks(db, ctx)
	})
	r.POST("/api/tasks", func(ctx *gin.Context) {
		createTask(db, ctx)
	})
	r.PUT("/api/tasks/:id", func(ctx *gin.Context) {
		updateTask(db, ctx)
	})
	r.DELETE("/api/tasks/:id", func(ctx *gin.Context) {
		deleteTask(db, ctx)
	})

	r.Static("/static", "./static")

	if err := r.Run(":8080"); err != nil {
		fmt.Println(err)
	}
}
