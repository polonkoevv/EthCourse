<template>
    <div class="flex items-center space-x-2">
      <span class="text-sm text-gray-600">{{ balance }} ETH</span>
      <img :src="avatarUrl" class="w-8 h-8 rounded-full">
    </div>
  </template>
  
  <script setup lang="ts">
  import { ref } from 'vue';
  
  const balance = ref('0.045');
  const avatarUrl = ref('https://storage.googleapis.com/uxpilot-auth.appspot.com/avatars/avatar-2.jpg');
  </script>