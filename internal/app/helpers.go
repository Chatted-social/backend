package app

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"strconv"
)

func Ok() map[string]bool {
	return map[string]bool{"ok": true}
}

func Err(msg string) map[string]string {
	return map[string]string{"error": msg}
}

func StringSliceToInt(slice []string) (r []int) {
	for _, s := range slice {
		if i, err := strconv.Atoi(s); err == nil {
			r = append(r, i)
		}
	}
	return
}

func NewFiber() *fiber.App {
	app := fiber.New()

	app.Use(cors.New())
	app.Use(logger.New())

	return app
}