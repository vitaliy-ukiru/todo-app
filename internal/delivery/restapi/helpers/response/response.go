package response

import "github.com/gofiber/fiber/v2"

// Response is type for response from RESTAPI.
// In Result must be only success result.
// Field Error for error message.
type Response struct {
	Ok     bool   `json:"ok"`
	Result any    `json:"result,omitempty"`
	Error  string `json:"error,omitempty"`
}

// JustOK is empty success response. It creates for reusing memory resources.
var JustOK = Response{Ok: true}

// Write writes response to fiber.Ctx.
func Write(c *fiber.Ctx, code int, resp Response) error {
	return c.Status(code).JSON(resp)
}

// Ok writes success response.
func Ok(c *fiber.Ctx, code int, result any) error {
	return Write(c, code, Response{Ok: true, Result: result})
}

// Err writes error response.
func Err(c *fiber.Ctx, code int, err string) error {
	return Write(c, code, Response{Error: err})
}

// WithError writes response to `c` with code `code` and error message from `err`.
func WithError(c *fiber.Ctx, code int, err error) error {
	return Err(c, code, err.Error())
}

func Wrap(c *fiber.Ctx, code int, msg string, err error) error {
	return c.Status(code).JSON(fiber.Map{
		"ok": false,
		"error": fiber.Map{
			"msg": msg,
			"err": err,
		},
	})
}
