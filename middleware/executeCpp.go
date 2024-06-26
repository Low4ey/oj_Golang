package middleware

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/low4ey/OJ/Golang-backend/models"
)

func RunCpp(codeBody string, testcases []models.TestCase) (int, string, error) {
	err := os.WriteFile("./files/solution.cpp", []byte(codeBody), 0644)
	if err != nil {
		return -1, "", fmt.Errorf("failed to write code to file: %v", err)
	}

	err = executeCppFile("./files/solution.cpp")
	if err != nil {
		return -1, compileError, fmt.Errorf("failed to execute C++ file: %v", err)
	}

	outcome, err := runExecutableWithTimeout("", "./files/a.out", testcases)
	if err != nil {
		if err == context.DeadlineExceeded {
			return outcome, timeExceeded, nil
		} else if strings.Contains(err.Error(), "exited with status") {
			return outcome, runtimeError, nil
		} else if strings.Contains(err.Error(), "exceeded memory limit") {
			return outcome, memoryExceeded, nil
		}
		return outcome, compileError, fmt.Errorf("failed to run executable: %v", err)
	}
	if outcome == len(testcases)-1 {
		lastIndex, isEqual, err := compareFile("./files/output.txt", "./files/expected_output.txt")
		if err != nil {
			return lastIndex, "", fmt.Errorf("failed to compare files: %v", err)
		}

		if !isEqual {
			return lastIndex, wrongAnswer, nil
		}
	}
	return outcome, correctAnswer, nil
}

func executeCppFile(filePath string) error {
	cmd := exec.Command("g++", filePath, "-o", "./files/a.out") // Compile the C++ file using g++
	errOutput := &bytes.Buffer{}                                // Buffer to capture the error output
	cmd.Stderr = errOutput                                      // Attach the buffer to cmd.Stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to compile C++ file: %v\n%s", err, errOutput.String())
	}

	return nil
}
