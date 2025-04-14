<template>
  <div class="bg-gray-100 dark:bg-gray-900 text-gray-800 dark:text-white container-app" :class="{ 'dark': isDarkMode }">
    <router-view @toggleDarkMode="applyTheme"/>
  </div>
</template>

<script setup lang="ts">
import {ref, watch} from "vue";
import {useRouter, useRoute} from "vue-router";

const isDarkMode = ref(false);
const router = useRouter();
const route = useRoute();

function applyTheme() {
  const savedTheme = localStorage.getItem("Theme");

  const prefersDark = window.matchMedia("(prefers-color-scheme: dark)").matches;
  const enabled = savedTheme === null
      ? prefersDark
      : savedTheme === "dark";


  isDarkMode.value = enabled;


  const html = document.documentElement;
  if (enabled) {
    html.classList.add("dark");
  } else {
    html.classList.remove("dark");
  }
}

function authHandler() {
  fetch('http://localhost:8081/api/auth', {
    method: 'GET',
    credentials: 'include',
  }).then(async res => {
    if (res.ok) {
      if (route.path === "/login"){
        await router.replace("/tasks");
      }
    } else if(res.status === 401) {
      await router.replace('/login');
    } else {
      throw new Error(res.statusText);
    }
  }).catch(err => {
    console.error(`error: ${err}`)
  })
}
watch(
    () => route.fullPath,
    () => {
      authHandler();
      applyTheme();
    },
    { immediate: true }
);
</script>

<style>
.container-app {
  width: 100%;
  height: 100%;
}
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  width: 100vw;
  height: 100vh;
}
</style>
