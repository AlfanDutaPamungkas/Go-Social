package main

import (
	"net/http"

	"github.com/AlfanDutaPamungkas/Go-Social/internal/store"
)

type CreateCommentPayload struct {
	Content string `json:"content" validate:"required"`
}

// CreateComment godoc
//
//	@Summary		Creates a new comment
//	@Description	Creates a new comment in post
//	@Tags			comments
//	@Accept			json
//	@Produce		json
//	@Param			postID	path		int						true	"Post ID"
//	@Param			body	body		CreateCommentPayload	true	"Created comment data"
//	@Success		201		{object}	store.Comment
//	@Failure		400		{object}	error	"Invalid request payload"
//	@Failure		404		{object}	error	"Post not found"
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/posts/{postID}/comment [post]
func (app *application) createCommentHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromCtx(r)
	post := getPostFromCtx(r)

	var payload CreateCommentPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	comment := &store.Comment{
		UserID:  user.ID,
		PostID:  post.ID,
		Content: payload.Content,
	}

	if err := app.store.Comments.Create(r.Context(), comment); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, comment); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
