package dispatch

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"synapsis-project/database/databasesModel"
	"synapsis-project/domain"
	"synapsis-project/helper"
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
				"code": fiber.StatusBadRequest,
				"msg":  fmt.Sprintf("[%d] %s", fiber.StatusBadRequest, err.Error()),
				"data": nil,
			})
		}

		result := h.useCase.ListProduct(*param)
		if result.ErrorMsg != nil {
			logrus.WithError(fmt.Errorf("failed to get list product : %v", result.ErrorMsg)).Error()
			return ctx.Status(result.HttpErrorCode).JSON(fiber.Map{
				"code": result.HttpErrorCode,
				"msg":  fmt.Sprintf("[%d] %s", result.HttpErrorCode, result.ErrorMsg),
				"data": nil,
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
				"code": fiber.StatusBadRequest,
				"msg":  fmt.Sprintf("[%d] %s", fiber.StatusBadRequest, err.Error()),
				"data": nil,
			})
		}

		claims, err := helper.ExtractData(ctx)
		if err != nil {
			logrus.WithError(fmt.Errorf("extract token: %v", err)).Error()
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"code": fiber.StatusUnauthorized,
				"msg":  fmt.Sprintf("[%d] %s", fiber.StatusUnauthorized, err.Error()),
				"data": nil,
			})
		}

		body.IdCustomer = claims.IdCustomer

		result := h.useCase.AddCart(*body)
		if result.ErrorMsg != nil {
			logrus.WithError(fmt.Errorf("failed to add product to cart : %v", result.ErrorMsg)).Error()
			return ctx.Status(result.HttpErrorCode).JSON(fiber.Map{
				"code": result.HttpErrorCode,
				"msg":  fmt.Sprintf("[%d] %s", result.HttpErrorCode, result.ErrorMsg),
				"data": nil,
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
		claims, err := helper.ExtractData(ctx)
		if err != nil {
			logrus.WithError(fmt.Errorf("extract token: %v", err)).Error()
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"code": fiber.StatusUnauthorized,
				"msg":  fmt.Sprintf("[%d] %s", fiber.StatusUnauthorized, err.Error()),
				"data": nil,
			})
		}

		result := h.useCase.ListCart(request.AddCartRequest{IdCustomer: claims.IdCustomer})
		if result.ErrorMsg != nil {
			logrus.WithError(fmt.Errorf("failed to get cart : %v", result.ErrorMsg)).Error()
			return ctx.Status(result.HttpErrorCode).JSON(fiber.Map{
				"code": result.HttpErrorCode,
				"msg":  fmt.Sprintf("[%d] %s", result.HttpErrorCode, result.ErrorMsg),
				"data": nil,
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
		claims, err := helper.ExtractData(ctx)
		if err != nil {
			logrus.WithError(fmt.Errorf("extract token: %v", err)).Error()
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"code": fiber.StatusUnauthorized,
				"msg":  fmt.Sprintf("[%d] %s", fiber.StatusUnauthorized, err.Error()),
				"data": nil,
			})
		}

		id := ctx.Params("id")
		if id == "" {
			logrus.WithError(fmt.Errorf("id required")).Error()
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"code": fiber.StatusBadRequest,
				"msg":  fmt.Sprintf("[%d] %s", fiber.StatusBadRequest, "id required"),
				"data": nil,
			})
		}

		result := h.useCase.DeleteCart(request.DeleteCartRequest{
			Id:         id,
			IdCustomer: claims.IdCustomer,
		})

		if result.ErrorMsg != nil {
			logrus.WithError(fmt.Errorf("failed to delete cart : %v", result.ErrorMsg)).Error()
			return ctx.Status(result.HttpErrorCode).JSON(fiber.Map{
				"code": result.HttpErrorCode,
				"msg":  fmt.Sprintf("[%d] %s", result.HttpErrorCode, result.ErrorMsg),
				"data": nil,
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
		claims, err := helper.ExtractData(ctx)
		if err != nil {
			logrus.WithError(fmt.Errorf("extract token: %v", err)).Error()
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"code": fiber.StatusUnauthorized,
				"msg":  fmt.Sprintf("[%d] %s", fiber.StatusUnauthorized, err.Error()),
				"data": nil,
			})
		}

		result := h.useCase.Order(request.OrderRequest{
			IdCustomer: claims.IdCustomer,
		})
		if result.ErrorMsg != nil {
			logrus.WithError(fmt.Errorf("failed to insert order : %v", result.ErrorMsg)).Error()
			return ctx.Status(result.HttpErrorCode).JSON(fiber.Map{
				"code": result.HttpErrorCode,
				"msg":  fmt.Sprintf("[%d] %s", result.HttpErrorCode, result.ErrorMsg),
				"data": nil,
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
				"code": fiber.StatusBadRequest,
				"msg":  fmt.Sprintf("[%d] %s", fiber.StatusBadRequest, err.Error()),
				"data": nil,
			})
		}

		result := h.useCase.AddCustomer(*body)
		if result.ErrorMsg != nil {
			logrus.WithError(fmt.Errorf("failed to add customer : %v", result.ErrorMsg)).Error()
			return ctx.Status(result.HttpErrorCode).JSON(fiber.Map{
				"code": result.HttpErrorCode,
				"msg":  fmt.Sprintf("[%d] %s", result.HttpErrorCode, result.ErrorMsg),
				"data": nil,
			})
		}

		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"code": fiber.StatusOK,
			"msg":  "success add customer",
			"data": result.Response,
		})
	}
}

func (h *Handler) Login() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		body := new(databasesModel.Customer)
		if err := ctx.BodyParser(body); err != nil {
			logrus.WithError(fmt.Errorf("parse body: %v", err)).Error()
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"code": fiber.StatusBadRequest,
				"msg":  fmt.Sprintf("[%d] %s", fiber.StatusBadRequest, err.Error()),
				"data": nil,
			})
		}

		result := h.useCase.Login(*body)
		if result.ErrorMsg != nil {
			logrus.WithError(fmt.Errorf("failed to login : %v", result.ErrorMsg)).Error()
			return ctx.Status(result.HttpErrorCode).JSON(fiber.Map{
				"code": result.HttpErrorCode,
				"msg":  fmt.Sprintf("[%d] %s", result.HttpErrorCode, result.ErrorMsg),
				"data": nil,
			})
		}

		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"code": fiber.StatusOK,
			"msg":  "success login",
			"data": result.Response,
		})
	}
}
