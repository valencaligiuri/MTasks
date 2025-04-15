<template>
  <div class="min-h-screen flex items-center justify-center">
    <div class="w-full max-w-md p-8 bg-white dark:bg-gray-800 rounded-lg shadow-lg">
      <h2 class="text-2xl font-bold text-left mb-6">Auth</h2>
      <form @submit.prevent="handleLogin">
        <div class="mb-6">
          <input
              placeholder="Password"
              type="password"
              id="password"
              v-model="password"
              required
              class="mt-2 p-2 w-full bg-gray-100 dark:bg-gray-700 border border-gray-300 dark:border-gray-600 dark:text-white rounded-md focus:outline-none focus:ring-2 focus:ring-blue-400 dark:focus:ring-blue-500"
          />
        </div>

        <button
            type="submit"
            class="w-full bg-blue-500 text-white py-2 rounded-md hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-400 dark:bg-blue-700 dark:hover:bg-blue-600"
        >
          Login
        </button>
      </form>
      <p class="mt-2 text-sm text-red-500 font-medium" v-if="err">{{ err }}</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import {capitalize, ref} from 'vue'
import {useRouter} from 'vue-router'

const password = ref('')
const err = ref('')
const router = useRouter()

function handleLogin() {
  const pwd = {
    "password": password.value
  }

  fetch(`http://${window.location.hostname}:8081/api/login`, {
    method: 'POST',
    headers: {'Content-Type': 'application/json'},
    credentials: "include",
    body: JSON.stringify(pwd)
  }).then(async res => {
    const data = await res.json()

    if (!res.ok) {
      err.value = capitalize(data.error);
    } else {
      await router.push('/tasks')
    }
  }).catch(e => {
    console.error(e)
    err.value = "failed connecting to server.."
  })
}
</script>

<style scoped>
/* Puedes agregar estilos adicionales aquí si es necesario */
</style>
