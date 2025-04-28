// internal/repository/person_repository.go
package repository

import (
	"gorm.io/gorm"
	"persons-api/internal/models"
)

// PersonRepository определяет методы для работы с сущностью Person в базе данных
type PersonRepository interface {
	Create(person *models.Person) error
	Get(id uint) (*models.Person, error)
	List(name, gender, country string, page, size int) ([]models.Person, error)
	Update(id uint, person *models.Person) error
	Delete(id uint) error
}

// personRepository реализует PersonRepository
type personRepository struct {
	db *gorm.DB
}

// NewPersonRepository создаёт новый репозиторий
func NewPersonRepository(db *gorm.DB) PersonRepository {
	return &personRepository{db: db}
}

// Create добавляет новую запись в базу данных
func (r *personRepository) Create(person *models.Person) error {
	return r.db.Create(person).Error
}

// Get получает запись по ID
func (r *personRepository) Get(id uint) (*models.Person, error) {
	var person models.Person
	err := r.db.First(&person, id).Error
	if err != nil {
		return nil, err
	}
	return &person, nil
}

// List возвращает список записей с фильтрацией и пагинацией
func (r *personRepository) List(name, gender, country string, page, size int) ([]models.Person, error) {
	var persons []models.Person
	query := r.db.Model(&models.Person{})

	if name != "" {
		query = query.Where("first_name ILIKE ?", "%"+name+"%")
	}
	if gender != "" {
		query = query.Where("gender = ?", gender)
	}
	if country != "" {
		query = query.Where("country = ?", country)
	}

	offset := (page - 1) * size
	return persons, query.Limit(size).Offset(offset).Find(&persons).Error
}

// Update обновляет запись
func (r *personRepository) Update(id uint, person *models.Person) error {
	return r.db.Model(&models.Person{}).Where("id = ?", id).Updates(person).Error
}

// Delete удаляет запись по ID
func (r *personRepository) Delete(id uint) error {
	return r.db.Delete(&models.Person{}, id).Error
}
