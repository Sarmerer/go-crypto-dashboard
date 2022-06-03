package droplet

type Droplet interface {
	Start(instances []PassivbotInstance) error
	Stop(instances []PassivbotInstance) error
	Restart(instances []PassivbotInstance) error

	GetInstances() []PassivbotInstance
	GetInstanceByID(id string) PassivbotInstance
	GetInstanceByPID(pid int) PassivbotInstance
	GetInstanceByPIDSignature(pidSignature string) PassivbotInstance
	GetInstanceByQuery(query string) []PassivbotInstance
}

type droplet struct {
	port int
}

func New(port int) Droplet {
	return &droplet{port: port}
}

func (d *droplet) Start(instances []PassivbotInstance) error {
	return nil
}

func (d *droplet) Stop(instances []PassivbotInstance) error {
	return nil
}

func (d *droplet) Restart(instances []PassivbotInstance) error {
	return nil
}

func (d *droplet) GetInstances() []PassivbotInstance {
	return nil
}

func (d *droplet) GetInstanceByID(id string) PassivbotInstance {
	return nil
}

func (d *droplet) GetInstanceByPID(pid int) PassivbotInstance {
	return nil
}

func (d *droplet) GetInstanceByPIDSignature(pidSignature string) PassivbotInstance {
	return nil
}

func (d *droplet) GetInstanceByQuery(query string) []PassivbotInstance {
	return nil
}
