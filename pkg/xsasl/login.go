package xsasl

import (
	"bytes"
)

// The LOGIN mechanism name.
const Login = "LOGIN"

var expectedLoginChallenge1 = []byte("Username:")
var expectedLoginChallenge2 = []byte("Password:")

type loginClient struct {
	Username string
	Password string
}

func (a *loginClient) Start() (mech string, ir []byte, err error) {
	mech = "LOGIN"
	ir = []byte(a.Username)
	return
}

func (a *loginClient) Next(challenge []byte) (response []byte, err error) {
	if bytes.Compare(challenge, expectedLoginChallenge2) != 0 {
		return nil, ErrUnexpectedServerChallenge
	} else {
		return []byte(a.Password), nil
	}
}

// A client implementation of the LOGIN authentication mechanism for SMTP,
// as described in http://www.iana.org/go/draft-murchison-sasl-login
//
// It is considered obsolete, and should not be used when other mechanisms are
// available. For plaintext password authentication use PLAIN mechanism.
func NewLoginClient(username, password string) Client {
	return &loginClient{username, password}
}

// Authenticates users with an username and a password.
type LoginAuthenticator func(username, password string) error

type loginServer struct {
	done         bool
	username     string
	authenticate LoginAuthenticator
}

func (a *loginServer) Next(response []byte) (challenge []byte, done bool, err error) {
	if a.done {
		err = ErrUnexpectedClientResponse
		return
	}

	// initial response is username, if not request it.
	if response == nil {
		return expectedLoginChallenge1, false, nil
	}

	if a.username == "" {
		a.username = string(response)
		return expectedLoginChallenge2, false, nil
	}

	password := string(response)

	err = a.authenticate(a.username, password)
	a.done = true
	return nil, a.done, err
}

// A server implementation of the LOGIN authentication mechanism.
func NewLoginServer(authenticator LoginAuthenticator) Server {
	return &loginServer{authenticate: authenticator}
}
