package backend

type Provider interface {
	ListAll(dir string) ([]string, error)
	Search(input string, dir string) ([]string, error)
	Flush()
}
