package admin

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"template/internal/models"
	"template/pgk/memcache"
	"time"
)

func CreateUser(c *fiber.Ctx) error {
	models.DB.Debug().Create(&models.UserModel{
		Model:    gorm.Model{},
		Username: c.FormValue("username"),
		Password: c.FormValue("password"),
		Status:   c.FormValue("status"),
	})

	go memcache.UserCache.Fetch()

	return c.Redirect("/admin/users")
}

func DeleteUser(c *fiber.Ctx) error {
	user := c.Params("user")

	userModel := memcache.UserCache.Get(user)
	models.DB.Delete(&userModel)

	go memcache.UserCache.Fetch()

	return c.Redirect("/admin/users")
}

func LoginIn(c *fiber.Ctx) error {
	login, password := c.FormValue("Username"), c.FormValue("Password")

	userFounded := memcache.UserCache.Get(login)
	if userFounded.Username == "" || userFounded.Password != password {
		return c.Redirect("/")
	}

	c.Cookie(&fiber.Cookie{Name: "USR", Value: c.FormValue("Username"), Expires: time.Now().Add(48 * time.Hour)})
	c.Cookie(&fiber.Cookie{Name: "Authenticated", Value: "true", Expires: time.Now().Add(48 * time.Hour)})

	return c.Redirect("/admin")
}
