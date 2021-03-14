package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/mp-hl-2021/splinter/internal/app/splinter/model"
	"log"
	"net/http"
	"reflect"
	"strings"
)

type Api struct {
	model.UserInterface
}

func NewApi(x model.UserInterface) Api {
	return Api{x}
}

type Context struct {
	*Api
	wr  http.ResponseWriter
	req *http.Request
}

type ErrorResponse struct {
	Error      string
	StatusCode int
}

func (ctx *Context) abortf(statusCode int, format string, a ...interface{}) {
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

func (ctx *Context) abort(err error) {
	ctx.abortf(http.StatusBadRequest, "%s", err.Error())
}

func (ctx *Context) token() model.Token {
	header := ctx.req.Header.Get("Authorization")
	if strings.HasPrefix(header, "Bearer ") {
		return model.Token(strings.TrimPrefix(header, "Bearer "))
	} else {
		return model.UnauthenticatedToken
	}
}

func (ctx *Context) decodeJson(v interface{}) {
	if ctx.req.Header.Get("Content-Type") != "application/json" {
		ctx.abortf(http.StatusBadRequest, "expected json data")
	}

	if err := json.NewDecoder(ctx.req.Body).Decode(v); err != nil {
		ctx.abortf(http.StatusBadRequest, "invalid json data: %e", err)
	}
}

func (ctx *Context) writeJson(v interface{}) {
	ctx.wr.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(ctx.wr).Encode(v); err != nil {
		log.Fatal(err)
	}
}

func (ctx *Context) ok() {
	ctx.wr.WriteHeader(http.StatusNoContent)
}

func (ctx *Context) do(err error) {
	if err != nil {
		ctx.abort(err)
	}
	ctx.ok()
}

func (a *Api) MakeHandler(name string) http.Handler {
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

func (a Api) Router() http.Handler {
	router := mux.NewRouter()

	router.Handle("/register", a.MakeHandler("Register")).Methods(http.MethodPost)
	router.Handle("/authenticate", a.MakeHandler("Authenticate")).Methods(http.MethodPost)

	router.Handle("/snippets", a.MakeHandler("Snippets")).Methods(http.MethodGet, http.MethodPost)
	router.Handle("/snippets/{id}", a.MakeHandler("SnippetById")).Methods(http.MethodGet, http.MethodDelete)
	router.Handle("/snippets/{id}/comments", a.MakeHandler("SnippetComments")).Methods(http.MethodGet, http.MethodPost)
	router.Handle("/snippets/lang/{lang}", a.MakeHandler("SnippetsByLang")).Methods(http.MethodGet)

	router.Handle("/users/{id}", a.MakeHandler("UserById")).Methods(http.MethodGet)
	router.Handle("/users/{id}/snippets", a.MakeHandler("UserSnippets")).Methods(http.MethodGet)

	router.Handle("/subscriptions", a.MakeHandler("Subscriptions")).Methods(http.MethodGet)
	router.Handle("/subscribe/user/{id}", a.MakeHandler("SubscribeToUser")).Methods(http.MethodPost)
	router.Handle("/unsubscribe/user/{id}", a.MakeHandler("UnsubscribeFromUser")).Methods(http.MethodPost)
	router.Handle("/subscribe/lang/{id}", a.MakeHandler("SubscribeToLang")).Methods(http.MethodPost)
	router.Handle("/unsubscribe/lang/{id}", a.MakeHandler("UnsubscribeFromLang")).Methods(http.MethodPost)

	return router
}

func (ctx *Context) HandleRegister() {
	request := model.RegisterRequest{}
	ctx.decodeJson(&request)
	response, err := ctx.Register(request)
	if err != nil {
		ctx.abort(err)
	}
	ctx.writeJson(response)
}

func (ctx *Context) HandleAuthenticate() {
	request := model.AuthenticateRequest{}
	ctx.decodeJson(&request)
	response, err := ctx.Authenticate(request)
	if err != nil {
		ctx.abort(err)
	}
	ctx.writeJson(response)
}

func (ctx *Context) HandleSnippetById() {
	id := mux.Vars(ctx.req)["id"]

	if ctx.req.Method == http.MethodGet {
		snippet, err := ctx.GetSnippetById(ctx.token(), id)
		if err != nil {
			ctx.abort(err)
		}
		ctx.writeJson(snippet)
	}

	if ctx.req.Method == http.MethodDelete {
		ctx.do(ctx.DeleteSnippetById(ctx.token(), id))
	}
}

func (ctx *Context) HandleSnippets() {
	if ctx.req.Method == http.MethodGet {
		snippets, err := ctx.GetSnippetsFeed(ctx.token())
		if err != nil {
			ctx.abort(err)
		}
		ctx.writeJson(snippets)
	}

	if ctx.req.Method == http.MethodPost {
		request := model.PostSnippetRequest{}
		ctx.decodeJson(&request)
		snippet, err := ctx.PostSnippet(ctx.token(), request)
		if err != nil {
			ctx.abort(err)
		}
		ctx.writeJson(snippet)
	}
}

func (ctx *Context) HandleUserById() {
	id := mux.Vars(ctx.req)["id"]
	user, err := ctx.GetUserById(ctx.token(), id)
	if err != nil {
		ctx.abort(err)
	}
	ctx.writeJson(user)
}

func (ctx *Context) HandleUserSnippets() {
	id := mux.Vars(ctx.req)["id"]
	snippets, err := ctx.GetSnippetsByUser(ctx.token(), id)
	if err != nil {
		ctx.abort(err)
	}
	ctx.writeJson(snippets)
}

func (ctx *Context) HandleSnippetsByLang() {
	lang := model.ProgrammingLanguage(mux.Vars(ctx.req)["lang"])
	snippets, err := ctx.GetSnippetsByLanguage(ctx.token(), lang)
	if err != nil {
		ctx.abort(err)
	}
	ctx.writeJson(snippets)
}

func (ctx *Context) HandleSnippetComments() {
	id := mux.Vars(ctx.req)["id"]

	if ctx.req.Method == http.MethodGet {
		comments, err := ctx.GetCommentsBySnippetId(ctx.token(), id)
		if err != nil {
			ctx.abort(err)
		}
		ctx.writeJson(comments)
	}

	if ctx.req.Method == http.MethodPost {
		request := model.PostCommentRequest{}
		ctx.decodeJson(&request)
		comment, err := ctx.PostComment(ctx.token(), id, request)
		if err != nil {
			ctx.abort(err)
		}
		ctx.writeJson(comment)
	}
}

func (ctx *Context) HandleSubscriptions() {
	subscriptions, err := ctx.GetSubscriptions(ctx.token())
	if err != nil {
		ctx.abort(err)
	}
	ctx.writeJson(subscriptions)
}

func (ctx *Context) HandleSubscribeToUser() {
	id := mux.Vars(ctx.req)["id"]
	ctx.do(ctx.SubscribeToUser(ctx.token(), id))
}

func (ctx *Context) HandleUnsubscribeFromUser() {
	id := mux.Vars(ctx.req)["id"]
	ctx.do(ctx.UnsubscribeFromUser(ctx.token(), id))
}

func (ctx *Context) HandleSubscribeToLang() {
	lang := model.ProgrammingLanguage(mux.Vars(ctx.req)["lang"])
	ctx.do(ctx.SubscribeToLanguage(ctx.token(), lang))
}

func (ctx *Context) HandleUnsubscribeFromLang() {
	lang := model.ProgrammingLanguage(mux.Vars(ctx.req)["lang"])
	ctx.do(ctx.UnsubscribeFromLanguage(ctx.token(), lang))
}
