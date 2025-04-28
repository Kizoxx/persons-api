// internal/services/person_service.go
package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"persons-api/internal/models"
	"persons-api/internal/repository"
)

// PersonService определяет методы бизнес-логики для работы с Person
type PersonService interface {
	Create(person *models.Person) error
	Get(id uint) (*models.Person, error)
	List(name, gender, country string, page, size int) ([]models.Person, error)
	Update(id uint, person *models.Person) error
	Delete(id uint) error
}

// personService реализует PersonService
type personService struct {
	repo repository.PersonRepository
}

// NewPersonService создаёт новый сервис
func NewPersonService(repo repository.PersonRepository) PersonService {
	return &personService{repo: repo}
}

// Create создаёт новую запись и обогащает данные
func (s *personService) Create(person *models.Person) error {
	// Обогащение данных (пример: получение пола, возраста, страны)
	if err := s.enrichPerson(person); err != nil {
		return fmt.Errorf("failed to enrich person: %w", err)
	}
	return s.repo.Create(person)
}

// Get получает запись по ID
func (s *personService) Get(id uint) (*models.Person, error) {
	return s.repo.Get(id)
}

// List возвращает список записей с фильтрацией и пагинацией
func (s *personService) List(name, gender, country string, page, size int) ([]models.Person, error) {
	return s.repo.List(name, gender, country, page, size)
}

// Update обновляет запись и переобогащает данные
func (s *personService) Update(id uint, person *models.Person) error {
	// Переобогащение данных
	if err := s.enrichPerson(person); err != nil {
		return fmt.Errorf("failed to enrich person: %w", err)
	}
	return s.repo.Update(id, person)
}

// Delete удаляет запись
func (s *personService) Delete(id uint) error {
	return s.repo.Delete(id)
}

// enrichPerson обогащает данные человека (возраст, пол, страна)
func (s *personService) enrichPerson(person *models.Person) error {
	// Получение пола через genderize.io
	genderResp, err := http.Get(fmt.Sprintf("https://api.genderize.io?name=%s", person.FirstName))
	if err != nil {
		return err
	}
	defer genderResp.Body.Close()
	var genderData struct {
		Gender string `json:"gender"`
	}
	if err := json.NewDecoder(genderResp.Body).Decode(&genderData); err != nil {
		return err
	}
	person.Gender = genderData.Gender

	// Получение возраста через agility.io (пример API)
	ageResp, err := http.Get(fmt.Sprintf("https://api.agify.io?name=%s", person.FirstName))
	if err != nil {
		return err
	}
	defer ageResp.Body.Close()
	var ageData struct {
		Age int `json:"age"`
	}
	if err := json.NewDecoder(ageResp.Body).Decode(&ageData); err != nil {
		return err
	}
	person.Age = ageData.Age

	// Получение страны через nationalize.io
	countryResp, err := http.Get(fmt.Sprintf("https://api.nationalize.io?name=%s", person.FirstName))
	if err != nil {
		return err
	}
	defer countryResp.Body.Close()
	var countryData struct {
		Country []struct {
			CountryID string `json:"country_id"`
		} `json:"country"`
	}
	if err := json.NewDecoder(countryResp.Body).Decode(&countryData); err != nil {
		return err
	}
	if len(countryData.Country) > 0 {
		person.Country = countryData.Country[0].CountryID
	}

	return nil
}
