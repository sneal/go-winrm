# WinRM for Go

### _NOTE_ This library does not have a stable API, use at your own risk!

This is a Go library to execute remote commands on Windows machines through the use of WinRM. This library doesn't support domain users (it doesn't support GSSAPI nor Kerberos).


## Getting Started

### Preparing the remote windows machine

WinRM is available on Windows Server 2008 and up. This project supports only basic authentication for local accounts (domain users are not supported).
The remote windows system must be prepared for winrm:

_For a PowerShell script to do what is described below in one go, check [Richard Downer's blog](http://www.frontiertown.co.uk/2011/12/overthere-control-windows-from-java/)_

On the remote host, open a Command Prompt (not a PowerShell prompt!) using the __Run as Administrator__ option and paste in the following lines:

		winrm quickconfig
		y
		winrm set winrm/config/service/Auth @{Basic="true"}
		winrm set winrm/config/service @{AllowUnencrypted="true"}
		winrm set winrm/config/winrs @{MaxMemoryPerShellMB="1024"}

__N.B.:__ The Windows Firewall needs to be running to run this command. See [Microsoft Knowledge Base article #2004640](http://support.microsoft.com/kb/2004640).

__N.B.:__ Do not disable Negotiate authentication as the `winrm` command itself uses this for internal authentication, and you risk getting a system where `winrm` doesn't work anymore.
	
__N.B.:__ The `MaxMemoryPerShellMB` option has no effects on some Windows 2008R2 systems because of a WinRM bug. Make sure to install the hotfix described [Microsoft Knowledge Base article #2842230](http://support.microsoft.com/kb/2842230) if you need to run commands that uses more than 150MB of memory.

For more information on WinRM, please refer to <a href="http://msdn.microsoft.com/en-us/library/windows/desktop/aa384426(v=vs.85).aspx">the online documentation at Microsoft's DevCenter</a>.

### Building the winrm go and executable

You can build winrm from source:

```sh
git clone https://github.com/sneal/go-winrm
cd go-winrm
make
```

_Note_: you need go 1.1+. Please check your installation with

```
go version
```

## Library Usage

**Warning the API might be subject to change.**

For the fast version (this doesn't allow to send input to the command):

```go
import (
  "os"
  "fmt"
  "github.com/sneal/go-winrm"
)

shell, err := client.CreateShell()
if err != nil {
  return err
}

var cmd *winrm.Command
cmd, err = shell.Execute(winrm.Powershell("Write-Host 'hello from PS'"), os.Stdout, os.Stderr)
if err != nil {
  return err
}

cmd.Wait()

if cmd.ExitCode() != 0 {
  fmt.Println("Command failed")
}


```

## Developing on WinRM

If you wish to work on `winrm` itself, you'll first need [Go](http://golang.org)
installed (version 1.1+ is _required_). Make sure you have Go properly installed,
including setting up your [GOPATH](http://golang.org/doc/code.html#GOPATH).

For some additional dependencies, Go needs [Mercurial](http://mercurial.selenic.com/)
and [Bazaar](http://bazaar.canonical.com/en/) to be installed.
Winrm itself doesn't require these, but a dependency of a dependency does.

Next, clone this repository into `$GOPATH/src/github.com/sneal/go-winrm` and
then just type `make`.

You can run tests by typing `make test`.

If you make any changes to the code, run `make format` in order to automatically
format the code according to Go standards.

When new dependencies are added to winrm you can use `make updatedeps` to
get the latest and subsequently use `make` to compile and generate the `winrm` binary.

