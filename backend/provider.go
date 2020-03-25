package backend

type AudioItem interface {
	Name() string
}

type Provider interface {
	ListAll() ([]AudioItem, error)
	Search(input string) ([]AudioItem, error)
	Filepath(index int) (string, error)
	Flush()
}
