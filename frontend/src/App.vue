<template>
  <router-view />
</template>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

html,
body {
  height: 100%;
}

body {
  font-family: var(--tm-font-ui);
  background: var(--tm-bg);
  color: var(--tm-text);
  -webkit-font-smoothing: antialiased;
  text-rendering: optimizeLegibility;
}

#app {
  height: 100%;
}

a {
  color: var(--tm-accent);
  text-decoration: none;
}

button {
  font-family: inherit;
}
</style>
