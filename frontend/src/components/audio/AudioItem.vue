<template>
  <div 
    class="flex items-center justify-between p-4 hover:bg-gray-50 rounded-lg transition-colors cursor-pointer" 
    @click="selectAudio"
  >
    <div class="flex items-center space-x-4">
      <div :class="['w-12', 'h-12', 'bg-gradient-to-br', 'from-blue-400', 'to-cyan-300', 'rounded', 'flex', 'items-center', 'justify-center']">
        <font-awesome-icon icon="music" class="text-white" />
      </div>
      <div>
        <p class="font-medium">{{ audio.title }}</p>
        <p class="text-sm text-gray-500">{{ audio.cid }}</p>
      </div>
    </div>
    <div class="flex items-center space-x-4">
      <button class="text-gray-600 hover:text-blue-600" @click.stop>
        <font-awesome-icon icon="share-nodes" />
      </button>
      <button class="text-gray-600 hover:text-blue-600" @click.stop>
        <font-awesome-icon icon="ellipsis" />
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Audio } from '../../types/audio';

const props = defineProps<{
  audio: Audio
}>();

const emit = defineEmits<{
  (e: 'select', audio: Audio): void
}>();

const selectAudio = () => {
  emit('select', props.audio);
};
</script>