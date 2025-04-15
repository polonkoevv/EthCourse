<template>
    <div class="bg-white rounded-2xl p-6 shadow-sm">
      <h3 class="text-lg font-bold mb-4">Recent Transactions</h3>
      <div class="space-y-4">
        <TransactionItem 
          v-for="(transaction, index) in transactions" 
          :key="index"
          :transaction="transaction"
        />
      </div>
    </div>
  </template>
  
  <script setup lang="ts">
  import { ref } from 'vue';
  import TransactionItem from './TransactionItem.vue';
  import type { Transaction } from '../../types/transaction';
  
  const transactions = ref<Transaction[]>([
    {
      type: 'exchange',
      title: 'Access Granted',
      address: '0x1234...5678',
      amount: '0.01',
      icon: 'arrow-right-arrow-left', // Упрощенный формат
      iconColor: 'text-blue-600'
    },
    {
      type: 'upload',
      title: 'File Upload',
      address: 'IPFS Hash: Qm...',
      amount: '0.005',
      icon: 'upload',
      iconColor: 'text-green-600'
    }
  ]);
  </script>