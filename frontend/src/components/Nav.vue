<script setup lang="ts">
import {useRouter} from "vue-router";

const router = useRouter();

import {ArrowLeftEndOnRectangleIcon} from "@heroicons/vue/24/solid";

function logout() {
  fetch(`http://${window.location.hostname}:8081/api/logout`, {
    method: 'GET',
    credentials: 'include',
  }).then(async res => {
    if (!res.ok) {
      console.error(`Error: ${res.status}`)
      return
    }
    await router.push("/login")
  }).catch(err => console.log(err))
}
</script>

<template>
  <nav
      class="flex items-center justify-between px-6 py-4 bg-white dark:bg-gray-800 shadow-md text-gray-800 dark:text-white">
    <h1 class="text-xl font-semibold">MTasks</h1>

    <button
        @click="logout"
        class="bg-red-500 hover:bg-red-600 text-white font-medium px-4 py-2 rounded transition flex items-center gap-2"
    >
      Logout
      <arrow-left-end-on-rectangle-icon class="h-6 w-6 inline"/>
    </button>

  </nav>


</template>

<style scoped>

</style>