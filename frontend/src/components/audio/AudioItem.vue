<template>
  <div 
    class="flex items-center justify-between p-4 hover:bg-gray-50 dark:hover:bg-gray-700 rounded-lg transition-colors cursor-pointer" 
    @click="selectAudio"
  >
    <div class="flex items-center space-x-4">
      <div :class="['w-12', 'h-12', 'bg-gradient-to-br', 'from-blue-400', 'to-cyan-300', 'rounded', 'flex', 'items-center', 'justify-center']">
        <font-awesome-icon icon="music" class="text-white" />
      </div>
      <div>
        <p class="font-medium dark:text-white">{{ audio.title }}</p>
        <p class="text-sm text-gray-500 dark:text-gray-400">{{ audio.cid }}</p>
      </div>
    </div>
    <div class="flex items-center space-x-4">
      <button 
        class="text-gray-600 dark:text-gray-400 hover:text-blue-600 dark:hover:text-blue-400" 
        @click.stop="shareAudio"
        title="Скопировать ссылку"
      >
        <font-awesome-icon icon="share-nodes" />
        <!-- Показать уведомление при копировании -->
        <span 
          v-if="showCopiedMessage" 
          class="absolute mt-6 -ml-12 bg-gray-800 text-white text-xs px-2 py-1 rounded"
        >
          Скопировано!
        </span>
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import type { Audio } from '../../types/audio';

const props = defineProps<{
  audio: Audio
}>();

const emit = defineEmits<{
  (e: 'select', audio: Audio): void
}>();

const showCopiedMessage = ref(false);

const selectAudio = () => {
  emit('select', props.audio);
};

const shareAudio = async () => {
  // Формируем ссылку на аудио трек
  const audioLink = props.audio.link;
  
  try {
    // Копируем ссылку в буфер обмена
    await navigator.clipboard.writeText(audioLink);
    
    // Показываем уведомление
    showCopiedMessage.value = true;
    
    // Скрываем уведомление через 2 секунды
    setTimeout(() => {
      showCopiedMessage.value = false;
    }, 2000);
  } catch (err) {
    console.error('Ошибка при копировании ссылки:', err);
    alert('Не удалось скопировать ссылку. Пожалуйста, попробуйте снова.');
  }
};
</script>

<style scoped>
/* Позиционирование для всплывающего уведомления */
button {
  position: relative;
}
</style>