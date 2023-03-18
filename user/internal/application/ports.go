package application

// MessageProducer is an interface for sending messages to a queue
type MessageProducer interface {
	//SendNewSignInSessionMessage(ctx context.Context, params core.SendNewSignInSessionParams) error
}

// Validator is an interface for validating structs using tags
type Validator interface {
	Validate(s interface{}) (err error)
}
