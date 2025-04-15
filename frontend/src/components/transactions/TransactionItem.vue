<template>
    <div class="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
      <div class="flex items-center space-x-3">
        <font-awesome-icon :icon="transaction.icon" :class="transaction.iconColor" />
        <div>
          <p class="font-medium">{{ transaction.title }}</p>
          <p class="text-sm text-gray-500">{{ transaction.address }}</p>
        </div>
      </div>
      <span class="text-sm text-gray-600">{{ transaction.amount }} ETH</span>
    </div>
  </template>
  
  <script setup lang="ts">
  import type { Transaction } from '../../types/transaction';
  
  defineProps<{
    transaction: Transaction
  }>();
  </script>