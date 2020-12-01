package delivery

import "github.com/gin-gonic/gin"

type HttpHandler interface {
	GetUserList(c *gin.Context)
	GetUser(c *gin.Context)
	CreateUser(c *gin.Context)
	ModifyUser(c *gin.Context)
	DeleteUser(c *gin.Context)
}
