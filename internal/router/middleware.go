package router

import (
	"fmt"
	"net/http"
	"strings"

	"firebase.google.com/go/v4/appcheck"
	"github.com/gofiber/fiber/v2"
	"github.com/prplx/cnvrt/internal/helpers"
	"github.com/prplx/cnvrt/internal/types"
	"github.com/prplx/cnvrt/internal/validator"
)

func requireAppCheck(appCheck *appcheck.Client, appCheckToken string) error {
	if _, err := appCheck.VerifyToken(appCheckToken); err != nil {
		return fmt.Errorf("AppCheck token verification failed: %v", err)
	}

	return nil
}

func checkAppCheckToken(ctx *fiber.Ctx) error {
	config := ctx.Locals("config").(*types.Config)
	if helpers.IsTest() || strings.Contains(ctx.Path(), "/ws") || strings.Contains(ctx.Path(), healthcheckEndpoint) {
		return ctx.Next()
	}

	var appCheckToken string

	if strings.Contains(ctx.Path(), "/ws") {
		appCheckToken = ctx.Query(firebaseAppCheckQuery)
	} else {
		appCheckToken = ctx.Get(config.Firebase.AppCheckHeader)
	}

	if err := requireAppCheck(appCheck, appCheckToken); err != nil {
		return ctx.Status(http.StatusForbidden).SendString("Forbidden")
	}

	return ctx.Next()
}

func checkFormFileLength(ctx *fiber.Ctx, config *types.Config) error {
	form, err := ctx.MultipartForm()
	if err != nil {
		return ctx.SendStatus(http.StatusBadRequest)
	}

	images := form.File["image"]
	if len(images) > config.App.MaxFileCount || len(images) == 0 {
		return ctx.SendStatus(http.StatusBadRequest)
	}

	return ctx.Next()
}

func checkQueryParams(ctx *fiber.Ctx, params ...string) error {
	v := validator.NewValidator()
	if validateRequestQueryParams(v, ctx, params...); !v.Valid() {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"errors": v.Errors,
		})
	}

	return ctx.Next()
}

func validateRequestQueryParams(v *validator.Validator, ctx *fiber.Ctx, requiredParams ...string) {
	for _, param := range requiredParams {
		v.Check(ctx.Query(param) != "", param, param+" is required")
	}
}
