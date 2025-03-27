// Variable global para almacenar el ID de la tarea a eliminar
let taskToDelete = null;

// Función para obtener las tareas desde el servidor
function getTasks() {
  fetch("/api/tasks")
    .then((response) => response.json())
    .then((data) => displayTasks(data))
    .catch((error) => console.error("Error al obtener tareas:", error));
}

// Función para mostrar las tareas en el HTML
function displayTasks(tasks) {
  const taskList = document.getElementById("taskList");
  const emptyMessage = document.getElementById("emptyMessage");
  taskList.innerHTML = ""; // Limpiar lista actual

  if (tasks.length === 0) {
    emptyMessage.style.display = "block";
  } else {
    emptyMessage.style.display = "none";
  }

  const currentTheme = localStorage.getItem("theme") || "light";

  tasks.forEach((task) => {
    const taskItem = document.createElement("li");
    taskItem.classList.add(
      "taskItem",
      "p-4",
      "rounded-lg",
      "shadow-sm",
      "flex",
      "justify-between",
      "items-center"
    );

    // Agregar clases según el tema actual
    if (currentTheme === "dark") {
      taskItem.classList.add("bg-gray-700", "text-white");
    } else {
      taskItem.classList.add("bg-gray-100", "text-gray-900");
    }

    // Contenedor para el icono de ojo y el contenido de la tarea
    const taskInfo = document.createElement("div");
    taskInfo.classList.add("flex", "items-center", "mr-4", "taskInfo"); // Flex para alinear el ojo y el contenido

    // Botón de ver tarea completa (icono de ojo)
    const viewButton = document.createElement("button");
    viewButton.classList.add(
      "mr-4",
      "text-gray-600",
      "hover:text-gray-900",
      "focus:outline-none"
    );
    viewButton.innerHTML = '<i class="fas fa-eye"></i>'; // Usando FontAwesome para el icono

    // Función para abrir el modal con la tarea completa
    viewButton.onclick = function () {
      openViewTaskModal(task);
    };

    // Contenedor para el contenido de la tarea (título y descripción)
    const taskContent = document.createElement("div");
    taskContent.classList.add("flex-grow", "taskContent");

    // Crear el título de la tarea
    const taskTitle = document.createElement("div");
    taskTitle.classList.add("taskTitle");
    taskTitle.textContent = task.title; // Asignar solo el título

    // Crear la descripción de la tarea
    const taskDescription = document.createElement("div");
    taskDescription.classList.add("taskDescription");
    taskDescription.textContent = task.description; // Asignar solo la descripción
    const completedCheckbox = document.createElement("input");
    completedCheckbox.type = "checkbox";
    completedCheckbox.checked = task.completed;
    completedCheckbox.classList.add("completedCheckbox", "dark:mark-white");

    // Añadir el título y la descripción a taskContent
    taskContent.appendChild(taskTitle);
    taskContent.appendChild(taskDescription);

    // Colocar el icono de ojo y el contenido de la tarea en taskInfo
    taskInfo.appendChild(viewButton); // El botón del ojo a la izquierda
    taskInfo.appendChild(completedCheckbox);
    taskInfo.appendChild(taskContent); // El contenido de la tarea (título + descripción)

    // Colocar el taskInfo (ojo + contenido) en taskItem
    taskItem.appendChild(taskInfo); // Colocar el contenedor de la tarea (ojo + contenido)

    // Contenedor de botones a la derecha
    const buttonContainer = document.createElement("div");
    buttonContainer.classList.add("flex", "items-center"); // Flex para los botones alineados

    // Botón de actualizar
    const updateButton = document.createElement("button");
    updateButton.textContent = "Update";
    updateButton.classList.add(
      "ml-4",
      "bg-blue-600",
      "text-white",
      "py-2",
      "px-4",
      "rounded-lg",
      "hover:bg-blue-700"
    );
    updateButton.onclick = function () {
      openUpdateModal(task);
    };

    // Botón de eliminar tarea
    const deleteButton = document.createElement("button");
    deleteButton.textContent = "Delete";
    deleteButton.classList.add(
      "ml-4",
      "bg-red-600",
      "text-white",
      "py-2",
      "px-4",
      "rounded-lg",
      "hover:bg-red-700"
    );
    deleteButton.onclick = function () {
      deleteTask(task.id);
    };

    // Checkbox para marcar como completada
    // Agregar margen a la izquierda para separarlo

    completedCheckbox.addEventListener("change", function () {
      if (this.checked) {
        // Guardar el ID de la tarea para eliminarla tras la felicitación
        taskToDelete = task.id;
        // Mostrar el modal de felicitación
        document.getElementById("congratsModal").style.display = "flex";
      }
    });

    // Colocar los botones de Update, Delete y el checkbox en buttonContainer
    buttonContainer.appendChild(updateButton);
    buttonContainer.appendChild(deleteButton);

    // Colocar el buttonContainer (botones) en el taskItem
    taskItem.appendChild(buttonContainer);

    // Añadir taskItem a la lista
    taskList.appendChild(taskItem);
  });
}

// Crear tarea

function createTask(title, description) {
  const newTask = { title, description, completed: false };

  fetch("/api/tasks", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(newTask),
  })
    .then((response) => response.json())
    .then(() => {
      document.getElementById("createModal").style.display = "none"; // Cerrar el modal
      getTasks(); // Recargar las tareas
    })
    .catch((error) => console.error("Error al crear tarea:", error));
}

// Función para eliminar tarea
function deleteTask(id) {
  fetch(`/api/tasks/${id}`, {
    method: "DELETE",
  })
    .then((response) => response.json())
    .then(() => getTasks()) // Volver a cargar las tareas
    .catch((error) => console.error("Error al eliminar tarea:", error));
}

// Función para mostrar el modal de actualización
function openUpdateModal(task) {
  const titleInput = document.getElementById("updateTaskTitle"); // Cambiado a nuevo id
  const descriptionInput = document.getElementById("updateTaskDescription"); // Cambiado a nuevo id

  // Rellenar los campos con los datos de la tarea
  titleInput.value = task.title;
  descriptionInput.value = task.description;

  // Mostrar el modal de actualización y asegurarse de cerrar el modal de creación
  document.getElementById("createModal").style.display = "none"; // Cerrar el modal de creación
  document.getElementById("updateModal").style.display = "flex"; // Mostrar el modal de actualización

  // Actualizar tarea al enviar el formulario
  document.getElementById("updateTaskForm").onsubmit = function (e) {
    e.preventDefault();
    const updatedTitle = titleInput.value.trim();
    const updatedDescription = descriptionInput.value.trim();

    // Evitar tareas con el mismo título
    fetch("/api/tasks")
      .then((response) => response.json())
      .then((data) => {
        if (updatedTitle.length > 50) {
          alert("The title must not exceed 50 characters!");
        } else if (updatedDescription.length > 300) {
          alert("The description must not exceed 300 characters!");
        } else {
          const existingTask = data.find(
            (task) => task.title.toLowerCase() === updatedTitle.toLowerCase()
          );
          if (existingTask && updatedTitle !== task.title) {
            alert("A task with this title already exists!");
          } else {
            // Si no existe y el título es válido, crear la tarea
            updateTask(task.id, updatedTitle, updatedDescription);
            document.getElementById("updateModal").style.display = "none"; // Cerrar el modal de actual
          }
        }
      })
      .catch((error) => console.error("Error al verificar tareas:", error));
  };
}

// Función para actualizar una tarea
function updateTask(id, title, description) {
  const updatedTask = { title, description };

  fetch(`/api/tasks/${id}`, {
    method: "PUT",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(updatedTask),
  })
    .then((response) => response.json())
    .then(() => {
      getTasks(); // Volver a cargar las tareas
      document.getElementById("updateModal").style.display = "none"; // Cerrar el modal de actualización
    })
    .catch((error) => console.error("Error al actualizar tarea:", error));
}

function openViewTaskModal(task) {
  // Obtener los elementos del modal
  const modal = document.getElementById("viewTaskModal");
  const modalTitle = document.getElementById("viewTaskTitle");
  const modalDescription = document.getElementById("viewTaskDescription");

  // Asignar los valores de la tarea al modal
  modalTitle.textContent = task.title;
  modalDescription.textContent = task.description;

  // Mostrar el modal
  modal.style.display = "flex";
}

document
  .getElementById("closeUpdateModalButton")
  .addEventListener("click", function () {
    document.getElementById("updateModal").style.display = "none";
  });

// Cerrar modal de felicitación y eliminar la tarea
document
  .getElementById("congratsCloseButton")
  .addEventListener("click", function () {
    if (taskToDelete !== null) {
      deleteTask(taskToDelete);
      taskToDelete = null;
    }
    document.getElementById("congratsModal").style.display = "none";
  });

// Inicializar la carga de tareas al cargar la página
document.addEventListener("DOMContentLoaded", function () {
  getTasks();
});

// Abrir el modal de creación de tarea
document
  .getElementById("openCreateModalButton")
  .addEventListener("click", function () {
    document.getElementById("createModal").style.display = "flex";
  });

// Cerrar el modal de creación de tarea
document
  .getElementById("closeCreateModalButton")
  .addEventListener("click", function () {
    document.getElementById("createModal").style.display = "none";
  });

// Función para cerrar el modal
document
  .getElementById("closeViewModalButton")
  .addEventListener("click", function () {
    document.getElementById("viewTaskModal").style.display = "none";
  });

// Función para crear una nueva tarea
document
  .getElementById("taskForm")
  .addEventListener("submit", function (event) {
    event.preventDefault();

    const title = document.getElementById("createTaskTitle").value.trim();
    const description = document
      .getElementById("createTaskDescription")
      .value.trim();

    // Verificar si ya existe una tarea con el mismo título
    fetch("/api/tasks")
      .then((response) => response.json())
      .then((data) => {
        const maxLength = 50;
        if (title.length > maxLength) {
          alert("The title must not exceed 50 characters!");
        } else if (description.length > 300) {
          alert("The description must not exceed 300 characters!");
        } else {
          const existingTask = data.find(
            (task) => task.title.toLowerCase() === title.toLowerCase()
          );
          if (existingTask) {
            alert("A task with this title already exists!");
          } else {
            // Si no existe y el título es válido, crear la tarea
            createTask(title, description);
          }
        }
      })
      .catch((error) => console.error("Error al verificar tareas:", error));
  });
