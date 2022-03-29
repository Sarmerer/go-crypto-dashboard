package hosting

type Provider interface {
	Start() error
	Stop() error
}
