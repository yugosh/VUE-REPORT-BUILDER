<template>
    <div class="report-canvas">
      <div class="controls">
        <button @click="addRow">Add Row</button>
        <button @click="generatePDF">Generate PDF</button>
      </div>
      <div class="a4-page">
        <div v-for="(row, rowIndex) in rows" :key="rowIndex" class="row">
          <div class="row-controls">
            <button @click="addColumn(rowIndex)">Add Column</button>
          </div>
          <div class="columns-container">
            <div 
              v-for="(column, columnIndex) in row.columns" 
              :key="columnIndex" 
              class="column-item" 
              :style="{ width: column.width + '%' }"
            >
              <div class="column-controls">
                <label>Column Width: 
                  <input type="range" min="10" max="100" v-model="column.width">
                  {{ column.width }}%
                </label>
                <button @click="addContent('text', rowIndex, columnIndex)">Add Text</button>
                <button @click="addContent('image', rowIndex, columnIndex)">Add Image</button>
                <button @click="addContent('field', rowIndex, columnIndex)">Add Field</button>
                <button @click="addContent('table', rowIndex, columnIndex)">Add Table</button>
              </div>
              <div class="content-container">
                <div 
                  v-for="(content, contentIndex) in column.contents" 
                  :key="contentIndex" 
                  class="content-item" 
                  :style="{ width: content.width + '%' }"
                >
                  <div class="content-controls">
                    <label>Content Width: 
                      <input type="range" min="10" max="100" v-model="content.width">
                      {{ content.width }}%
                    </label>
                  </div>
                  <div v-if="content.type === 'text'">Text Content</div>
                  <div v-if="content.type === 'image'">Image Content</div>
                  <div v-if="content.type === 'field'">Field Content</div>
                  <div v-if="content.type === 'table'">Table Content</div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </template>
  
  <script>
  import jsPDF from 'jspdf';
  import 'jspdf-autotable';
  
  export default {
    data() {
      return {
        rows: [],
      };
    },
    methods: {
      addRow() {
        this.rows.push({ columns: [] });
      },
      addColumn(rowIndex) {
        const row = this.rows[rowIndex];
        const newWidth = 100 / (row.columns.length + 1);
        row.columns.forEach(column => {
          column.width = newWidth;
        });
        row.columns.push({ width: newWidth, contents: [] });
      },
      addContent(type, rowIndex, columnIndex) {
        this.rows[rowIndex].columns[columnIndex].contents.push({ type, width: 100 });
      },
      generatePDF() {
        const doc = new jsPDF('p', 'pt', 'a4');
        let yPos = 20; // Initial Y position in the PDF
  
        this.rows.forEach((row) => {
          row.columns.forEach((column) => {
            const columnXPos = (column.width / 100) * 794; // Calculate the X position based on the column width
            let columnYPos = yPos;
  
            column.contents.forEach(content => {
              if (content.type === 'text') {
                doc.text('Text Content', columnXPos, columnYPos);
              } else if (content.type === 'image') {
                doc.text('Image Content', columnXPos, columnYPos);
              } else if (content.type === 'field') {
                doc.text('Field Content', columnXPos, columnYPos);
              } else if (content.type === 'table') {
                doc.autoTable({
                  startY: columnYPos,
                  head: [['Column 1', 'Column 2']],
                  body: [
                    ['Data 1', 'Data 2'],
                    ['Data 3', 'Data 4']
                  ],
                });
              }
              columnYPos += 40; // Increment Y position after each content
            });
          });
          yPos += 100; // Increment Y position after each row
        });
  
        doc.save('report.pdf');
      },
    },
  };
  </script>
  
  <style>
  .report-canvas {
    padding: 20px;
    background-color: #f4f4f4;
  }
  
  .controls {
    margin-bottom: 20px;
  }
  
  .a4-page {
    width: 794px; /* A4 width */
    height: 1123px; /* A4 height */
    background-color: white;
    border: 1px solid #ccc;
    padding: 20px;
    box-sizing: border-box;
    display: flex;
    flex-direction: column;
  }
  
  .row {
    margin-bottom: 20px;
    border: 1px solid #ddd;
    padding: 10px;
    background-color: #e0e0e0;
    border-radius: 5px;
  }
  
  .row-controls {
    margin-bottom: 10px;
  }
  
  .columns-container {
    display: flex;
    gap: 10px;
    flex-wrap: wrap;
  }
  
  .column-item {
    background-color: #dcdcdc;
    padding: 10px;
    border-radius: 5px;
    border: 1px solid #ccc;
    flex-grow: 1;
  }
  
  .column-controls {
    margin-bottom: 10px;
  }
  
  .content-container {
    display: flex;
    gap: 10px;
    flex-wrap: wrap;
  }
  
  .content-item {
    background-color: #f0f0f0;
    padding: 10px;
    border-radius: 5px;
    border: 1px solid #bbb;
    flex-grow: 1;
  }
  
  .content-controls {
    margin-bottom: 10px;
  }
  </style>
  