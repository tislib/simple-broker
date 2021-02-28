package model

type RecordData = map[string]interface{}

type Record struct {
	MessageId  string      `json:"messageId"`
	BrokerId   string      `json:"brokerId"`
	ConsumerId string      `json:"consumerId"`
	Message    *RecordData `json:"message"`
}
