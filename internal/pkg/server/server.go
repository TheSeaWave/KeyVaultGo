package server

import (
	"log"
	"net/http"

	"dz1/internal/pkg/storage" // замените на фактический путь к storage.go

	"github.com/gin-gonic/gin"
)

// Импортируйте интерфейс StorageInterface вместо конкретного типа
type Server struct {
	router  *gin.Engine
	storage storage.StorageInterface // Используем интерфейс
}

// NewServer создает новый сервер с указанным хранилищем
func NewServer(s storage.StorageInterface) *Server { // Используем интерфейс
	server := &Server{
		router:  gin.Default(),
		storage: s,
	}
	server.routes()
	return server
}

// routes задает маршруты сервера
func (s *Server) routes() {
	s.router.GET("/health", s.handleHealth)
	s.router.GET("/scalar/get/:key", s.handleGetScalar)
	s.router.PUT("/scalar/set/:key", s.handleSetScalar)
}

// Health endpoint для проверки состояния сервера
func (s *Server) handleHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}

// GET /scalar/get/:key - получение значения скаляра по ключу
func (s *Server) handleGetScalar(c *gin.Context) {
	key := c.Param("key")
	value := s.storage.Get(key)
	if value == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "key not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"key": key, "value": *value})
}

// PUT /scalar/set/:key - установка значения скаляра по ключу
func (s *Server) handleSetScalar(c *gin.Context) {
	key := c.Param("key")

	// Парсим JSON body для получения значения
	var body struct {
		Value string `json:"Value"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// Устанавливаем значение в хранилище
	s.storage.Set(key, body.Value)
	c.Status(http.StatusOK)
}

// Run запускает сервер на указанном адресе
func (s *Server) Run(addr string) error {
	log.Printf("Starting server on %s", addr)
	return s.router.Run(addr)
}
