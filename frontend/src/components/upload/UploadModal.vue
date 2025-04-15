<template>
  <Teleport to="body">
    <div v-if="isOpen" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-[60]" @click.self="closeModal">
      <div class="bg-white rounded-xl w-full max-w-lg mx-4">
        <div class="p-6 border-b border-neutral-100">
          <div class="flex justify-between items-center">
            <h3 class="text-xl">Загрузить аудио</h3>
            <button class="text-neutral-400 hover:text-neutral-600" @click="closeModal">
              <font-awesome-icon icon="xmark" class="text-xl" />
            </button>
          </div>
        </div>

        <div class="p-6">
          <!-- Drag & Drop Zone -->
          <div 
            class="border-2 border-dashed border-neutral-200 rounded-xl p-8 text-center"
            :class="{ 'border-blue-500 bg-blue-50': isDragging }"
            @dragover.prevent="isDragging = true"
            @dragleave.prevent="isDragging = false"
            @drop.prevent="onDrop"
          >
            <div class="space-y-4">
              <div class="w-16 h-16 bg-neutral-100 rounded-full mx-auto flex items-center justify-center">
                <font-awesome-icon icon="cloud-arrow-up" class="text-2xl text-neutral-400" />
              </div>
              <div class="space-y-2">
                <h4>Перетащите аудио файлы сюда</h4>
                <p class="text-sm text-neutral-500">Поддерживаются форматы MP3, WAV, FLAC</p>
              </div>
              <div class="relative">
                <button class="bg-neutral-800 text-white px-6 py-2 rounded-lg hover:bg-neutral-700">
                  Выбрать файлы
                </button>
                <input 
                  type="file" 
                  class="absolute inset-0 w-full h-full opacity-0 cursor-pointer" 
                  accept="audio/*"
                  @change="onFileSelected"
                >
              </div>
            </div>
          </div>

          <!-- Upload Progress -->
          <div v-if="uploadProgress > 0" class="mt-6 space-y-4">
            <div class="flex items-center space-x-4">
              <font-awesome-icon icon="music" class="text-neutral-400" />
              <div class="flex-1">
                <div class="flex justify-between text-sm mb-1">
                  <span>{{ selectedFile?.name || 'Загрузка файла...' }}</span>
                  <span>{{ uploadProgress }}%</span>
                </div>
                <div class="h-2 bg-neutral-100 rounded-full">
                  <div 
                    class="h-full bg-neutral-800 rounded-full" 
                    :style="{ width: `${uploadProgress}%` }"
                  ></div>
                </div>
              </div>
            </div>
          </div>

          <!-- Form Fields -->
          <div v-if="selectedFile" class="mt-6 space-y-4">
            <div>
              <label class="block text-sm text-neutral-700 mb-1">Название трека</label>
              <input 
                v-model="trackTitle" 
                type="text" 
                class="w-full px-3 py-2 border border-neutral-200 rounded-lg focus:outline-none focus:border-neutral-400"
              >
            </div>
          </div>
        </div>

        <div class="p-6 border-t border-neutral-100 flex justify-end space-x-4">
          <button 
            class="px-4 py-2 text-neutral-600 hover:text-neutral-800"
            @click="closeModal"
          >
            Отмена
          </button>
          <button 
            class="px-6 py-2 bg-neutral-800 text-white rounded-lg hover:bg-neutral-700"
            :disabled="!canUpload"
            :class="{ 'opacity-50 cursor-not-allowed': !canUpload }"
            @click="uploadFile"
          >
            Загрузить
          </button>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue';
import type { Audio } from '../../types/audio';
import AxiosEntity from '../../scripts/axios';
const props = defineProps<{
  isOpen: boolean
}>();

const emit = defineEmits<{
  (e: 'close'): void,
  (e: 'upload-success', audio: Audio): void
}>();

// Состояние
const isDragging = ref(false);
const selectedFile = ref<File | null>(null);
const uploadProgress = ref(0);
const trackTitle = ref('');
const artistName = ref('');
const description = ref('');

// Вычисляемые свойства
const canUpload = computed(() => {
  return selectedFile.value && trackTitle.value.trim() !== '';
});

// Методы
const closeModal = () => {
  emit('close');
  resetForm();
};

const resetForm = () => {
  selectedFile.value = null;
  uploadProgress.value = 0;
  trackTitle.value = '';
  artistName.value = '';
  description.value = '';
  isDragging.value = false;
};

const onDrop = (event: DragEvent) => {
  isDragging.value = false;
  
  if (!event.dataTransfer?.files.length) return;
  
  const file = event.dataTransfer.files[0];
  if (file.type.startsWith('audio/')) {
    selectedFile.value = file;
    
    // Автоматически устанавливаем название из имени файла
    if (!trackTitle.value) {
      const filename = file.name.split('.').slice(0, -1).join('.');
      trackTitle.value = filename;
    }
  } else {
    alert('Пожалуйста, выберите аудиофайл');
  }
};

const onFileSelected = (event: Event) => {
  const target = event.target as HTMLInputElement;
  if (!target.files?.length) return;
  
  const file = target.files[0];
  if (file.type.startsWith('audio/')) {
    selectedFile.value = file;
    
    // Автоматически устанавливаем название из имени файла
    if (!trackTitle.value) {
      const filename = file.name.split('.').slice(0, -1).join('.');
      trackTitle.value = filename;
    }
  } else {
    alert('Пожалуйста, выберите аудиофайл');
  }
};

const uploadFile = async () => {
  if (!selectedFile.value || !canUpload.value) return;
  
  try {
    // Имитация загрузки с прогрессом
    for (let i = 0; i <= 100; i += 5) {
      uploadProgress.value = i;
      await new Promise(resolve => setTimeout(resolve, 100));
    }
    
    // В реальном приложении здесь был бы код загрузки на сервер
    // const formData = new FormData();
    // formData.append('file', selectedFile.value);
    // formData.append('title', trackTitle.value);
    // formData.append('artist', artistName.value);
    // formData.append('description', description.value);
    // const response = await axios.post('/api/upload', formData);

    const response = await AxiosEntity.UploadMusic(selectedFile.value, trackTitle.value);
    console.log(response.data);

    // Создаем объект URL для аудиофайла (в реальном приложении это был бы URL с сервера)
    const audioURL = URL.createObjectURL(selectedFile.value);
    
    // Создаем новый объект аудио
    const newAudio: Audio = {
      title: trackTitle.value,
      gradient: 'from-green-400 to-emerald-300',
      link: audioURL,
      id: 0,
      cid: ""
    };
    
    // Отправляем событие об успешной загрузке
    emit('upload-success', newAudio);
    
    // Закрываем модальное окно
    closeModal();
  } catch (error) {
    console.error('Ошибка загрузки:', error);
    alert('Произошла ошибка при загрузке файла');
  }
};

// Сброс формы при закрытии
watch(() => props.isOpen, (newVal) => {
  if (!newVal) {
    resetForm();
  }
});
</script> 