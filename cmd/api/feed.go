package main

import (
	"net/http"

	"github.com/AlfanDutaPamungkas/Go-Social/internal/store"
)

// GetUserFeed godoc
//
//	@Summary		Get user feed
//	@Description	Retrieves a paginated feed of posts for the user, with filtering options
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			limit	query		int			false	"Number of posts to retrieve (1-20)"	default(20)
//	@Param			offset	query		int			false	"Pagination offset (>=0)"				default(0)
//	@Param			sort	query		string		false	"Sort order (asc or desc)"				default(desc)	Enums(asc, desc)
//	@Param			tags	query		[]string	false	"Filter by up to 5 tags"
//	@Param			search	query		string		false	"Search query (max 100 chars)"
//	@Param			since	query		string		false	"Start date (RFC3339 format)"
//	@Param			until	query		string		false	"End date (RFC3339 format)"
//	@Success		200		{array}		store.Post
//	@Failure		400		{object}	error	"Invalid request parameters"
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/feed [get]
func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {
	p := &store.PaginatedFeedQuery{
		Limit: 20,
		Offset: 0,
		Sort: "desc",
	}

	p, err := p.Parse(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(p); err != nil{
		app.badRequestResponse(w, r, err)
		return
	}
	
	feed, err := app.store.Posts.GetUserFeed(r.Context(), int64(21), p)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, feed); err != nil {
		app.internalServerError(w, r, err)
	}
}
