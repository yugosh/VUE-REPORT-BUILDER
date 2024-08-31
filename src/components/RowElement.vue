<template>
    <div class="row-element">
      <div class="controls">
        <button @click="addRowElement('text')">Add Text</button>
        <button @click="addRowElement('image')">Add Image</button>
        <button @click="addRowElement('table')">Add Table</button>
        <button @click="addRowElement('row')">Add Inner Row</button> <!-- Option to add row within a row -->
      </div>
      <div class="row-content">
        <draggable v-model="localRowContent" group="report-elements" class="draggable-row" @end="updateRowContent">
          <template #item="{ element }">
            <div v-if="element.type === 'row'" class="inner-row">
              <RowElement :row="element" />
            </div>
            <div v-else class="element" :style="element.css">
              <span class="move-handle">::</span>
              <component 
                :is="getElementComponent(element)" 
                :element="element" 
                @update-content="updateElementContent(element, $event)"
              />
            </div>
          </template>
        </draggable>
      </div>
    </div>
  </template>
  
  <script>
  import draggable from 'vuedraggable';
  import TextComponent from './TextComponent.vue';
  import ImageComponent from './ImageComponent.vue';
  import TableComponent from './TableComponent.vue';
  import RowElement from './RowElement.vue'; // Impor RowElement

  export default {
    name: 'RowElement',
    components: {
      draggable,
      TextComponent,
      ImageComponent,
      TableComponent,
      RowElement,
    },
    props: {
      row: Object,
    },
    data() {
      return {
        localRowContent: [...this.row.content],
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
            return 'RowElement';
          default:
            return null;
        }
      },
      addRowElement(type) {
        const newElement = {
          id: this.generateUUID(),
          type: type,
          content: '',
          css: {
            margin: '10px',
            padding: '5px',
            border: '1px solid #ddd',
          },
        };
  
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
          case 'row':
            newElement.content = [];
            newElement.css = { display: 'flex', 'flex-direction': 'row' };
            break;
        }
  
        this.localRowContent.push(newElement);
        this.updateRowContent();
      },
      updateRowContent() {
        this.$emit('update:row', { ...this.row, content: this.localRowContent });
      },
      updateElementContent(element, newContent) {
        element.content = newContent;
        this.updateRowContent();
      },
      generateUUID() {
        return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function(c) {
          var r = (Math.random() * 16) | 0,
            v = c == 'x' ? r : (r & 0x3) | 0x8;
          return v.toString(16);
        });
      },
    },
    watch: {
      row: {
        handler(newVal) {
          this.localRowContent = [...newVal.content];
        },
        deep: true,
      },
    },
  };
  </script>
  
  <style scoped>
  .row-element {
    display: flex;
    flex-direction: column;
    gap: 10px;
    border: 1px dashed #ddd;
    padding: 5px;
    margin-bottom: 10px;
  }
  
  .row-content {
    display: flex;
    flex-direction: row;
    gap: 10px;
  }
  
  .inner-row {
    display: flex;
    flex-direction: row;
    gap: 10px;
  }
  
  .element {
    cursor: move;
  }
  
  .controls {
    margin-bottom: 10px;
  }
  
  .controls button {
    margin-right: 5px;
    padding: 5px 10px;
    font-size: 12px;
    cursor: pointer;
  }
  
  .move-handle {
    cursor: move;
    display: inline-block;
    margin-right: 5px;
  }
  </style>
  