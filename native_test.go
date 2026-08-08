package main

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestNativeCatalogAndPackCoverage(t *testing.T) {
	catalog, err := loadCatalog()
	if err != nil {
		t.Fatal(err)
	}
	if problems := validateCatalog(catalog); len(problems) != 0 {
		t.Fatalf("catalog problems: %v", problems)
	}
	if len(catalog.Skills) != 24 || len(catalog.Packs) != 7 {
		t.Fatalf("unexpected catalog size: %d skills, %d packs", len(catalog.Skills), len(catalog.Packs))
	}
	complete, err := catalog.resolvePack("complete")
	if err != nil {
		t.Fatal(err)
	}
	if len(complete) != len(catalog.Skills) {
		t.Fatalf("complete pack has %d skills, want %d", len(complete), len(catalog.Skills))
	}
}

func TestEmbeddedFrontmatterValidationNormalizesWindowsNewlines(t *testing.T) {
	normalized := normalizeEmbeddedText([]byte("---\r\nname: clarify-outcome\r\ndescription: example\r\n---\r\n"))
	if !strings.HasPrefix(normalized, "---\nname: clarify-outcome\ndescription:") {
		t.Fatalf("CRLF frontmatter was not normalized: %q", normalized)
	}
}

func TestNativeInstallIsIdempotentAndAgentAware(t *testing.T) {
	catalog, err := loadCatalog()
	if err != nil {
		t.Fatal(err)
	}
	project := t.TempDir()
	events := []string{}
	emit := func(key string, values ...any) { events = append(events, key) }
	if err := installNative(catalog, []string{"clarify-outcome"}, []string{"codex", "claude-code"}, "project", project, false, false, emit); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(project, ".agents", "skills", "clarify-outcome", "agents", "openai.yaml")); err != nil {
		t.Fatalf("Codex metadata missing: %v", err)
	}
	if _, err := os.Stat(filepath.Join(project, ".claude", "skills", "clarify-outcome", "agents")); !os.IsNotExist(err) {
		t.Fatalf("Claude materialization contains Codex metadata: %v", err)
	}
	events = nil
	if err := installNative(catalog, []string{"clarify-outcome"}, []string{"codex", "claude-code"}, "project", project, false, false, emit); err != nil {
		t.Fatal(err)
	}
	if strings.Join(events, ",") != "unchanged,unchanged" {
		t.Fatalf("unexpected idempotent events: %v", events)
	}
}

func TestNativeDryRunAndConflictProtection(t *testing.T) {
	catalog, err := loadCatalog()
	if err != nil {
		t.Fatal(err)
	}
	project := t.TempDir()
	emit := func(string, ...any) {}
	if err := installNative(catalog, []string{"clarify-outcome"}, []string{"codex"}, "project", project, true, false, emit); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(project, ".agents")); !os.IsNotExist(err) {
		t.Fatalf("dry-run wrote target: %v", err)
	}
	if _, err := os.Stat(filepath.Join(project, ".anywork")); !os.IsNotExist(err) {
		t.Fatalf("dry-run wrote state or lock files: %v", err)
	}
	if err := installNative(catalog, []string{"clarify-outcome"}, []string{"codex"}, "project", project, false, false, emit); err != nil {
		t.Fatal(err)
	}
	skillFile := filepath.Join(project, ".agents", "skills", "clarify-outcome", "SKILL.md")
	file, err := os.OpenFile(skillFile, os.O_APPEND|os.O_WRONLY, 0)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := file.WriteString("\nlocal edit\n"); err != nil {
		t.Fatal(err)
	}
	_ = file.Close()
	if err := installNative(catalog, []string{"clarify-outcome"}, []string{"codex"}, "project", project, false, false, emit); err == nil || !strings.Contains(err.Error(), "Refusing to overwrite") {
		t.Fatalf("expected conflict error, got %v", err)
	}
}

func TestNativeUninstallIsRecoverable(t *testing.T) {
	catalog, err := loadCatalog()
	if err != nil {
		t.Fatal(err)
	}
	project := t.TempDir()
	emit := func(string, ...any) {}
	if err := installNative(catalog, []string{"clarify-outcome"}, []string{"codex"}, "project", project, false, false, emit); err != nil {
		t.Fatal(err)
	}
	if err := uninstallNative([]string{"clarify-outcome"}, []string{"codex"}, "project", project, false, false, emit); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(project, ".agents", "skills", "clarify-outcome")); !os.IsNotExist(err) {
		t.Fatalf("managed target remains after uninstall: %v", err)
	}
	matches, err := filepath.Glob(filepath.Join(project, ".agents", ".anywork", "backups", "*", "codex", "clarify-outcome", "SKILL.md"))
	if err != nil || len(matches) != 1 {
		t.Fatalf("expected one recoverable backup, got %v (%v)", matches, err)
	}
}

func setCrashPoint(t *testing.T, point string) {
	t.Helper()
	transactionFaultHook = func(current string) error {
		if current == point {
			return errSimulatedCrash
		}
		return nil
	}
	t.Cleanup(func() { transactionFaultHook = nil })
}

func TestRecoverRollsBackNewInstallAfterTargetRenameCrash(t *testing.T) {
	catalog, err := loadCatalog()
	if err != nil {
		t.Fatal(err)
	}
	project := t.TempDir()
	setCrashPoint(t, "after-target-rename")
	emit := func(string, ...any) {}
	err = installNative(catalog, []string{"clarify-outcome"}, []string{"codex"}, "project", project, false, false, emit)
	if !errors.Is(err, errSimulatedCrash) {
		t.Fatalf("expected simulated crash, got %v", err)
	}
	transactionFaultHook = nil
	target := filepath.Join(project, ".agents", "skills", "clarify-outcome")
	if _, err := os.Stat(target); err != nil {
		t.Fatalf("target rename did not happen before crash: %v", err)
	}
	if err := recoverNative("project", project, emit); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(target); !os.IsNotExist(err) {
		t.Fatalf("rollback left newly installed target: %v", err)
	}
	recovered, err := filepath.Glob(filepath.Join(project, ".agents", ".anywork", "recovery", "*", "codex", "clarify-outcome", "SKILL.md"))
	if err != nil || len(recovered) != 1 {
		t.Fatalf("rollback did not preserve displaced content: %v (%v)", recovered, err)
	}
	statePath, _ := anyworkStatePath("project", project)
	state, err := loadState(statePath)
	if err != nil || len(state.Installations) != 0 {
		t.Fatalf("state was not rolled back: %+v (%v)", state, err)
	}
}

func TestRecoverRestoresOriginalAfterBackupCrash(t *testing.T) {
	catalog, err := loadCatalog()
	if err != nil {
		t.Fatal(err)
	}
	project := t.TempDir()
	emit := func(string, ...any) {}
	if err := installNative(catalog, []string{"clarify-outcome"}, []string{"codex"}, "project", project, false, false, emit); err != nil {
		t.Fatal(err)
	}
	skillFile := filepath.Join(project, ".agents", "skills", "clarify-outcome", "SKILL.md")
	file, err := os.OpenFile(skillFile, os.O_APPEND|os.O_WRONLY, 0)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := file.WriteString("\nimportant local work\n"); err != nil {
		t.Fatal(err)
	}
	_ = file.Close()
	setCrashPoint(t, "after-backup")
	err = installNative(catalog, []string{"clarify-outcome"}, []string{"codex"}, "project", project, false, true, emit)
	if !errors.Is(err, errSimulatedCrash) {
		t.Fatalf("expected simulated crash, got %v", err)
	}
	transactionFaultHook = nil
	if err := recoverNative("project", project, emit); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(skillFile)
	if err != nil || !strings.Contains(string(data), "important local work") {
		t.Fatalf("original content was not restored: %v", err)
	}
}

func TestRecoverFinishesTransactionWhoseStateWasCommitted(t *testing.T) {
	catalog, err := loadCatalog()
	if err != nil {
		t.Fatal(err)
	}
	project := t.TempDir()
	emit := func(string, ...any) {}
	setCrashPoint(t, "after-state-write")
	err = installNative(catalog, []string{"clarify-outcome"}, []string{"codex"}, "project", project, false, false, emit)
	if !errors.Is(err, errSimulatedCrash) {
		t.Fatalf("expected simulated crash, got %v", err)
	}
	transactionFaultHook = nil
	if err := recoverNative("project", project, emit); err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(project, ".agents", "skills", "clarify-outcome")
	if _, err := os.Stat(target); err != nil {
		t.Fatalf("committed target was rolled back: %v", err)
	}
	statePath, _ := anyworkStatePath("project", project)
	if _, err := os.Stat(journalPathForState(statePath)); !os.IsNotExist(err) {
		t.Fatalf("recovery left transaction journal: %v", err)
	}
}

func TestRecoverRestoresInterruptedUninstall(t *testing.T) {
	catalog, err := loadCatalog()
	if err != nil {
		t.Fatal(err)
	}
	project := t.TempDir()
	emit := func(string, ...any) {}
	if err := installNative(catalog, []string{"clarify-outcome"}, []string{"codex"}, "project", project, false, false, emit); err != nil {
		t.Fatal(err)
	}
	setCrashPoint(t, "after-backup")
	err = uninstallNative([]string{"clarify-outcome"}, []string{"codex"}, "project", project, false, false, emit)
	if !errors.Is(err, errSimulatedCrash) {
		t.Fatalf("expected simulated crash, got %v", err)
	}
	transactionFaultHook = nil
	if err := recoverNative("project", project, emit); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(project, ".agents", "skills", "clarify-outcome", "SKILL.md")); err != nil {
		t.Fatalf("interrupted uninstall was not restored: %v", err)
	}
}

func TestSagaRollsBackEveryAppliedOperationInReverse(t *testing.T) {
	catalog, err := loadCatalog()
	if err != nil {
		t.Fatal(err)
	}
	project := t.TempDir()
	emit := func(string, ...any) {}
	applied := 0
	transactionFaultHook = func(point string) error {
		if point == "after-target-rename" {
			applied++
			if applied == 2 {
				return errSimulatedCrash
			}
		}
		return nil
	}
	t.Cleanup(func() { transactionFaultHook = nil })
	err = installNative(catalog, []string{"clarify-outcome", "plan-work"}, []string{"codex"}, "project", project, false, false, emit)
	if !errors.Is(err, errSimulatedCrash) {
		t.Fatalf("expected crash after second operation, got %v", err)
	}
	transactionFaultHook = nil
	if err := recoverNative("project", project, emit); err != nil {
		t.Fatal(err)
	}
	for _, skill := range []string{"clarify-outcome", "plan-work"} {
		if _, err := os.Stat(filepath.Join(project, ".agents", "skills", skill)); !os.IsNotExist(err) {
			t.Fatalf("saga left %s applied: %v", skill, err)
		}
	}
	statePath, _ := anyworkStatePath("project", project)
	state, err := loadState(statePath)
	if err != nil || len(state.Installations) != 0 {
		t.Fatalf("saga state was not restored: %+v (%v)", state, err)
	}
}

func TestOrdinaryApplyFailureRollsBackAutomatically(t *testing.T) {
	catalog, err := loadCatalog()
	if err != nil {
		t.Fatal(err)
	}
	project := t.TempDir()
	emit := func(string, ...any) {}
	transactionFaultHook = func(point string) error {
		if point == "after-target-rename" {
			return errors.New("injected write failure")
		}
		return nil
	}
	t.Cleanup(func() { transactionFaultHook = nil })
	err = installNative(catalog, []string{"clarify-outcome"}, []string{"codex"}, "project", project, false, false, emit)
	if err == nil || !strings.Contains(err.Error(), "injected write failure") {
		t.Fatalf("expected injected failure, got %v", err)
	}
	transactionFaultHook = nil
	target := filepath.Join(project, ".agents", "skills", "clarify-outcome")
	if _, err := os.Stat(target); !os.IsNotExist(err) {
		t.Fatalf("automatic rollback left target: %v", err)
	}
	statePath, _ := anyworkStatePath("project", project)
	if _, err := os.Stat(journalPathForState(statePath)); !os.IsNotExist(err) {
		t.Fatalf("automatic rollback left journal: %v", err)
	}
}

func TestSharedTargetTracksConsumersAndPreventsPrematureRemoval(t *testing.T) {
	catalog, err := loadCatalog()
	if err != nil {
		t.Fatal(err)
	}
	project := t.TempDir()
	emit := func(string, ...any) {}
	if err := installNative(catalog, []string{"clarify-outcome"}, []string{"codex", "gemini-cli"}, "project", project, false, false, emit); err != nil {
		t.Fatal(err)
	}
	statePath, _ := anyworkStatePath("project", project)
	state, err := loadState(statePath)
	if err != nil || len(state.Installations) != 1 {
		t.Fatalf("shared target should have one ownership record: %+v (%v)", state, err)
	}
	for _, record := range state.Installations {
		if strings.Join(record.Consumers, ",") != "codex,gemini-cli" {
			t.Fatalf("unexpected consumers: %v", record.Consumers)
		}
	}
	if err := uninstallNative([]string{"clarify-outcome"}, []string{"codex"}, "project", project, false, false, emit); err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(project, ".agents", "skills", "clarify-outcome")
	if _, err := os.Stat(target); err != nil {
		t.Fatalf("first consumer removed shared target: %v", err)
	}
	state, err = loadState(statePath)
	if err != nil {
		t.Fatal(err)
	}
	for _, record := range state.Installations {
		if strings.Join(record.Consumers, ",") != "gemini-cli" {
			t.Fatalf("consumer removal was not recorded: %v", record.Consumers)
		}
	}
	if err := uninstallNative([]string{"clarify-outcome"}, []string{"gemini-cli"}, "project", project, false, false, emit); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(target); !os.IsNotExist(err) {
		t.Fatalf("last consumer did not remove target: %v", err)
	}
}

func TestAgentSkillsFamilyAlwaysKeepsCodexMetadata(t *testing.T) {
	catalog, err := loadCatalog()
	if err != nil {
		t.Fatal(err)
	}
	project := t.TempDir()
	emit := func(string, ...any) {}
	if err := installNative(catalog, []string{"clarify-outcome"}, []string{"gemini-cli"}, "project", project, false, false, emit); err != nil {
		t.Fatal(err)
	}
	metadata := filepath.Join(project, ".agents", "skills", "clarify-outcome", "agents", "openai.yaml")
	if _, err := os.Stat(metadata); err != nil {
		t.Fatalf("agent-skills family lost Codex UI metadata: %v", err)
	}
	if err := installNative(catalog, []string{"clarify-outcome"}, []string{"codex"}, "project", project, false, false, emit); err != nil {
		t.Fatal(err)
	}
	statePath, _ := anyworkStatePath("project", project)
	state, err := loadState(statePath)
	if err != nil {
		t.Fatal(err)
	}
	for _, record := range state.Installations {
		if strings.Join(record.Consumers, ",") != "codex,gemini-cli" {
			t.Fatalf("late Codex consumer was not added: %v", record.Consumers)
		}
	}
	codexHash, err := digestEmbeddedSkill(catalog.Skills["clarify-outcome"], "codex")
	if err != nil {
		t.Fatal(err)
	}
	geminiHash, err := digestEmbeddedSkill(catalog.Skills["clarify-outcome"], "gemini-cli")
	if err != nil {
		t.Fatal(err)
	}
	if codexHash != geminiHash {
		t.Fatal("shared target family produced consumer-specific hashes")
	}
}

func TestLegacyAgentKeyedStateMigratesToConsumers(t *testing.T) {
	directory := t.TempDir()
	statePath := filepath.Join(directory, "state.json")
	target := filepath.Join(directory, "skills", "clarify-outcome")
	legacy := stateFile{SchemaVersion: 1, Installations: map[string]installationRecord{
		"codex:" + target:  {Agent: "codex", Skill: "clarify-outcome", Target: target, Hash: "same", Version: "1"},
		"gemini:" + target: {Agent: "gemini-cli", Skill: "clarify-outcome", Target: target, Hash: "same", Version: "1"},
	}}
	data, err := json.Marshal(legacy)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(statePath, data, 0o600); err != nil {
		t.Fatal(err)
	}
	state, err := loadState(statePath)
	if err != nil || len(state.Installations) != 1 {
		t.Fatalf("legacy state did not migrate: %+v (%v)", state, err)
	}
	for _, record := range state.Installations {
		if strings.Join(record.Consumers, ",") != "codex,gemini-cli" {
			t.Fatalf("legacy consumers not merged: %v", record.Consumers)
		}
	}
}

func TestClaudeConfigDirExpansionAndRelativePathRejection(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	t.Setenv("CLAUDE_CONFIG_DIR", "~/.claude-custom")
	root, err := agentTargetRoot("claude-code", "user", home)
	if err != nil {
		t.Fatal(err)
	}
	if root != filepath.Join(home, ".claude-custom", "skills") {
		t.Fatalf("unexpected expanded root: %s", root)
	}
	t.Setenv("CLAUDE_CONFIG_DIR", "relative/config")
	if _, err := agentTargetRoot("claude-code", "user", home); err == nil {
		t.Fatal("relative CLAUDE_CONFIG_DIR was accepted")
	}
}

func TestProjectSymlinkEscapeIsRejected(t *testing.T) {
	project := t.TempDir()
	outside := t.TempDir()
	if err := os.Symlink(outside, filepath.Join(project, ".agents")); err != nil {
		t.Skipf("symlink unavailable: %v", err)
	}
	if _, err := agentTargetRoot("codex", "project", project); err == nil || !strings.Contains(err.Error(), "escapes") {
		t.Fatalf("project symlink escape was accepted: %v", err)
	}
}

func TestAdjacentTransactionAreaSymlinkEscapeIsRejected(t *testing.T) {
	catalog, err := loadCatalog()
	if err != nil {
		t.Fatal(err)
	}
	project := t.TempDir()
	outside := t.TempDir()
	if err := os.MkdirAll(filepath.Join(project, ".agents"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(project, ".agents", ".anywork")); err != nil {
		t.Skipf("symlink unavailable: %v", err)
	}
	emit := func(string, ...any) {}
	err = installNative(catalog, []string{"clarify-outcome"}, []string{"codex"}, "project", project, false, false, emit)
	if err == nil || !strings.Contains(err.Error(), "escapes") {
		t.Fatalf("transaction area symlink escape was accepted: %v", err)
	}
	entries, err := os.ReadDir(outside)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 0 {
		t.Fatalf("installer wrote through transaction symlink: %v", entries)
	}
}

func TestLockReleaseUsesOwnerNonceAndDeadOwnerIsReclaimed(t *testing.T) {
	statePath := filepath.Join(t.TempDir(), ".anywork", "state.json")
	release, err := acquireMutationLock(statePath)
	if err != nil {
		t.Fatal(err)
	}
	lockPath := filepath.Join(filepath.Dir(statePath), "mutation.lock")
	if err := os.Remove(lockPath); err != nil {
		t.Fatal(err)
	}
	replacement := lockOwner{PID: os.Getpid(), Nonce: "replacement-owner", Created: time.Now().UTC().Format(time.RFC3339Nano)}
	if err := writeJSON(lockPath, replacement, 0o600); err != nil {
		t.Fatal(err)
	}
	release()
	if _, err := os.Stat(lockPath); err != nil {
		t.Fatalf("old owner removed replacement lock: %v", err)
	}
	if err := os.Remove(lockPath); err != nil {
		t.Fatal(err)
	}
	dead := lockOwner{PID: 2147483647, Nonce: "dead-owner", Created: time.Now().UTC().Format(time.RFC3339Nano)}
	if err := writeJSON(lockPath, dead, 0o600); err != nil {
		t.Fatal(err)
	}
	newRelease, err := acquireMutationLock(statePath)
	if err != nil {
		t.Fatalf("dead owner lock was not reclaimed: %v", err)
	}
	newRelease()
}

func TestStateReplacementRemainsReadable(t *testing.T) {
	statePath := filepath.Join(t.TempDir(), "state.json")
	for index := 0; index < 20; index++ {
		state := emptyState()
		state.LastTransaction = strings.Repeat("x", index)
		if err := writeState(statePath, state); err != nil {
			t.Fatal(err)
		}
		loaded, err := loadState(statePath)
		if err != nil || loaded.LastTransaction != state.LastTransaction {
			t.Fatalf("replacement %d unreadable: %+v (%v)", index, loaded, err)
		}
	}
}

func TestNativeMutationLockRejectsConcurrentWriter(t *testing.T) {
	statePath := filepath.Join(t.TempDir(), ".anywork", "state.json")
	release, err := acquireMutationLock(statePath)
	if err != nil {
		t.Fatal(err)
	}
	defer release()
	if _, err := acquireMutationLock(statePath); err == nil {
		t.Fatal("second writer acquired the same mutation lock")
	}
}

func TestNativeCLIParsingAcceptsFlagsAfterPack(t *testing.T) {
	options, err := parseMutationOptions([]string{"essential", "--agent", "codex", "--scope", "project", "--dry-run"})
	if err != nil {
		t.Fatal(err)
	}
	if options.Pack != "essential" || len(options.Agents) != 1 || options.Agents[0] != "codex" || options.Scope != "project" || !options.DryRun {
		t.Fatalf("unexpected options: %+v", options)
	}
}
