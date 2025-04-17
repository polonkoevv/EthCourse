<template>
  <header class="fixed dark:bg-gray-800 w-full bg-white/80 backdrop-blur-lg border-b border-gray-100 z-50">
    <div class="container mx-auto px-4">
      <div class="flex items-center justify-between h-16">
        <div class="flex items-center space-x-4">
          <div class="text-blue-600 text-xl font-bold">AudioChain</div>
          <nav class="hidden md:flex space-x-6">
            <!-- <span class="text-gray-600 hover:text-blue-600 cursor-pointer">Discover</span>
            <span class="text-gray-600 hover:text-blue-600 cursor-pointer">Library</span>
            <span class="text-gray-600 hover:text-blue-600 cursor-pointer">Upload</span> -->
          </nav>
        </div>
        <div class="flex items-center space-x-4">
          <ThemeToggle @click="$emit('toggle-theme')" :isDarkMode="isDarkMode" />
          <button 
            class="text-gray-600 hover:text-blue-600 relative"
            @click="openWalletModal"
          >
            <div class="flex items-center justify-center">
              <!-- Показываем разные иконки в зависимости от статуса подключения -->
              <div v-if="!isWalletConnected" class="h-8 w-8 flex items-center justify-center">
                <font-awesome-icon icon="user" class="text-lg" />
              </div>
              <div v-else class="h-8 w-8 bg-blue-500 rounded-full flex items-center justify-center text-white">
                <font-awesome-icon icon="wallet" class="text-sm" />
              </div>
              
              <!-- Индикатор подключенного кошелька -->
              <div 
                v-if="isWalletConnected" 
                class="absolute -top-1 -right-1 h-3 w-3 bg-green-500 rounded-full border border-white"
              ></div>
            </div>
          </button>
        </div>
      </div>
    </div>
    
    <!-- Модальное окно подключения кошелька -->
    <ConnectWalletModal 
      :is-open="isWalletModalOpen" 
      @close="isWalletModalOpen = false"
      @wallet-connected="onWalletConnected"
      @wallet-disconnected="onWalletDisconnected"
    />
  </header>
</template>

<script setup lang="ts">
// Объявление типа для window.ethereum


import ThemeToggle from '../ui/ThemeToggle.vue';
import ConnectWalletModal from '../wallet/ConnectWalletModal.vue';
import { ref } from 'vue';

defineProps<{
  isDarkMode: boolean
}>();

const emit = defineEmits<{
  (e: 'toggle-theme'): void,
  (e: 'wallet-connected', address: string): void,
  (e: 'wallet-disconnected'): void
}>();

const isWalletModalOpen = ref(false);
const isWalletConnected = ref(false);
const walletAddress = ref('');

const toggleTheme = () => {
  emit('toggle-theme');
};

const openWalletModal = () => {
  isWalletModalOpen.value = true;
};

const onWalletConnected = (address: string) => {
  isWalletConnected.value = true;
  walletAddress.value = address;
  emit('wallet-connected', address);
};

const onWalletDisconnected = () => {
  isWalletConnected.value = false;
  walletAddress.value = '';
  emit('wallet-disconnected');
};
</script>