package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gilwong00/url-shortner/pkg/models"
)

func (h *Handler) AddTag(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	defer r.Body.Close()
	payload, err := ParseJSONBody[models.AddTagRequest](r)
	if err != nil {
		GenericErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	result, err := h.store.Get(ctx, payload.ShortID).Result()
	if err != nil {
		GenericErrorResponse(w, http.StatusNotFound, err.Error())
		return
	}
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(result), &data); err != nil {
		// If data is not an object, assume its a string
		// TODO: create a data model to store shortened urls and tags
		data["data"] = result
	}
	var tags []string
	if existingTags, ok := data["tags"].([]string); ok {
		tags = existingTags
	}
	if hasDupe := hasDuplicateTag(tags, payload.Tag); hasDupe {
		GenericErrorResponse(w, http.StatusBadRequest, "tag already exists")
		return
	}
	// append new tag
	tags = append(tags, payload.Tag)
	data["tags"] = tags
	updated, err := json.Marshal(data)
	if err != nil {
		GenericErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	if err = h.store.Set(ctx, payload.ShortID, updated, 0).Err(); err != nil {
		GenericErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	resp := models.AddTagResponse{
		Tags: tags,
	}
	w.Header().Set(contentType, appJSON)
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

func hasDuplicateTag(tags []string, tag string) bool {
	tagMap := make(map[string]struct{})
	for _, t := range tags {
		tagMap[t] = struct{}{}
	}
	if _, ok := tagMap[tag]; ok {
		return true
	}
	return false
}
