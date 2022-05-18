package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/bazelbuild/rules_go/go/tools/bazel"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func Test_stdliblist(t *testing.T) {
	testDir := t.TempDir()
	outJSON := filepath.Join(testDir, "out.json")

	runFilePath, err := bazel.RunfilesPath()
	if err != nil {
		t.Errorf("cannot file runfile path %v", err)
	}
	test_args := []string{
		fmt.Sprintf("-out=%s", outJSON),
		fmt.Sprintf("-sdk=%s", abs(filepath.Join(filepath.Dir(runFilePath), "go_sdk"))),
	}

	err = stdliblist(test_args)
	if err != nil {
		t.Errorf("calling stdliblist got err: %v", err)
	}
	f, err := os.Open(outJSON)
	if err != nil {
		t.Errorf("cannot open output json: %v", err)
	}
	defer func() { _ = f.Close() }()
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
