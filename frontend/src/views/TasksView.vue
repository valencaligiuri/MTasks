<template>
  <p>Doing <a href="" @click="logout">Logout</a></p>
  <button
    @click="toggleDarkMode"
    class="p-2 bg-blue-500 text-white rounded dark:bg-red-500"
  >
    Toggle Dark Mode
  </button>
</template>

<script setup lang="ts">

// Define el evento que vas a emitir al padre
const emit = defineEmits<{
  (e: 'toggleDarkMode'): void
}>()

const toggleDarkMode = () => {
  const Theme = (localStorage.getItem('Theme') === "dark") ? "light" : "dark"


  localStorage.setItem('Theme', Theme)
  emit('toggleDarkMode')
}

const logout = async() => {
  try {
    const response = await fetch('http://localhost:8081/api/logout', {
      method: 'GET',
      credentials: 'include',
    })
    if (!response.ok) {
      console.error(`Error: ${response.status}`)
    }

    console.log(response.json())
  } catch (e){
    console.error(e)
  }

}

</script>
