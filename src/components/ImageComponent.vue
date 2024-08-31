<template>
    <img :src="element.content" :style="element.css" alt="Image Component" />
  </template>
  
  <script>
  export default {
    name: 'ImageComponent',
    props: {
      element: Object,
    },
  };
  </script>
  