<template>
  <router-link to="/auth">Go to Login</router-link>
  <FirstComponent />
  <button @click="toggleDarkMode" class="p-2 bg-blue-500 text-white rounded dark:bg-red-500">
    Toggle Dark Mode
  </button>
</template>

<script setup lang="ts">
import FirstComponent from "../components/FirstComponent.vue"

// Define el evento que vas a emitir al padre
const emit = defineEmits<{
  (e: 'toggleDarkMode'): void
}>()

const toggleDarkMode = () => {
  const isDarkMode = !(localStorage.getItem('darkMode') === "true")
  localStorage.setItem('darkMode', JSON.stringify(isDarkMode))
  emit('toggleDarkMode')
}
</script>
