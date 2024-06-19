// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen/v2 version v2.1.0 DO NOT EDIT.
package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/oapi-codegen/runtime"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

const (
	BearerAuthScopes = "bearerAuth.Scopes"
)

// Post defines model for Post.
type Post struct {
	Author  *string `json:"author,omitempty"`
	Content *string `json:"content,omitempty"`
	Id      *int64  `json:"id,omitempty"`
}

// AuthSuccess defines model for AuthSuccess.
type AuthSuccess struct {
	ExpirationDate *time.Time `json:"expiration_date,omitempty"`
	Token          *string    `json:"token,omitempty"`
}

// PostLoginJSONBody defines parameters for PostLogin.
type PostLoginJSONBody struct {
	Password *string `json:"password,omitempty"`
	Username *string `json:"username,omitempty"`
}

// GetPostsParams defines parameters for GetPosts.
type GetPostsParams struct {
	// FirstId ID of the first post needed
	FirstId int `form:"first_id" json:"first_id"`

	// MaxPosts Maximum number of posts to retrieve
	MaxPosts int `form:"max_posts" json:"max_posts"`
}

// PostPostsJSONBody defines parameters for PostPosts.
type PostPostsJSONBody struct {
	Content *string `json:"content,omitempty"`
}

// PutPostsPostIdJSONBody defines parameters for PutPostsPostId.
type PutPostsPostIdJSONBody struct {
	Content *string `json:"content,omitempty"`
}

// PostRegisterJSONBody defines parameters for PostRegister.
type PostRegisterJSONBody struct {
	Password *string `json:"password,omitempty"`
	Username *string `json:"username,omitempty"`
}

// PutUpdateJSONBody defines parameters for PutUpdate.
type PutUpdateJSONBody struct {
	DateOfBirth *openapi_types.Date  `json:"date_of_birth,omitempty"`
	Email       *openapi_types.Email `json:"email,omitempty"`
	Name        *string              `json:"name,omitempty"`
	PhoneNumber *string              `json:"phone_number,omitempty"`
	Surname     *string              `json:"surname,omitempty"`
}

// PostLoginJSONRequestBody defines body for PostLogin for application/json ContentType.
type PostLoginJSONRequestBody PostLoginJSONBody

// PostPostsJSONRequestBody defines body for PostPosts for application/json ContentType.
type PostPostsJSONRequestBody PostPostsJSONBody

// PutPostsPostIdJSONRequestBody defines body for PutPostsPostId for application/json ContentType.
type PutPostsPostIdJSONRequestBody PutPostsPostIdJSONBody

// PostRegisterJSONRequestBody defines body for PostRegister for application/json ContentType.
type PostRegisterJSONRequestBody PostRegisterJSONBody

// PutUpdateJSONRequestBody defines body for PutUpdate for application/json ContentType.
type PutUpdateJSONRequestBody PutUpdateJSONBody

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// User login
	// (POST /login)
	PostLogin(ctx echo.Context) error
	// Retrieve list of posts with pagination
	// (GET /posts)
	GetPosts(ctx echo.Context, params GetPostsParams) error
	// Create new post
	// (POST /posts)
	PostPosts(ctx echo.Context) error
	// Delete post
	// (DELETE /posts/{post_id})
	DeletePostsPostId(ctx echo.Context, postId int64) error
	// Retrieve a post by ID
	// (GET /posts/{post_id})
	GetPostsPostId(ctx echo.Context, postId int64) error
	// Update post
	// (PUT /posts/{post_id})
	PutPostsPostId(ctx echo.Context, postId int64) error
	// Send a post like
	// (POST /posts/{post_id}/like)
	PostPostsPostIdLike(ctx echo.Context, postId int64) error
	// Send a post view
	// (POST /posts/{post_id}/view)
	PostPostsPostIdView(ctx echo.Context, postId int64) error
	// Register a new user
	// (POST /register)
	PostRegister(ctx echo.Context) error
	// Update user data
	// (PUT /update)
	PutUpdate(ctx echo.Context) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// PostLogin converts echo context to params.
func (w *ServerInterfaceWrapper) PostLogin(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PostLogin(ctx)
	return err
}

// GetPosts converts echo context to params.
func (w *ServerInterfaceWrapper) GetPosts(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{})

	// Parameter object where we will unmarshal all parameters from the context
	var params GetPostsParams
	// ------------- Required query parameter "first_id" -------------

	err = runtime.BindQueryParameter("form", true, true, "first_id", ctx.QueryParams(), &params.FirstId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter first_id: %s", err))
	}

	// ------------- Required query parameter "max_posts" -------------

	err = runtime.BindQueryParameter("form", true, true, "max_posts", ctx.QueryParams(), &params.MaxPosts)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter max_posts: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetPosts(ctx, params)
	return err
}

// PostPosts converts echo context to params.
func (w *ServerInterfaceWrapper) PostPosts(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PostPosts(ctx)
	return err
}

// DeletePostsPostId converts echo context to params.
func (w *ServerInterfaceWrapper) DeletePostsPostId(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "post_id" -------------
	var postId int64

	err = runtime.BindStyledParameterWithOptions("simple", "post_id", ctx.Param("post_id"), &postId, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter post_id: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.DeletePostsPostId(ctx, postId)
	return err
}

// GetPostsPostId converts echo context to params.
func (w *ServerInterfaceWrapper) GetPostsPostId(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "post_id" -------------
	var postId int64

	err = runtime.BindStyledParameterWithOptions("simple", "post_id", ctx.Param("post_id"), &postId, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter post_id: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetPostsPostId(ctx, postId)
	return err
}

// PutPostsPostId converts echo context to params.
func (w *ServerInterfaceWrapper) PutPostsPostId(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "post_id" -------------
	var postId int64

	err = runtime.BindStyledParameterWithOptions("simple", "post_id", ctx.Param("post_id"), &postId, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter post_id: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PutPostsPostId(ctx, postId)
	return err
}

// PostPostsPostIdLike converts echo context to params.
func (w *ServerInterfaceWrapper) PostPostsPostIdLike(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "post_id" -------------
	var postId int64

	err = runtime.BindStyledParameterWithOptions("simple", "post_id", ctx.Param("post_id"), &postId, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter post_id: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PostPostsPostIdLike(ctx, postId)
	return err
}

// PostPostsPostIdView converts echo context to params.
func (w *ServerInterfaceWrapper) PostPostsPostIdView(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "post_id" -------------
	var postId int64

	err = runtime.BindStyledParameterWithOptions("simple", "post_id", ctx.Param("post_id"), &postId, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter post_id: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PostPostsPostIdView(ctx, postId)
	return err
}

// PostRegister converts echo context to params.
func (w *ServerInterfaceWrapper) PostRegister(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PostRegister(ctx)
	return err
}

// PutUpdate converts echo context to params.
func (w *ServerInterfaceWrapper) PutUpdate(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PutUpdate(ctx)
	return err
}

// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "")
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.POST(baseURL+"/login", wrapper.PostLogin)
	router.GET(baseURL+"/posts", wrapper.GetPosts)
	router.POST(baseURL+"/posts", wrapper.PostPosts)
	router.DELETE(baseURL+"/posts/:post_id", wrapper.DeletePostsPostId)
	router.GET(baseURL+"/posts/:post_id", wrapper.GetPostsPostId)
	router.PUT(baseURL+"/posts/:post_id", wrapper.PutPostsPostId)
	router.POST(baseURL+"/posts/:post_id/like", wrapper.PostPostsPostIdLike)
	router.POST(baseURL+"/posts/:post_id/view", wrapper.PostPostsPostIdView)
	router.POST(baseURL+"/register", wrapper.PostRegister)
	router.PUT(baseURL+"/update", wrapper.PutUpdate)

}
