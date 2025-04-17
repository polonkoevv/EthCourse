import { ethers } from 'ethers';
import AudioChainABI from '../../build/contracts/AudioChain.json';

// Адрес смарт-контракта (после деплоя в Ganache)
const CONTRACT_ADDRESS = '0x...'; 

// Вспомогательная функция для получения экземпляра контракта
const getContract = async (needSigner = false) => {
  // @ts-ignore - для MetaMask
  if (!window.ethereum) throw new Error('MetaMask не установлен');
  
  // @ts-ignore
  const provider = new ethers.providers.Web3Provider(window.ethereum);
  
  if (needSigner) {
    const signer = provider.getSigner();
    return new ethers.Contract(CONTRACT_ADDRESS, AudioChainABI, signer);
  }
  
  return new ethers.Contract(CONTRACT_ADDRESS, AudioChainABI, provider);
};

// Публикация аудио в блокчейне
export const publishAudio = async (title: string, artist: string, ipfsHash: string, price: string) => {
  try {
    const contract = await getContract(true);
    const priceInWei = ethers.utils.parseEther(price);
    
    // Вызываем метод контракта
    const tx = await contract.publishAudio(title, artist, ipfsHash, priceInWei);
    
    // Ожидаем завершения транзакции
    const receipt = await tx.wait();
    
    // Получаем ID аудио из событий транзакции
    const event = receipt.events.find(event => event.event === 'AudioPublished');
    const audioId = event.args.id.toString();
    
    return {
      txHash: receipt.transactionHash,
      audioId
    };
  } catch (error) {
    console.error('Ошибка при публикации аудио:', error);
    throw error;
  }
};

// Покупка аудио
export const purchaseAudio = async (audioId: string, price: string) => {
  try {
    const contract = await getContract(true);
    const priceInWei = ethers.utils.parseEther(price);
    
    // Вызываем метод покупки
    const tx = await contract.purchaseAudio(audioId, {
      value: priceInWei
    });
    
    // Ожидаем завершения транзакции
    const receipt = await tx.wait();
    
    return {
      txHash: receipt.transactionHash,
      success: true
    };
  } catch (error) {
    console.error('Ошибка при покупке аудио:', error);
    throw error;
  }
};

// Проверка доступа к аудио
export const checkAudioAccess = async (audioId: string, userAddress: string) => {
  try {
    const contract = await getContract();
    const hasAccess = await contract.hasAccess(userAddress, audioId);
    return hasAccess;
  } catch (error) {
    console.error('Ошибка при проверке доступа:', error);
    return false;
  }
};

// Получение списка всех аудио
export const getAllAudios = async () => {
  try {
    const contract = await getContract();
    const count = await contract.getAudioCount();
    
    const audios = [];
    for (let i = 1; i <= count; i++) {
      const audio = await contract.audios(i);
      audios.push({
        id: audio.id.toString(),
        title: audio.title,
        artist: audio.artist,
        ipfsHash: audio.ipfsHash,
        price: ethers.utils.formatEther(audio.price),
        owner: audio.owner,
        isForSale: audio.isForSale
      });
    }
    
    return audios;
  } catch (error) {
    console.error('Ошибка при получении списка аудио:', error);
    throw error;
  }
};