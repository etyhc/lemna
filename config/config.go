package config

type Serialization interface {
	String() string
	Init(string) error
}

type Config interface {
	Publish(Serialization) error
	Subscribe(string, Serialization) (<-chan Serialization, error)
}
