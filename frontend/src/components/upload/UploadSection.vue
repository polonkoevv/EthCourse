<template>
  <section class="bg-gradient-to-r from-blue-500 to-cyan-500 rounded-2xl p-8 mb-8">
    <div class="flex flex-col items-center text-white">
      <h1 class="text-3xl font-bold mb-4">Загрузите аудио в IPFS</h1>
      <p class="text-blue-50 mb-6">Безопасное, децентрализованное хранилище для ваших аудиофайлов</p>
      <button 
        class="bg-white text-blue-600 px-6 py-3 rounded-lg hover:bg-blue-50 transition-colors"
        @click="openUploadModal"
      >
        <font-awesome-icon icon="cloud-arrow-up" class="mr-2" />Загрузить аудио
      </button>
    </div>
    
    <!-- Модальное окно загрузки -->
    <UploadModal 
      :is-open="isModalOpen" 
      @close="isModalOpen = false"
      @upload-success="onUploadSuccess"
    />
  </section>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import UploadModal from './UploadModal.vue';
import type { Audio } from '../../types/audio';

const isModalOpen = ref(false);
const emit = defineEmits<{
  (e: 'add-audio', audio: Audio): void
}>();

const openUploadModal = () => {
  isModalOpen.value = true;
};

const onUploadSuccess = (audio: Audio) => {
  emit('add-audio', audio);
};
</script>