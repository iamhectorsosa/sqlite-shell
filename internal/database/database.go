package database

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func ExecCmd(path, query string) (headers []string, rows [][]string, err error) {
	resolvedPath, err := resolvePath(path)
	if err != nil {
		return nil, nil, formatErrors("resolving path", err)
	}

	cmd := exec.Command("sqlite3", "-csv", "-header", resolvedPath, query)
	var out, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &stderr
	if err := cmd.Run(); err != nil {
		return nil, nil, formatErrors(
			"executing command",
			fmt.Errorf("exec: %v", err),
			fmt.Errorf("sqlite3: %s", stderr.String()),
		)
	}

	headers, rows, err = parseCSV(out.String())
	if err != nil {
		return nil, nil, formatErrors("parsing data", err)
	}

	return headers, rows, nil
}

func resolvePath(path string) (string, error) {
	path = os.ExpandEnv(path)

	if len(path) > 2 && path[:2] == "~/" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("reading home dir: %v", err)
		}
		return filepath.Join(home, path[2:]), nil
	}

	if filepath.IsAbs(path) {
		return path, nil
	}

	return path, nil
}

func parseCSV(input string) ([]string, [][]string, error) {
	reader := csv.NewReader(strings.NewReader(input))
	records, err := reader.ReadAll()
	if err != nil {
		return nil, nil, fmt.Errorf("reading CSV: %v", err)
	}
	if len(records) == 0 {
		return nil, nil, nil
	}

	headers := records[0]
	rows := records[1:]
	return headers, rows, nil
}

func formatErrors(context string, errs ...error) error {
	errMsgs := make([]string, 0, len(errs))
	for _, err := range errs {
		if err != nil {
			errMsgs = append(errMsgs, err.Error())
		}
	}
	if len(errMsgs) > 0 {
		return fmt.Errorf("%s: %s", context, strings.Join(errMsgs, "; "))
	}
	return nil
}
