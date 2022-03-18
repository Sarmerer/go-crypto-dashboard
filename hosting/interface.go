package hosting

type Hosting interface {
	Start() error
	Stop() error
	
}
