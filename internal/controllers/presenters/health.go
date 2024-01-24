package presenters

type HealthCheckStatus string

const (
	HealthCheckStatusUp   HealthCheckStatus = "up"
	HealthCheckStatusDown HealthCheckStatus = "down"
)

type HealthCheckResponse struct {
	Status HealthCheckStatus `json:"status" example:"down"`
}
