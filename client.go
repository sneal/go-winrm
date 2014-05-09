package winrm

import (
	"github.com/sneal/go-winrm/soap"
	"strings"
)

type Client struct {
	Parameters
	username string
	password string
	useHTTPS bool
	http     HttpPost
}

// NewClient will create a new remote client on url, connecting with user and password
// This function doesn't connect (connection happens only when CreateShell is called)
func NewClient(host string, user string, password string) (client *Client) {
	params := DefaultParameters()
	params.url = winRMUrl(host)
	client = &Client{Parameters: *params, username: user, password: password, http: Http_post}
	return
}

// NewClient will create a new remote client on url, connecting with user and password
// This function doesn't connect (connection happens only when CreateShell is called)
func NewClientWithParameters(host string, user string, password string, params *Parameters) (client *Client) {
	params.url = winRMUrl(host)
	client = &Client{Parameters: *params, username: user, password: password, http: Http_post}
	return
}

func winRMUrl(host string) string {
	if !strings.Contains(host, ":") {
		host += ":5985"
	}
	return "http://" + host + "/wsman"
}

// CreateShell will create a WinRM Shell, which is required before running
// commands.
func (client *Client) CreateShell() (shell *Shell, err error) {
	request := NewOpenShellRequest(client.Parameters.url, &client.Parameters)
	defer request.Free()

	response, err := client.sendRequest(request)
	if err == nil {
		var shellId string
		if shellId, err = ParseOpenShellResponse(response); err == nil {
			shell = &Shell{client: client, ShellId: shellId}
		}
	}
	return
}

func (client *Client) sendRequest(request *soap.SoapMessage) (response string, err error) {
	return client.http(client, request)
}

