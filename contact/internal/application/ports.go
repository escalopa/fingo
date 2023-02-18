package application

type CodeRepository interface {
	Save(code string, userID string) error
	Get(code string) (string, error)
}
