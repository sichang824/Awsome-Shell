package exec

import (
	"bytes"
	"fmt"
	"os"
	osexec "os/exec"
)

// Run runs command and returns combined output and error.
func Run(name string, args ...string) (stdout, stderr string, err error) {
	cmd := osexec.Command(name, args...)
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	err = cmd.Run()
	return outBuf.String(), errBuf.String(), err
}

// RunInherit runs command with stdin/stdout/stderr connected to current process.
func RunInherit(name string, args ...string) error {
	cmd := osexec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// RunDir runs command in the given directory.
func RunDir(dir, name string, args ...string) (stdout, stderr string, err error) {
	cmd := osexec.Command(name, args...)
	cmd.Dir = dir
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	err = cmd.Run()
	return outBuf.String(), errBuf.String(), err
}

// DockerComposeExec runs: docker compose exec -T <service> <args>
func DockerComposeExec(service string, args ...string) (stdout, stderr string, err error) {
	a := []string{"compose", "exec", "-T", service}
	a = append(a, args...)
	return Run("docker", a...)
}

// DockerComposeExecTTY runs: docker compose exec <service> <args> (with TTY for interactive)
func DockerComposeExecTTY(service string, args ...string) error {
	a := []string{"compose", "exec", service}
	a = append(a, args...)
	return RunInherit("docker", a...)
}

// DockerComposeExecTTYWithEnv runs: docker compose exec -e KEY=val <service> <args>
func DockerComposeExecTTYWithEnv(service string, env map[string]string, args ...string) error {
	a := []string{"compose", "exec"}
	for k, v := range env {
		a = append(a, "-e", k+"="+v)
	}
	a = append(a, service)
	a = append(a, args...)
	cmd := osexec.Command("docker", a...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// MustDockerOut runs Docker compose exec and returns stdout; exits on error.
func MustDockerOut(service string, args ...string) string {
	out, stderr, err := DockerComposeExec(service, args...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", stderr)
		os.Exit(1)
	}
	return out
}
