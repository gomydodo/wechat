package wxofficial

import (
	"errors"
)

var (
	RequestIllegalError   = errors.New("illegal request")
	ArticleCountOverError = errors.New("article count over 10")
)
