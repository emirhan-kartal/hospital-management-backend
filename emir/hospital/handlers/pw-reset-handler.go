package handlers

import (
	"emir/hospital/models"
	"emir/hospital/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type PWHandler struct {
	Redis *redis.Client
	DB    *gorm.DB
}
type ResetPasswordInitiateData struct {
	Tel_no string `json:"tel_no"`
}

// @Summary Initiate password reset
// @Description Send a reset code to the user's phone number
// @Tags password-reset
// @Accept json
// @Produce json
// @Param phone body ResetPasswordInitiateData true "Phone number"
// @Success 200 {string} string "Code sent to your phone number. Code(For Development):<code>"
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Router /reset-password/initiate [post]
// @Security BearerAuth
func (h *PWHandler) ResetPasswordInitiate(c *fiber.Ctx) error {
	var phone ResetPasswordInitiateData
	if err := c.BodyParser(&phone); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error()})
	}
	if c.Cookies("reset_password") != "" {
		return c.SendStatus(fiber.StatusConflict)
	}

	var dbUser models.User
	if err := h.DB.Where("tel_no = ?", phone.Tel_no).First(&dbUser).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	sixDigitCode := utils.RandCode(6)

	ctx := c.Context()
	h.Redis.Set(ctx, phone.Tel_no, sixDigitCode, time.Minute*3)
	c.Cookie(&fiber.Cookie{
		Name:    "reset_password",
		Value:   string(phone.Tel_no),
		Expires: time.Now().Add(time.Minute * 3),
	},
	)
	return c.Send([]byte("Code sent to your phone number. Code(For Development):" + sixDigitCode))
}

type ResetPasswordFinalizeCode struct {
	Code   string `json:"code"`
	Tel_no string `json:"tel_no"`
}

// @Summary Finalize password reset
// @Description Validate the reset code and get a token to change the password
// @Tags password-reset
// @Accept json
// @Produce json
// @Param resetPassword body ResetPasswordFinalizeCode true "Reset code and phone number"
// @Success 200 {string} string "Password reset successful. go to /change-password with this token:<token>"
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /reset-password/finalize [post]
// @Security BearerAuth
func (h *PWHandler) ResetPasswordFinalize(c *fiber.Ctx) error {

	var resetPassword ResetPasswordFinalizeCode
	if err := c.BodyParser(&resetPassword); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	ctx := c.Context()
	redisCode, err := h.Redis.Get(ctx, resetPassword.Tel_no).Result()
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Code not found",
		})
	}
	if redisCode != resetPassword.Code {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid code",
		})
	}
	randomString := utils.RandString(30)
	h.Redis.Set(ctx, randomString, resetPassword.Tel_no, time.Minute*3)
	return c.SendString("Password reset successful. go to /change-password with this token:" + randomString)
}

type ResetPasswordData struct {
	Password       string `json:"password"`
	RepeatPassword string `json:"repeat_password"`
	ValidateCode   string `json:"validate_code"`
}

// @Summary Reset password
// @Description Change the user's password using the validation token
// @Tags password-reset
// @Accept json
// @Produce json
// @Param data body ResetPasswordData true "Password reset data"
// @Success 200 {string} string "Password changed successfully"
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /reset-password [post]
// @Security BearerAuth
func (h *PWHandler) ResetPassword(c *fiber.Ctx) error {
	var data ResetPasswordData
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	if data.Password != data.RepeatPassword {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Passwords do not match",
		})
	}
	ctx := c.Context()
	phone, err := h.Redis.Get(ctx, data.ValidateCode).Result()
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Code not found",
		})
	}
	var dbUser models.User
	if err := h.DB.Where("tel_no = ?", phone).First(&dbUser).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{

			"error": err.Error(),
		})
	}
	newHashedPw := string(hashedPassword)
	if dbUser.Password == newHashedPw {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "New password cannot be the same as the old password",
		})
	}
	dbUser.Password = newHashedPw
	dbUser.TokenVersion++
	if err := h.DB.Save(&dbUser).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	h.Redis.Del(ctx, data.ValidateCode)
	h.Redis.Del(ctx, phone)

	return c.SendString("Password changed successfully")
}
