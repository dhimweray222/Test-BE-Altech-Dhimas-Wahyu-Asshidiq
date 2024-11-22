package main

import (
	"log"
	"os"
	"test-backend-altech/utils"
	"time"

	"test-backend-altech/config"
	"test-backend-altech/controller"
	"test-backend-altech/repository"
	"test-backend-altech/repository/query"
	"test-backend-altech/service"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"

	_ "github.com/joho/godotenv/autoload"
)

var logger = utils.NewLogger()

func main() {

	time.Local = time.UTC
	serverConfig := config.NewServerConfig()
	db := config.NewPostgresDatabase()
	store := repository.NewStore(db)
	validate := validator.New()
	authorQuery := query.NewAuthor()
	authorRepository := repository.NewAuthorRepository(store, authorQuery)
	authorMemberService := service.NewAuthorService(authorRepository)
	authorController := controller.NewAuthorController(validate, authorMemberService)

	bookQuery := query.NewBook()
	bookRepository := repository.NewBookRepository(store, bookQuery)

	cache, errCache := config.NewRedisCache(&config.RedisConfig{
		Host: os.Getenv("REDIS_HOST"),
	})
	bookMemberService := service.NewBookService(bookRepository, cache)
	bookController := controller.NewBookController(validate, bookMemberService, cache)

	if errCache != nil {
		log.Fatalf("Failed to connect to cache: %v", errCache)
	}

	app := fiber.New(fiber.Config{BodyLimit: 10 * 1024 * 1024})
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "*",
		AllowHeaders:     "*",
		AllowCredentials: false,
	}))

	authorController.Route(app)
	bookController.Route(app)
	err := app.Listen(serverConfig.Host)
	if err != nil {
		log.Fatal(err)
		logger.Error(err)
	}
}
