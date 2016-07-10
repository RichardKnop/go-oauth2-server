package config

import (
	"github.com/RichardKnop/logging"
)

var logger *logging.Logger

func init() {
	logger = logging.New(nil, nil, new(logging.ColouredFormatter))
}
