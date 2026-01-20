<template>
  <div id="app">
    <ToDoList />
  </div>
</template>

<script>
import ToDoList from './components/ToDoList.vue'

export default {
  name: 'App',
  components: {
    ToDoList
  }
}
</script>

<style lang="scss">
@import url('https://use.fontawesome.com/releases/v5.7.2/css/all.css');

body {
  margin: 0;
  padding: 0;
  width: 100%;
  min-height: 100vh;
  display: flex;
  justify-content: center;
  align-items: flex-start;
  background: repeating-linear-gradient(to bottom right, #21656D, #C3D98D);
}

#app {
  font-family: Arial, Helvetica, sans-serif;
  font-size: 1rem;
  color: #333;
  width: 100%;
  max-width: 768px;
}
</style>
