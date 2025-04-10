<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-100 dark:bg-gray-900">
    <div class="w-full max-w-md p-8 bg-white dark:bg-gray-800 rounded-lg shadow-lg">
      <h2 class="text-2xl font-bold text-left mb-6 text-gray-800 dark:text-white">Login</h2>
      <form @submit.prevent="handleLogin">
        <div class="mb-6">
          <input placeholder="Password" type="password" id="password" v-model="password" required
                 class="mt-2 p-2 w-full border border-gray-300 dark:border-gray-600 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-400 dark:focus:ring-blue-500"/>
        </div>

        <button type="submit"
                class="w-full bg-blue-500 text-white py-2 rounded-md hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-400 dark:bg-blue-700 dark:hover:bg-blue-600">
          Login
        </button>
      </form>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      email: '',
      password: ''
    };
  },
  methods: {
    handleLogin() {
      // Aquí puedes agregar la lógica para autenticar al usuario
      console.log("Email:", this.email);
      console.log("Password:", this.password);
    }
  }
}
</script>

<style scoped>
/* Puedes agregar estilos adicionales aquí si es necesario */
</style>
