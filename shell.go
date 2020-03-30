package shell

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"

	"github.com/mattn/go-shellwords"
)

var (
	powershellPattern = regexp.MustCompile(`\.ps1($| )`)
)

// GetTokens retrieves the individual components of the given command
func GetTokens(command string) ([]string, error) {
	return shellwords.Parse(command)
}

// ExecuteShellCommand runs a command via the shell
func ExecuteShellCommand(command string) error {
	tokens, err := GetTokens(command)
	if err != nil {
		return err
	}

	cmd := exec.Command(tokens[0], tokens[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// ExecuteShellCommandAsync runs a command via the shell asynchronously
func ExecuteShellCommandAsync(command string) error {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		if powershellPattern.FindString(command) != "" {
			cmd = exec.Command("powershell", "&", fmt.Sprintf("'%s'", command))
		} else {
			cmd = exec.Command("cmd.exe", "/c", command)
		}
	} else if runtime.GOOS == "darwin" || runtime.GOOS == "linux" {
		cmd = exec.Command("sh", "-c", command)
	} else {
		return errors.New("platform not supported at this current time")
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// ExecuteCommand runs a command with given args
func ExecuteCommand(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// IsFlag returns whether the given token matches either flag
// Specify flag arg as empty if that type of flag shouldn't be considered be considered
func IsFlag(token string, shortFlag string, longFlag string) bool {
	return (shortFlag != "" && token == shortFlag) || (longFlag != "" && token == longFlag)
}

// GetFlagValue returns the associated value for the specified flags given a token
func GetFlagValue(token string, shortFlag string, longFlag string) string {
	if strings.HasPrefix(token, shortFlag) || strings.HasPrefix(token, longFlag) {
		return token[strings.Index(token, "=")+1:]
	}
	return ""
}
