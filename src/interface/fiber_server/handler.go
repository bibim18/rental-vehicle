package fiber_server

import (
	"context"

	"rental-vehicle-system/src/entity/price_model"

	"github.com/cockroachdb/errors"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
)

const contextLocalKey = "fiber-otel-tracer"
const OK string = "OK"

type ErrorResponse struct {
	Error     string `json:"error"`
	ErrorCode int    `json:"error_code"`
	IssueId   string `json:"issue_id"`
}

type ErrorMapInfo struct {
	StatusCode int
	ErrorCode  int
}

var errorList = map[error]ErrorMapInfo{
	price_model.ErrInvalidVehicle: {400, 1},
}

func (f FiberServer) sendError(c *fiber.Ctx, status int, err error, errCode int, issueId string) error {
	c.Locals("error", err)

	return c.Status(status).JSON(ErrorResponse{
		Error:     err.Error(),
		ErrorCode: errCode,
		IssueId:   issueId,
	})
}

func (f FiberServer) errorHandler(c *fiber.Ctx, err error, issueIds ...string) error {
	unwrapErr := errors.Unwrap(err)
	if unwrapErr == nil {
		unwrapErr = err
	}

	var issueId string
	if len(issueIds) > 0 {
		issueId = issueIds[0]
	}

	for iErr, code := range errorList {
		if errors.Is(err, iErr) {
			return f.sendError(c, code.StatusCode, err, code.ErrorCode, issueId)
		}
	}

	return f.sendError(c, 500, err, 0, issueId)
}

func getSpanContext(c *fiber.Ctx) context.Context {
	ctx, ok := c.Locals(contextLocalKey).(context.Context)

	if !ok {
		log.Warn("Failed to get span context from fiber context")
		return c.Context()
	}

	return ctx
}
