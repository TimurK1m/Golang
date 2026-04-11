package v1

import (
	"Secure/internal/entity"
	"Secure/internal/usecase"
	"Secure/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/logger"
)

type UserRoutes struct {
	t usecase.UserInterface 
	l logger.Interface
}

func NewUserRoutes(handler *gin.RouterGroup, t usecase.UserInterface, l logger.Interface) {
	r := &UserRoutes{t, l}
	h := handler.Group("/users")
	{
		h.POST("/", r.RegisterUser)
		h.POST("/login", r.LoginUser)
		protected := h.Group("/")
		protected.Use(utils.JWTAuthMiddleware())
			{
				protected.GET("/protected/hello", r.ProtectedFunc)
				protected.GET("/me", r.GetMe)
				
			}
			
		admin := h.Group("/")
		admin.Use(utils.JWTAuthMiddleware(), utils.RoleMiddleware("admin"))
		{
			admin.PATCH("/promote/:id", r.PromoteUser)
		}
	}
}

func (r *UserRoutes) PromoteUser(c *gin.Context) {
	id := c.Param("id")

	err := r.t.PromoteUser(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user promoted"})
}

func (r *UserRoutes) GetMe(c *gin.Context) {
	userID := c.GetString("userID")

	user, err := r.t.GetMe(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"email": user.Email,
		"username": user.Username,
		"role": user.Role,
	})
}

func (r *UserRoutes) ProtectedFunc(c *gin.Context) {
	c.JSON(200, gin.H{"message": "OK"})
}


func (r *UserRoutes) RegisterUser(c * gin.Context){
	var createUserDTO entity.CreateUserDTO
	if err := c.ShouldBindJSON(&createUserDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := utils.HashPassword(createUserDTO.Password)
	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
		return
	}
	user := entity.User{
		Username: createUserDTO.Username,
		Email: createUserDTO.Email,
		Password: hashedPassword,
		Role: "user",
	}
	createdUser, sessionID,err := r.t.RegisterUser(&user)
	if err!=nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"massage": "User registered successfully",
		"session_id": sessionID,
		"user": createdUser,
	})
}

func (r *UserRoutes) LoginUser(c *gin.Context) {
	var input entity.LoginUserDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, err := r.t.LoginUser(&input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error":
		err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

