package main

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/AlfanDutaPamungkas/Go-Social/internal/store"
	"github.com/go-chi/chi/v5"
)

type userKey string

const userCtx userKey = "user"

// GetUser godoc
//
//	@Summary		Fetches a user profile
//	@Description	Fetches a user profile by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@param			userID	path		int	true	"User ID"
//	@Success		200		{object}	store.User
//	@Failure		400		{object}	error
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/{userID}/ [get]
func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromCtx(r)

	if err := app.jsonResponse(w, http.StatusOK, user); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// FollowUser godoc
//
//	@Summary		Follows a user
//	@Description	Follows a user profile by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@param			userID	path	int	true	"User ID"
//	@Success		204
//	@Failure		400	{object}	error	"User payload missing"
//	@Failure		404	{object}	error	"User not found"
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/{userID}/follow [put]
func (app *application) followUserHandler(w http.ResponseWriter, r *http.Request) {
	followerUser := getUserFromCtx(r)

	type FollowUser struct {
		UserID int64 `json:"user_id"`
	}

	//! revert back to auth userID from ctx
	var payload FollowUser
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := app.store.Followers.Follow(r.Context(), followerUser.ID, payload.UserID); err != nil {
		switch {
		case errors.Is(err, store.ErrConflict):
			app.conflictResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// UnfollowUser godoc
//
//	@Summary		Unfollows a user
//	@Description	Unfollows a user profile by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@param			userID	path	int	true	"User ID"
//	@Success		204
//	@Failure		400	{object}	error	"User payload missing"
//	@Failure		404	{object}	error	"User not found"
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/{userID}/follow [delete]
func (app *application) unfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	unfollowedUser := getUserFromCtx(r)

	type FollowUser struct {
		UserID int64 `json:"user_id"`
	}

	//! revert back to auth userID from ctx
	var payload FollowUser
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := app.store.Followers.Unfollow(r.Context(), unfollowedUser.ID, payload.UserID); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (app *application) userContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
		if err != nil {
			app.badRequestResponse(w, r, err)
			return
		}

		user, err := app.store.Users.GetByID(ctx, userID)
		if err != nil {
			switch {
			case errors.Is(err, store.ErrNotFound):
				app.notFoundResponse(w, r, err)
			default:
				app.internalServerError(w, r, err)
			}
			return
		}

		ctx = context.WithValue(ctx, userCtx, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getUserFromCtx(r *http.Request) *store.User {
	user, _ := r.Context().Value(userCtx).(*store.User)
	return user
}
