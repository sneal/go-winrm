package winrm

import (
	"errors"
	"io"
	"io/ioutil"
)

// Command represents a given command running on a Shell. This structure allows to get access
// to the various stdout, stderr and stdin pipes.
type Command struct {
	client    *Client
	shell     *Shell
	commandId string
	exitCode  int
	finished  bool

	Stdout io.Writer
	Stderr io.Writer

	done chan bool
}

func newCommand(shell *Shell, commandId string, stdout io.Writer, stderr io.Writer) *Command {
	command := &Command{shell: shell, client: shell.client, commandId: commandId, done: make(chan bool)}
	command.Stdout = stdout
	command.Stderr = stderr

	if command.Stdout == nil {
		command.Stdout = ioutil.Discard
	}
	if command.Stderr == nil {
		command.Stderr = ioutil.Discard
	}

	go fetchOutput(command)

	return command
}

func fetchOutput(command *Command) {
	for {
		select {
		case <-command.done:
			break
		default:
			finished, _ := command.slurpAllOutput()
			if finished {
				command.done <- true
				break
			}
		}
	}
}

// Close will terminate the running command
func (command *Command) Close() (err error) {
	if err = command.check(); err != nil {
		return err
	}

	request := NewSignalRequest(command.client.Parameters.url, command.shell.ShellId, command.commandId, &command.client.Parameters)
	defer request.Free()

	_, err = command.client.sendRequest(request)
	return err
}

func (command *Command) slurpAllOutput() (finished bool, err error) {
	if err = command.check(); err != nil {
		return true, err
	}

	request := NewGetOutputRequest(command.client.Parameters.url, command.shell.ShellId, command.commandId, "stdout stderr", &command.client.Parameters)
	defer request.Free()

	response, err := command.client.sendRequest(request)
	if err != nil {
		return true, err
	}

	var exitCode int
	finished, exitCode, err = ParseSlurpOutputErrResponse(response, command.Stdout, command.Stderr)
	if err != nil {
		return true, err
	}
	if finished {
		command.exitCode = exitCode
	}

	return
}

func (command *Command) sendInput(data []byte) (err error) {
	if err = command.check(); err != nil {
		return err
	}

	request := NewSendInputRequest(command.client.Parameters.url, command.shell.ShellId, command.commandId, data, &command.client.Parameters)
	defer request.Free()

	_, err = command.client.sendRequest(request)
	return
}

func (command *Command) check() (err error) {
	if command.commandId == "" {
		return errors.New("Command has already been closed")
	}
	if command.shell == nil {
		return errors.New("Command has no associated shell")
	}
	if command.client == nil {
		return errors.New("Command has no associated client")
	}
	return
}

// ExitCode returns command exit code when it is finished. Before that the result is always 0.
func (command *Command) ExitCode() int {
	return command.exitCode
}

// Calling this function will block the current goroutine until the remote command terminates.
func (command *Command) Wait() {
	// block until finished
	<-command.done
}
