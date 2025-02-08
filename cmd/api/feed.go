package main

import (
	"net/http"

	"github.com/AlfanDutaPamungkas/Go-Social/internal/store"
)

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
