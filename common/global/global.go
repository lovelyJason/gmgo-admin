package global

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/songzhibin97/gkit/cache/local_cache"
)

var (
	BlackCache local_cache.Cache
	Trans      ut.Translator
)
