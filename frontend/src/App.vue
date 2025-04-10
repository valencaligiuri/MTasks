<template>
  <div :class="{ dark: isDarkMode.value }">
    <router-view @toggleDarkMode="applyTheme" />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from "vue";
import { useRoute } from "vue-router";

const isDarkMode = ref(false);
const route = useRoute();

function applyTheme() {
  const savedTheme = localStorage.getItem("darkMode");
  const prefersDark = window.matchMedia("(prefers-color-scheme: dark)").matches;
  const enabled = savedTheme === null ? prefersDark : savedTheme === "true";

  isDarkMode.value = enabled;

  const html = document.documentElement;
  if (enabled) {
    html.classList.add("dark");
  } else {
    html.classList.remove("dark");
  }
}
onMounted(() => {
  applyTheme();
});

watch(
  () => route.fullPath,
  () => {
    applyTheme();
  }
);
</script>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
}
</style>
