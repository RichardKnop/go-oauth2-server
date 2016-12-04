package response

import (
	"github.com/RichardKnop/jsonhal"
)

// LookupUintIDResponse ...
type LookupUintIDResponse struct {
	jsonhal.Hal
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// LookupStringIDResponse ...
type LookupStringIDResponse struct {
	jsonhal.Hal
	ID        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
