<template>
    <div class="bg-white rounded-2xl p-6 shadow-sm">
      <h3 class="text-lg font-bold mb-4">Последние транзакции</h3>
      <div class="transactions-container overflow-y-auto max-h-80 pr-1">
        <div class="space-y-4">
          <TransactionItem 
            v-for="(transaction, index) in transactions" 
            :key="index"
            :transaction="transaction"
          />
        </div>
        <div v-if="transactions.length === 0" class="text-center py-8 text-gray-500">
          <p>Подключите свой кошелек и получайте транзакции.</p>
        </div>
      </div>
    </div>
  </template>
  
  <script setup lang="ts">
  import { onMounted, ref } from 'vue';
  import TransactionItem from './TransactionItem.vue';
  import type { Transaction } from '../../types/transaction';
  import axiosEntity from '../../scripts/axios';

  // Функция для получения адреса кошелька из localStorage
  const getWalletAddressFromStorage = (): string => {
    const wallet = localStorage.getItem('wallet');
    if (wallet) {
      const walletData = JSON.parse(wallet);
      return walletData.address || '';
    }
    return '';
  };

  const transactions = ref<Transaction[]>([]);

  const shortenAddress = (address: string, startChars: number = 6, endChars: number = 4): string => {
  if (!address) return '';
  if (address.length <= startChars + endChars) return address;
  
  return `${address.substring(0, startChars)}...${address.substring(address.length - endChars)}`;
};
const formatEthBalance = (weiAmount: string | number, decimals: number = 4): string => {
  // Проверка входного параметра
  if (!weiAmount && weiAmount !== 0) return '0';
  
  // Преобразование в BigNumber для избежания проблем с точностью
  const amount = typeof weiAmount === 'string' ? weiAmount : weiAmount.toString();
  
  // Преобразование из Wei в ETH (1 ETH = 10^18 Wei)
  const ethValue = parseFloat(amount) / 1e18;
  
  // Форматирование числа с заданным количеством десятичных знаков
  const formattedValue = ethValue.toFixed(decimals);
  
  // Удаление незначащих нулей в конце (опционально)
  const trimmedValue = parseFloat(formattedValue).toString();
  
  // Добавление ETH в конце
  return `${trimmedValue}`;
};

  onMounted(async () => {
    try {
      // Получаем адрес кошелька из localStorage
      const walletAddress = getWalletAddressFromStorage();
      
      if (walletAddress) {
        // Используем метод из axiosEntity для получения истории транзакций
        const response = await axiosEntity.GetTransactionHistoryFromChain(walletAddress);
        
        // Если есть данные, обновляем массив transactions
        if (response.data && Array.isArray(response.data)) {
          let tempTrans = ref<Transaction[]>([]);
          for (let i = 0; i < response.data.length; i++) {
            var tt = {
              type: '',
              title: '',
              address: '',
              amount: '',
              icon: '',
              iconColor: '',
              date: response.data[i].Timestamp
            }
            if (response.data[i].TransactionType === "outgoing") {
              tt.type = "upload";
              tt.title = "Uploaded";
              tt.address = shortenAddress(response.data[i].To);
              tt.amount = formatEthBalance(response.data[i].Value);
              tt.icon = 'upload';
              tt.iconColor = 'text-green-600'
            } else {
              tt.type = "exchange";
              tt.title = "Access Granted";
              tt.address = shortenAddress(response.data[i].From);
              tt.amount = formatEthBalance(response.data[i].Value);
              tt.icon = 'exchange';
              tt.iconColor = 'text-blue-600';
            }
            tempTrans.value.push(tt);
          }
          transactions.value = tempTrans.value;
        }


        
        console.log('Получены транзакции:', response.data);
      } else {
        console.log('Адрес кошелька не найден в localStorage');
      }
    } catch (error) {
      console.error('Ошибка при получении истории транзакций:', error);
    }
  });
  </script>
  
  <style scoped>
  /* Стили для полосы прокрутки */
  .transactions-container {
    scrollbar-width: thin;
    scrollbar-color: #cbd5e0 #f7fafc;
  }
  
  .transactions-container::-webkit-scrollbar {
    width: 4px;
  }
  
  .transactions-container::-webkit-scrollbar-track {
    background: #f7fafc;
    border-radius: 4px;
  }
  
  .transactions-container::-webkit-scrollbar-thumb {
    background-color: #cbd5e0;
    border-radius: 4px;
  }
  </style>