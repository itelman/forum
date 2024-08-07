package check

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type ReactionRequest struct {
	UserID       string `json:"user_id"`
	ID           string `json:"id"`
	ReactionType string `json:"reactionType"`
}

func (app *application) handleLike(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/like/" {
		app.notFound(w, r)
		return
	}
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		app.clientError(w, r, http.StatusMethodNotAllowed)
		return
	}
	var request ReactionRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	num, err := strconv.Atoi(request.ID)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	err = app.post_reactions.Insert(request.ID, request.UserID, "1")
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	err = app.snippets.UpdateReactions(num, app.post_reactions.Likes, app.post_reactions.Dislikes)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/sneep?id=%s", request.ID), http.StatusSeeOther)
}

func (app *application) handleDislike(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/dislike/" {
		app.notFound(w, r)
		return
	}
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		app.clientError(w, r, http.StatusMethodNotAllowed)
		return
	}
	var request ReactionRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	num, err := strconv.Atoi(request.ID)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	err = app.post_reactions.Insert(request.ID, request.UserID, "0")
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// UpdateReactions expects integers as arguments
	err = app.snippets.UpdateReactions(num, app.post_reactions.Likes, app.post_reactions.Dislikes)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/sneep?id=%s", request.ID), http.StatusSeeOther)
}
