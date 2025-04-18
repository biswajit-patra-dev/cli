package logout

import (
	"testing"

	"github.com/debricked/cli/internal/auth"
	"github.com/debricked/cli/internal/auth/testdata"
	"github.com/stretchr/testify/assert"
)

func TestNewLogoutCmd(t *testing.T) {
	authenticator := auth.NewDebrickedAuthenticator("")
	cmd := NewLogoutCmd(authenticator)
	commands := cmd.Commands()
	nbrOfCommands := 0
	assert.Len(t, commands, nbrOfCommands)
}

func TestPreRun(t *testing.T) {
	mockAuthenticator := testdata.MockAuthenticator{}
	cmd := NewLogoutCmd(mockAuthenticator)
	cmd.PreRun(cmd, nil)
}

func TestRunE(t *testing.T) {
	a := testdata.MockAuthenticator{}
	runE := RunE(a)

	err := runE(nil, []string{})

	assert.NoError(t, err)
}

func TestRunEError(t *testing.T) {
	a := testdata.ErrorMockAuthenticator{}
	runE := RunE(a)

	err := runE(nil, []string{})

	assert.Error(t, err)
}
