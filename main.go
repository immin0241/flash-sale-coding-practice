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

type BuyForm struct {
	ProductID int `json:"product_id"`
	Amount    int `json:"amount"`
}

type TransactionErr struct {
	Code    int
	Message string
	Err     error
}

func (e *TransactionErr) Error() string {
	return e.Message
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

	sqlDB, err := db.DB()

	sqlDB.SetMaxIdleConns(100)
	sqlDB.SetMaxOpenConns(200)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)
	sqlDB.SetConnMaxIdleTime(5 * time.Minute)

	migrate(db)

	if err != nil {
		log.Fatal("DB Connection failure: exiting.")
	}

	app := fiber.New()

	authConfig := basicauth.Config{
		Users: map[string]string{
			"admin":   "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user1":   "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user2":   "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user3":   "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user4":   "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user5":   "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user6":   "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user7":   "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user8":   "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user9":   "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user10":  "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user11":  "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user12":  "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user13":  "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user14":  "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user15":  "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user16":  "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user17":  "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user18":  "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user19":  "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user20":  "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user21":  "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user22":  "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user23":  "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user24":  "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user25":  "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user26":  "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user27":  "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user28":  "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user29":  "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user30":  "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user31":  "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user32":  "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user33":  "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user34":  "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user35":  "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user36":  "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user37":  "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user38":  "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user39":  "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user40":  "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user41":  "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user42":  "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user43":  "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user44":  "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user45":  "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user46":  "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user47":  "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user48":  "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user49":  "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user50":  "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user51":  "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user52":  "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user53":  "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user54":  "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user55":  "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user56":  "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user57":  "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user58":  "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user59":  "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user60":  "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user61":  "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user62":  "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user63":  "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user64":  "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user65":  "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user66":  "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user67":  "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user68":  "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user69":  "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user70":  "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user71":  "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user72":  "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user73":  "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user74":  "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user75":  "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user76":  "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user77":  "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user78":  "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user79":  "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user80":  "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user81":  "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user82":  "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user83":  "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user84":  "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user85":  "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user86":  "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user87":  "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user88":  "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user89":  "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user90":  "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user91":  "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user92":  "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user93":  "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user94":  "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user95":  "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user96":  "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user97":  "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user98":  "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user99":  "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
			"user100": "$2a$12$5ZRgwbTLZJDaE7itOdFLiORvgtVaCEgbqH5f8WXLlJckZWrc2UwRi",
		},
	}

	api := app.Group("/api")

	products := api.Group("/products")
	transactions := api.Group("/transactions")

	products.Get("/", func(c fiber.Ctx) error {
		product, _ := gorm.G[Product](db).Find(c)

		return c.JSON(product)
	})

	products.Get("/:id", basicauth.New(authConfig), func(c fiber.Ctx) error {
		userId, err := extractUserId(c)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Internal Server Error",
			})
		}

		if userId != "admin" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Unauthenticated access.",
			})
		}

		product, _ := gorm.G[Product](db).Where("id = ?", c.Params("id")).Find(c.Context())

		return c.JSON(product)
	})

	transactions.Post("/buy", basicauth.New(authConfig), func(c fiber.Ctx) error {
		payload := new(BuyForm)

		if err := c.Bind().Body(payload); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid JSON.",
			})
		}

		userId, err := extractUserId(c)

		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Cannot authenticate.",
			})
		}

		var transactionResult Order
		var existingOrder Order

		err = db.Transaction(func(tx *gorm.DB) error {
			result := tx.Model(&Product{}).
				Where("id = ? AND stock - ? >= 0", payload.ProductID, payload.Amount).
				Update("stock", gorm.Expr("stock - ?", payload.Amount))

			if result.Error != nil {
				return &TransactionErr{Code: fiber.StatusInternalServerError, Message: "서버 오류가 발생했습니다.", Err: result.Error}
			}
			if result.RowsAffected == 0 {
				return &TransactionErr{Code: fiber.StatusGone, Message: "올바르지 않은 상품이거나 재고가 없습니다.", Err: result.Error}
			}

			if err := tx.Where("orderer = ? AND product_id = ?", userId, payload.ProductID).Limit(1).Find(&existingOrder).Error; err != nil {
				return &TransactionErr{Code: fiber.StatusInternalServerError, Message: "DB 조회 중 오류가 발생했습니다.", Err: err}
			}

			if existingOrder.ID != 0 {
				return &TransactionErr{Code: fiber.StatusBadRequest, Message: "1인 1개만 구매할 수 있습니다.", Err: nil}
			}

			var updatedProduct Product
			if err := tx.First(&updatedProduct, payload.ProductID).Error; err != nil {
				return err
			}

			newTx := Order{
				OrderDate: time.Now(),
				Orderer:   userId,
				ProductID: uint(payload.ProductID),
				Amount:    payload.Amount,
				Product:   updatedProduct,
			}

			if err := tx.Create(&newTx).Error; err != nil {
				return &TransactionErr{Code: fiber.StatusInternalServerError, Message: "서버 오류가 발생했습니다.", Err: result.Error}
			}

			transactionResult = newTx
			return nil
		})

		if err != nil {
			var transErr *TransactionErr

			if errors.As(err, &transErr) {
				return c.Status(transErr.Code).JSON(fiber.Map{
					"error": transErr.Message,
				})
			}

			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Internal Server Error",
			})
		}

		return c.JSON(transactionResult)
	})

	transactions.Get("/log", basicauth.New(authConfig), func(c fiber.Ctx) error {
		userId, err := extractUserId(c)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Internal Server Error",
			})
		}

		orders, _ := gorm.G[Order](db).Where("orderer = ?", userId).Find(c.Context())

		return c.JSON(orders)
	})

	log.Fatal(app.Listen(":8000"))
}
