<template>
  <div :class="{ 'dark': isDarkMode }">
    <router-view @toggleDarkMode="applyTheme"/>
  </div>
</template>

<script setup lang="ts">
import {ref, onMounted, watch} from "vue";
import {useRouter} from "vue-router";

const isDarkMode = ref(false);
const router = useRouter();

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
    if (!res.ok) {
      console.log(`Status: ${JSON.stringify(await res.json())}`)
      await router.push('/auth');
    } else {
      console.log(`Status: ${JSON.stringify(await res.json())}`)
    }
  }).catch(err => {
    console.error(`eh ${err}`)
  })

}

onMounted(() => {
  authHandler()
  applyTheme();
});

watch(
    () => router.fullPath,
    () => {
      authHandler();
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
