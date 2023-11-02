// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023-Present The UDS Authors

// Package test provides e2e tests for UDS.
package test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUseCLI(t *testing.T) {
	t.Log("E2E: Use CLI")

	t.Run("run copy", func(t *testing.T) {
		t.Parallel()

		baseFilePath := "base"
		copiedFilePath := "copy"

		e2e.CleanFiles(baseFilePath, copiedFilePath)
		t.Cleanup(func() {
			e2e.CleanFiles(baseFilePath, copiedFilePath)
		})

		err := os.WriteFile(baseFilePath, []byte{}, 0600)
		require.NoError(t, err)

		stdOut, stdErr, err := e2e.RunTasksWithFile("run", "copy")
		require.NoError(t, err, stdOut, stdErr)

		require.FileExists(t, copiedFilePath)
	})

	t.Run("run copy-exec", func(t *testing.T) {
		t.Parallel()

		baseFilePath := "exectest"
		copiedFilePath := "exec"

		e2e.CleanFiles(baseFilePath, copiedFilePath)
		t.Cleanup(func() {
			e2e.CleanFiles(baseFilePath, copiedFilePath)
		})

		err := os.WriteFile(baseFilePath, []byte{}, 0600)
		require.NoError(t, err)

		stdOut, stdErr, err := e2e.RunTasksWithFile("run", "copy-exec")
		require.NoError(t, err, stdOut, stdErr)

		require.FileExists(t, copiedFilePath)
		execFileInfo, err := os.Stat(copiedFilePath)
		require.NoError(t, err)
		require.True(t, execFileInfo.Mode()&0111 != 0)
	})

	t.Run("run copy-verify", func(t *testing.T) {
		t.Parallel()

		baseFilePath := "data"
		copiedFilePath := "verify"

		e2e.CleanFiles(baseFilePath, copiedFilePath)
		t.Cleanup(func() {
			e2e.CleanFiles(baseFilePath, copiedFilePath)
		})

		err := os.WriteFile(baseFilePath, []byte("test"), 0600)
		require.NoError(t, err)

		stdOut, stdErr, err := e2e.RunTasksWithFile("run", "copy-verify")
		require.NoError(t, err, stdOut, stdErr)

		require.FileExists(t, copiedFilePath)
	})

	t.Run("run copy-symlink", func(t *testing.T) {
		t.Parallel()

		baseFilePath := "symtest"
		copiedFilePath := "symcopy"
		symlinkName := "testlink"

		e2e.CleanFiles(baseFilePath, copiedFilePath, symlinkName)
		t.Cleanup(func() {
			e2e.CleanFiles(baseFilePath, copiedFilePath, symlinkName)
		})

		err := os.WriteFile(baseFilePath, []byte{}, 0600)
		require.NoError(t, err)

		stdOut, stdErr, err := e2e.RunTasksWithFile("run", "copy-symlink")
		require.NoError(t, err, stdOut, stdErr)

		require.FileExists(t, symlinkName)
	})

	t.Run("run remote", func(t *testing.T) {
		t.Parallel()

		downloadedFile := "checksums.txt"

		e2e.CleanFiles(downloadedFile)
		t.Cleanup(func() {
			e2e.CleanFiles(downloadedFile)
		})

		stdOut, stdErr, err := e2e.RunTasksWithFile("run", "remote")
		require.NoError(t, err, stdOut, stdErr)

		require.FileExists(t, downloadedFile)
	})

	t.Run("run template-file", func(t *testing.T) {
		t.Parallel()

		baseFilePath := "raw"
		copiedFilePath := "templated"

		e2e.CleanFiles(baseFilePath, copiedFilePath)
		t.Cleanup(func() {
			e2e.CleanFiles(baseFilePath, copiedFilePath)
		})

		err := os.WriteFile(baseFilePath, []byte("###ZARF_VAR_REPLACE_ME###"), 0600)
		require.NoError(t, err)

		stdOut, stdErr, err := e2e.RunTasksWithFile("run", "template-file")
		require.NoError(t, err, stdOut, stdErr)

		require.FileExists(t, copiedFilePath)

		templatedContentsBytes, err := os.ReadFile(copiedFilePath)
		require.NoError(t, err)
		require.Equal(t, "replaced\n", string(templatedContentsBytes))
	})

	t.Run("run action", func(t *testing.T) {
		t.Parallel()

		stdOut, stdErr, err := e2e.RunTasksWithFile("run", "action")
		require.NoError(t, err, stdOut, stdErr)
		require.Contains(t, stdErr, "specific test string")
	})

	t.Run("run cmd-set-variable", func(t *testing.T) {
		t.Parallel()

		stdOut, stdErr, err := e2e.RunTasksWithFile("run", "cmd-set-variable")
		require.NoError(t, err, stdOut, stdErr)
		require.Contains(t, stdErr, "unique-value")
	})

	t.Run("run reference", func(t *testing.T) {
		t.Parallel()

		stdOut, stdErr, err := e2e.RunTasksWithFile("run", "reference")
		require.NoError(t, err, stdOut, stdErr)
		require.Contains(t, stdErr, "other-task")
	})

	t.Run("run recursive", func(t *testing.T) {
		t.Parallel()

		stdOut, stdErr, err := e2e.RunTasksWithFile("run", "recursive")
		require.Error(t, err, stdOut, stdErr)
		require.Contains(t, stdErr, "task loop detected")
	})
}