package main

import (
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
)

func createResponseError(c *fiber.Ctx, code int, message string) error {

	errorRes, _ := json.Marshal(Response{
		Message: message,
		Code:    code,
		Data:    nil,
	})

	return c.SendString(string(errorRes))
}

func createResponseSuccess(c *fiber.Ctx, data interface{}, arg ...string) error {
	message := "Success"
	if len(arg) != 0 && arg[0] != "" {
		message = arg[0]
	}
	res, _ := json.Marshal(Response{
		Message: message,
		Code:    fiber.StatusOK,
		Data:    data,
	})

	return c.SendString(string(res))
}
