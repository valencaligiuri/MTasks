<template>
  <Nav />

  <!-- Botón Modo Oscuro -->
  <button
      @click="toggleDarkMode"
      class="fixed bottom-4 left-4 w-12 h-12 rounded-full bg-blue-600 text-white dark:bg-blue-500 dark:text-gray-900 flex items-center justify-center shadow-lg hover:bg-blue-700 dark:hover:bg-blue-400 transition"
  >
    <sun-icon class="h-6 w-6 block dark:hidden" />
    <moon-icon class="h-6 w-6 hidden dark:block" />
  </button>

  <!-- Botón Crear Tarea -->
  <button
      @click="openCreateModal"
      class="fixed bottom-4 right-4 w-12 h-12 rounded-full bg-blue-600 text-white dark:bg-blue-500 dark:text-gray-900 flex items-center justify-center shadow-lg hover:bg-blue-700 dark:hover:bg-blue-400 transition"
  >
    <plus-icon class="h-6 w-6" />
  </button>

  <!-- Contenedor de tareas -->
  <div class="px-4 py-8 sm:px-8 md:px-16 lg:px-32 xl:px-48 max-w-screen-xl mx-auto">
    <div class="bg-white dark:bg-gray-800 rounded-2xl shadow-md border border-gray-200 dark:border-gray-700 p-6">
      <h1 class="text-3xl font-bold mb-6 text-gray-900 dark:text-white relative inline-block after:block after:w-full after:h-[2px] after:bg-blue-500 after:mt-1">My Tasks</h1>

      <Task v-for="task in tasks" :key="task.id" :task="task" />

      <p
          v-if="tasks.length === 0"
          class="text-center text-gray-500 dark:text-gray-400 italic mt-4"
      >
        No tasks available!
      </p>
    </div>
  </div>



  <!-- CreateTaskModal.vue -->
  <div
      v-if="showCreateModal"
      class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-40"
  >
    <div class="bg-white dark:bg-gray-800 p-6 rounded-lg w-full max-w-xl shadow-xl flex flex-col min-h-[500px]">
      <h2 class="text-xl font-bold mb-4 text-gray-900 dark:text-white">Create Task</h2>

      <!-- Formulario -->
      <form @submit.prevent="createTask" class="flex-grow flex flex-col">
        <!-- Inputs -->
        <input
            type="text"
            placeholder="Title"
            v-model="newTaskTitle"
            class="w-full mb-3 px-3 py-2 rounded border border-gray-300 dark:border-gray-600
           bg-white dark:bg-gray-700 text-gray-900 dark:text-white"
        />
        <textarea
            placeholder="Description"
            v-model="newTaskDescription"
            class="w-full mb-3 px-3 py-2 rounded border border-gray-300 dark:border-gray-600
           bg-white dark:bg-gray-700 text-gray-900 dark:text-white resize-none flex-grow"
        ></textarea>

        <!-- Mensaje de error -->
        <p v-if="err" class="mb-3 text-sm text-red-500 font-medium text-center">
          {{ err }}
        </p>

        <!-- Botones -->
        <div class="flex flex-col space-y-2 mt-auto">
          <button
              type="submit"
              class="px-4 py-2 rounded bg-blue-600 hover:bg-blue-700 text-white dark:bg-blue-500 dark:hover:bg-blue-600"
          >
            Create
          </button>
          <button
              type="button"
              @click="closeCreateModal"
              class="px-4 py-2 rounded bg-gray-300 hover:bg-gray-400 dark:bg-gray-600 dark:hover:bg-gray-500"
          >
            Cancel
          </button>
        </div>
      </form>
    </div>
  </div>




</template>

<script setup lang="ts">

import Nav from "../components/Nav.vue";
import {SunIcon, MoonIcon, PlusIcon} from '@heroicons/vue/24/solid';
import Task from "../components/Task.vue";
import {ref, watch} from 'vue'
import {useTasks} from "../composables/useTasks.ts";
import {useRoute} from "vue-router";

const route = useRoute()

// Define el evento que vas a emitir al padre

let {tasks, fetchTasks} = useTasks();
fetchTasks()

const emit = defineEmits<{
  (e: 'toggleDarkMode'): void
}>()


const showCreateModal = ref(false)
const newTaskTitle = ref('')
const newTaskDescription = ref('')
let err = ref('')

function openCreateModal() {
  showCreateModal.value = true
}

function closeCreateModal() {
  showCreateModal.value = false
  newTaskTitle.value = ''
  newTaskDescription.value = ''
}

function createTask() {

  if (!newTaskTitle.value.trim()) return;


  fetch("http://localhost:8081/api/tasks", {
    method: "POST",
    headers: {'Content-Type': 'application/json'},
    credentials: 'include',
    body: JSON.stringify({
      title: newTaskTitle.value,
      description: newTaskDescription.value,
    }),
  }).then(async res => {
    if (!res.ok) {
      const data = await res.json()
      err.value = data.error
      return
    } else {
      await fetchTasks()
      closeCreateModal()

    }


  }).catch(err => console.error(err));
}


const toggleDarkMode = () => {
  const Theme = (localStorage.getItem('Theme') === "dark") ? "light" : "dark"


  localStorage.setItem('Theme', Theme)
  emit('toggleDarkMode')

  watch(
      () => route.fullPath,
      async () => {
        await fetchTasks()
      },
      {immediate: true}
  )


}


</script>
