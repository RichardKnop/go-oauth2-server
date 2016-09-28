package response

import (
	"github.com/RichardKnop/jsonhal"
)

// ListResponse ...
type ListResponse struct {
	jsonhal.Hal
	Count uint `json:"count"`
	Page  uint `json:"page"`
}

// NewListResponse creates new ListResponse instance
func NewListResponse(count, page int, self, first, last, previous, next, embedName string, items interface{}) *ListResponse {
	response := &ListResponse{
		Count: uint(count),
		Page:  uint(page),
	}

	response.SetLink("self", self, "")
	response.SetLink("first", first, "")
	response.SetLink("last", last, "")
	response.SetLink("prev", previous, "")
	response.SetLink("next", next, "")

	response.SetEmbedded(embedName, jsonhal.Embedded(items))

	return response
}
