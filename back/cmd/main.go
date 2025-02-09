package main

import (
	"database/sql"
	"fmt"
	"log"

	"net/http"

	"back/config"
	myHttp "back/delivery/http"
	"back/repository"
	"back/usecase"

	_ "github.com/lib/pq"
)

func main() {
	// Загружаем конфигурацию
	cfg := config.LoadConfig()

	// Подключение к базе данных
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()

	fmt.Println(cfg)

	// Инициализация репозитория, юзкейса и маршрутов
	repo := repository.NewEntityRepo(db)
	api := usecase.NewUsecase(repo)

	mux := http.NewServeMux()
	myHttp.SetupRoutes(api, mux)

	// Получаем порт из конфигурации или переменной окружения
	port := cfg.Port
	if port == "" {
		port = "8080"
	}

	// Запуск HTTP-сервера
	log.Printf("Server is running on port %s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
