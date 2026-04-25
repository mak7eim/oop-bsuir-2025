package response

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	apperrors "lab7-oop/shared/errors"
)

type Envelope struct {
	Data  any                  `json:"data,omitempty"`
	Error *apperrors.AppError `json:"error,omitempty"`
}

func OK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Envelope{Data: data})
}

func Created(c *gin.Context, data any) {
	c.JSON(http.StatusCreated, Envelope{Data: data})
}

func Error(c *gin.Context, err error) {
	var appErr *apperrors.AppError
	if errors.As(err, &appErr) {
		c.JSON(statusFromCode(appErr.Code), Envelope{Error: appErr})
		return
	}

	switch {
	case errors.Is(err, apperrors.ErrNotFound):
		c.JSON(http.StatusNotFound, Envelope{Error: apperrors.From(err)})
	case errors.Is(err, apperrors.ErrUnauthorized), errors.Is(err, apperrors.ErrInvalidCredentials):
		c.JSON(http.StatusUnauthorized, Envelope{Error: apperrors.From(err)})
	case errors.Is(err, apperrors.ErrForbidden):
		c.JSON(http.StatusForbidden, Envelope{Error: apperrors.From(err)})
	case errors.Is(err, apperrors.ErrConflict):
		c.JSON(http.StatusConflict, Envelope{Error: apperrors.From(err)})
	default:
		c.JSON(http.StatusBadRequest, Envelope{Error: apperrors.From(err)})
	}
}

func statusFromCode(code string) int {
	switch code {
	case "not_found":
		return http.StatusNotFound
	case "unauthorized", "invalid_credentials":
		return http.StatusUnauthorized
	case "forbidden":
		return http.StatusForbidden
	case "conflict":
		return http.StatusConflict
	default:
		return http.StatusBadRequest
	}
}
