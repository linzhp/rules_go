package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bazelbuild/rules_go/go/tools/bazel"
)

func Test_stdliblist(t *testing.T) {
	testDir := t.TempDir()
	f, _ := ioutil.TempFile(testDir, "*")

	// test files are at run file directory, but this test is run at
	// {runfile directory}/bazel.TestWorkspace()
	// since -sdk is assumed to be a relative path to execRoot
	// (go.sdk.root_file.dirname), thus setting wd to
	// {runfile directory} so that go_sdk is discoverable
	// {runfile directory} is the parent directory of bazel.RunfilesPath()
	runFilesPath, err := bazel.RunfilesPath()
	if err != nil {
		t.Error("failed to find runfiles path")
	}
	test_args := []string{
		fmt.Sprintf("-out=%s", f.Name()),
		fmt.Sprintf("-sdk=%s", "go_sdk"),
		fmt.Sprintf("-wd=%s", filepath.Dir(filepath.Clean(runFilesPath))),
	}

	err = stdliblist(test_args)
	if err != nil {
		t.Errorf("calling stdliblist got err: %v", err)
	}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		var result flatPackage
		jsonLineStr := scanner.Text()
		if err := json.Unmarshal([]byte(jsonLineStr), &result); err != nil {
			t.Errorf("cannot parse result line %s \n to goListPackage{}: %v\n", err)
		}
		if !strings.HasPrefix(result.ID, "@io_bazel_rules_go//stdlib") {
			t.Errorf("ID should be prefixed with @io_bazel_rules_go//stdlib :%s", jsonLineStr)
		}
		if !strings.HasPrefix(result.ExportFile, "__BAZEL_OUTPUT_BASE__") {
			t.Errorf("export file should be prefixed with __BAZEL_OUTPUT_BASE__ :%s", jsonLineStr)
		}
		for _, gofile := range result.GoFiles {
			if !strings.HasPrefix(gofile, "__BAZEL_OUTPUT_BASE__/go_sdk") {
				t.Errorf("All go files should be prefixed with __BAZEL_OUTPUT_BASE__/go_sdk :%s", jsonLineStr)
			}
		}
	}
}
