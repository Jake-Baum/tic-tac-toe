package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
)

type parseError struct {
	message string
}

func (p parseError) Error() string {
	return p.message
}

func CreateMessageResponse(message string, args ...any) gin.H {
	return gin.H{"message": fmt.Sprintf(message, args...)}
}

func IntPathParam(c *gin.Context, name string) (int, error) {
	s := c.Param(name)

	if param, err := strconv.Atoi(s); err != nil {
		return 0, parseError{message: fmt.Sprintf("%s is not an integer", s)}
	} else {
		return param, nil
	}
}
