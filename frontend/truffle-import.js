// Этот файл использует динамический импорт для загрузки ES модуля
async function getTruffleConfig() {
  const config = await import('./truffle-config.js');
  return config.default;
}

module.exports = getTruffleConfig(); 