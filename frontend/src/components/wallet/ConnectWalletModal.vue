<template>
  <Teleport to="body">
    <div v-if="isOpen" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-[60]" @click.self="closeModal">
      <div class="bg-white rounded-xl w-full max-w-md mx-4 shadow-lg">
        <div class="p-6 border-b border-neutral-100">
          <div class="flex justify-between items-center">
            <h3 class="text-xl font-bold">Подключить кошелек</h3>
            <button class="text-neutral-400 hover:text-neutral-600" @click="closeModal">
              <font-awesome-icon icon="xmark" class="text-xl" />
            </button>
          </div>
        </div>

        <div class="p-6">
          <!-- Шаг 1: Выбор провайдера кошелька -->
          <div v-if="step === 'provider'" class="space-y-4">
            <p class="text-gray-600 mb-6">Выберите провайдер для подключения кошелька</p>
            
            <button 
              @click="selectProvider('metamask')" 
              class="w-full py-3 px-4 bg-white border border-gray-200 rounded-xl hover:bg-gray-50 transition-colors mb-4 flex items-center"
            >
              <img src="../../assets/images/metamask-fox.svg" alt="MetaMask" class="w-8 h-8 mr-3">
              <div class="text-left flex-1">
                <h4 class="font-medium">MetaMask</h4>
                <p class="text-sm text-gray-500">Популярный Ethereum кошелек</p>
              </div>
              <font-awesome-icon icon="arrow-right" class="text-gray-400" />
            </button>
            
            <!-- Можно добавить другие провайдеры (Coinbase, WalletConnect и т.д.) -->
          </div>
          
          <!-- Шаг 2: Выбор конкретного аккаунта -->
          <div v-else-if="step === 'accounts'" class="space-y-4">
            <p class="text-gray-600 mb-6">Выберите аккаунт для подключения</p>
            
            <div class="max-h-64 overflow-y-auto space-y-2">
              <button 
                v-for="account in availableAccounts" 
                :key="account.address"
                @click="selectAccount(account.address)"
                class="w-full py-3 px-4 bg-white border border-gray-200 rounded-xl hover:bg-gray-50 transition-colors flex items-center"
                :class="{'border-blue-500 bg-blue-50': account.address === selectedAccountAddress}"
              >
                <div class="w-8 h-8 bg-blue-100 rounded-full flex items-center justify-center mr-3">
                  <span class="text-xs font-medium text-blue-600">{{ account.index + 1 }}</span>
                </div>
                <div class="text-left flex-1">
                  <h4 class="font-medium truncate">{{ account.address }}</h4>
                  <p class="text-sm text-gray-500">Баланс: {{ account.balance || '—' }} ETH</p>
                </div>
                <font-awesome-icon 
                  v-if="account.address === selectedAccountAddress" 
                  icon="check" 
                  class="text-blue-500 ml-2" 
                />
              </button>
            </div>
            
            <div class="flex justify-between mt-6">
              <button 
                @click="step = 'provider'" 
                class="px-4 py-2 text-gray-600 hover:text-gray-800"
              >
                <font-awesome-icon icon="arrow-left" class="mr-2" />
                Назад
              </button>
              
              <button 
                @click="connectSelectedAccount" 
                class="px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed"
                :disabled="!selectedAccountAddress"
              >
                Подключить
              </button>
            </div>
          </div>
          
          <!-- Шаг 3: Успешное подключение -->
          <div v-else-if="step === 'connected'" class="text-center py-4">
            <div class="w-16 h-16 bg-green-100 rounded-full mx-auto flex items-center justify-center mb-4">
              <font-awesome-icon icon="check" class="text-green-500 text-xl" />
            </div>
            <h4 class="text-lg font-medium mb-2">Кошелек подключен</h4>
            <p class="text-gray-600 mb-2">{{ shortenAddress(walletAddress) }}</p>
            <p class="text-sm text-gray-500 mb-6">Баланс: {{ walletBalance }} ETH</p>
            
            <div class="flex justify-center space-x-4">
              <button 
                @click="step = 'accounts'" 
                class="px-4 py-2 text-blue-500 border border-blue-200 rounded-lg hover:bg-blue-50"
              >
                Сменить аккаунт
              </button>
              
              <button 
                @click="disconnectWallet" 
                class="px-4 py-2 text-red-500 border border-red-200 rounded-lg hover:bg-red-50"
              >
                Отключить кошелек
              </button>
            </div>
          </div>
          
          <!-- Отображение ошибок -->
          <div v-if="error" class="mt-4 p-3 bg-red-50 text-red-600 rounded-lg text-sm">
            {{ error }}
          </div>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue';
import { ethers } from 'ethers';
// Типы для аккаунтов
interface AccountInfo {
  address: string;
  balance?: string;
  index: number;
}

const props = defineProps<{
  isOpen: boolean
}>();

const emit = defineEmits<{
  (e: 'close'): void,
  (e: 'wallet-connected', address: string): void,
  (e: 'wallet-disconnected'): void
}>();

// Состояния модального окна
const step = ref<'provider' | 'accounts' | 'connected'>('provider');
const error = ref('');
const walletAddress = ref('');
const walletBalance = ref('0');
const availableAccounts = ref<AccountInfo[]>([]);
const selectedAccountAddress = ref('');
const selectedProvider = ref<'metamask' | ''>('');

// Проверка наличия MetaMask в окне браузера
const checkIfMetaMaskIsInstalled = (): boolean => {
  // @ts-ignore
  const { ethereum } = window;
  return Boolean(ethereum && ethereum.isMetaMask);
};

// Сокращение адреса кошелька
const shortenAddress = (address: string): string => {
  if (!address) return '';
  return `${address.slice(0, 6)}...${address.slice(-4)}`;
};

// Получение баланса конкретного аккаунта
const getAccountBalance = async (address: string): Promise<string> => {
  try {
    // @ts-ignore
    const provider = new ethers.providers.Web3Provider(window.ethereum);
    const balance = await provider.getBalance(address);
    // @ts-ignore
    return ethers.utils.formatEther(balance).substring(0, 6);
  } catch (err) {
    console.error('Ошибка при получении баланса:', err);
    return '?';
  }
};

// Получение всех доступных аккаунтов
const fetchAvailableAccounts = async (): Promise<void> => {
  try {
    // @ts-ignore
    const accounts = await window.ethereum.request({ method: 'eth_accounts' });
    
    availableAccounts.value = await Promise.all(
      accounts.map(async (address: string, index: number) => {
        const balance = await getAccountBalance(address);
        return { address, balance, index };
      })
    );
    
    // Если не получилось получить аккаунты через eth_accounts (пользователь не авторизован),
    // запрашиваем доступные аккаунты без авторизации
    if (availableAccounts.value.length === 0) {
      try {
        // @ts-ignore
        const accounts = await window.ethereum.request({ method: 'eth_requestAccounts' });
        availableAccounts.value = await Promise.all(
          accounts.map(async (address: string, index: number) => {
            const balance = await getAccountBalance(address);
            return { address, balance, index };
          })
        );
      } catch (err) {
        console.error('Пользователь отклонил запрос на подключение:', err);
      }
    }
  } catch (err) {
    console.error('Ошибка при получении аккаунтов:', err);
    error.value = 'Не удалось получить список аккаунтов. Пожалуйста, убедитесь, что MetaMask разблокирован.';
  }
};

// Выбор провайдера кошелька
const selectProvider = async (provider: 'metamask') => {
  selectedProvider.value = provider;
  error.value = '';
  
  if (provider === 'metamask') {
    if (!checkIfMetaMaskIsInstalled()) {
      error.value = 'MetaMask не установлен. Пожалуйста, установите расширение MetaMask для вашего браузера.';
      return;
    }
    
    try {
      await fetchAvailableAccounts();
      
      if (availableAccounts.value.length > 0) {
        step.value = 'accounts';
        // Предварительно выбираем первый аккаунт
        selectedAccountAddress.value = availableAccounts.value[0].address;
      } else {
        error.value = 'Нет доступных аккаунтов в MetaMask. Пожалуйста, создайте аккаунт или разблокируйте MetaMask.';
      }
    } catch (err: any) {
      console.error('Ошибка при подключении к MetaMask:', err);
      error.value = err.message || 'Произошла ошибка при подключении к MetaMask.';
    }
  }
};

// Выбор конкретного аккаунта
const selectAccount = (address: string) => {
  selectedAccountAddress.value = address;
};

// Подключение выбранного аккаунта
const connectSelectedAccount = async () => {
  if (!selectedAccountAddress.value) return;
  
  try {
    // Запрашиваем переключение на выбранный аккаунт
    // @ts-ignore
    await window.ethereum.request({
      method: 'wallet_requestPermissions',
      params: [{ eth_accounts: { accounts: [selectedAccountAddress.value] } }]
    });
    
    walletAddress.value = selectedAccountAddress.value;
    walletBalance.value = availableAccounts.value.find(
      account => account.address === selectedAccountAddress.value
    )?.balance || '0';
    
    step.value = 'connected';
    emit('wallet-connected', selectedAccountAddress.value);
    
    // Подписка на изменение аккаунтов
    // @ts-ignore
    window.ethereum.on('accountsChanged', handleAccountsChanged);
  } catch (err: any) {
    console.error('Ошибка при подключении аккаунта:', err);
    error.value = err.message || 'Произошла ошибка при подключении аккаунта.';
  }
};

// Обработчик изменения аккаунтов
const handleAccountsChanged = async (accounts: string[]) => {
  if (accounts.length === 0) {
    // Пользователь отключил кошелек
    disconnectWallet();
  } else {
    // Пользователь сменил аккаунт
    walletAddress.value = accounts[0];
    walletBalance.value = await getAccountBalance(accounts[0]);
    
    // Обновляем список доступных аккаунтов
    await fetchAvailableAccounts();
    
    // Если аккаунт сменился и мы находимся на шаге подключенного кошелька
    if (step.value === 'connected') {
      emit('wallet-connected', accounts[0]);
    }
  }
};

// Отключение кошелька
const disconnectWallet = () => {
  step.value = 'provider';
  walletAddress.value = '';
  walletBalance.value = '0';
  selectedAccountAddress.value = '';
  emit('wallet-disconnected');
  
  // Отписываемся от событий
  // @ts-ignore
  if (window.ethereum) {
    // @ts-ignore
    window.ethereum.removeListener('accountsChanged', handleAccountsChanged);
  }
};

// Закрытие модального окна
const closeModal = () => {
  emit('close');
};

// Проверяем, подключен ли кошелек при монтировании компонента
onMounted(async () => {
  try {
    if (checkIfMetaMaskIsInstalled()) {
      // @ts-ignore
      const accounts = await window.ethereum.request({ method: 'eth_accounts' });
      if (accounts.length > 0) {
        walletAddress.value = accounts[0];
        walletBalance.value = await getAccountBalance(accounts[0]);
        await fetchAvailableAccounts();
        step.value = 'connected';
        emit('wallet-connected', accounts[0]);
        
        // Подписка на изменение аккаунтов
        // @ts-ignore
        window.ethereum.on('accountsChanged', handleAccountsChanged);
      }
    }
  } catch (err) {
    console.error('Ошибка при проверке подключения кошелька:', err);
  }
});

// Сброс состояния при закрытии модального окна
watch(() => props.isOpen, (newVal) => {
  if (newVal && !walletAddress.value) {
    // Если модальное окно открыто и кошелек не подключен,
    // показываем первый шаг
    step.value = 'provider';
    error.value = '';
  }
});
</script>

<style scoped>
.max-h-64 {
  max-height: 16rem;
}
</style> 