package delete

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"

	resp "url-shortner/url-shortner/internal/lib/api/response"
	"url-shortner/url-shortner/internal/lib/logger/sl"
)

type Request struct {
	Alias string `json:"alias" validate:"required"`
}

type Response struct {
	resp.Response
	Result string `json:"result,omitempty"`
}

type URLDeleter interface {
	DeleteURL(alias string) (string, error)
}

func New(log *slog.Logger, urlDeleter URLDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.delete.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", sl.Err(err))

			//render.JSON(w, r, resp.Error("invalid request"))
			render.JSON(w, r, resp.ValidationError(validateErr))

			return
		}

		result, err := urlDeleter.DeleteURL(req.Alias)
		if err != nil {
			log.Error("failed to delete url", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to delete url"))

			return
		}

		log.Info("url deletes", slog.String("result:", result))

		responseOK(w, r, result)

	}
}

func responseOK(w http.ResponseWriter, r *http.Request, result string) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		Result:   result,
	})
}
