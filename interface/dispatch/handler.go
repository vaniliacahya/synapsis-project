package dispatch

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"synapsis-project/database/databasesModel"
	"synapsis-project/domain"
	"synapsis-project/structures/request"
)

type Handler struct {
	useCase domain.UseCase
	data    domain.Data
}

func New(useCase domain.UseCase, data domain.Data) domain.Handler {
	return &Handler{
		useCase: useCase,
		data:    data,
	}
}

func (h *Handler) ListProduct() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		param := new(request.ListProductRequest)
		if err := ctx.QueryParser(param); err != nil {
			logrus.WithError(fmt.Errorf("parse param: %v", err)).Error()
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"code":   fiber.StatusBadRequest,
				"reason": fmt.Sprintf("[%d] %s", fiber.StatusBadRequest, err.Error()),
				"data":   nil,
			})
		}

		result := h.useCase.ListProduct(*param)
		if result.ErrorMsg != nil {
			logrus.WithError(fmt.Errorf("failed to get list product : %v", result.ErrorMsg)).Error()
			return ctx.Status(result.HttpErrorCode).JSON(fiber.Map{
				"code":   result.HttpErrorCode,
				"reason": fmt.Sprintf("[%d] %s", result.HttpErrorCode, result.ErrorMsg),
				"data":   nil,
			})
		}

		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"code": fiber.StatusOK,
			"msg":  "success get list product",
			"data": result.Response,
		})
	}
}

func (h *Handler) AddCart() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		body := new(request.AddCartRequest)
		if err := ctx.BodyParser(body); err != nil {
			logrus.WithError(fmt.Errorf("parse body: %v", err)).Error()
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"code":   fiber.StatusBadRequest,
				"reason": fmt.Sprintf("[%d] %s", fiber.StatusBadRequest, err.Error()),
				"data":   nil,
			})
		}

		result := h.useCase.AddCart(*body)
		if result.ErrorMsg != nil {
			logrus.WithError(fmt.Errorf("failed to add product to cart : %v", result.ErrorMsg)).Error()
			return ctx.Status(result.HttpErrorCode).JSON(fiber.Map{
				"code":   result.HttpErrorCode,
				"reason": fmt.Sprintf("[%d] %s", result.HttpErrorCode, result.ErrorMsg),
				"data":   nil,
			})
		}

		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"code": fiber.StatusOK,
			"msg":  "success add product to cart",
			"data": result.Response,
		})
	}
}

func (h *Handler) ListCart() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		param := new(request.AddCartRequest)
		if err := ctx.QueryParser(param); err != nil {
			logrus.WithError(fmt.Errorf("parse param: %v", err)).Error()
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"code":   fiber.StatusBadRequest,
				"reason": fmt.Sprintf("[%d] %s", fiber.StatusBadRequest, err.Error()),
				"data":   nil,
			})
		}

		result := h.useCase.ListCart(request.AddCartRequest{IdCustomer: param.IdCustomer})
		if result.ErrorMsg != nil {
			logrus.WithError(fmt.Errorf("failed to get cart : %v", result.ErrorMsg)).Error()
			return ctx.Status(result.HttpErrorCode).JSON(fiber.Map{
				"code":   result.HttpErrorCode,
				"reason": fmt.Sprintf("[%d] %s", result.HttpErrorCode, result.ErrorMsg),
				"data":   nil,
			})
		}

		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"code": fiber.StatusOK,
			"msg":  "success get cart",
			"data": result.Response,
		})
	}
}

func (h *Handler) DeleteCart() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		body := new(request.DeleteCartRequest)
		if err := ctx.BodyParser(body); err != nil {
			logrus.WithError(fmt.Errorf("parse body: %v", err)).Error()
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"code":   fiber.StatusBadRequest,
				"reason": fmt.Sprintf("[%d] %s", fiber.StatusBadRequest, err.Error()),
				"data":   nil,
			})
		}

		body.Id = ctx.Params("id")
		if body.Id == "" {
			logrus.WithError(fmt.Errorf("id required")).Error()
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"code":   fiber.StatusBadRequest,
				"reason": fmt.Sprintf("[%d] %s", fiber.StatusBadRequest, "id required"),
				"data":   nil,
			})
		}

		result := h.useCase.DeleteCart(*body)
		if result.ErrorMsg != nil {
			logrus.WithError(fmt.Errorf("failed to delete cart : %v", result.ErrorMsg)).Error()
			return ctx.Status(result.HttpErrorCode).JSON(fiber.Map{
				"code":   result.HttpErrorCode,
				"reason": fmt.Sprintf("[%d] %s", result.HttpErrorCode, result.ErrorMsg),
				"data":   nil,
			})
		}

		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"code": fiber.StatusOK,
			"msg":  "success delete cart",
			"data": result.Response,
		})
	}
}

func (h *Handler) Order() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		body := new(request.OrderRequest)
		if err := ctx.BodyParser(body); err != nil {
			logrus.WithError(fmt.Errorf("parse body: %v", err)).Error()
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"code":   fiber.StatusBadRequest,
				"reason": fmt.Sprintf("[%d] %s", fiber.StatusBadRequest, err.Error()),
				"data":   nil,
			})
		}

		result := h.useCase.Order(*body)
		if result.ErrorMsg != nil {
			logrus.WithError(fmt.Errorf("failed to insert order : %v", result.ErrorMsg)).Error()
			return ctx.Status(result.HttpErrorCode).JSON(fiber.Map{
				"code":   result.HttpErrorCode,
				"reason": fmt.Sprintf("[%d] %s", result.HttpErrorCode, result.ErrorMsg),
				"data":   nil,
			})
		}

		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"code": fiber.StatusOK,
			"msg":  "success insert order",
			"data": result.Response,
		})
	}
}

func (h *Handler) Register() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		body := new(databasesModel.Customer)
		if err := ctx.BodyParser(body); err != nil {
			logrus.WithError(fmt.Errorf("parse body: %v", err)).Error()
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"code":   fiber.StatusBadRequest,
				"reason": fmt.Sprintf("[%d] %s", fiber.StatusBadRequest, err.Error()),
				"data":   nil,
			})
		}

		result := h.useCase.AddCustomer(*body)
		if result.ErrorMsg != nil {
			logrus.WithError(fmt.Errorf("failed to add customer : %v", result.ErrorMsg)).Error()
			return ctx.Status(result.HttpErrorCode).JSON(fiber.Map{
				"code":   result.HttpErrorCode,
				"reason": fmt.Sprintf("[%d] %s", result.HttpErrorCode, result.ErrorMsg),
				"data":   nil,
			})
		}

		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"code": fiber.StatusOK,
			"msg":  "success add customer",
			"data": result.Response,
		})
	}
}
