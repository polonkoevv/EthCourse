<template>
    <div class="bg-white rounded-2xl p-6 shadow-sm">
      <div class="flex items-center space-x-4 mb-6">
        <div class="w-24 h-24 bg-gradient-to-br from-blue-400 to-cyan-300 rounded-lg flex items-center justify-center">
          <font-awesome-icon icon="music" class="text-3xl text-white" />
        </div>
        <div>
          <h3 class="text-xl font-bold">Currently Playing</h3>
          <p class="text-gray-600">{{ currentTrack.name }}</p>
          <p class="text-sm text-gray-500">{{ currentTrack.artist }}</p>
        </div>
      </div>
      <div class="space-y-4">
        <div 
          class="h-2 bg-gray-200 rounded-full relative cursor-pointer" 
          ref="progressBarRef"
          @click="seekAudio"
          @mousemove="updateSeekPreview"
          @mouseleave="isPreviewVisible = false"
        >
          <div 
            class="h-full bg-blue-500 rounded-full" 
            :style="{ width: `${progressPercentage}%` }"
          ></div>
          <div 
            class="absolute top-1/2 -translate-y-1/2 w-4 h-4 bg-white rounded-full border-2 border-blue-500 shadow-md"
            :style="{ left: `calc(${progressPercentage}% - 0.5rem)` }"
            v-show="progressPercentage > 0"
          ></div>
          <div 
            v-if="isPreviewVisible"
            class="absolute top-0 -translate-y-8 bg-gray-800 text-white px-2 py-1 text-xs rounded"
            :style="{ left: `calc(${previewPosition}% - 1.5rem)` }"
          >
            {{ formatTime(previewTime) }}
          </div>
        </div>
        <div class="flex justify-between text-sm text-gray-500">
          <span>{{ formatTime(currentTime) }}</span>
          <span>{{ formatTime(duration) }}</span>
        </div>
        <div class="flex items-center justify-center space-x-6">
          <button 
            class="text-gray-600 hover:text-blue-600 text-xl" 
            @click="playPrevious"
            :disabled="!hasPrevious"
            :class="{'opacity-50 cursor-not-allowed': !hasPrevious}"
          >
            <font-awesome-icon icon="backward-step" />
          </button>
          <button 
            class="w-12 h-12 bg-blue-600 rounded-full flex items-center justify-center text-white hover:bg-blue-700"
            @click="togglePlay"
          >
            <font-awesome-icon :icon="isPlaying ? 'pause' : 'play'" />
          </button>
          <button 
            class="text-gray-600 hover:text-blue-600 text-xl" 
            @click="playNext"
            :disabled="!hasNext"
            :class="{'opacity-50 cursor-not-allowed': !hasNext}"
          >
            <font-awesome-icon icon="forward-step" />
          </button>
        </div>
      </div>
      <audio 
        ref="audioElement"
        :src="currentTrack.link" 
        @timeupdate="updateProgress" 
        @loadedmetadata="onAudioLoaded"
        @ended="onAudioEnded"
      ></audio>
    </div>
  </template>
  
  <script setup lang="ts">
  import { ref, computed, onMounted, watch } from 'vue';
  import type { Track, Audio } from '../../types/audio';
  
  const props = defineProps<{
    selectedAudio?: Audio
    audioList: Audio[]
  }>();
  
  const emit = defineEmits<{
    (e: 'change-track', index: number): void
  }>();
  
  const isPlaying = ref(false);
  const currentTime = ref(0);
  const duration = ref(0);
  const audioElement = ref<HTMLAudioElement | null>(null);
  const progressBarRef = ref<HTMLDivElement | null>(null);
  const isPreviewVisible = ref(false);
  const previewPosition = ref(0);
  const previewTime = ref(0);
  const currentIndex = ref(-1);
  const currentTrack = ref<Track>({
    name: 'Выберите трек',
    artist: 'Нет исполнителя',
    link: ''
  });
  
  const progressPercentage = computed(() => {
    if (duration.value === 0) return 0;
    return (currentTime.value / duration.value) * 100;
  });
  
  const hasNext = computed(() => {
    return currentIndex.value < props.audioList.length - 1 && currentIndex.value !== -1;
  });
  
  const hasPrevious = computed(() => {
    return currentIndex.value > 0;
  });
  
  const togglePlay = () => {
    if (!audioElement.value || !currentTrack.value.link) return;
    
    if (isPlaying.value) {
      audioElement.value.pause();
    } else {
      audioElement.value.play();
    }
    isPlaying.value = !isPlaying.value;
  };
  
  const updateProgress = () => {
    if (!audioElement.value) return;
    currentTime.value = audioElement.value.currentTime;
  };
  
  const onAudioLoaded = () => {
    if (!audioElement.value) return;
    duration.value = audioElement.value.duration;
  };
  
  const onAudioEnded = () => {
    isPlaying.value = false;
    currentTime.value = 0;
    
    if (hasNext.value) {
      playNext();
    }
  };
  
  const seekAudio = (event: MouseEvent) => {
    if (!audioElement.value || !progressBarRef.value || duration.value === 0) return;
    
    const rect = progressBarRef.value.getBoundingClientRect();
    const clickPosition = event.clientX - rect.left;
    const progressBarWidth = rect.width;
    const seekPercentage = (clickPosition / progressBarWidth);
    
    audioElement.value.currentTime = seekPercentage * duration.value;
  };
  
  const updateSeekPreview = (event: MouseEvent) => {
    if (!progressBarRef.value || duration.value === 0) return;
    
    isPreviewVisible.value = true;
    
    const rect = progressBarRef.value.getBoundingClientRect();
    const hoverPosition = event.clientX - rect.left;
    const progressBarWidth = rect.width;
    const hoverPercentage = (hoverPosition / progressBarWidth) * 100;
    
    previewPosition.value = Math.max(0, Math.min(100, hoverPercentage));
    previewTime.value = (previewPosition.value / 100) * duration.value;
  };
  
  const playNext = () => {
    if (!hasNext.value) return;
    
    const nextIndex = currentIndex.value + 1;
    emit('change-track', nextIndex);
  };
  
  const playPrevious = () => {
    if (!hasPrevious.value) return;
    
    const prevIndex = currentIndex.value - 1;
    emit('change-track', prevIndex);
  };
  
  const skipAudio = (seconds: number) => {
    if (!audioElement.value) return;
    
    const newTime = audioElement.value.currentTime + seconds;
    audioElement.value.currentTime = Math.max(0, Math.min(duration.value, newTime));
  };
  
  const formatTime = (seconds: number): string => {
    if (isNaN(seconds)) return '0:00';
    const mins = Math.floor(seconds / 60);
    const secs = Math.floor(seconds % 60);
    return `${mins}:${secs.toString().padStart(2, '0')}`;
  };
  
  watch(() => props.selectedAudio, (newAudio) => {
    if (newAudio) {
      currentIndex.value = props.audioList.findIndex(
        audio => audio.title === newAudio.title && audio.link === newAudio.link
      );
      
      currentTrack.value = {
        name: newAudio.title,
        artist: 'Неизвестный исполнитель',
        link: newAudio.link
      };
      
      if (audioElement.value) {
        isPlaying.value = false;
        currentTime.value = 0;
        
        const playWhenReady = () => {
          if (audioElement.value) {
            audioElement.value.play()
              .then(() => { isPlaying.value = true; })
              .catch(error => console.error('Ошибка воспроизведения:', error));
            
            audioElement.value.removeEventListener('loadedmetadata', playWhenReady);
          }
        };
        
        audioElement.value.addEventListener('loadedmetadata', playWhenReady);
        
        audioElement.value.src = newAudio.link;
        audioElement.value.load();
      }
    }
  }, { immediate: true });
  
  const setupKeyboardControls = () => {
    window.addEventListener('keydown', (event) => {
      if (!audioElement.value) return;
      
      switch (event.code) {
        case 'Space':
          event.preventDefault();
          togglePlay();
          break;
        case 'ArrowRight':
          if (event.shiftKey) {
            event.preventDefault();
            playNext();
          } else {
            event.preventDefault();
            skipAudio(5);
          }
          break;
        case 'ArrowLeft':
          if (event.shiftKey) {
            event.preventDefault();
            playPrevious();
          } else {
            event.preventDefault();
            skipAudio(-5);
          }
          break;
      }
    });
  };
  
  onMounted(() => {
    if (props.selectedAudio) {
      currentTrack.value = {
        name: props.selectedAudio.title,
        artist:'Неизвестный исполнитель',
        link: props.selectedAudio.link
      };
      
      currentIndex.value = props.audioList.findIndex(
        audio => audio.title === props.selectedAudio.title && audio.link === props.selectedAudio.link
      );
    }
    
    setupKeyboardControls();
  });
  </script>

<style scoped>
.cursor-pointer {
  cursor: pointer;
}

button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
</style>