package scan

import (
	"testing"

	"github.com/debricked/cli/pkg/client"
	"github.com/debricked/cli/pkg/client/testdata"
	"github.com/debricked/cli/pkg/scan"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestNewScanCmd(t *testing.T) {
	var c client.IDebClient = testdata.NewDebClientMock()
	cmd := NewScanCmd(&c)

	viperKeys := viper.AllKeys()
	flags := cmd.Flags()
	flagAssertions := map[string]string{
		RepositoryFlag:    "r",
		CommitFlag:        "c",
		BranchFlag:        "b",
		CommitAuthorFlag:  "a",
		RepositoryUrlFlag: "u",
		IntegrationFlag:   "i",
	}
	for name, shorthand := range flagAssertions {
		flag := flags.Lookup(name)
		assert.NotNil(t, flag)
		assert.Equalf(t, shorthand, flag.Shorthand, "failed to assert that %s flag shorthand %s was set correctly", name, shorthand)

		match := false
		for _, key := range viperKeys {
			if key == name {
				match = true
			}
		}
		assert.Truef(t, match, "failed to assert that %s was present", name)
	}
}

func TestRunE(t *testing.T) {
	var s scan.IScanner = &scannerMock{}
	runE := RunE(&s)

	err := runE(nil, []string{"."})

	assert.NoError(t, err)
}

func TestRunENoPath(t *testing.T) {
	var s scan.IScanner = &scannerMock{}
	runE := RunE(&s)

	err := runE(nil, []string{})

	assert.NoError(t, err)
}

func TestRunEFailPipelineErr(t *testing.T) {
	var s scan.IScanner
	mock := &scannerMock{}
	mock.setErr(scan.FailPipelineErr)
	s = mock
	runE := RunE(&s)
	cmd := &cobra.Command{}

	err := runE(cmd, nil)

	assert.Error(t, err, scan.FailPipelineErr)
	assert.True(t, cmd.SilenceUsage, "failed to assert that usage was silenced")
	assert.True(t, cmd.SilenceErrors, "failed to assert that errors were silenced")
}

func TestRunEError(t *testing.T) {
	runE := RunE(nil)
	err := runE(nil, []string{"."})

	assert.ErrorContains(t, err, "⨯ scanner was nil")
}

type scannerMock struct {
	err error
}

func (s *scannerMock) Scan(_ scan.IOptions) error {
	return s.err
}

func (s *scannerMock) setErr(err error) {
	s.err = err
}
