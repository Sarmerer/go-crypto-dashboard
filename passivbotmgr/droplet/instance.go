package droplet

type PassivbotInstance interface {
	GetID() string

	GetPID() int
	GetPIDSignature() string

	GetRunCommand() []string
	GetArgs() []string
	GetFlags() []string
	GetStatus() string
	IsRunning() bool
	Match(query string) bool
}

type passivbotInstance struct {
	User       string
	Symbol     string
	ConfigPath string
	Flags      map[string]string
}
