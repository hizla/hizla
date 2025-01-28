package auth

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Handler struct {
	DB *pgxpool.Pool
}

func NewAuthHandler(db *pgxpool.Pool) *Handler {
	return &Handler{DB: db}
}

var store = InitStore()

func InitStore() *session.Store {
	store := session.New()
	return store
}

func (h *Handler) SignUp(c *fiber.Ctx) error {
	user := new(User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not hash password"})
	}

	_, err = h.DB.Exec(context.Background(),
		"INSERT INTO users(email, password) VALUES($1, $2)",
		user.Email, string(hashedPassword))
	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "User already exists"})
	}

	return c.JSON(fiber.Map{"message": "User created successfully"})
}

func (h *Handler) SignIn(c *fiber.Ctx) error {
	session, _ := store.Get(c)
	var input User
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	var user User
	err := h.DB.QueryRow(context.Background(),
		"SELECT id, email, password FROM users WHERE email = $1", input.Email).
		Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	session.Set("user_id", user.ID)
	if err := session.Save(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save session"})
	}

	return c.JSON(fiber.Map{"message": "Logged in successfully"})
}

func (h *Handler) Logout(c *fiber.Ctx) error {
	session, _ := store.Get(c)
	session.Delete("user_id")
	if err := session.Save(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save session"})
	}
	return c.JSON(fiber.Map{"message": "Logged out successfully"})
}

func (h *Handler) Profile(c *fiber.Ctx) error {
	session, _ := store.Get(c)
	userID := session.Get("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var user User
	err := h.DB.QueryRow(context.Background(),
		"SELECT id, email FROM users WHERE id = $1", userID).
		Scan(&user.ID, &user.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "User not found"})
	}

	return c.JSON(user)
}
