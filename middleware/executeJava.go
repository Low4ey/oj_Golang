package middleware

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/low4ey/OJ/Golang-backend/models"
)

func RunJava(codeBody string, testcases []models.TestCase) (int, string, error) {
	err := os.WriteFile("./files/Solution.java", []byte(codeBody), 0644)
	if err != nil {
		return -1, "", fmt.Errorf("failed to write code to file: %v", err)
	}

	compErr := compileJava("./files/Solution.java")
	if compErr != nil {
		return -1, "", fmt.Errorf("failed to compile Java code: %v", compErr)
	}
	outcome, err := runExecutableWithTimeout("java", "./files/Solution.java", testcases)
	if err != nil {
		if err == context.DeadlineExceeded {
			return outcome, timeExceeded, nil
		} else if strings.Contains(err.Error(), "exited with status") {
			return outcome, runtimeError, nil
		} else if strings.Contains(err.Error(), "exceeded memory limit") {
			return outcome, memoryExceeded, nil
		}
		return outcome, runtimeError, fmt.Errorf("failed to run Java code: %v", err)
	}

	if outcome == -1 {
		return -1, "", fmt.Errorf("no testcases executed during runtime")
	}

	lastTestCaseExecuted, isEqual, err := compareFile("./files/output.txt", "./files/expected_output.txt")
	if err != nil {
		return lastTestCaseExecuted, "", fmt.Errorf("failed to compare files: %v", err)
	}

	if !isEqual {
		return lastTestCaseExecuted, wrongAnswer, nil
	}

	return outcome, correctAnswer, nil
}

func compileJava(fileAddress string) error {
	cmd := exec.Command("javac", fileAddress)

	cmd.Stdout = &bytes.Buffer{}
	cmd.Stderr = &bytes.Buffer{}

	err := cmd.Run()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
				return fmt.Errorf("javac exited with status: %d", status.ExitStatus())
			}
		}
		return fmt.Errorf("failed to compile Java code: %v", err)
	}

	return nil
}
