<template>
  <button 
    class="w-10 h-10 rounded-full flex items-center justify-center focus:outline-none transition-colors" 
    :class="[
      isDarkMode 
        ? 'bg-gray-700 text-yellow-300 hover:bg-gray-600' 
        : 'bg-blue-100 text-blue-600 hover:bg-blue-200'
    ]"
    @click="$emit('click')"
    aria-label="Toggle dark mode"
  >
    <font-awesome-icon :icon="isDarkMode ? 'sun' : 'moon'" class="text-lg" />
  </button>
</template>

<script setup lang="ts">
defineProps<{
  isDarkMode: boolean
}>();

defineEmits<{
  (e: 'click'): void
}>();
</script>