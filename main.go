// package main

// import(
// 	"log"
// 	"net/http"
// 	"os"
// 	"github.com/gofiber/fiber/v2"
// 	"github.com/joho/godotenv"
// 	"gorm.io/gorm"
// 	"fmt"
// 	"book-api/src/storage"
// 	"book-api/src/models"
// )

// type Book struct{
// 	Author string `json:"author"`
// 	Title string `json:"title"`
// 	Publisher string `json:"publisher"`
// }

// type Repository struct{
// 	DB *gorm.DB
// }

// func (r *Repository) GetBooks(c *fiber.Ctx) error {
// 	var books []Book
// 	r.DB.Find(&books)
// 	return c.JSON(books)

// 	fmt.Println("id is: ", id)
// }


// func (r *Repository) GetBook(c *fiber.Ctx) error {
// 	var book Book
// 	id := c.Params("id")
// 	r.DB.Find(&book, id)
// 	return c.JSON(book)
// }	

// func (r *Repository) NewBook(c *fiber.Ctx) error {
// 	book := new(Book)
// 	if err := c.BodyParser(book); err != nil {
// 		return c.Status(500).SendString(err.Error())
// 	}
// 	r.DB.Create(&book)
// 	return c.JSON(book)
// }

// func (r *Repository) UpdateBook(c *fiber.Ctx) error {
// 	var book Book
// 	id := c.Params("id")
// 	r.DB.Find(&book, id)
// 	if err := c.BodyParser(&book); err != nil {
// 		return c.Status(500).SendString(err.Error())
// 	}
// 	r.DB.Save(&book)
// 	return c.JSON(book)

// }

// func (r *Repository) DeleteBook(c *fiber.Ctx) error {
// 	id := c.Params("id")
// 	var book Book
// 	r.DB.First(&book, id)
// 	if book.Title == "" {
// 		return c.Status(500).SendString("No book found with given ID")
// 	}
// 	r.DB.Delete(&book)
// 	return c.SendString("Book successfully deleted")

// }

// func (r *Repository) SetupRoutes(app *fiber.App)  {
// 	api := app.Group("/api")
// 	api.Get("/get_books", r.GetBooks)
// 	api.Get("/get_book/:id", r.GetBook)
// 	api.Post("/create_book", r.NewBook)
// 	api.Put("/update_book/:id", r.UpdateBook)
// 	api.Delete("/delete_book/:id", r.DeleteBook)
// }

// func main()  {
// 	err := godotenv.Load(".env")
// 	if  err != nil {
// 		log.Fatalf(err)
// 	}

// 	config := storage.Config{
// 		Host: os.Getenv("DB_HOST"),
// 		Port: os.Getenv("DB_PORT"),
// 		User: os.Getenv("DB_USER"),
// 		DBName: os.Getenv("DB_NAME"),
// 		Password: os.Getenv("DB_PASSWORD"),
// 		SSLMode: os.Getenv("DB_SSLMODE"),
// 	}

// 	db, err := storage.NewConnection(config)
// 	if err != nil {
// 		log.Fatalf(err)
// 	}

// 	err = models.Migrate(db)
// 	if err != nil {
// 		log.Fatalf("could not migrate DB")
// 	}

// 	r := Repository{
// 		DB: db,
// 	}

// 	app := fiber.New()
// 	r.SetupRoutes(app)
// 	app.Listen(":8080")
// }

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
	"book-api/src/storage"
	"book-api/src/models"

)

type Book struct {
	Author    string `json:"author"`
	Title     string `json:"title"`
	Publisher string `json:"publisher"`
}

type Repository struct {
	DB *gorm.DB
}

func (r *Repository) CreateBook(context *fiber.Ctx) error {
	book := Book{}

	err := context.BodyParser(&book)

	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "request failed"})
		return err
	}

	err = r.DB.Create(&book).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not create book"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "book has been added"})
	return nil
}

func (r *Repository) DeleteBook(context *fiber.Ctx) error {
	bookModel := models.Books{}
	id := context.Params("id")
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
		return nil
	}

	err := r.DB.Delete(bookModel, id)

	if err.Error != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "could not delete book",
		})
		return err.Error
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "book delete successfully",
	})
	return nil
}

func (r *Repository) GetBooks(context *fiber.Ctx) error {
	bookModels := &[]models.Books{}

	err := r.DB.Find(bookModels).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not get books"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "books fetched successfully",
		"data":    bookModels,
	})
	return nil
}

func (r *Repository) GetBookByID(context *fiber.Ctx) error {

	id := context.Params("id")
	bookModel := &models.Books{}
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
		return nil
	}

	fmt.Println("the ID is", id)

	err := r.DB.Where("id = ?", id).First(bookModel).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not get the book"})
		return err
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "book id fetched successfully",
		"data":    bookModel,
	})
	return nil
}

func (r *Repository) SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Post("/create_books", r.CreateBook)
	api.Delete("delete_book/:id", r.DeleteBook)
	api.Get("/get_books/:id", r.GetBookByID)
	api.Get("/books", r.GetBooks)
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	config := &storage.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASS"),
		User:     os.Getenv("DB_USER"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
		DBName:   os.Getenv("DB_NAME"),
	}

	db, err := storage.NewConnection(config)

	if err != nil {
		log.Fatal("could not load the database")
	}
	err = models.MigrateBooks(db)
	if err != nil {
		log.Fatal("could not migrate db")
	}

	r := Repository{
		DB: db,
	}
	app := fiber.New()
	r.SetupRoutes(app)
	app.Listen(":8080")
}