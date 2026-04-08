package services

type HealthService interface {
}

type healthService struct {
}

func NewHealthService() HealthService {
	return &healthService{}
}
