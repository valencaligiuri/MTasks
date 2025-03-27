// Selecciona el botón y el ícono dentro del botón
const changeThemeButton = document.getElementById("changeThemeButton");
const themeIcon = document.getElementById("themeIcon");

// Configuración de Tailwind para el modo oscuro
tailwind.config = {
  darkMode: "class", // Asegúrate de que Tailwind use clases para el modo oscuro
};

// Función para establecer el tema
function setTheme(theme) {
  if (theme === "dark") {
    document.documentElement.classList.add("dark"); // Activa el modo oscuro en el <html>
    themeIcon.classList.remove("fa-sun");
    themeIcon.classList.add("fa-moon");
  } else {
    document.documentElement.classList.remove("dark"); // Remueve el modo oscuro del <html>
    themeIcon.classList.remove("fa-moon");
    themeIcon.classList.add("fa-sun");
  }
  localStorage.setItem("theme", theme); // Guarda el tema seleccionado
}

// Alternar entre los temas claro y oscuro
function toggleTheme() {
    const currentTheme = localStorage.getItem('theme') || 'light';
    const newTheme = currentTheme === 'light' ? 'dark' : 'light';
  
    // Cambiar el tema global
    setTheme(newTheme);
  
    // Actualizar las tareas ya existentes
    updateTaskItemsTheme(newTheme);
  
    // Guardar el nuevo tema en localStorage
    localStorage.setItem('theme', newTheme);
    console.log(`${newTheme.charAt(0).toUpperCase() + newTheme.slice(1)} tema activado`);
  }

  function updateTaskItemsTheme(theme) {
    // Obtener todos los taskItem existentes
    const taskItems = document.querySelectorAll(".taskItem");
  
    taskItems.forEach((taskItem) => {
      // Limpiar las clases de fondo y texto previas
      taskItem.classList.remove("bg-gray-100", "text-gray-900", "bg-gray-700", "text-white");
  
      // Aplicar nuevas clases según el tema
      if (theme === "dark") {
        taskItem.classList.add("bg-gray-700", "text-white");
      } else {
        taskItem.classList.add("bg-gray-100", "text-gray-900");
      }
  
    });
  }
  

// Inicializa el tema al cargar la página
document.addEventListener("DOMContentLoaded", () => {
  const savedTheme = localStorage.getItem("theme") || "light"; // Recupera el tema guardado o 'light' por defecto
  setTheme(savedTheme); // Aplica el tema al cargar la página
});

// Añadir el evento de clic al botón para alternar el tema
changeThemeButton.addEventListener("click", toggleTheme);
