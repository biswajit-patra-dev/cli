package pip

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type ExecPathMock struct {
}

func (_ ExecPathMock) LookPath(file string) (string, error) {

	if file == "python3" {
		return "", errors.New("executable file not found in $PATH")
	}

	return "python", nil
}

func TestCreateVenvCmd(t *testing.T) {
	venvName := "test-file.venv"
	cmd, err := CmdFactory{
		execPath: ExecPath{},
	}.MakeCreateVenvCmd(venvName)
	assert.NoError(t, err)
	assert.NotNil(t, cmd)
	args := cmd.Args
	assert.Contains(t, args, "python3")
	assert.Contains(t, args, "-m")
	assert.Contains(t, args, "venv")
	assert.Contains(t, args, venvName)
	assert.Contains(t, args, "--clear")

	execPathMock := ExecPathMock{}
	venvName = "test-file.venv"
	cmd2, err := CmdFactory{
		execPath: execPathMock,
	}.MakeCreateVenvCmd(venvName)

	assert.NoError(t, err)
	assert.NotNil(t, cmd2)
	args = cmd2.Args
	assert.Contains(t, args, "python")
	assert.Contains(t, args, "-m")
	assert.Contains(t, args, "venv")
	assert.Contains(t, args, venvName)
	assert.Contains(t, args, "--clear")
}

func TestMakeInstallCmd(t *testing.T) {
	fileName := "test-file"
	pipCommand := "pip"
	cmd, err := CmdFactory{
		execPath: ExecPath{},
	}.MakeInstallCmd(pipCommand, fileName)
	assert.NoError(t, err)
	assert.NotNil(t, cmd)
	args := cmd.Args
	assert.Contains(t, args, "pip")
	assert.Contains(t, args, "install")
	assert.Contains(t, args, "-r")
	assert.Contains(t, args, fileName)
}

func TestMakeCatCmd(t *testing.T) {
	fileName := "test-file"
	cmd, _ := CmdFactory{}.MakeCatCmd(fileName)
	assert.NotNil(t, cmd)
	args := cmd.Args
	assert.Contains(t, args, "cat")
	assert.Contains(t, args, fileName)
}
func TestMakeListCmd(t *testing.T) {
	mockCommand := "mock-cmd"
	cmd, _ := CmdFactory{
		execPath: ExecPath{},
	}.MakeListCmd(mockCommand)
	assert.NotNil(t, cmd)
	args := cmd.Args
	assert.Contains(t, args, "pip")
	assert.Contains(t, args, "list")
}

func TestMakeShowCmd(t *testing.T) {
	input := []string{"package1", "package2"}
	mockCommand := "pip"
	cmd, _ := CmdFactory{
		execPath: ExecPath{},
	}.MakeShowCmd(mockCommand, input)
	assert.NotNil(t, cmd)
	args := cmd.Args
	assert.Contains(t, args, "pip")
	assert.Contains(t, args, "show")
	assert.Contains(t, args, "package1")
	assert.Contains(t, args, "package2")
}
