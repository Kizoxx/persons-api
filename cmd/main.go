package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	_ "net/http"
	"persons-api/internal/api"
	"persons-api/internal/config"
	"persons-api/internal/logger"
	"persons-api/internal/models"
	"persons-api/internal/repository"
	"persons-api/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "persons-api/cmd/docs" // Пустой импорт для Swagger-документации
)

// @title Persons API
// @version 1.0
// @description This is a sample Persons API.
// @host localhost:8080
// @BasePath /api/v1
func main() {
	// 1. Загружаем конфигурацию
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(fmt.Errorf("cannot load config: %w", err))
	}

	// 2. Инициализация логгера
	logger.Init(cfg.LogLevel)

	// 3. Подключаемся к базе данных
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logrus.Fatalf("Не удалось подключиться к БД: %v", err)
	}

	// Миграция схемы (создание таблицы persons)
	if err := db.AutoMigrate(&models.Person{}); err != nil {
		logrus.Fatalf("Ошибка миграции: %v", err)
	}

	// 4. Создаём репозиторий и сервис
	personRepo := repository.NewPersonRepository(db)
	personService := services.NewPersonService(personRepo)

	// 5. Инициализируем Gin-роутер
	router := gin.New()
	router.Use(gin.Recovery()) // обработка паник

	// Добавляем Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 6. Регистрируем контроллеры
	api.NewPersonHandler(router, personService)

	// 7. Запускаем сервер
	addr := fmt.Sprintf("%s:%d", cfg.ServerHost, cfg.ServerPort)
	logrus.Infof("Запуск сервера на %s", addr)
	if err := router.Run(addr); err != nil {
		logrus.Fatalf("Не удалось запустить сервер: %v", err)
	}
}
