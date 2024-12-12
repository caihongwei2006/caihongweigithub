<template>
    <div class="person">
      <h2>你好啊！</h2>
      <p>Name: {{ name }}</p>
      <p>Age: {{ age }}</p>
      <p>Tel: {{ tel }}</p>
      <button @click="showData">Show Data</button>
    </div>
  </template>
  
  <script setup lang="ts">
  import { ref } from 'vue'
  
  // Define component data
  const name = ref('CHW')
  const age = ref(18)
  const tel = ref('123123')
  
  function showData() {
    alert(`Data: ${name.value}, ${age.value}, ${tel.value}`)
  }
  </script>
  
  <style scoped>
  .person {
    background-color: skyblue;
    padding: 10px;
    width: 200px;
  }
  </style>
  