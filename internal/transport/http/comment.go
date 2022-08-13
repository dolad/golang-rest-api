package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/dolad/rest-api/internal/comment"
	"github.com/gorilla/mux"
)

// Get Comment by id
func (h *Handler) GetComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	i, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		sendErrorResponse(w, "Unable to pass uint from ID", err)
		return
	}
	comment, err := h.Service.GetComment(uint(i))
	if err != nil {
		sendErrorResponse(w, "Error Retrieving Comment By ID", err)
		return
	}
	if err := json.NewEncoder(w).Encode(comment); err != nil {
		panic(err)
	}
}

// GetAllComment
func (h *Handler) GetAllComments(w http.ResponseWriter, r *http.Request) {
	comments, err := h.Service.GetAllComments()
	if err != nil {
		sendErrorResponse(w, "Failed to post new comment", err)
		return
	}
	fmt.Fprintf(w, "%+v", comments)
}

// PostComment
func (h *Handler) PostComment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	var comment comment.Comment
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		sendErrorResponse(w, "Failed to decode JSON Body", err)
	}

	comment, err := h.Service.PostComment(comment)
	if err != nil {
		sendErrorResponse(w, "Failed to post new comment", err)
	}
	if err := json.NewEncoder(w).Encode(comment); err != nil {
		panic(err)
	}
}

// UpdateComment by Id
func (h *Handler) UpdateComment(w http.ResponseWriter, r *http.Request) {

	var cmt comment.Comment

	if err := json.NewDecoder(r.Body).Decode(&cmt); err != nil {
		sendErrorResponse(w, "Failed to decode JSON Body", err)
		return
	}
	vars := mux.Vars(r)
	id := vars["id"]
	commentID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		sendErrorResponse(w, "Failed to parse uint from ID", err)
		return
	}
	cmt, err = h.Service.UpdateComment(uint(commentID), cmt)

	if err != nil {
		sendErrorResponse(w, "Failed to update new comment", err)
		return
	}

	if err := json.NewEncoder(w).Encode(cmt); err != nil {
		panic(err)
	}
}

// DeleteComment by Id
func (h *Handler) DeleteComment(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	commentID, err := strconv.ParseUint(id, 10, 64)

	if err != nil {
		sendErrorResponse(w, "Failed to parse uint from ID", err)
		return
	}

	err = h.Service.DeleteComment(uint(commentID))

	if err != nil {
		sendErrorResponse(w, "Failed to delete comment by comment ID", err)
		return
	}

	if err = sendOkResponse(w, Response{Message: "Successfully Deleted"}); err != nil {
		panic(err)
	}
}
