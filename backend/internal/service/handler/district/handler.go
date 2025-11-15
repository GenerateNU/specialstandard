package district

import (
	"specialstandard/internal/storage"
)

type Handler struct {
	districtRepository storage.DistrictRepository
}

func NewHandler(districtRepository storage.DistrictRepository) *Handler {
	return &Handler{
		districtRepository: districtRepository,
	}
}