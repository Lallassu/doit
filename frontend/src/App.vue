<template>
    <div id="app">
        <ToDoList />
    </div>
</template>

<script>
import ToDoList from './components/ToDoList.vue'

export default {
    name: 'app',
    components: {
        "ToDoList": ToDoList,
    }
}
</script>


<style lang="scss">
@import url('https://use.fontawesome.com/releases/v5.7.2/css/all.css');

body {
    margin: 0;
    padding: 0;
    width: 100%;
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
#app .bm-burger-button {
    height: 20px;
    width: 20px;
    top: 15px;       
    left: 92%;
}
#app .bm-menu {
    z-index: 10;
    background: linear-gradient(to bottom right, #21656D, #0D4826);
}

#app .show-archived {
    color: #aaa;
    padding: 10px;
    font-size: 0.8rem;
    text-transform: uppercase;
    overflow: hidden;
    cursor: pointer;
}

#app .list-menu-item {
    margin-left: 20px;
    color: #fff;
    font-size: 1.3rem;
}
</style>
