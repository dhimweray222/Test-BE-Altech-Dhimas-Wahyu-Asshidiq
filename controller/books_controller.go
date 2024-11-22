package controller

import (
	"test-backend-altech/config"
	"test-backend-altech/exception"
	web "test-backend-altech/model/web"
	req "test-backend-altech/model/web/req"
	"test-backend-altech/service"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type BookController interface {
	Route(app *fiber.App)
}

type bookController struct {
	validate    *validator.Validate
	bookService service.BookService
	cache       config.Cache
}

func NewBookController(validate *validator.Validate, bookService service.BookService, cache config.Cache) BookController {
	return &bookController{
		validate:    validate,
		bookService: bookService,
		cache:       cache,
	}
}
func (controller *bookController) Route(app *fiber.App) {
	api := app.Group("/books")
	api.Post("/",
		controller.CreateBook,
	)
	api.Get("/:book_id",
		controller.FindByID,
	)
	api.Put("/:book_id",
		controller.UpdateBook)

	api.Get("/",
		controller.FindAllBook,
	)

	api.Delete("/:book_id",
		controller.DeleteBook,
	)
}
func (controller *bookController) CreateBook(ctx *fiber.Ctx) error {
	var request req.BookRequest
	err := ctx.BodyParser(&request)
	if err != nil {
		return exception.ErrValidateBadRequest(err.Error(), request)
	}

	err = controller.validate.Struct(request)
	if err != nil {
		return exception.ErrValidateBadRequest(err.Error(), request)
	}
	bookResponse, err := controller.bookService.CreateBook(ctx.Context(), request)
	if err != nil {
		return exception.ErrorHandler(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(web.WebResponse{
		Code:    fiber.StatusOK,
		Status:  true,
		Message: "success",
		Data:    bookResponse,
	})
}

func (controller *bookController) FindByID(ctx *fiber.Ctx) error {

	bookId := ctx.Params("book_id")

	book, err := controller.bookService.FindByID(ctx.Context(), bookId)
	if err != nil {
		return err
	}
	if book.Id == "" {
		return exception.ErrNotFound("Book not found")
	}

	return ctx.Status(fiber.StatusOK).JSON(web.WebResponse{
		Code:    fiber.StatusOK,
		Status:  true,
		Message: "success",
		Data:    book,
	})
}

func (controller *bookController) UpdateBook(ctx *fiber.Ctx) error {
	id := ctx.Params("book_id")

	var request req.BookRequest
	err := ctx.BodyParser(&request)
	if err != nil {
		return exception.ErrValidateBadRequest(err.Error(), request)
	}

	err = controller.validate.Struct(request)
	if err != nil {
		return exception.ErrValidateBadRequest(err.Error(), request)
	}
	bookResponse, err := controller.bookService.UpdateBook(ctx.Context(), request, id)
	if err != nil {
		return exception.ErrorHandler(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(web.WebResponse{
		Code:    fiber.StatusOK,
		Status:  true,
		Message: "success",
		Data:    bookResponse,
	})
}

func (controller *bookController) FindAllBook(ctx *fiber.Ctx) error {
	books, err := controller.bookService.FindAllBook(ctx.Context(), controller.cache)
	if err != nil {
		return exception.ErrorHandler(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(web.WebResponse{
		Code:    fiber.StatusOK,
		Status:  true,
		Message: "success",
		Data:    books,
	})
}

func (controller *bookController) DeleteBook(ctx *fiber.Ctx) error {
	id := ctx.Params("book_id")
	data, err := controller.bookService.DeleteBook(ctx.Context(), id)
	if err != nil {
		return exception.ErrorHandler(ctx, err)
	}
	return ctx.Status(fiber.StatusOK).JSON(web.WebResponse{
		Code:    fiber.StatusOK,
		Status:  true,
		Message: "success",
		Data:    data,
	})
}
