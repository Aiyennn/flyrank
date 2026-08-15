package health

type HealthService struct {
}

func NewHealthService() *HealthService {
	return &HealthService{}
}

func (s *HealthService) healthCheck() map[string]string {
	return map[string]string{
		"status":  "ok",
		"service": "healthy",
	}
}