<template>
  <div class="task-item p-4 border-b flex flex-col gap-2">
    <div class="flex flex-wrap items-start">
      <!-- Contenedor del título y descripción -->
      <div class="flex-1 min-w-0 text-left break-words overflow-hidden mr-5">
        <p class="text-lg font-semibold truncate" :class="{ 'line-through text-gray-500': task.completed }">
          {{ capitalizeSentence(task.title) }}
        </p>
        <!-- Descripción visible solo en pantallas más grandes -->
        <p class="text-sm text-gray-600 dark:text-gray-300 hidden sm:block">
          {{ capitalizeSentence(task.description) }}
        </p>
      </div>

      <!-- Contenedor de los botones, en forma de columna en pantallas pequeñas -->
      <div class="flex gap-2 flex-wrap justify-start min-w-0">
        <button
            @click="markAsDone"
            :disabled="task.completed"
            class="bg-blue-500 hover:bg-blue-600 text-white px-2 py-1 rounded text-sm disabled:opacity-50"
        >
          <CheckIcon class="h-5 w-5"/>
        </button>
        <button
            @click="openEditModal"
            class="bg-yellow-400 hover:bg-yellow-500 text-white px-2 py-1 rounded text-sm"
        >
          <PencilIcon class="h-5 w-5"/>
        </button>
        <button
            @click="deleteTask"
            class="bg-red-500 hover:bg-red-600 text-white px-2 py-1 rounded text-sm"
        >
          <TrashIcon class="h-5 w-5"/>
        </button>
        <!-- Botón de vista solo lectura -->
        <button
            @click="openViewModal"
            class="bg-green-500 hover:bg-green-600 text-white px-2 py-1 rounded text-sm"
        >
          <EyeIcon class="h-5 w-5"/>
        </button>
      </div>
    </div>

    <!-- Completion Modal -->
    <div
        v-if="showCompletionModal"
        class="fixed inset-0 bg-black bg-opacity-50 flex justify-center items-center z-50"
    >
      <div class="bg-white dark:bg-gray-800 p-6 rounded-lg shadow-md text-center w-72">
        <h2 class="text-xl font-semibold mb-4 text-green-600">Task Completed!</h2>
        <p class="text-gray-700 dark:text-gray-200">Great job! 🎉</p>
        <button
            @click="closeCompletionModal"
            class="mt-4 bg-green-500 hover:bg-green-600 text-white px-4 py-2 rounded"
        >
          Close
        </button>
      </div>
    </div>

    <!-- Edit Modal -->
    <div v-if="showEditModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-40">
      <div class="bg-white dark:bg-gray-800 p-6 rounded-lg w-full max-w-xl shadow-xl flex flex-col min-h-[500px]">
        <h2 class="text-xl font-bold mb-4 text-gray-900 dark:text-white">Edit Task</h2>

        <!-- Inputs -->
        <input
            type="text"
            placeholder="Title"
            v-model="editedTask.title"
            class="w-full mb-3 px-3 py-2 rounded border border-gray-300 dark:border-gray-600
             bg-white dark:bg-gray-700 text-gray-900 dark:text-white"
        />
        <textarea
            placeholder="Description"
            v-model="editedTask.description"
            class="w-full mb-3 px-3 py-2 rounded border border-gray-300 dark:border-gray-600
             bg-white dark:bg-gray-700 text-gray-900 dark:text-white flex-grow resize-none"
        ></textarea>

        <!-- Error Message -->
        <p v-if="err" class="mb-3 text-sm text-red-500 font-medium text-center">
          {{ err }}
        </p>

        <!-- Buttons -->
        <div class="flex flex-col space-y-2 mt-auto">
          <button
              @click="updateTask"
              class="px-4 py-2 rounded bg-blue-600 hover:bg-blue-700 text-white dark:bg-blue-500 dark:hover:bg-blue-600"
          >
            Update
          </button>
          <button
              @click="closeEditModal"
              class="px-4 py-2 rounded bg-gray-300 hover:bg-gray-400 dark:bg-gray-600 dark:hover:bg-gray-500"
          >
            Cancel
          </button>
        </div>
      </div>
    </div>

    <!-- View Modal -->
    <div
        v-if="showViewModal"
        class="fixed inset-0 bg-black bg-opacity-50 flex justify-center items-center z-50"
    >
      <div
          class="bg-white dark:bg-gray-800 p-6 rounded-lg shadow-md text-center w-full max-w-md min-h-[300px] flex flex-col justify-between">
        <div>
          <p
              class="text-lg text-left font-semibold border-b border-gray-300 dark:border-gray-700 pb-2"
              :class="{ 'line-through text-gray-500': task.completed }"
          >
            {{ capitalizeSentence(task.title) }}
          </p>
          <p class="text-sm text-left text-gray-600 dark:text-gray-300 mt-2">
            {{ capitalizeSentence(task.description) }}
          </p>
        </div>

        <div class="mt-6">
          <button
              @click="closeViewModal"
              class="bg-gray-500 hover:bg-gray-600 text-white px-4 py-2 rounded text-sm"
          >
            Close
          </button>
        </div>
      </div>
    </div>
  </div>

</template>

<script setup lang="ts">
import {ref} from 'vue'
import {CheckIcon, PencilIcon, TrashIcon, EyeIcon} from '@heroicons/vue/24/solid'
import {useTasks} from "../composables/useTasks.ts";

let {fetchTasks} = useTasks();

const props = defineProps<{
  task: {
    id: number
    title: string
    description: string
    completed: boolean
  }
}>()

const emit = defineEmits<{
  (e: 'edit', task: any): void
  (e: 'delete', id: number): void
}>()

const showCompletionModal = ref(false)
const showEditModal = ref(false)
const showViewModal = ref(false) // Nuevo: Modal de vista previa
const editedTask = ref({...props.task})
const err = ref('')

const markAsDone = async () => {
  fetch(`http://${window.location.hostname}:8081/api/tasks/${props.task.id}`, {
    method: 'DELETE',
    headers: {'Content-Type': 'application/json'},
    credentials: 'include',
  }).then(res => {

  }).catch(e => {
    console.error(e)
  })

  props.task.completed = true
  showCompletionModal.value = true
}

const closeCompletionModal = () => {
  fetchTasks()
  showCompletionModal.value = false
}

const openEditModal = () => {
  editedTask.value = {...props.task}
  showEditModal.value = true
}

const closeEditModal = () => {
  showEditModal.value = false
}

const updateTask = () => {
  fetch(`http://${window.location.hostname}:8081/api/tasks/` + props.task.id, {
    method: 'PUT',
    headers: {'Content-Type': 'application/json'},
    credentials: 'include',
    body: JSON.stringify(editedTask.value),
  }).then(async res => {
    if (!res.ok) {
      const data = await res.json()
      err.value = data.error
      return
    } else {
      closeEditModal()
      await fetchTasks()
    }
  }).catch(e => {
    console.error(e)
  })
  emit('edit', editedTask.value)
}

const deleteTask = () => {
  fetch(`http://${window.location.hostname}:8081/api/tasks/${props.task.id}`, {
    method: 'DELETE',
    headers: {'Content-Type': 'application/json'},
    credentials: 'include'
  }).then(async res => {
    fetchTasks()
  }).catch(
      err => {
        console.error(err)
      }
  )
  emit('delete', props.task.id)
}

// Funciones nuevas para el modal de vista previa
const openViewModal = () => {
  showViewModal.value = true
}

const closeViewModal = () => {
  showViewModal.value = false
}

function capitalizeSentence(text: string) {
  return text.charAt(0).toUpperCase() + text.slice(1).toLowerCase();
}

</script>


<style scoped>
.task-item {
  transition: background-color 0.2s;
}
</style>
