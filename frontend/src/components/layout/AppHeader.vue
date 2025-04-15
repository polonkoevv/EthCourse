<template>
  <header class="fixed w-full bg-white/80 backdrop-blur-lg border-b border-gray-100 z-50">
    <div class="container mx-auto px-4">
      <div class="flex items-center justify-between h-16">
        <div class="flex items-center space-x-4">
          <div class="text-blue-600 text-xl font-bold">AudioChain</div>
          <nav class="hidden md:flex space-x-6">
            <span class="text-gray-600 hover:text-blue-600 cursor-pointer">Discover</span>
            <span class="text-gray-600 hover:text-blue-600 cursor-pointer">Library</span>
            <span class="text-gray-600 hover:text-blue-600 cursor-pointer">Upload</span>
          </nav>
        </div>
        <div class="flex items-center space-x-4">
          <ThemeToggle @click="$emit('toggle-theme')" :isDarkMode="isDarkMode" />
          <WalletInfo />
        </div>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
import ThemeToggle from '../ui/ThemeToggle.vue';
import WalletInfo from '../ui/WalletInfo.vue';

defineProps<{
  isDarkMode: boolean
}>();

defineEmits<{
  (e: 'toggle-theme'): void
}>();
</script>