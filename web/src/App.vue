<template>
  <div id="app">
    <Deposit/>
  </div>
</template>

<script>
import Deposit from './components/Deposit.vue'

export default {
  name: 'App',
  components: {
    Deposit
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
  margin-top: 60px;
}
</style>
