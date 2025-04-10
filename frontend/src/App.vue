<template>
  <div :class="{'dark': isDarkMode}">
    <router-view @toggleDarkMode="checkTheme"/>
  </div>
</template>

<script>
export default {
  data() {
    return {
      isDarkMode: false,
    }
  },
  mounted() {
    this.checkTheme()
  },
  watch: {
    $route() {
      this.checkTheme()
    }
  },
  methods: {
    checkTheme() {
      const savedTheme = localStorage.getItem("darkMode");
      if (savedTheme) {
        this.isDarkMode = savedTheme === 'true';
      } else {
        this.isDarkMode = window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches;
      }
    }
  }
}
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
