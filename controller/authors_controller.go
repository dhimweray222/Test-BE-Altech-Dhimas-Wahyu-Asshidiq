package controller

import (
	"test-backend-altech/exception"
	web "test-backend-altech/model/web"
	req "test-backend-altech/model/web/req"
	"test-backend-altech/service"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type AuthorController interface {
	Route(app *fiber.App)
}

type authorController struct {
	validate      *validator.Validate
	authorService service.AuthorService
}

func NewAuthorController(validate *validator.Validate, authorService service.AuthorService) AuthorController {
	return &authorController{
		validate:      validate,
		authorService: authorService,
	}
}
func (controller *authorController) Route(app *fiber.App) {
	api := app.Group("/authors")

	api.Post("/",
		controller.CreateAuthor,
	)
	api.Get("/:author_id",
		controller.FindByID,
	)
	api.Put("/:author_id",
		controller.UpdateAuthor)

	api.Get("/",
		controller.FindAllAuthor,
	)

	api.Delete("/:author_id",
		controller.DeleteAuthor,
	)
}
func (controller *authorController) CreateAuthor(ctx *fiber.Ctx) error {
	var request req.AuthorRequest
	err := ctx.BodyParser(&request)
	if err != nil {
		return exception.ErrValidateBadRequest(err.Error(), request)
	}

	err = controller.validate.Struct(request)
	if err != nil {
		return exception.ErrValidateBadRequest(err.Error(), request)
	}
	authorResponse, err := controller.authorService.CreateAuthor(ctx.Context(), request)
	if err != nil {
		return exception.ErrorHandler(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(web.WebResponse{
		Code:    fiber.StatusOK,
		Status:  true,
		Message: "success",
		Data:    authorResponse,
	})
}

func (controller *authorController) FindByID(ctx *fiber.Ctx) error {

	authorId := ctx.Params("author_id")

	author, err := controller.authorService.FindByID(ctx.Context(), authorId)
	if err != nil {
		return err
	}
	if author.Id == "" {
		return exception.ErrNotFound("Author not found")
	}

	return ctx.Status(fiber.StatusOK).JSON(web.WebResponse{
		Code:    fiber.StatusOK,
		Status:  true,
		Message: "success",
		Data:    author,
	})
}

func (controller *authorController) UpdateAuthor(ctx *fiber.Ctx) error {
	id := ctx.Params("author_id")

	var request req.AuthorRequest
	err := ctx.BodyParser(&request)
	if err != nil {
		return exception.ErrValidateBadRequest(err.Error(), request)
	}

	err = controller.validate.Struct(request)
	if err != nil {
		return exception.ErrValidateBadRequest(err.Error(), request)
	}
	authorResponse, err := controller.authorService.UpdateAuthor(ctx.Context(), request, id)
	if err != nil {
		return exception.ErrorHandler(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(web.WebResponse{
		Code:    fiber.StatusOK,
		Status:  true,
		Message: "success",
		Data:    authorResponse,
	})
}

func (controller *authorController) FindAllAuthor(ctx *fiber.Ctx) error {
	authors, err := controller.authorService.FindAllAuthor(ctx.Context())
	if err != nil {
		return exception.ErrorHandler(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(web.WebResponse{
		Code:    fiber.StatusOK,
		Status:  true,
		Message: "success",
		Data:    authors,
	})
}

func (controller *authorController) DeleteAuthor(ctx *fiber.Ctx) error {
	id := ctx.Params("author_id")
	data, err := controller.authorService.DeleteAuthor(ctx.Context(), id)
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
