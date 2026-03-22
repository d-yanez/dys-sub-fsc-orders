package services

func BuildIdempotencyKey(eventType, orderID string) string {
	return "falabella:" + eventType + ":" + orderID
}
