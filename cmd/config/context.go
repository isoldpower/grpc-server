package config

import (
	"errors"
	"fmt"
	"golang-grpc/internal/util"
	"os"
	"path/filepath"
	"strings"
)

type ProcessContext struct {
	SourceDir string
	TargetDir string
	RootDir   string
}

func resolveExecutionProcessContext() string {
	var executionFile string
	if exe, err := os.Executable(); err == nil && !strings.HasPrefix(exe, os.TempDir()) {
		if sym, err := filepath.EvalSymlinks(exe); err == nil {
			executionFile = sym
		} else {
			executionFile = exe
		}
	}

	return filepath.Dir(executionFile)
}

func handleTargetError(error error) error {
	if errors.Is(error, os.ErrNotExist) {
		return fmt.Errorf("execution target does not exist")
	}

	return error
}

// NewProcessContext is a constructor function to safely initialize new ProcessContext object
func NewProcessContext() *ProcessContext {
	return &ProcessContext{}
}

// ResolveProcessContext The purpose of this function is to handle the context of the run command,
// specifically the sourceDir path and targetDir path, which represent
// the path to source project and the path to execution directory respectively
func (c *ProcessContext) ResolveProcessContext(args []string) error {
	var errorMessage error

	executionProcessContext := resolveExecutionProcessContext()
	sourceDir := util.ResolvePath(filepath.Dir(""), "")
	targetDir := util.ResolvePath(executionProcessContext, "")
	rootDir := util.RootPath()

	if stat, err := os.Stat(targetDir); err != nil {
		errorMessage = handleTargetError(err)
	} else if !stat.IsDir() {
		errorMessage = fmt.Errorf("%s is not a directory", targetDir)
	} else {
		c.TargetDir = targetDir
		c.SourceDir = sourceDir
		c.RootDir = rootDir
	}

	return errorMessage
}
