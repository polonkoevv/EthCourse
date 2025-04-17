<template>
  <div class="min-h-screen flex flex-col transition-colors duration-300" :class="[
    isDarkMode ? 'dark bg-gray-900' : 'bg-gray-50',
  ]">
    <AppHeader 
      @toggle-theme="toggleTheme" 
      :isDarkMode="isDarkMode" 
      @wallet-connected="onWalletConnected"
      @wallet-disconnected="onWalletDisconnected"
    />
    
    <main class="pt-16 flex-grow">
      <div class="container mx-auto px-4 py-8">
        <UploadSection @add-audio="addNewAudio" />
        
        <section class="grid grid-cols-1 md:grid-cols-4 gap-6">
          <AudioPlayer 
            class="md:col-span-2 max-h-100" 
            :selectedAudio="selectedAudio"
            :audioList="audioList" 
            @change-track="changeTrack"
          />
          <TransactionsPanel class="md:col-span-2 max-h-100" />
        </section>
        
        <div class="container mx-auto mt-6">
          <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div>
              <AudioLibrary 
                :audioFiles="audioList" 
                @play-audio="playAudio" 
                @load-audio-list="loadAudioList"
              />
            </div>
            <div>
              <UserAudioLibrary 
                :audioFiles="audioList" 
                :walletAddress="walletAddress"
                :isWalletConnected="isWalletConnected"
                @play-audio="playAudio"
              />
            </div>
          </div>
        </div>
      </div>
    </main>
    
    <AppFooter />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue';
import AppHeader from './components/layout/AppHeader.vue';
import AppFooter from './components/layout/AppFooter.vue';
import UploadSection from './components/upload/UploadSection.vue';
import AudioPlayer from './components/audio/AudioPlayer.vue';
import TransactionsPanel from './components/transactions/TransactionsPanel.vue';
import AudioLibrary from './components/audio/AudioLibrary.vue';
import UserAudioLibrary from './components/audio/UserAudioLibrary.vue';
import type { Audio } from './types/audio';

const isDarkMode = ref(false);
const selectedAudio = ref<Audio | undefined>(undefined);
const audioList = ref<Audio[]>([]);
const currentIndex = ref(-1);
const walletAddress = ref('');
const isWalletConnected = ref(false);

// Загрузка сохраненной темы из localStorage при старте
onMounted(() => {
  const savedTheme = localStorage.getItem('darkMode');
  if (savedTheme === 'true') {
    isDarkMode.value = true;
    document.documentElement.classList.add('dark');
  }
});

// Применяет класс к html элементу при изменении темы
watch(() => isDarkMode.value, (newValue) => {
  if (newValue) {
    document.documentElement.classList.add('dark');
  } else {
    document.documentElement.classList.remove('dark');
  }
});

const toggleTheme = (): void => {
  isDarkMode.value = !isDarkMode.value;
  // Сохраняем выбор темы в localStorage
  localStorage.setItem('darkMode', isDarkMode.value.toString());
};

const playAudio = (audio: Audio): void => {
  selectedAudio.value = audio;
  
  // Обновить текущий индекс
  currentIndex.value = audioList.value.findIndex(
    a => a.title === audio.title && a.link === audio.link
  );
};

const loadAudioList = (list: Audio[]): void => {
  audioList.value = list;
};

const changeTrack = (index: number): void => {
  if (index >= 0 && index < audioList.value.length) {
    selectedAudio.value = audioList.value[index];
    currentIndex.value = index;
  }
};

const addNewAudio = (audio: Audio): void => {
  // Добавляем новый трек в начало списка
  audioList.value = [audio, ...audioList.value];
  
  // Автоматически выбираем новый трек для воспроизведения
  selectedAudio.value = audio;
  currentIndex.value = 0;
};

const onWalletConnected = (address: string): void => {
  isWalletConnected.value = true;
  walletAddress.value = address;
  console.log('Кошелек подключен:', address);
};

const onWalletDisconnected = (): void => {
  isWalletConnected.value = false;
  walletAddress.value = '';
  console.log('Кошелек отключен');
};
</script>

<style>
/* Добавляем стили для компонентов, чтобы контролировать переполнение */
.max-h-96 {
  max-height: 24rem; /* 384px */
  overflow: hidden;
}

/* Добавляем общие стили для темной темы */
.dark {
  color-scheme: dark;
}

.dark .bg-white {
  background-color: #1f2937;
  color: #f3f4f6;
}

.dark .text-gray-600 {
  color: #d1d5db;
}

.dark .text-gray-500 {
  color: #9ca3af;
}

.dark .border-gray-100 {
  border-color: #374151;
}

.dark .bg-gray-200 {
  background-color: #4b5563;
}

.dark .bg-gray-50 {
  background-color: #374151;
}

.dark .shadow-sm {
  box-shadow: 0 1px 2px 0 rgba(0, 0, 0, 0.3);
}

.dark .text-blue-600 {
  color: #60a5fa;
}

.dark .hover\:text-blue-600:hover {
  color: #93c5fd;
}

.dark .hover\:bg-gray-50:hover {
  background-color: #2d3748;
}

.dark .hover\:bg-blue-700:hover {
  background-color: #1d4ed8;
}

.dark .text-neutral-500 {
  color: #9ca3af;
}

/* Добавляем плавную анимацию при переключении темы */
.transition-colors {
  transition-property: background-color, border-color, color, fill, stroke;
  transition-timing-function: cubic-bezier(0.4, 0, 0.2, 1);
  transition-duration: 300ms;
}
</style>