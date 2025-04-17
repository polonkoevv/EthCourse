package service

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	shell "github.com/ipfs/go-ipfs-api"
	"github.com/polonkoevv/ethcourse/internal/model"
	"github.com/polonkoevv/ethcourse/internal/storage/postgres"
)

type Service struct {
	sh *shell.Shell
	pg *postgres.Postgres
}

func NewService(sh *shell.Shell, pg *postgres.Postgres) *Service {
	return &Service{sh: sh, pg: pg}
}

func (s *Service) UploadFile(ctx context.Context, name string, file *os.File, walletAddress, signature string, uploadedAt time.Time) (string, error) {
	cid, err := s.sh.Add(file)
	if err != nil {
		return "", err
	}

	err = s.sh.Pin(cid)
	if err != nil {
		return "", err
	}

	music := model.Music{
		Title:      name,
		CID:        cid,
		OwnerAddr:  walletAddress,
		Signature:  signature,
		UploadedAt: uploadedAt,
	}

	_, err = s.pg.CreateMusic(ctx, music)
	if err != nil {
		return "", err
	}
	return cid, nil
}

func (s *Service) GetAllMusic(ctx context.Context) ([]model.Music, error) {
	return s.pg.GetAllMusic(ctx)
}

func (s *Service) GetTransactionHistoryFromChain(address string, startBlock, endBlock uint64, rpcURL string) ([]model.BlockchainTransaction, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к ноде: %v", err)
	}

	ctx := context.Background()
	targetAddress := common.HexToAddress(address)
	fmt.Printf("Ищем транзакции для адреса: %s\n", targetAddress.Hex())
	var transactions []model.BlockchainTransaction

	// Если конечный блок не указан, используем текущий блок
	if endBlock == 0 {
		header, err := client.HeaderByNumber(ctx, nil)
		if err != nil {
			return nil, fmt.Errorf("ошибка получения текущего блока: %v", err)
		}
		endBlock = header.Number.Uint64()
	}

	startBlock = 0

	fmt.Printf("Сканирование всех блоков с %d по %d...\n", startBlock, endBlock)

	type rpcTransaction struct {
		Hash        string `json:"hash"`
		From        string `json:"from"`
		To          string `json:"to"`
		Value       string `json:"value"`
		Gas         string `json:"gas"`
		GasPrice    string `json:"gasPrice"`
		Input       string `json:"input"`
		BlockNumber string `json:"blockNumber"`
	}

	for blockNum := startBlock; blockNum <= endBlock; blockNum++ {
		fmt.Printf("Сканирование блока %d...\n", blockNum)

		block, err := client.BlockByNumber(ctx, big.NewInt(int64(blockNum)))
		if err != nil {
			fmt.Printf("Ошибка получения блока %d: %v\n", blockNum, err)
			continue
		}

		blockTime := time.Unix(int64(block.Time()), 0)

		for _, tx := range block.Transactions() {
			txHash := tx.Hash().Hex()

			var result rpcTransaction
			err := client.Client().Call(&result, "eth_getTransactionByHash", txHash)
			if err != nil {
				fmt.Printf("Ошибка получения транзакции %s: %v\n", txHash, err)
				continue
			}

			fromAddrLower := strings.ToLower(result.From)
			targetAddrLower := strings.ToLower(targetAddress.Hex())
			var toAddrLower string
			if result.To != "" {
				toAddrLower = strings.ToLower(result.To)
			}

			isIncoming := toAddrLower == targetAddrLower
			isOutgoing := fromAddrLower == targetAddrLower

			if isIncoming || isOutgoing {
				txType := "incoming"
				if isOutgoing {
					txType = "outgoing"
				}

				valueInt, _ := new(big.Int).SetString(result.Value[2:], 16)
				gasInt, _ := strconv.ParseUint(result.Gas[2:], 16, 64)
				gasPriceInt, _ := new(big.Int).SetString(result.GasPrice[2:], 16)

				transaction := model.BlockchainTransaction{
					Hash:            txHash,
					BlockNumber:     blockNum,
					From:            result.From,
					To:              result.To,
					Value:           valueInt.String(),
					Gas:             gasInt,
					GasPrice:        gasPriceInt.String(),
					Input:           result.Input[2:], // Убираем "0x" префикс
					Timestamp:       blockTime,
					TransactionType: txType,
				}

				transactions = append(transactions, transaction)
				fmt.Printf("НАЙДЕНА %s транзакция в блоке %d: %s\n", txType, blockNum, txHash)
			}
		}
	}

	// Альтернативный метод: получаем все транзакции через eth_getBlockByNumber с полными деталями
	if len(transactions) == 0 {
		fmt.Println("Стандартный метод не нашел транзакций. Пробуем альтернативный метод...")

		type rpcBlock struct {
			Transactions []rpcTransaction `json:"transactions"`
			Timestamp    string           `json:"timestamp"`
		}

		for blockNum := startBlock; blockNum <= endBlock; blockNum++ {
			blockNumHex := fmt.Sprintf("0x%x", blockNum)

			var block rpcBlock
			err := client.Client().Call(&block, "eth_getBlockByNumber", blockNumHex, true)
			if err != nil {
				fmt.Printf("Ошибка получения блока %d: %v\n", blockNum, err)
				continue
			}

			// Преобразуем timestamp из hex в int
			timestampHex := block.Timestamp[2:] // убираем "0x" префикс
			timestamp, _ := strconv.ParseInt(timestampHex, 16, 64)
			blockTime := time.Unix(timestamp, 0)

			for _, tx := range block.Transactions {
				fromAddrLower := strings.ToLower(tx.From)
				targetAddrLower := strings.ToLower(targetAddress.Hex())
				var toAddrLower string
				if tx.To != "" {
					toAddrLower = strings.ToLower(tx.To)
				}

				isIncoming := toAddrLower == targetAddrLower
				isOutgoing := fromAddrLower == targetAddrLower

				if isIncoming || isOutgoing {
					txType := "incoming"
					if isOutgoing {
						txType = "outgoing"
					}

					// Конвертируем hex в int для значений
					valueInt, _ := new(big.Int).SetString(tx.Value[2:], 16)
					gasInt, _ := strconv.ParseUint(tx.Gas[2:], 16, 64)
					gasPriceInt, _ := new(big.Int).SetString(tx.GasPrice[2:], 16)

					transaction := model.BlockchainTransaction{
						Hash:            tx.Hash,
						BlockNumber:     blockNum,
						From:            tx.From,
						To:              tx.To,
						Value:           valueInt.String(),
						Gas:             gasInt,
						GasPrice:        gasPriceInt.String(),
						Input:           tx.Input[2:], // Убираем "0x" префикс
						Timestamp:       blockTime,
						TransactionType: txType,
					}

					transactions = append(transactions, transaction)
					fmt.Printf("НАЙДЕНА %s транзакция в блоке %d: %s\n", txType, tx.Hash, transaction.Value)
				}
			}
		}
	}

	fmt.Printf("Сканирование завершено. Найдено транзакций: %d\n", len(transactions))
	return transactions, nil
}

// Вывод транзакций
func (s *Service) PrintBlockchainTransactions(transactions []model.BlockchainTransaction) {
	fmt.Println("История транзакций из блокчейна:")
	fmt.Println("--------------------------------------------------")

	for i, tx := range transactions {
		fmt.Printf("Транзакция #%d (%s):\n", i+1, tx.TransactionType)
		fmt.Printf("  Хеш: %s\n", tx.Hash)
		fmt.Printf("  Блок: %d\n", tx.BlockNumber)
		fmt.Printf("  Дата: %s\n", tx.Timestamp.Format("02.01.2006 15:04:05"))
		fmt.Printf("  От: %s\n", tx.From)
		fmt.Printf("  Кому: %s\n", tx.To)
		fmt.Printf("  Сумма (Wei): %s\n", tx.Value)
		fmt.Printf("  Газ: %d, Цена газа: %s Wei\n", tx.Gas, tx.GasPrice)
		if len(tx.Input) > 10 {
			fmt.Printf("  Данные: 0x%s...\n", tx.Input[:10])
		} else {
			fmt.Printf("  Данные: 0x%s\n", tx.Input)
		}
		fmt.Println("--------------------------------------------------")
	}
}

// VerifySignature проверяет подпись Ethereum и восстанавливает адрес кошелька
func (s *Service) VerifySignature(message, signature string) (bool, string, error) {
	// Очищаем префикс '0x' если он есть
	cleanSignature := signature
	if strings.HasPrefix(cleanSignature, "0x") {
		cleanSignature = cleanSignature[2:]
	} else {
		// Если префикса нет, добавляем его для hexutil.Decode
		signature = "0x" + signature
	}

	// Декодирование подписи
	signatureBytes, err := hexutil.Decode(signature)
	if err != nil {
		// Пробуем альтернативный способ декодирования
		signatureBytes, err = hex.DecodeString(cleanSignature)
		if err != nil {
			return false, "", fmt.Errorf("ошибка декодирования подписи: %w", err)
		}
	}

	// Проверка длины подписи
	if len(signatureBytes) != 65 {
		return false, "", fmt.Errorf("недопустимая длина подписи: %d", len(signatureBytes))
	}

	// Создаем хеш сообщения с Ethereum-префиксом
	messageBytes := []byte(message)
	prefixedMessage := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(messageBytes), message)
	messageHash := crypto.Keccak256([]byte(prefixedMessage))

	// Извлекаем r, s, v из подписи
	r := signatureBytes[:32]
	ss := signatureBytes[32:64]
	v := signatureBytes[64]

	// Для метода Ecrecover нужно v = 0 или v = 1
	normalizedV := byte(0)
	if v >= 27 {
		normalizedV = v - 27
	} else {
		normalizedV = v
	}

	// Создаем новую сигнатуру с корректным v
	recoveryID := normalizedV
	if recoveryID != 0 && recoveryID != 1 {
		// Если значение не стандартное, пробуем разные варианты
		// Возможно это подпись с EIP-155 (chainId)
		chainId := uint64((v - 35) / 2)
		if chainId > 0 && chainId <= 38 { // Проверка на разумные пределы chainId
			normalizedV = byte((v - 35) % 2)
			recoveryID = normalizedV
		} else {
			// Пробуем простое преобразование
			recoveryID = v % 2
		}
	}

	// Используем go-ethereum низкоуровневую функцию secp256k1.RecoverPubkey
	// Этот метод более прямолинейный, чем Ecrecover
	pubKeyBytes, err := crypto.SigToPub(messageHash, append(append(r, ss...), recoveryID))
	if err != nil {
		return false, "", fmt.Errorf("ошибка восстановления публичного ключа: %w", err)
	}

	// Получаем адрес из публичного ключа
	recoveredAddress := crypto.PubkeyToAddress(*pubKeyBytes).Hex()
	return true, recoveredAddress, nil
}

// SaveAudioMetadata сохраняет метаданные аудио в базу данных
func (s *Service) SaveAudioMetadata(ctx context.Context, audio *model.Audio) (int64, error) {
	// Здесь должен быть код для сохранения в базу данных
	// В простом случае можно использовать Repository для сохранения
	return 0, nil
}
