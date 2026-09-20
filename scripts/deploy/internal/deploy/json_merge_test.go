package deploy

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestRunnerMergeJSONPreservesUnmanagedSettings(t *testing.T) {
	root := t.TempDir()
	src := filepath.Join(root, "cursor.json")
	dst := filepath.Join(root, "home", "cli-config.json")
	writeConfig(t, src, `{"version":1,"editor":{"vimMode":false},"permissions":{"allow":["Shell(ls)"],"deny":[]},"approvalMode":"allowlist"}`)
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		t.Fatal(err)
	}
	writeConfig(t, dst, `{"version":1,"editor":{"vimMode":true},"model":{"modelId":"saved"},"authInfo":{"userId":"private"},"permissions":{"allow":["Shell(rm)"],"deny":["Read(secret.key)"]}}`)
	if err := os.Chmod(dst, 0600); err != nil {
		t.Fatal(err)
	}
	config := filepath.Join(root, "deploy.json")
	writeConfig(t, config, `{"items":[{"source":"cursor.json","destination":"home/cli-config.json","mergeJSON":true}]}`)

	var out bytes.Buffer
	if err := runFromDir(t, root, func() error { return NewRunner(&out).Run(config, Options{NoColor: true}) }); err != nil {
		t.Fatal(err)
	}
	contents, err := os.ReadFile(dst)
	if err != nil {
		t.Fatal(err)
	}
	var got struct {
		Model struct {
			ModelID string `json:"modelId"`
		} `json:"model"`
		AuthInfo struct {
			UserID string `json:"userId"`
		} `json:"authInfo"`
		Permissions struct {
			Allow []string `json:"allow"`
			Deny  []string `json:"deny"`
		} `json:"permissions"`
		ApprovalMode string `json:"approvalMode"`
	}
	if err := json.Unmarshal(contents, &got); err != nil {
		t.Fatal(err)
	}
	if got.Model.ModelID != "saved" || got.AuthInfo.UserID != "private" || got.ApprovalMode != "allowlist" {
		t.Fatalf("unmanaged settings changed: %+v", got)
	}
	if len(got.Permissions.Allow) != 1 || got.Permissions.Allow[0] != "Shell(ls)" {
		t.Fatalf("permissions not replaced: %+v", got.Permissions)
	}
	if len(got.Permissions.Deny) != 1 || got.Permissions.Deny[0] != "Read(secret.key)" {
		t.Fatalf("existing deny rules not preserved: %+v", got.Permissions)
	}
	info, err := os.Stat(dst)
	if err != nil || info.Mode().Perm() != 0600 {
		t.Fatalf("destination mode changed: %v %v", info, err)
	}
}

func TestRunnerMergeJSONRejectsInvalidDestinationWithoutWriting(t *testing.T) {
	root := t.TempDir()
	writeConfig(t, filepath.Join(root, "source.json"), `{"permissions":{"allow":[],"deny":[]}}`)
	dst := filepath.Join(root, "destination.json")
	writeConfig(t, dst, `invalid`)
	config := filepath.Join(root, "deploy.json")
	writeConfig(t, config, `{"items":[{"source":"source.json","destination":"destination.json","mergeJSON":true}]}`)
	var out bytes.Buffer
	if err := runFromDir(t, root, func() error { return NewRunner(&out).Run(config, Options{NoColor: true}) }); err == nil {
		t.Fatal("expected invalid destination to fail")
	}
	assertFileContent(t, dst, "invalid")
}

func TestRunnerMergeJSONDryRunDoesNotWrite(t *testing.T) {
	root := t.TempDir()
	writeConfig(t, filepath.Join(root, "source.json"), `{"permissions":{"allow":[],"deny":[]}}`)
	dst := filepath.Join(root, "destination.json")
	writeConfig(t, dst, `{"model":"saved"}`)
	config := filepath.Join(root, "deploy.json")
	writeConfig(t, config, `{"items":[{"source":"source.json","destination":"destination.json","mergeJSON":true}]}`)
	var out bytes.Buffer
	if err := runFromDir(t, root, func() error { return NewRunner(&out).Run(config, Options{DryRun: true, NoColor: true}) }); err != nil {
		t.Fatal(err)
	}
	assertFileContent(t, dst, `{"model":"saved"}`)
}

func TestRunnerMergeJSONCreatesPrivateDestination(t *testing.T) {
	root := t.TempDir()
	writeConfig(t, filepath.Join(root, "source.json"), `{"permissions":{"allow":["Shell(ls)"],"deny":[]}}`)
	config := filepath.Join(root, "deploy.json")
	writeConfig(t, config, `{"items":[{"source":"source.json","destination":"nested/destination.json","mergeJSON":true}]}`)
	var out bytes.Buffer
	if err := runFromDir(t, root, func() error { return NewRunner(&out).Run(config, Options{NoColor: true}) }); err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(filepath.Join(root, "nested", "destination.json"))
	if err != nil || info.Mode().Perm() != 0600 {
		t.Fatalf("new destination mode: %v %v", info, err)
	}
}
