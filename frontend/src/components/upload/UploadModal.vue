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

          <!-- Статус публикации -->
          <div v-if="publishingStatus" class="mt-4">
            <div v-if="publishingStatus === 'pending'" class="text-blue-600">
              <font-awesome-icon icon="spinner" spin class="mr-2" />
              {{ publishingStatusMessage }}
            </div>
            <div v-else-if="publishingStatus === 'success'" class="text-green-600">
              <font-awesome-icon icon="check-circle" class="mr-2" />
              {{ publishingStatusMessage }}
            </div>
            <div v-else-if="publishingStatus === 'error'" class="text-red-600">
              <font-awesome-icon icon="exclamation-circle" class="mr-2" />
              {{ publishingStatusMessage }}
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

declare global {
  interface Window {
    ethereum: any;
  }
}

import { ref, computed, watch } from 'vue';
import type { Audio } from '../../types/audio';
import AxiosEntity from '../../scripts/axios';
import { ethers } from 'ethers';
import AudioChainABI from '../../../build/contracts/AudioChain.json';

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
const price = ref('');

// Добавляем переменные для взаимодействия с Web3
const contractAddress = ref(''); // Адрес смарт-контракта в Ganache
const provider = ref(null);
const signer = ref(null);
const contract = ref(null);
const transactionHash = ref('');
const publishingStatus = ref('');
const publishingStatusMessage = ref('');

const isWalletModalOpen = ref(false);
const isWalletConnected = ref(false);
const walletAddress = ref('');

const GetCurrentWallet = async () => {
  try {
    // Проверяем, доступен ли объект ethereum в окне браузера
    if (window.ethereum) {
      const provider = new ethers.providers.Web3Provider(window.ethereum);
      
      // Запрашиваем доступ к аккаунтам
      const accounts = await provider.send('eth_requestAccounts', []);
      
      if (accounts.length > 0) {
        const address = accounts[0];
        console.log(accounts);
        isWalletConnected.value = true;
        walletAddress.value = address;
        return address;
      } else {
        isWalletConnected.value = false;
        walletAddress.value = '';
        return null;
      }
    } else {
      console.log('MetaMask не установлен');
      return null;
    }
  } catch (error) {
    console.error('Ошибка при получении кошелька:', error);
    return null;
  }
}

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
  price.value = '';
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

// Инициализация Web3 и контракта
const initWeb3 = async () => {
  try {
    // Проверяем наличие MetaMask
    if (window.ethereum) {
      provider.value = new ethers.providers.Web3Provider(window.ethereum);
      await provider.value.send("eth_requestAccounts", []);
      signer.value = provider.value.getSigner();
      
      // Инициализация контракта
      contract.value = new ethers.Contract(
        contractAddress.value,
        AudioChainABI.abi,
        signer.value
      );
      
      return true;
    } else {
      console.error('MetaMask не установлен');
      return false;
    }
  } catch (error) {
    console.error('Ошибка инициализации Web3:', error);
    return false;
  }
};

// Модифицируем функцию загрузки файла
const uploadFile = async () => {
  if (!selectedFile.value || !canUpload.value) return;
  
  try {
    // Получаем адрес текущего кошелька
    const walletAddress = await GetCurrentWallet();
    if (!walletAddress) {
      publishingStatus.value = 'error';
      publishingStatusMessage.value = 'Не удалось подключиться к кошельку';
      return;
    }
    
    // Создаем сообщение для подписи
    const messageToSign = JSON.stringify({
      action: 'audio_upload',
      title: trackTitle.value,
      artist: artistName.value || 'Unknown',
      filename: selectedFile.value.name,
      filesize: selectedFile.value.size,
      timestamp: Date.now(),
      wallet: walletAddress
    });
    
    // Подписываем сообщение
    publishingStatus.value = 'pending';
    publishingStatusMessage.value = 'Подписываем сообщение...';
    
    const provider = new ethers.providers.Web3Provider(window.ethereum);
    const signer = provider.getSigner();
    const signature = await signer.signMessage(messageToSign);
    
    console.log('Подписанное сообщение:', signature);
    
    // Загрузка файла в IPFS
    uploadProgress.value = 0;
    publishingStatus.value = 'pending';
    publishingStatusMessage.value = 'Загрузка файла...';
    
    // Реальная загрузка файла на сервер с подписью
    const formData = new FormData();
    formData.append('file', selectedFile.value);
    formData.append('title', trackTitle.value);
    formData.append('artist', artistName.value || 'Unknown');
    formData.append('message', messageToSign);
    formData.append('signature', signature);
    formData.append('walletAddress', walletAddress);
    
    // Имитация прогресса загрузки
    const uploadInterval = setInterval(() => {
      if (uploadProgress.value < 95) {
        uploadProgress.value += 5;
      }
    }, 200);
    
    try {
      const response = await fetch('http://localhost:8000/upload', {
        method: 'POST',
        body: formData
      });
      
      clearInterval(uploadInterval);
      uploadProgress.value = 100;
      
      const data = await response.json();
      
      if (!data.success) {
        throw new Error(data.message || 'Ошибка загрузки файла');
      }
      
      const ipfsHash = data.cid;
      
      // Успешная загрузка
      publishingStatus.value = 'success';
      publishingStatusMessage.value = 'Файл успешно загружен!';
      
      // Отправляем событие об успешной загрузке
      emit('upload-success', {
        id: data.audioId || Date.now(),
        title: trackTitle.value,
        artist: artistName.value || 'Unknown',
        ipfsHash,
        signature,
        walletAddress
      });
      
      // Закрываем модальное окно
      setTimeout(() => {
        closeModal();
      }, 2000);
      
    } catch (error) {
      clearInterval(uploadInterval);
      console.error('Ошибка загрузки:', error);
      publishingStatus.value = 'error';
      publishingStatusMessage.value = `Ошибка загрузки: ${error.message}`;
    }
    
  } catch (error) {
    console.error('Ошибка:', error);
    publishingStatus.value = 'error';
    publishingStatusMessage.value = `Ошибка: ${error.message}`;
  }
};

// Сброс формы при закрытии
watch(() => props.isOpen, (newVal) => {
  if (!newVal) {
    resetForm();
  }
});
</script> 