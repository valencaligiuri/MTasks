<template>
  <router-link to="/auth">Go to Login</router-link>
  <FirstComponent></FirstComponent>
  <button @click="toggleDarkMode" class="p-2 bg-blue-500 text-white rounded dark:bg-red-500">Toggle Dark Mode</button>
</template>
<script setup>
import FirstComponent from "@/components/FirstComponent.vue";

</script>

<script>
export default {
  emits: ["toggleDarkMode"],
  methods: {
    toggleDarkMode() {
      let isDarkMode = !(localStorage.getItem('darkMode') === "true"); // Alternar el valor
      localStorage.setItem('darkMode', isDarkMode); // Guardar la preferencia en localStorage
      if (typeof document !== "undefined") {
        this.$emit('toggleDarkMode');
      }
    },
  },
}
</script>