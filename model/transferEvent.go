package model

// TransferEvent is the Event
type TransferEvent struct {
	Sender    string `json:"sender"`
	Recipient string `json:"recipient"`
	Amount    int    `json:"amount"`
}

func NewTransferEvent(sender, recipient string, amount int) *TransferEvent {
	return &TransferEvent{
		Sender:    sender,
		Recipient: recipient,
		Amount:    amount,
	}
}
