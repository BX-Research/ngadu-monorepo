package routes

import (
	"time"

	"github.com/BX-Research/ngadu-monorepo/backend/app/services"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/limiter"
)

// func registerUserRoutes(app *fiber.App, lookup func(context.Context, string) (userInfo, error)) {
// 	app.Get("/api/users/:id", func(c fiber.Ctx) error {
// 		c.Set("Cache-Control", "no-store")
// 		fail := func(code int, message string) error {
// 			return c.Status(code).JSON(fiber.Map{"error": message}, fiber.MIMEApplicationJSON)
// 		}
// 		id := c.Params("id")
// 		if !userIDPattern.MatchString(id) {
// 			return fail(fiber.StatusBadRequest, "invalid user ID")
// 		}
// 		ctx, cancel := context.WithTimeout(c.Context(), 3*time.Second)
// 		defer cancel()
// 		user, err := lookup(ctx, id)
// 		if errors.Is(err, pgx.ErrNoRows) {
// 			return fail(fiber.StatusNotFound, "user not found")
// 		}
// 		if err != nil {
// 			return fail(fiber.StatusServiceUnavailable, "user lookup unavailable")
// 		}
// 		return c.JSON(user, fiber.MIMEApplicationJSON)
// 	})
// 	app.All("/api/users/:id", methodNotAllowed)
// }

func UserRoutes(app fiber.Router, users services.UserFinder) {
	userLimiter := limiter.New(
		limiter.Config{
			Max:        10,
			Expiration: 1 * time.Minute,
		})

	r := app.Group("/users")

	r.Use(userLimiter)

	r.Get("/:id", services.FindUser(users))
}
