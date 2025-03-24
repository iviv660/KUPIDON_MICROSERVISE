package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"service1/internal/entity"
	"service1/internal/usecase"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	usecase *usecase.UserUsecase
}

func NewUserHandler(usecase usecase.UserUsecase) (*UserHandler, *gin.Engine) {
	h := UserHandler{usecase: &usecase}
	router := gin.New()
	router.POST("/users", h.CreateUser)
	router.GET("/users/:id", h.GetByID)
	router.GET("/users/search", h.Search)
	router.PUT("/users/:id", h.Update)
	router.DELETE("/users/:id", h.Delete)

	return &h, router
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	// 1️⃣ Получаем JSON как строку
	jsonData := c.PostForm("json")
	if jsonData == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing JSON data"})
		return
	}

	// 2️⃣ Распаковываем JSON в структуру
	var req struct {
		Name        string `json:"name"`
		City        string `json:"city"`
		Gender      string `json:"gender"`
		Description string `json:"description"`
		Age         int    `json:"age"`
		TelegramID  int64  `json:"telegram_id"`
	}

	if err := json.Unmarshal([]byte(jsonData), &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON: " + err.Error()})
		return
	}

	// 3️⃣ Проверяем обязательные поля
	if req.TelegramID == 0 || req.Name == "" || req.Age <= 0 || req.Description == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
		return
	}

	// 4️⃣ Получаем файл
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}

	// 5️⃣ Проверяем размер файла (макс. 5MB)
	const maxFileSize = 100 * 1024 * 1024
	if fileHeader.Size > maxFileSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File size exceeds the limit of 5MB"})
		return
	}

	// 6️⃣ Открываем файл
	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to open file: " + err.Error()})
		return
	}
	defer file.Close()

	// 7️⃣ Вызываем бизнес-логику
	userID, err := h.usecase.Create(
		c.Request.Context(),
		req.Name,
		req.Description,
		fileHeader.Filename,
		req.Gender,
		req.City,
		req.Age,
		file,
		fileHeader.Size,
		req.TelegramID,
	)

	if err != nil {
		log.Printf("Error in creating user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	if userID == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user: invalid user ID"})
		return
	}

	// 8️⃣ Возвращаем успешный ответ
	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"user_id": userID,
	})
}

// @Summary Search user
// @Description Search for users based on filters such as age, city, and gender
// @Tags filter_users
// @Accept json
// @Produce json
// @Param min_age query int false "Minimum age of user"
// @Param max_age query int false "Maximum age of user"
// @Param description query string false "Description filter"
// @Param city query string false "City filter"
// @Param gender query string false "Gender filter"
// @Success 200 {array} usecase.User "List of users"
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Router /users/search [get]
func (h *UserHandler) Search(c *gin.Context) {
	var req struct {
		MinAge *int   `form:"min_age,omitempty"` // Используем form-теги для query параметров
		MaxAge *int   `form:"max_age,omitempty"`
		City   string `form:"city,omitempty"`
		Gender string `form:"gender,omitempty"`
	}

	// Привязываем параметры запроса
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	var filter entity.UserFilter
	filter.MinAge = req.MinAge
	filter.MaxAge = req.MaxAge
	filter.Gender = req.Gender
	filter.City = req.City
	users, err := h.usecase.Search(c.Request.Context(), filter)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

// @Summary Get user by ID
// @Description Get a user by their unique ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} usecase.User "User details"
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Router /users/{id} [get]
func (h *UserHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")

	// Проверяем, что ID – это валидное положительное число
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Получаем пользователя из usecase
	user, err := h.usecase.GetByID(c.Request.Context(), id)
	if err != nil {
		log.Printf("Error fetching user %d: %v", id, err)

		// Если юзер не найден, возвращаем 404
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		// В остальных случаях – внутренняя ошибка сервера
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Успешный ответ
	c.JSON(http.StatusOK, user)
}

// @Summary Update user
// @Description Update an existing user's details
// @Tags users
// @Accept json
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "User ID"
// @Param name formData string true "Name of the user"
// @Param age formData int true "Age of the user"
// @Param description formData string true "Description of the user"
// @Param telegram_id formData int64 true "Telegram ID"
// @Param file formData file true "User's new photo"
// @Success 200 {string} string "User updated"
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Router /users/{id} [put]
func (h *UserHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var req struct {
		TelegramID  int64  `json:"telegram_id"`
		Name        string `json:"name"`
		Age         int    `json:"age"`
		City        string `json:"city,omitempty"`
		Gender      string `json:"gender,omitempty"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}
	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to open file"})
		return
	}
	defer file.Close()
	err = h.usecase.Update(c.Request.Context(), req.Name, req.Description, fileHeader.Filename, req.Gender, req.City, req.Age, id, file, fileHeader.Size, req.TelegramID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

// @Summary Delete user
// @Description Delete a user by their ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {string} string "User deleted"
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Router /users/{id} [delete]
func (h *UserHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err = h.usecase.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.Status(http.StatusOK)
}
