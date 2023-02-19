package application

type CodeRepository interface {
	Save(code string, userID string) error
	Get(code string) (string, error)
	Close() error
}

type EmailSender interface {
	SendVerificationCode(email string, code string) error
	Close() error
}

type CodeGenerator interface {
	GenerateCode() string
}
