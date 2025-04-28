package api

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"persons-api/internal/models"
	"persons-api/internal/services"
	"strconv"
	"time"
)

// PersonHandler обрабатывает HTTP-запросы для сущности Person
type PersonHandler struct {
	service services.PersonService
}

// NewPersonHandler регистрирует маршруты для Person
func NewPersonHandler(r *gin.Engine, service services.PersonService) {
	handler := &PersonHandler{service: service}

	// Группа /api/v1/people
	v1 := r.Group("/api/v1")
	people := v1.Group("/people")
	{
		people.GET("", handler.ListPersons)         // Список людей с фильтрами и пагинацией
		people.GET("/:id", handler.GetPerson)       // Получение по ID
		people.POST("", handler.CreatePerson)       // Создание нового (автообогащение)
		people.PUT("/:id", handler.UpdatePerson)    // Обновление
		people.DELETE("/:id", handler.DeletePerson) // Удаление по ID
	}
}

// ListPersons godoc
// @Summary Получить список людей
// @Description Список людей с возможностью фильтрации по имени, полу, стране. Пагинация (page, size).
// @Tags people
// @Accept json
// @Produce json
// @Param name query string false "Фильтр по имени (частичный поиск)"
// @Param gender query string false "Фильтр по полу (male/female)"
// @Param country query string false "Фильтр по стране (код страны)"
// @Param page query int false "Номер страницы (начинается с 1)"
// @Param size query int false "Размер страницы"
// @Success 200 {array} models.PersonSwagger
// @Failure 400 {object} object
// @Router /api/v1/people [get]
func (h *PersonHandler) ListPersons(c *gin.Context) {
	name := c.Query("name")
	gender := c.Query("gender")
	country := c.Query("country")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	persons, err := h.service.List(name, gender, country, page, size)
	if err != nil {
		logrus.Errorf("Ошибка при получении списка людей: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, persons)
}

// GetPerson godoc
// @Summary Получить человека по ID
// @Description Получить детали человека по его ID.
// @Tags people
// @Accept json
// @Produce json
// @Param id path int true "ID человека"
// @Success 200 {object} models.PersonSwagger
// @Failure 404 {object} object
// @Router /api/v1/people/{id} [get]
func (h *PersonHandler) GetPerson(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		logrus.Warnf("Неверный ID: %s, ошибка: %v", idParam, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	person, err := h.service.Get(uint(id))
	if err != nil {
		logrus.Warnf("Человек с ID %d не найден: %v", id, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Person not found"})
		return
	}

	// Конвертируем Person в PersonSwagger для ответа
	var deletedAt *time.Time
	if person.DeletedAt.Valid {
		deletedAt = &person.DeletedAt.Time
	}
	personSwagger := models.PersonSwagger{
		ID:         person.ID,
		CreatedAt:  person.CreatedAt,
		UpdatedAt:  person.UpdatedAt,
		DeletedAt:  deletedAt,
		FirstName:  person.FirstName,
		LastName:   person.LastName,
		Patronymic: person.Patronymic,
		Age:        person.Age,
		Gender:     person.Gender,
		Country:    person.Country,
	}
	c.JSON(http.StatusOK, personSwagger)
}

// CreatePerson godoc
// @Summary Создать нового человека
// @Description Создает нового человека. Имя передается в запросе, остальное поле (возраст, пол, страна) заполняются автоматически.
// @Tags people
// @Accept json
// @Produce json
// @Param person body models.PersonSwagger true "Новый человек (только имя обязательно)"
// @Success 201 {object} models.PersonSwagger
// @Failure 400 {object} object
// @Router /api/v1/people [post]
func (h *PersonHandler) CreatePerson(c *gin.Context) {
	var personSwagger models.PersonSwagger
	if err := c.ShouldBindJSON(&personSwagger); err != nil {
		logrus.Errorf("Ошибка десериализации JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	logrus.Infof("Десериализованный personSwagger: %+v", personSwagger)

	if personSwagger.FirstName == "" {
		logrus.Warn("Поле first_name пустое")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name is required"})
		return
	}

	// Конвертируем PersonSwagger в Person
	person := models.Person{
		FirstName:  personSwagger.FirstName,
		LastName:   personSwagger.LastName,
		Patronymic: personSwagger.Patronymic,
		Age:        personSwagger.Age,
		Gender:     personSwagger.Gender,
		Country:    personSwagger.Country,
	}

	if err := h.service.Create(&person); err != nil {
		logrus.Errorf("Ошибка создания человека: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Обновляем personSwagger с новым ID и датами
	var deletedAt *time.Time
	if person.DeletedAt.Valid {
		deletedAt = &person.DeletedAt.Time
	}
	personSwagger.ID = person.ID
	personSwagger.CreatedAt = person.CreatedAt
	personSwagger.UpdatedAt = person.UpdatedAt
	personSwagger.DeletedAt = deletedAt

	logrus.Infof("Человек успешно создан: %+v", personSwagger)
	c.JSON(http.StatusCreated, personSwagger)
}

// UpdatePerson godoc
// @Summary Обновить человека
// @Description Обновляет данные человека. Переобогащает по имени.
// @Tags people
// @Accept json
// @Produce json
// @Param id path int true "ID человека"
// @Param person body models.PersonSwagger true "Обновленные данные человека"
// @Success 200 {object} models.PersonSwagger
// @Failure 400 {object} object
// @Router /api/v1/people/{id} [put]
func (h *PersonHandler) UpdatePerson(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		logrus.Warnf("Неверный ID: %s, ошибка: %v", idParam, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var personSwagger models.PersonSwagger
	if err := c.ShouldBindJSON(&personSwagger); err != nil {
		logrus.Errorf("Ошибка десериализации JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	logrus.Infof("Десериализованный personSwagger для обновления: %+v", personSwagger)

	// Конвертируем PersonSwagger в Person
	person := models.Person{
		FirstName:  personSwagger.FirstName,
		LastName:   personSwagger.LastName,
		Patronymic: personSwagger.Patronymic,
		Age:        personSwagger.Age,
		Gender:     personSwagger.Gender,
		Country:    personSwagger.Country,
	}

	if err := h.service.Update(uint(id), &person); err != nil {
		logrus.Errorf("Ошибка обновления человека с ID %d: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Обновляем personSwagger с актуальными данными
	var deletedAt *time.Time
	if person.DeletedAt.Valid {
		deletedAt = &person.DeletedAt.Time
	}
	personSwagger.ID = person.ID
	personSwagger.CreatedAt = person.CreatedAt
	personSwagger.UpdatedAt = person.UpdatedAt
	personSwagger.DeletedAt = deletedAt

	logrus.Infof("Человек с ID %d успешно обновлён: %+v", id, personSwagger)
	c.JSON(http.StatusOK, personSwagger)
}

// DeletePerson godoc
// @Summary Удалить человека
// @Description Удаляет человека по ID.
// @Tags people
// @Accept json
// @Produce json
// @Param id path int true "ID человека"
// @Success 204
// @Failure 400 {object} object
// @Router /api/v1/people/{id} [delete]
func (h *PersonHandler) DeletePerson(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		logrus.Warnf("Неверный ID: %s, ошибка: %v", idParam, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := h.service.Delete(uint(id)); err != nil {
		logrus.Errorf("Ошибка удаления человека с ID %d: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logrus.Infof("Человек с ID %d успешно удалён", id)
	c.Status(http.StatusNoContent)
}
