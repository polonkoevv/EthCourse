import { Web3Storage } from 'web3.storage';

// API-ключ Web3.Storage
// В реальном приложении следует получать его из переменных окружения
const WEB3_STORAGE_TOKEN = 'ваш_api_ключ';

// Инициализация клиента Web3.Storage
const client = new Web3Storage({ token: WEB3_STORAGE_TOKEN });

/**
 * Загружает файл в IPFS через Web3.Storage
 * @param file - Файл для загрузки
 * @param onProgress - Функция обратного вызова для отслеживания прогресса (опционально)
 * @returns Объект с CID (IPFS-хешем) загруженного файла
 */
export const uploadToIPFS = async (
  file: File,
  onProgress?: (progress: number) => void
): Promise<{ ipfsHash: string }> => {
  try {
    // Создаем массив файлов для загрузки
    const files = [
      new File([file], file.name, { type: file.type })
    ];
    
    // Имитация прогресса загрузки (Web3.Storage не имеет нативной поддержки прогресса)
    let progress = 0;
    const progressInterval = setInterval(() => {
      progress += 0.1;
      if (progress >= 1) {
        progress = 0.95; // Держим на 95% до полного завершения
        clearInterval(progressInterval);
      }
      if (onProgress) onProgress(progress);
    }, 500);
    
    // Загружаем файл
    const cid = await client.put(files, {
      name: file.name,
      wrapWithDirectory: false
    });
    
    // Очищаем интервал и устанавливаем прогресс на 100%
    clearInterval(progressInterval);
    if (onProgress) onProgress(1);
    
    return { ipfsHash: cid };
  } catch (error) {
    console.error('Ошибка при загрузке в IPFS:', error);
    throw new Error('Не удалось загрузить файл в IPFS');
  }
};

/**
 * Получает URL для доступа к файлу по IPFS-хешу
 * @param ipfsHash - IPFS-хеш (CID) файла
 * @returns URL для доступа к файлу
 */
export const getIPFSUrl = (ipfsHash: string): string => {
  return `https://localhost:8080/ipfs/${ipfsHash}`;
};
