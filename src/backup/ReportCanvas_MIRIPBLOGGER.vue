<template>
    <div class="report-canvas">
      <!-- Tombol untuk Menambahkan Elemen -->
      <div class="controls">
        <button @click="addElement('text')">Add Text</button>
        <button @click="addElement('image')">Add Image</button>
        <button @click="addElement('table')">Add Table</button>
        <button @click="addElement('row')">Add Row</button>
      </div>
  
      <!-- Header Area -->
      <div class="header-area">
        <h3>Header</h3>
        <draggable v-model="reportData.header.elements" group="report-elements" class="draggable-header">
          <template #item="{ element }">
            <RowElement v-if="element.type === 'row'" :row="element" />
            <div v-else class="element" :style="element.css">
              <span class="move-handle">::</span>
              <component :is="getElementComponent(element)" :element="element" />
            </div>
          </template>
        </draggable>
      </div>
  
      <!-- Body Area -->
      <div class="body-area">
        <h3>Body</h3>
        <draggable v-model="reportData.body.elements" group="report-elements" class="draggable-body">
          <template #item="{ element }">
            <RowElement v-if="element.type === 'row'" :row="element" />
            <div v-else class="element" :style="element.css">
              <span class="move-handle">::</span>
              <component :is="getElementComponent(element)" :element="element" />
            </div>
          </template>
        </draggable>
      </div>
  
      <!-- Footer Area -->
      <div class="footer-area">
        <h3>Footer</h3>
        <draggable v-model="reportData.footer.elements" group="report-elements" class="draggable-footer">
          <template #item="{ element }">
            <RowElement v-if="element.type === 'row'" :row="element" />
            <div v-else class="element" :style="element.css">
              <span class="move-handle">::</span>
              <component :is="getElementComponent(element)" :element="element" />
            </div>
          </template>
        </draggable>
      </div>
    </div>
  </template>
  
  <script>
  import draggable from 'vuedraggable';
  import RowElement from './RowElement.vue'; // Import komponen RowElement
  import TextComponent from './TextComponent.vue';
import ImageComponent from './ImageComponent.vue';
import TableComponent from './TableComponent.vue';

  export default {
    name: 'ReportCanvas',
    components: {
      draggable,
      RowElement, // Tambahkan RowElement sebagai komponen anak
      TextComponent,
    ImageComponent,
    TableComponent,
    },
    data() {
      return {
        reportData: {
          id: this.generateUUID(),
          product_name: 'Product ABC',
          created_time: new Date().toISOString(),
          last_update: new Date().toISOString(),
          header: {
            elements: [], // Elemen untuk header
          },
          body: {
            elements: [], // Elemen untuk body
          },
          footer: {
            elements: [], // Elemen untuk footer
          },
        },
      };
    },
    methods: {
      getElementComponent(element) {
        switch (element.type) {
          case 'text':
            return 'TextComponent';
          case 'image':
            return 'ImageComponent';
          case 'table':
            return 'TableComponent';
          case 'row':
            return null; // Row ditangani secara khusus
          default:
            return null;
        }
      },
      generateUUID() {
        return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function(c) {
          var r = (Math.random() * 16) | 0,
            v = c == 'x' ? r : (r & 0x3) | 0x8;
          return v.toString(16);
        });
      },
      addElement(type) {
        const newElement = {
          id: this.generateUUID(),
          type: type,
          content: '',
          css: {
            position: 'relative',
            margin: '10px',
            padding: '5px',
            border: '1px solid #ddd',
          },
        };
  
        if (type === 'row') {
          newElement.content = []; // Row memiliki array content
          newElement.css = { ...newElement.css, display: 'flex', 'align-items': 'center' };
        } else {
          switch (type) {
            case 'text':
              newElement.content = 'New Text';
              newElement.css = { ...newElement.css, 'font-size': '16px', 'color': '#000' };
              break;
            case 'image':
              newElement.content = 'https://via.placeholder.com/150';
              newElement.css = { ...newElement.css, 'width': '150px', 'height': '150px' };
              break;
            case 'table':
              newElement.content = {
                headers: ['Column 1', 'Column 2'],
                rows: [['Data 1', 'Data 2'], ['Data 3', 'Data 4']],
              };
              newElement.css = { ...newElement.css, 'width': '100%', 'border-collapse': 'collapse' };
              break;
          }
        }
  
        // Default: Menambah elemen ke dalam body
        this.reportData.body.elements.push(newElement);
        this.updateLastModified();
      },
      updateLastModified() {
        this.reportData.last_update = new Date().toISOString();
      },
    },
  };
  </script>
  
  <style scoped>
  .report-canvas {
    width: 210mm;
    height: 297mm;
    padding: 20px;
  }
  
  .controls {
    margin-bottom: 20px;
  }
  
  .controls button {
    margin-right: 10px;
    padding: 5px 10px;
    font-size: 14px;
    cursor: pointer;
  }
  
  .header-area, .body-area, .footer-area {
    margin-bottom: 20px;
    padding: 10px;
    border: 1px dashed #ccc;
  }
  
  .element {
    cursor: move;
    margin-bottom: 10px;
  }
  
  .row-element {
    display: flex;
    flex-direction: column;
    gap: 10px;
    border: 1px dashed #ddd;
    padding: 5px;
    margin-bottom: 10px;
  }
  
  .move-handle {
    cursor: move;
    display: inline-block;
    margin-right: 5px;
  }
  </style>
  