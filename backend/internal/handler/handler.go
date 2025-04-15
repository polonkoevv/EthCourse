package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/polonkoevv/ethcourse/internal/service"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins: []string{"https://foo.com"}, // Используйте это для продакшена
		AllowedOrigins: []string{"http://localhost:5173"}, // Разрешаем запросы с вашего Vue-сервера
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Максимальное время (в секундах) кеширования результатов preflight-запросов
	}))
	r.Post("/upload", h.UploadFile)
	r.Get("/music", h.GetAllMusic)
	return r
}

func (h *Handler) UploadFile(w http.ResponseWriter, r *http.Request) {
	// Максимальный размер файла (32 МБ)
	r.ParseMultipartForm(32 << 20)

	// Получение файла из формы
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Ошибка получения файла: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Проверка типа файла
	// Читаем первые 512 байт для определения MIME-типа
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		http.Error(w, "Ошибка чтения файла: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Сбрасываем позицию чтения в начало файла
	_, err = file.Seek(0, 0)
	if err != nil {
		http.Error(w, "Ошибка при обработке файла: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Определяем MIME-тип
	fileType := http.DetectContentType(buffer)

	// Проверяем, является ли файл аудиофайлом
	isAudioFile := false

	// Проверка по MIME-типу
	audioMimeTypes := []string{
		"audio/mpeg", "audio/mp3", "audio/wav", "audio/wave",
		"audio/ogg", "audio/flac", "audio/aac", "audio/mp4",
		"audio/x-m4a", "audio/webm", "audio/x-matroska",
	}

	for _, mimeType := range audioMimeTypes {
		if strings.HasPrefix(fileType, mimeType) {
			isAudioFile = true
			break
		}
	}

	// Если MIME-тип не определен как аудио, проверяем расширение файла
	if !isAudioFile {
		ext := strings.ToLower(filepath.Ext(handler.Filename))
		audioExtensions := []string{
			".mp3", ".wav", ".ogg", ".flac", ".aac",
			".m4a", ".wma", ".opus", ".webm", ".mka",
		}

		for _, audioExt := range audioExtensions {
			if ext == audioExt {
				isAudioFile = true
				break
			}
		}
	}

	// Если это не аудиофайл, возвращаем ошибку
	if !isAudioFile {
		http.Error(w, "Недопустимый тип файла. Разрешены только аудиофайлы.", http.StatusBadRequest)
		return
	}

	fmt.Printf("Загружен аудиофайл: %s, тип: %s\n", handler.Filename, fileType)

	// Временное сохранение файла
	tempFilePath := filepath.Join(os.TempDir(), handler.Filename)
	tempFile, err := os.Create(tempFilePath)
	if err != nil {
		http.Error(w, "Ошибка создания временного файла: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer tempFile.Close()

	_, err = io.Copy(tempFile, file)
	if err != nil {
		http.Error(w, "Ошибка сохранения файла: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Добавление файла в IPFS
	tempFile.Seek(0, 0)
	cid, err := h.service.UploadFile(context.Background(), handler.Filename, tempFile)
	if err != nil {
		http.Error(w, "Ошибка добавления в IPFS: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Файл успешно загружен в IPFS с CID: %s", cid)
}

func (h *Handler) GetAllMusic(w http.ResponseWriter, r *http.Request) {
	music, err := h.service.GetAllMusic(context.Background())
	if err != nil {
		http.Error(w, "Ошибка получения музыки: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(music)
}
