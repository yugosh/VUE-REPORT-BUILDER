<template>
  <div id="app">
    <h1>Report Builder</h1>
    <ReportCanvas />
  </div>
</template>

<script>
import ReportCanvas from './components/ReportCanvas.vue';

export default {
  name: 'App',
  components: {
    ReportCanvas,
  },
};
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
