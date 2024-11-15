package internal

import (
	"context"
	"log"
	"time"
	"os"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"quiz.com/quiz/internal/collection"
	"quiz.com/quiz/internal/controller"
	"quiz.com/quiz/internal/service"
	"quiz.com/quiz/internal/middleware"
	"github.com/joho/godotenv"
)

type App struct {
	httpServer  *fiber.App
	database    *mongo.Database
	secretKey   []byte

	quizService *service.QuizService
	netService  *service.NetService
	userService *service.UserService
}

func (a *App) Init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading environment variables")
	}

	a.secretKey = []byte(os.Getenv("JWT_SECRET"))
	if len(a.secretKey) == 0 {
		log.Fatal("JWT_SECRET is not set in the environment")
	}

	dbURI := os.Getenv("DB_STRING")

	a.setupDb(dbURI)
	a.setupServices()
	a.setupHttp()

	log.Fatal(a.httpServer.Listen(":3000"))
}

func (a *App) setupHttp() {
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "https://vestibulados-front-qwvn.vercel.app", // No trailing slash
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true, // Allow credentials like cookies or Authorization headers
	}))

	quizController := controller.Quiz(a.quizService)
	userController := controller.User(a.userService)
	wsController := controller.Ws(a.netService)

	// Auth Middleware
	app.Use("/api/protected", middleware.JWTAuthMiddleware)

	// User Routes
	app.Get("/api/protected/me", userController.DetailUser)
	app.Post("/api/register", userController.Register)
	app.Post("/api/login", userController.Login)
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"teste": "teste",
		})
	// Quiz Routes
	app.Post("/api/quiz/create", quizController.CreateQuiz)

	// Websocket Routes
	app.Get("/ws", websocket.New(wsController.Ws))

	a.httpServer = app
}

func (a *App) setupServices() {
	a.quizService = service.Quiz(collection.Quiz(a.database.Collection("quizzes")))
	a.netService = service.Net(a.quizService)
	a.userService = service.User(collection.User(a.database.Collection("users"), a.secretKey))
}

func (a *App) setupDb(dbURI string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbURI))
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	a.database = client.Database("quiz")
}
