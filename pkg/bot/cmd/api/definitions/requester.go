package definitions

type Requester interface {
	Serialize() ([]byte, error)
	Validate() error
}
