package handler

import (
	"context"
	"fmt"

	"github.com/Talonmortem/wb-test-task1/internal/postgres"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	DB     postgres.DB
	Schema string
}

func NewHandler(db *postgres.DB, schema string) *Handler {
	return &Handler{
		DB:     *db,
		Schema: schema,
	}
}

func (h *Handler) HandleRequest(c *gin.Context) {
	method := c.Request.Method
	path := c.Param("path")
	var idParams []int
	var params interface{}
	if method == "POST" || method == "PUT" {
		if err := c.ShouldBindJSON(&params); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
	}

	//Преобразоване метода HTTP в SQL-запрос
	var suffix string
	switch method {
	case "GET":
		suffix = "get"
	case "POST":
		suffix = "ins"
	case "PUT":
		suffix = "upd"
	case "DELETE":
		suffix = "del"
	default:
		c.JSON(405, gin.H{"error": "Method not allowed"})
		return
	}

	functionName := fmt.Sprintf("%s.%s_%s", h.Schema, path, suffix)
	query := fmt.Sprintf("SELECT * FROM %s($1, $2, $3)", functionName)

	//Выполнение SQL-запроса

	results, err := h.DB.ExecuteFunction(context.Background(), query, idParams, params)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, results)
}
