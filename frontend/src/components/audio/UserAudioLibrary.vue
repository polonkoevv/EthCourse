<template>
  <section class="mt-8 bg-white rounded-2xl p-6 shadow-sm">
    <h2 class="text-xl font-bold mb-6">Ваши загрузки</h2>
    <div v-if="!isWalletConnected" class="text-center py-8 text-gray-500">
      Подключите кошелек, чтобы увидеть ваши загруженные треки.
    </div>
    <div v-else-if="userAudioFiles.length === 0" class="text-center py-8 text-gray-500">
      У вас пока нет загруженных треков. Загрузите свой первый трек, нажав на кнопку "Загрузить аудио" выше.
    </div>
    <div v-else class="space-y-4">
      <AudioItem 
        v-for="(audio, index) in userAudioFiles"
        :key="index"
        :audio="audio"
        @select="selectAudio"
      />
    </div>
  </section>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue';
import AudioItem from './AudioItem.vue';
import type { Audio } from '../../types/audio';
import AxiosEntity from '../../scripts/axios';

const props = defineProps<{
  audioFiles?: Audio[],
  walletAddress?: string,
  isWalletConnected?: boolean,
}>();

const emit = defineEmits<{
  (e: 'play-audio', audio: Audio): void
}>();

const audioFiles = ref<Audio[]>([]);
const isWalletConnected = ref(false);
const walletAddress = ref('');

// Отфильтрованные аудиофайлы, принадлежащие пользователю
const userAudioFiles = computed(() => {
  return audioFiles.value.filter(audio => {
    // Сравниваем адреса в нижнем регистре для игнорирования различий в написании (0x...)
    return audio.owner_addr && walletAddress.value && 
           audio.owner_addr.toLowerCase() === walletAddress.value.toLowerCase();
  });
});

// Установка данных из props
watch(() => props.audioFiles, (newFiles) => {
  if (newFiles?.length) {
    audioFiles.value = newFiles;
  }
}, { immediate: true });

watch(() => props.walletAddress, (newAddress) => {
  if (newAddress) {
    walletAddress.value = newAddress;
  }
}, { immediate: true });

watch(() => props.isWalletConnected, (newConnected) => {
  isWalletConnected.value = newConnected || false;
}, { immediate: true });

// Функция воспроизведения аудио
const selectAudio = (audio: Audio): void => {
  emit('play-audio', audio);
};

// Загрузка данных при монтировании компонента
onMounted(async () => {
  if (!props.audioFiles?.length) {
    try {
      const response = await AxiosEntity.GetAllMusic();
      audioFiles.value = response.data.map((audio: any) => ({
        id: audio.id,
        title: audio.title,
        cid: audio.cid,
        link: audio.link,
        owner_addr: audio.owner_addr,
        signature: audio.signature,
        uploaded_at: audio.uploaded_at
      }));
    } catch (error) {
      console.error('Ошибка при загрузке аудиофайлов:', error);
    }
  }
});
</script> 