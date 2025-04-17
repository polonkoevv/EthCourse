<template>
    <section class="mt-8 bg-white rounded-2xl p-6 shadow-sm">
      <h2 class="text-xl font-bold mb-6">Библиотека</h2>
      <div v-if="audioFiles.length === 0" class="text-center py-8 text-gray-500">
        Библиотека пуста. Загрузите свой первый трек, нажав на кнопку "Загрузить аудио" выше.
      </div>
      <div v-else class="space-y-4">
        <AudioItem 
          v-for="(audio, index) in audioFiles"
          :key="index"
          :audio="audio"
          @select="selectAudio"
        />
      </div>
    </section>
  </template>
  
  <script setup lang="ts">
  import { ref, onMounted, watch } from 'vue';
  import AudioItem from './AudioItem.vue';
  import type { Audio } from '../../types/audio';
  import AxiosEntity from '../../scripts/axios';

  const props = defineProps<{
    audioFiles?: Audio[]
  }>();

  const emit = defineEmits<{
    (e: 'play-audio', audio: Audio): void,
    (e: 'load-audio-list', list: Audio[]): void
  }>();

  const audioFiles = ref<Audio[]>([]);

  

  // Установка демо-треков при отсутствии внешнего списка
  

  const selectAudio = (audio: Audio): void => {
    emit('play-audio', audio);
  };
  const defaultAudioFiles = audioFiles;



  onMounted(async () => {
    const response = await AxiosEntity.GetAllMusic()
    console.log(response.data);
    audioFiles.value = response.data.map((audio: Audio) => ({
      id: audio.id,
      title: audio.title,
      cid: audio.cid,
      gradient: audio.gradient,
      link: audio.link
    }));
    
    if (!props.audioFiles?.length) {
      emit('load-audio-list', defaultAudioFiles.value);
    }
    emit('load-audio-list', audioFiles.value);
  });


  // Следим за внешним списком, чтобы инициализировать при его появлении
  watch(() => props.audioFiles, (newFiles) => {
    if (newFiles?.length) {
      // Внешний список имеет приоритет над демо-треками
      // Ничего не делаем, т.к. внешние файлы уже установлены
    } 
    else if (defaultAudioFiles.value.length) {
      // Если внешних треков нет, но есть демо, используем их
      emit('load-audio-list', defaultAudioFiles.value);
    }
  }, { immediate: true });
</script>