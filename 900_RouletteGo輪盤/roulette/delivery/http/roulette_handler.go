package http

import (
	"RouletteGo/domain"
	swagger "RouletteGo/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// RouletteHandler ...
type RouletteHandler struct {
	RouletteUsecase domain.RouletteUsecase
}

// NewRouletteHandler ...
func NewRouletteHandler(server *gin.Engine, rouletteUsecase domain.RouletteUsecase) {
	handler := &RouletteHandler{
		RouletteUsecase: rouletteUsecase,
	}

	server.GET("/api/v1/Roulette/:rouletteID", handler.GetRouletteByBetID)

	server.GET("/test", handler.GetTest)
}

func (d *RouletteHandler) GetTest(c *gin.Context) {
	result := "{'msg':'test ok!'}"
	c.JSON(http.StatusOK, result)
}

// GetRouletteByBetID ...
func (d *RouletteHandler) GetRouletteByBetID(c *gin.Context) {
	rouletteID := c.Param("")

	anDigimon, err := d.RouletteUsecase.GetByID(c, rouletteID)
	if err != nil {
		logrus.Error(err)
		c.JSON(500, &swagger.ModelError{
			Code:    3000,
			Message: "Internal error. Query digimon error",
		})
		return
	}

	c.JSON(200, &swagger.RoulletteInfo{
		Id:   anDigimon.ID,
		Name: anDigimon.Name,
	})
}
