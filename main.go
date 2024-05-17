package main

import (
	"Student_REST_API/models"
	"Student_REST_API/storage"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

type Student struct {
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Class   string `json:"class"`
	Subject string `json:"subject"`
}

func (r *Repository) CreateStudent(context *fiber.Ctx) error {
	student := Student{}

	err := context.BodyParser(&student)
	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{"message": "request failed"})
		return err
	}

	err = r.DB.Create(&student).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "could not create a student"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{"message": "student create successfully"})
	return nil
}

func (r *Repository) GetStudent(context *fiber.Ctx) error {
	studentModels := &[]models.Students{}

	err := r.DB.Find(studentModels).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "could not get the student"})
		return err
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{"message": "Student fetched successfully", "data": studentModels})
	return nil
}

func (r *Repository) DeleteStudent(context *fiber.Ctx) error {
	studentModel := models.Students{}
	id := context.Params("id")
	fmt.Println("The ID is", id)
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{"message": "id cannot be empty"})
		return nil
	}

	err := r.DB.Delete(studentModel, id)
	if err.Error != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "could not delete student"})
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{"message": "student deleted successfully"})
	return nil
}

func (r *Repository) GetStudentByID(context *fiber.Ctx) error {
	id := context.Params("id")
	studentModel := &models.Students{}
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{"message": "id cannot be empty"})
		return nil
	}
	fmt.Println("The ID is", id)
	err := r.DB.Where("id= ?", id).First(studentModel).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "could not get the student"})
		return err
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{"message": "Student ID fetched successfully", "data": studentModel})
	return nil
}

func (r *Repository) SetUpRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Post("/student/v1/students/create", r.CreateStudent)
	api.Get("/GET/student/v1/students", r.GetStudent)
	api.Delete("/student/v1/students/delete/:id", r.DeleteStudent)
	api.Get("/GET/student/v1/students/:id", r.GetStudentByID)

}

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	config := &storage.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}
	db, err := storage.NewConnection(config)
	if err != nil {
		log.Fatal("Could not load the database")
	}
	err = models.MigrateStudents(db)

	if err != nil {
		log.Fatal("Could not migrate the database")
	}
	r := Repository{
		DB: db,
	}

	app := fiber.New()
	r.SetUpRoutes(app)
	app.Listen(":8080")
	fmt.Println("I am building REST API For Students")
}
