package main

import (
	"encoding/base64"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/basicauth"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

type Product struct {
	gorm.Model
	Name   string
	Stock  int
	Detail string
}

type Order struct {
	gorm.Model
	OrderDate time.Time
	Orderer   string
	Amount    int
	ProductID uint
	Product   Product
}

func extractUserId(c fiber.Ctx) (string, error) {
	authHeader := c.Get(fiber.HeaderAuthorization)

	if !strings.HasPrefix(authHeader, "Basic ") {
		return "", errors.New("Missing or Invalid auth header")
	}

	encoded := strings.TrimPrefix(authHeader, "Basic ")

	payload, err := base64.StdEncoding.DecodeString(encoded)

	if err != nil {
		return "", errors.New("Invalid auth header.")
	}

	pair := strings.SplitN(string(payload), ":", 2)

	userId := pair[0]

	return userId, nil
}

func migrate(db *gorm.DB) {
	var exists bool
	checkTypeSQL := `SELECT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'user_type')`

	if err := db.Raw(checkTypeSQL).Scan(&exists).Error; err != nil {
		log.Fatal("Cannot check enum existence: ", err)
	}

	if !exists {
		createTypeSQL := `CREATE TYPE public."user_type" AS ENUM ('admin', 'user');`
		if err := db.Exec(createTypeSQL).Error; err != nil {
			log.Fatal("Failed to create USER_TYPE enum", err)
		}
		log.Println("user_type enum created.")
	}

	if err := db.AutoMigrate(&Order{}, &Product{}); err != nil {
		log.Fatal("AutoMigrate failed:", err)
	}

	seed_products := []Product{
		{Name: "신발", Stock: 100, Detail: "예쁜 신발이에요."},
		{Name: "양말", Stock: 100, Detail: "예쁜 양말이에요."},
		{Name: "버선", Stock: 100, Detail: "예쁜 버선이에요."},
	}

	for _, product := range seed_products {
		err := db.Where(Product{Name: product.Name}).FirstOrCreate(&product).Error
		if err != nil {
			log.Printf("Could not seed item")
		}
	}

	log.Println("Migration completed.")
}

func main() {
	dsn := "host=localhost user=postgres password=password dbname=flash-sale-coding-practice port=5432 sslmode=disable TimeZone=Asia/Seoul"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	migrate(db)

	if err != nil {
		log.Fatal("DB Connection failure: exiting.")
	}

	app := fiber.New()

	authConfig := basicauth.Config{
		Users: map[string]string{
			"admin":  "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user1":  "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user2":  "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user3":  "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user4":  "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user5":  "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user6":  "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user7":  "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user8":  "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user9":  "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user10": "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
		},
	}

	api := app.Group("/api")

	products := api.Group("/products")
	transactions := api.Group("/transactions")

	products.Get("/", func(c fiber.Ctx) error {
		product, _ := gorm.G[Product](db).Find(c)

		return c.JSON(product)
	})

	products.Get("/:id", func(c fiber.Ctx) error {
		product, _ := gorm.G[Product](db).Where("Id = ?", c.Params("id")).First(c.Context())

		return c.JSON(product)
	})

	transactions.Get("/buy/:id", basicauth.New(authConfig), func(c fiber.Ctx) error {
		return c.SendString("WIP")
	})

	transactions.Get("/log", basicauth.New(authConfig), func(c fiber.Ctx) error {
		userId, _ := extractUserId(c)

		return c.SendString(userId)
	})

	log.Fatal(app.Listen(":8000"))
}
