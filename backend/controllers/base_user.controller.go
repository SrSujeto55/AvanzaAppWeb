package controllers

import (
	"backend/models"
	"backend/tools"
	"backend/tools/db"
	va "backend/tools/validation"

	"github.com/gofiber/fiber/v2"
)


func GetBaseUser(c *fiber.Ctx) error {
	isLogged, id := tools.GetCurrentUserId(c)
	if !isLogged {
		return ApiError(c, "Not logged in", 400)
	}


	var user models.BaseUser
	err := db.Orm().First(&user, id).Error
	if err != nil {
		return ApiError(c, "DB error", 500)
	}
	return ApiSuccess(c, fiber.Map{
		"alias":       user.Alias,
		"birthDate":    user.BirthDate,
		"phone":       user.PhoneNumber,
		"photo":       user.Photo,
		"description": user.Description,
	})
}

func UpdateBaseUser(c *fiber.Ctx) error {
	isLogged, id := tools.GetCurrentUserId(c)
	if !isLogged {
		return ApiError(c, "Not logged in", 400)
	}

	vaErr := va.Check(c, va.Rmap {
		"alias":     "required",
		"phone":     "required",
	})
	if vaErr != nil {
		return ApiError(c, "Wrong data", 400)
	}

	var user models.BaseUser
	getUserErr := db.Orm().First(&user, id).Error
	if getUserErr != nil {
		return ApiError(c, "DB error", 500)
	}

	user.Alias = c.FormValue("alias")
	user.PhoneNumber = c.FormValue("phone")
	user.Description = c.FormValue("description")

	if db.Orm().Save(&user).Error != nil {
		return ApiError(c, "DB error", 500)
	}

	return ApiSuccess(c, nil)
}

func UpdateBaseUserPhoto(c *fiber.Ctx) error {
	isLogged, id := tools.GetCurrentUserId(c)
	if !isLogged {
		return ApiError(c, "Not logged in", 400)
	}

	vaErr := va.Check(c, va.Rmap {
		"photo": "required",
	})
	if vaErr != nil {
		return ApiError(c, "Wrong data", 400)
	}

	var user models.BaseUser
	getUserErr := db.Orm().First(&user, id).Error
	if getUserErr != nil {
		return ApiError(c, "DB error", 500)
	}

	user.Photo = c.FormValue("photo")

	if db.Orm().Save(&user).Error != nil {
		return ApiError(c, "DB error", 500)
	}

	return ApiSuccess(c, nil)
}

func DeleteBaseUser(c *fiber.Ctx) error {
	isLogged, id := tools.GetCurrentUserId(c)
	if !isLogged {
		return ApiError(c, "Not logged in", 400)
	}

	if db.Orm().Delete(&models.BaseUser{}, id).Error != nil {
		return ApiError(c, "Db error", 500)
	}
	return ApiSuccess(c, fiber.Map{"id": id})
}


// For a given base user returns "trainer" if this base
// user belongs to a trainer, returns "athlete" otherwise.
func UserKind(u models.BaseUser) string {
	var trainer models.Trainer
	trainerErr := db.Orm().Where("base_user_id = ?", u.ID).First(&trainer).Error

	if trainerErr == nil {
		return "trainer"
	} else {
		return "athlete"
	}
}

