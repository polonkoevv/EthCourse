<template>
  <div class="min-h-screen bg-gray-50" :class="{ 'dark': isDarkMode }">
    <AppHeader 
      @toggle-theme="toggleTheme" 
      :isDarkMode="isDarkMode" 
      @wallet-connected="onWalletConnected"
      @wallet-disconnected="onWalletDisconnected"
    />
    
    <main class="pt-16">
      <div class="container mx-auto px-4 py-8">
        <UploadSection @add-audio="addNewAudio" />
        
        <section class="grid grid-cols-1 md:grid-cols-3 gap-6">
          <AudioPlayer 
            class="md:col-span-2" 
            :selectedAudio="selectedAudio"
            :audioList="audioList" 
            @change-track="changeTrack"
          />
          <TransactionsPanel />
        </section>
        
        <AudioLibrary 
          :audio-files="audioList"
          @play-audio="playAudio"
          @load-audio-list="loadAudioList" 
        />
      </div>
    </main>
    
    <AppFooter />
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import AppHeader from './components/layout/AppHeader.vue';
import AppFooter from './components/layout/AppFooter.vue';
import UploadSection from './components/upload/UploadSection.vue';
import AudioPlayer from './components/audio/AudioPlayer.vue';
import TransactionsPanel from './components/transactions/TransactionsPanel.vue';
import AudioLibrary from './components/audio/AudioLibrary.vue';
import type { Audio } from './types/audio';

const isDarkMode = ref(false);
const selectedAudio = ref<Audio | undefined>(undefined);
const audioList = ref<Audio[]>([]);
const currentIndex = ref(-1);
const walletAddress = ref('');

const toggleTheme = (): void => {
  isDarkMode.value = !isDarkMode.value;
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
  walletAddress.value = address;
  console.log('Кошелек подключен:', address);
};

const onWalletDisconnected = (): void => {
  walletAddress.value = '';
  console.log('Кошелек отключен');
};
</script>