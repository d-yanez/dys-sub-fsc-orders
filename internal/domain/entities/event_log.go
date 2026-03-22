package entities

type EventLog struct {
	ID             string
	IdempotencyKey string
	Processed      bool
	Attempts       int
	Status         string
}
