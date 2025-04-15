<template>
  <button class="text-gray-600 hover:text-blue-600" @click="$emit('click')">
    <font-awesome-icon v-if="!isDarkMode" :icon="['far', 'moon']" />
    <font-awesome-icon v-else :icon="['far', 'sun']" />
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