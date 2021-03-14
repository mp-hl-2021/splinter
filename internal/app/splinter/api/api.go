package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/mp-hl-2021/splinter/internal/app/splinter/model"
	"log"
	"net/http"
	"reflect"
)

type RestApi struct {
	model.UserInterface
}

type Context struct {
	*RestApi
	wr  http.ResponseWriter
	req *http.Request
}

type ErrorResponse struct {
	Error      string
	StatusCode int
}

func (ctx *Context) Abort(statusCode int, format string, a ...interface{}) {
	ctx.wr.WriteHeader(statusCode)
	t := ErrorResponse{
		fmt.Errorf(format, a...).Error(),
		statusCode,
	}

	if err := json.NewEncoder(ctx.wr).Encode(t); err != nil {
		log.Fatal(err)
	}

	panic(t)
}

func (ctx *Context) DecodeJson(v interface{}) {
	if ctx.req.Header.Get("Content-Type") != "application/json" {
		ctx.Abort(http.StatusBadRequest, "Expected json data")
	}

	if err := json.NewDecoder(ctx.req.Body).Decode(v); err != nil {
		ctx.Abort(http.StatusBadRequest, "Invalid json data: %e", err)
	}
}

func (ctx *Context) WriteJson(v interface{}) {
	ctx.wr.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(ctx.wr).Encode(v); err != nil {
		log.Fatal(err)
	}
}

func (a *RestApi) MakeHandler(name string) http.Handler {
	return http.HandlerFunc(func(wr http.ResponseWriter, req *http.Request) {
		ctx := Context{
			a, wr, req,
		}

		defer func() {
			if rec := recover(); rec != nil {
				switch rec.(type) {
				case ErrorResponse:
					return
				default:
					panic(rec)
				}
			}
		}()

		reflect.ValueOf(&ctx).MethodByName(fmt.Sprintf("Handle%s", name)).Call([]reflect.Value{})
	})
}

func (a *RestApi) Router() http.Handler {
	router := mux.NewRouter()

	router.Handle("/register", a.MakeHandler("Register")).Methods(http.MethodPost)
	router.Handle("/authenticate", a.MakeHandler("Authenticate")).Methods(http.MethodPost)

	return router
}

func (ctx *Context) HandleRegister() {
	request := model.RegisterRequest{}
	ctx.DecodeJson(&request)
	response, err := ctx.Register(request)
	if err != nil {
		ctx.Abort(http.StatusBadRequest, "%e", err)
	}
	ctx.WriteJson(response)
}

func (ctx *Context) HandleAuthenticate() {
	request := model.AuthenticateRequest{}
	ctx.DecodeJson(&request)
	response, err := ctx.Authenticate(request)
	if err != nil {
		ctx.Abort(http.StatusBadRequest, "%e", err)
	}
	ctx.WriteJson(response)
}
