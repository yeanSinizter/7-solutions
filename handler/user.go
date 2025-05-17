package handler

import (
	"7-solutions/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	Usecase usecase.UserUsecase
}

func NewUserHandler(r *gin.Engine, uc usecase.UserUsecase, auth gin.HandlerFunc) {
	h := &UserHandler{Usecase: uc}

	r.POST("/register", h.Register)
	r.POST("/login", h.Login)

	authGroup := r.Group("/users", auth)
	authGroup.GET("/", h.List)
	authGroup.GET("/:id", h.Get)
	authGroup.PUT("/:id", h.Update)
	authGroup.DELETE("/:id", h.Delete)
}

func (h *UserHandler) Register(c *gin.Context) {
	var req struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := h.Usecase.Register(c.Request.Context(), req.Name, req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "registered"})
}

func (h *UserHandler) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, err := h.Usecase.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *UserHandler) Get(c *gin.Context) {
	id := c.Param("id")
	user, err := h.Usecase.GetUser(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) List(c *gin.Context) {
	users, err := h.Usecase.ListUsers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

func (h *UserHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := h.Usecase.UpdateUser(c.Request.Context(), id, req.Name, req.Email)
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "updated"})
}

func (h *UserHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	err := h.Usecase.DeleteUser(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}
