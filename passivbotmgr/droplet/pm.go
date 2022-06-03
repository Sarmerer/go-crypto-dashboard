package droplet

import (
	"errors"
	"os/exec"
	"runtime"
)

type ProcessManager interface {
	Start(command []string) (int, error)
	Stop(pid int) error
	Restart(pid int) error

	GetPIDBySignature(signature string) int
	IsRunning(pid int) bool
	Info(pid int) (string, error)
}

type processManager struct{}

func NewProcessManager() ProcessManager {
	return &processManager{}
}

func (pm *processManager) Start(command []string) (int, error) {
	if len(command) == 0 {
		return -1, errors.New("no command provided")
	}

	cmd := exec.Command(command[0], command[1:]...)
	err := cmd.Start()
	if err != nil {
		return -1, err
	}

	return cmd.Process.Pid, nil
}

func (pm *processManager) Stop(pid int) error {
	return nil
}

func (pm *processManager) Restart(pid int) error {
	return nil
}

func (pm *processManager) GetPIDBySignature(signature string) int {
	return 0
}

func (pm *processManager) IsRunning(pid int) bool {
	return false
}

func (pm *processManager) Info(pid int) (string, error) {
	// get process info
	// if goos == linux
	// 	get process info from /proc/pid/stat
	// else if goos == windows
	// 	get process info from tasklist

	if runtime.GOOS == "windows" {
		cmd := exec.Command("tasklist /fi \"pid eq " + pid + "\"")
		out, err := cmd.Output()
		if err != nil {
			return "", err
		}

		return string(out), nil
	}

	if runtime.GOOS == "linux" {
		cmd := exec.Command("ps -p " + pid)
		out, err := cmd.Output()
		if err != nil {
			return "", err
		}

		return string(out), nil
	}
}
