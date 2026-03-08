package installer

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
)

func checkKubectl() error {
	_, err := exec.LookPath("kubectl")
	if err != nil {
		return fmt.Errorf("kubectl not found in PATH — install it from https://kubernetes.io/docs/tasks/tools/")
	}
	return nil
}

func kubectl(kubeCtx string, args ...string) (string, error) {
	cmdArgs := args
	if kubeCtx != "" {
		cmdArgs = append([]string{"--context", kubeCtx}, args...)
	}

	cmd := exec.CommandContext(context.Background(), "kubectl", cmdArgs...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%w: %s", err, stderr.String())
	}
	return stdout.String(), nil
}
