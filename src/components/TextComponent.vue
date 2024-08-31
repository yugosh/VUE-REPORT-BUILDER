<template>
    <div :style="element.css">
      <input
        v-if="isEditing"
        v-model="localContent"
        @blur="finishEditing"
        @keyup.enter="finishEditing"
        @keyup.esc="cancelEditing"
      />
      <span v-else @dblclick="editContent">{{ localContent }}</span>
    </div>
  </template>
  
  <script>
  export default {
    name: 'TextComponent',
    props: {
      element: Object,
    },
    data() {
      return {
        isEditing: false,
        localContent: this.element.content, // Membuat salinan dari konten untuk diedit
      };
    },
    methods: {
      editContent() {
        this.isEditing = true;
        this.$nextTick(() => {
          this.$el.querySelector('input').focus();
        });
      },
      finishEditing() {
        this.isEditing = false;
        this.$emit('update-content', this.localContent); // Emit perubahan ke parent component
      },
      cancelEditing() {
        this.isEditing = false;
        this.localContent = this.element.content; // Kembalikan ke konten asli jika dibatalkan
      }
    },
  };
  </script>
  
  <style scoped>
  input {
    width: 100%;
    border: none;
    background: transparent;
    outline: none;
  }
  </style>
  