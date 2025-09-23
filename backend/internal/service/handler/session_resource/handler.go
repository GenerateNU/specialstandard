package sessionresource

import (
	"specialstandard/internal/storage"
	"specialstandard/internal/xvalidator"
)

type Handler struct {
	sessionResourceRepository storage.SessionResourceRepository
	validator                 *xvalidator.XValidator
}
