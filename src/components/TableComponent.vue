<template>
    <table :style="element.css">
      <thead>
        <tr>
          <th v-for="header in element.content.headers" :key="header">{{ header }}</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="row in element.content.rows" :key="row[0]">
          <td v-for="cell in row" :key="cell">{{ cell }}</td>
        </tr>
      </tbody>
    </table>
  </template>
  
  <script>
  export default {
    name: 'TableComponent',
    props: {
      element: Object,
    },
  };
  </script>
  