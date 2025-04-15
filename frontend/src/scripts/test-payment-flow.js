import { ethers } from "ethers";
import fs from "fs";

async function main() {
  // Подключаемся к Ganache
  const provider = new ethers.providers.JsonRpcProvider("http://127.0.0.1:8545");
  
  // Получаем адрес контракта
  const contractData = JSON.parse(fs.readFileSync("../contracts/MusicPayment.json"));
  const contractAddress = contractData.MusicPayment;
  
  // Загружаем ABI
  const contractABI = JSON.parse(fs.readFileSync("./src/contracts/MusicPayment.json")).abi;
  
  // Получаем аккаунты из Ganache
  const accounts = await provider.listAccounts();
  const owner = accounts[0];
  const artist = accounts[1];
  const listener = accounts[2];
  
  // Подключаемся к контракту от имени разных пользователей
  const ownerContract = new ethers.Contract(contractAddress, contractABI, provider.getSigner(0));
  const artistContract = new ethers.Contract(contractAddress, contractABI, provider.getSigner(1));
  const listenerContract = new ethers.Contract(contractAddress, contractABI, provider.getSigner(2));
  
  console.log("Тестирование полного цикла платежей");
  console.log("----------------------------------");
  
  // 1. Артист загружает трек
  const songCID = "QmTestCIDForSongUpload123456";
  const uploadPrice = await ownerContract.uploadPrice();
  
  console.log(`Артист ${artist} загружает трек ${songCID}`);
  console.log(`Комиссия за загрузку: ${ethers.utils.formatEther(uploadPrice)} ETH`);
  
  const uploadTx = await artistContract.uploadSong(songCID, { value: uploadPrice });
  await uploadTx.wait();
  console.log("Трек успешно загружен");
  
  // 2. Артист устанавливает цену за прослушивание
  const songPrice = ethers.utils.parseEther("0.005");
  console.log(`Артист устанавливает цену: ${ethers.utils.formatEther(songPrice)} ETH`);
  
  const setPriceTx = await artistContract.setSongPrice(songCID, songPrice);
  await setPriceTx.wait();
  console.log("Цена установлена");
  
  // 3. Слушатель покупает доступ
  console.log(`Слушатель ${listener} приобретает доступ к треку`);
  
  const purchaseTx = await listenerContract.purchaseSong(songCID, { value: songPrice });
  await purchaseTx.wait();
  console.log("Доступ успешно приобретен");
  
  // 4. Проверка доступа
  const hasAccess = await ownerContract.hasSongAccess(listener, songCID);
  console.log(`Доступ слушателя к треку: ${hasAccess ? "Разрешен" : "Запрещен"}`);
  
  console.log("----------------------------------");
  console.log("Тестирование успешно завершено");
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });