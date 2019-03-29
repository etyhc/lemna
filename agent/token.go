package agent

type Token interface {
	GetSessionID(string) (int32, error)
}
