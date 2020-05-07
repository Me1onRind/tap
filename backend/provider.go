package backend

type Provider interface {
	ListAll() ([]string, error)
	ListDirs() []string
	CurrDir() string
	Search(input string) ([]string, error)
	Filepath(name string) string
	SetDir(dir string) error
	Flush()
}
