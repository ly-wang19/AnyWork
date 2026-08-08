package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

var supportedAgents = []string{"codex", "claude-code", "gemini-cli", "github-copilot", "opencode"}

var (
	// transactionFaultHook is nil in production. Tests use it to model a process
	// disappearing between durable journal and filesystem/state transitions.
	transactionFaultHook func(string) error
	errSimulatedCrash    = errors.New("simulated process crash")
)

type installationRecord struct {
	Agent       string   `json:"agent"`
	Consumers   []string `json:"consumers,omitempty"`
	Skill       string   `json:"skill"`
	Target      string   `json:"target"`
	Hash        string   `json:"hash"`
	Version     string   `json:"version"`
	InstalledAt string   `json:"installed_at"`
}

type stateFile struct {
	SchemaVersion   int                           `json:"schema_version"`
	LastTransaction string                        `json:"last_transaction,omitempty"`
	Installations   map[string]installationRecord `json:"installations"`
}

type lockOwner struct {
	PID     int    `json:"pid"`
	Nonce   string `json:"nonce"`
	Created string `json:"created"`
}

type journalOperation struct {
	Action         string `json:"action"`
	Agent          string `json:"agent"`
	Skill          string `json:"skill"`
	Target         string `json:"target,omitempty"`
	Staging        string `json:"staging,omitempty"`
	Backup         string `json:"backup,omitempty"`
	Recovery       string `json:"recovery,omitempty"`
	ExpectedHash   string `json:"expected_hash,omitempty"`
	OriginalExists bool   `json:"original_exists,omitempty"`
	Status         string `json:"status"`
}

type transactionJournal struct {
	SchemaVersion int                `json:"schema_version"`
	ID            string             `json:"id"`
	Kind          string             `json:"kind"`
	Phase         string             `json:"phase"`
	StatePath     string             `json:"state_path"`
	CreatedAt     string             `json:"created_at"`
	BeforeState   stateFile          `json:"before_state"`
	AfterState    stateFile          `json:"after_state"`
	Operations    []journalOperation `json:"operations"`
}

type eventEmitter func(string, ...any)

func normalizeAgent(agent string) string {
	switch agent {
	case "claude", "claude_code":
		return "claude-code"
	case "gemini":
		return "gemini-cli"
	case "copilot", "github_copilot":
		return "github-copilot"
	case "open-code":
		return "opencode"
	default:
		return agent
	}
}

func isSupportedAgent(agent string) bool {
	for _, candidate := range supportedAgents {
		if agent == candidate {
			return true
		}
	}
	return false
}

func targetFamily(agent string) string {
	switch normalizeAgent(agent) {
	case "claude-code", "cursor", "cursor-agent":
		return "portable-skills"
	default:
		return "agents-skills"
	}
}

func expandHomePath(value, home string) (string, error) {
	if value == "~" {
		return home, nil
	}
	if strings.HasPrefix(value, "~/") || strings.HasPrefix(value, `~\`) {
		return filepath.Join(home, value[2:]), nil
	}
	if strings.HasPrefix(value, "~") {
		return "", fmt.Errorf("user-specific home expansion is not supported: %s", value)
	}
	return value, nil
}

func absoluteConfiguredPath(name, value, home string) (string, error) {
	expanded, err := expandHomePath(value, home)
	if err != nil {
		return "", fmt.Errorf("%s: %w", name, err)
	}
	if !filepath.IsAbs(expanded) {
		return "", fmt.Errorf("%s must be an absolute path (or start with ~/): %s", name, value)
	}
	return filepath.Clean(expanded), nil
}

func pathWithin(base, candidate string) bool {
	relative, err := filepath.Rel(base, candidate)
	if err != nil {
		return false
	}
	return relative == "." || (relative != ".." && !strings.HasPrefix(relative, ".."+string(os.PathSeparator)) && !filepath.IsAbs(relative))
}

// resolveProjectChild follows every existing symlink in the requested project
// path and child prefix. A .agents/.claude/.anywork symlink may point elsewhere
// inside the project, but never outside it.
func resolveProjectChild(projectDirectory string, components ...string) (string, error) {
	base, err := filepath.Abs(projectDirectory)
	if err != nil {
		return "", err
	}
	resolvedBase, err := filepath.EvalSymlinks(base)
	if err != nil {
		return "", fmt.Errorf("resolve project directory %s: %w", base, err)
	}
	candidate := filepath.Join(append([]string{resolvedBase}, components...)...)
	existing := candidate
	missing := []string{}
	for {
		_, statErr := os.Lstat(existing)
		if statErr == nil {
			break
		}
		if !errors.Is(statErr, os.ErrNotExist) {
			return "", statErr
		}
		parent := filepath.Dir(existing)
		if parent == existing {
			return "", fmt.Errorf("cannot resolve project child: %s", candidate)
		}
		missing = append(missing, filepath.Base(existing))
		existing = parent
	}
	resolvedExisting, err := filepath.EvalSymlinks(existing)
	if err != nil {
		return "", err
	}
	resolvedCandidate := resolvedExisting
	for index := len(missing) - 1; index >= 0; index-- {
		resolvedCandidate = filepath.Join(resolvedCandidate, missing[index])
	}
	if !pathWithin(resolvedBase, resolvedCandidate) {
		return "", fmt.Errorf("project path escapes through a symlink: %s", candidate)
	}
	return filepath.Clean(resolvedCandidate), nil
}

func resolveProspectivePath(candidate string) (string, error) {
	existing := filepath.Clean(candidate)
	missing := []string{}
	for {
		_, statErr := os.Lstat(existing)
		if statErr == nil {
			break
		}
		if !errors.Is(statErr, os.ErrNotExist) {
			return "", statErr
		}
		parent := filepath.Dir(existing)
		if parent == existing {
			return "", fmt.Errorf("cannot resolve path: %s", candidate)
		}
		missing = append(missing, filepath.Base(existing))
		existing = parent
	}
	resolved, err := filepath.EvalSymlinks(existing)
	if err != nil {
		return "", err
	}
	for index := len(missing) - 1; index >= 0; index-- {
		resolved = filepath.Join(resolved, missing[index])
	}
	return filepath.Clean(resolved), nil
}

func validateContainedPath(base, candidate string) error {
	resolvedBase, err := resolveProspectivePath(base)
	if err != nil {
		return err
	}
	resolvedCandidate, err := resolveProspectivePath(candidate)
	if err != nil {
		return err
	}
	if !pathWithin(resolvedBase, resolvedCandidate) {
		return fmt.Errorf("managed path escapes through a symlink: %s", candidate)
	}
	return nil
}

func validateOperationContainment(operation journalOperation) error {
	if operation.Target == "" {
		return nil
	}
	skillsRoot := filepath.Dir(operation.Target)
	agentRoot := filepath.Dir(skillsRoot)
	for _, candidate := range []string{operation.Target, operation.Staging, operation.Backup, operation.Recovery} {
		if candidate != "" {
			if err := validateContainedPath(agentRoot, candidate); err != nil {
				return err
			}
		}
	}
	return nil
}

func agentTargetRoot(agent, scope, projectDirectory string) (string, error) {
	agent = normalizeAgent(agent)
	if !isSupportedAgent(agent) {
		return "", fmt.Errorf("unsupported agent: %s", agent)
	}
	if scope == "project" {
		if agent == "claude-code" {
			return resolveProjectChild(projectDirectory, ".claude", "skills")
		}
		return resolveProjectChild(projectDirectory, ".agents", "skills")
	}
	if scope != "user" {
		return "", fmt.Errorf("unsupported scope: %s", scope)
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve user home: %w", err)
	}
	if agent == "claude-code" {
		claudeRoot := os.Getenv("CLAUDE_CONFIG_DIR")
		if claudeRoot == "" {
			claudeRoot = filepath.Join(home, ".claude")
		} else if claudeRoot, err = absoluteConfiguredPath("CLAUDE_CONFIG_DIR", claudeRoot, home); err != nil {
			return "", err
		}
		return filepath.Join(claudeRoot, "skills"), nil
	}
	return filepath.Join(home, ".agents", "skills"), nil
}

func anyworkStatePath(scope, projectDirectory string) (string, error) {
	if scope == "project" {
		return resolveProjectChild(projectDirectory, ".anywork", "state.json")
	}
	if scope != "user" {
		return "", fmt.Errorf("unsupported scope: %s", scope)
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	configured := os.Getenv("ANYWORK_HOME")
	if configured != "" {
		configured, err = absoluteConfiguredPath("ANYWORK_HOME", configured, home)
		if err != nil {
			return "", err
		}
		return filepath.Join(configured, "state.json"), nil
	}
	return filepath.Join(home, ".anywork", "state.json"), nil
}

func emptyState() stateFile {
	return stateFile{SchemaVersion: 1, Installations: map[string]installationRecord{}}
}

func cloneState(source stateFile) stateFile {
	copy := stateFile{SchemaVersion: source.SchemaVersion, LastTransaction: source.LastTransaction, Installations: map[string]installationRecord{}}
	for key, record := range source.Installations {
		record.Consumers = append([]string(nil), record.Consumers...)
		copy.Installations[key] = record
	}
	return copy
}

func loadState(filePath string) (stateFile, error) {
	state := emptyState()
	data, err := os.ReadFile(filePath)
	if errors.Is(err, os.ErrNotExist) {
		return state, nil
	}
	if err != nil {
		return stateFile{}, fmt.Errorf("read state %s: %w", filePath, err)
	}
	if err := json.Unmarshal(data, &state); err != nil {
		return stateFile{}, fmt.Errorf("decode state %s: %w", filePath, err)
	}
	if state.SchemaVersion != 1 || state.Installations == nil {
		return stateFile{}, fmt.Errorf("invalid state file: %s", filePath)
	}
	return normalizeOwnershipState(state)
}

func installationKey(target, skill string) string {
	return "target:" + filepath.Clean(target) + "|skill:" + skill
}

func containsString(values []string, candidate string) bool {
	for _, value := range values {
		if value == candidate {
			return true
		}
	}
	return false
}

func addConsumer(record installationRecord, agent string) installationRecord {
	if !containsString(record.Consumers, agent) {
		record.Consumers = append(record.Consumers, agent)
		sort.Strings(record.Consumers)
	}
	record.Agent = ""
	return record
}

func removeConsumer(record installationRecord, agent string) installationRecord {
	remaining := make([]string, 0, len(record.Consumers))
	for _, consumer := range record.Consumers {
		if consumer != agent {
			remaining = append(remaining, consumer)
		}
	}
	record.Consumers = remaining
	record.Agent = ""
	return record
}

// normalizeOwnershipState migrates alpha state keyed by agent:target to one
// ownership record per resolved target and skill. Consumers prevent one agent
// from uninstalling a shared .agents/skills materialization still used by others.
func normalizeOwnershipState(state stateFile) (stateFile, error) {
	normalized := stateFile{SchemaVersion: state.SchemaVersion, LastTransaction: state.LastTransaction, Installations: map[string]installationRecord{}}
	for _, record := range state.Installations {
		if record.Target == "" || record.Skill == "" {
			return stateFile{}, errors.New("invalid state ownership record")
		}
		consumers := append([]string(nil), record.Consumers...)
		if record.Agent != "" && !containsString(consumers, record.Agent) {
			consumers = append(consumers, normalizeAgent(record.Agent))
		}
		record.Agent = ""
		record.Consumers = nil
		for _, consumer := range consumers {
			record = addConsumer(record, normalizeAgent(consumer))
		}
		key := installationKey(record.Target, record.Skill)
		if previous, exists := normalized.Installations[key]; exists {
			if previous.Hash != record.Hash || previous.Target != record.Target || previous.Skill != record.Skill {
				return stateFile{}, fmt.Errorf("conflicting ownership records for %s", record.Target)
			}
			for _, consumer := range record.Consumers {
				previous = addConsumer(previous, consumer)
			}
			record = previous
		}
		normalized.Installations[key] = record
	}
	return normalized, nil
}

func writeJSON(filePath string, value any, mode fs.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(filePath), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	temporary, err := os.CreateTemp(filepath.Dir(filePath), ".anywork-write-*.tmp")
	if err != nil {
		return err
	}
	temporaryPath := temporary.Name()
	defer os.Remove(temporaryPath)
	if err := temporary.Chmod(mode); err != nil {
		temporary.Close()
		return err
	}
	if _, err := temporary.Write(data); err != nil {
		temporary.Close()
		return err
	}
	if err := temporary.Sync(); err != nil {
		temporary.Close()
		return err
	}
	if err := temporary.Close(); err != nil {
		return err
	}
	if err := atomicReplace(temporaryPath, filePath); err != nil {
		return err
	}
	return syncDirectory(filepath.Dir(filePath))
}

func writeState(filePath string, state stateFile) error {
	return writeJSON(filePath, state, 0o600)
}

func randomNonce() (string, error) {
	buffer := make([]byte, 16)
	if _, err := rand.Read(buffer); err != nil {
		return "", err
	}
	return hex.EncodeToString(buffer), nil
}

func readLockOwner(lockPath string) (lockOwner, []byte, error) {
	data, err := os.ReadFile(lockPath)
	if err != nil {
		return lockOwner{}, nil, err
	}
	var owner lockOwner
	if err := json.Unmarshal(data, &owner); err != nil {
		return lockOwner{}, data, err
	}
	if owner.PID <= 0 || owner.Nonce == "" || owner.Created == "" {
		return lockOwner{}, data, errors.New("invalid lock owner")
	}
	return owner, data, nil
}

func removeLockIfUnchanged(lockPath string, expected []byte) bool {
	current, err := os.ReadFile(lockPath)
	if err != nil || !stringEqualBytes(current, expected) {
		return false
	}
	return os.Remove(lockPath) == nil
}

func stringEqualBytes(left, right []byte) bool {
	if len(left) != len(right) {
		return false
	}
	var different byte
	for index := range left {
		different |= left[index] ^ right[index]
	}
	return different == 0
}

func acquireMutationLock(statePath string) (func(), error) {
	directory := filepath.Dir(statePath)
	if err := os.MkdirAll(directory, 0o755); err != nil {
		return nil, err
	}
	lockPath := filepath.Join(directory, "mutation.lock")
	nonce, err := randomNonce()
	if err != nil {
		return nil, err
	}
	owner := lockOwner{PID: os.Getpid(), Nonce: nonce, Created: time.Now().UTC().Format(time.RFC3339Nano)}
	data, err := json.Marshal(owner)
	if err != nil {
		return nil, err
	}
	data = append(data, '\n')
	open := func() error {
		file, openErr := os.OpenFile(lockPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
		if openErr != nil {
			return openErr
		}
		if _, writeErr := file.Write(data); writeErr != nil {
			file.Close()
			_ = os.Remove(lockPath)
			return writeErr
		}
		if syncErr := file.Sync(); syncErr != nil {
			file.Close()
			_ = os.Remove(lockPath)
			return syncErr
		}
		return file.Close()
	}
	if err := open(); errors.Is(err, os.ErrExist) {
		existingOwner, existingData, readErr := readLockOwner(lockPath)
		reclaim := false
		if readErr == nil {
			reclaim = !processAlive(existingOwner.PID)
		} else if info, statErr := os.Stat(lockPath); statErr == nil {
			reclaim = time.Since(info.ModTime()) > time.Hour
		}
		if reclaim && removeLockIfUnchanged(lockPath, existingData) {
			err = open()
		}
		if err != nil {
			return nil, fmt.Errorf("mutation lock %s: %w", lockPath, err)
		}
	} else if err != nil {
		return nil, fmt.Errorf("mutation lock %s: %w", lockPath, err)
	}
	return func() {
		current, _, readErr := readLockOwner(lockPath)
		if readErr == nil && current.Nonce == nonce && current.PID == os.Getpid() {
			_ = os.Remove(lockPath)
		}
	}, nil
}

func hashDirectory(directory string) (string, error) {
	entries := []string{}
	err := filepath.WalkDir(directory, func(filePath string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.Type()&os.ModeSymlink != 0 {
			return fmt.Errorf("symlink is not allowed in a managed skill: %s", filePath)
		}
		if entry.IsDir() {
			return nil
		}
		entries = append(entries, filePath)
		return nil
	})
	if err != nil {
		return "", err
	}
	sort.Strings(entries)
	hash := sha256.New()
	for _, filePath := range entries {
		relative, err := filepath.Rel(directory, filePath)
		if err != nil {
			return "", err
		}
		relative = filepath.ToSlash(relative)
		data, err := os.ReadFile(filePath)
		if err != nil {
			return "", err
		}
		hash.Write([]byte(relative))
		hash.Write([]byte{0})
		hash.Write(data)
		hash.Write([]byte{0})
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func materializeSkill(skill skillRecord, agent, destination string) error {
	family := targetFamily(agent)
	return fs.WalkDir(officialContent, skill.Path, func(filePath string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		relative := strings.TrimPrefix(filePath, skill.Path)
		relative = strings.TrimPrefix(relative, "/")
		if relative == "" {
			return os.MkdirAll(destination, 0o755)
		}
		if family == "portable-skills" && (relative == "agents" || strings.HasPrefix(relative, "agents/")) {
			if entry.IsDir() {
				return fs.SkipDir
			}
			return nil
		}
		target := filepath.Join(destination, filepath.FromSlash(relative))
		if !pathWithin(filepath.Clean(destination), filepath.Clean(target)) {
			return fmt.Errorf("embedded path escapes skill destination: %s", relative)
		}
		if entry.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		data, err := officialContent.ReadFile(filePath)
		if err != nil {
			return err
		}
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return err
		}
		return os.WriteFile(target, data, 0o644)
	})
}

func transactionArea(root string) string {
	return filepath.Join(filepath.Dir(root), ".anywork")
}

func transactionPaths(root, transactionID, agent, skillID string) (string, string, string) {
	area := transactionArea(root)
	return filepath.Join(area, "staging", transactionID, agent, skillID),
		filepath.Join(area, "backups", transactionID, agent, skillID),
		filepath.Join(area, "recovery", transactionID, agent, skillID)
}

func journalPathForState(statePath string) string {
	return filepath.Join(filepath.Dir(statePath), "transaction.json")
}

func writeJournal(journal *transactionJournal) error {
	return writeJSON(journalPathForState(journal.StatePath), journal, 0o600)
}

func loadJournal(statePath string) (*transactionJournal, error) {
	journalPath := journalPathForState(statePath)
	data, err := os.ReadFile(journalPath)
	if err != nil {
		return nil, err
	}
	var journal transactionJournal
	if err := json.Unmarshal(data, &journal); err != nil {
		return nil, fmt.Errorf("decode transaction journal %s: %w", journalPath, err)
	}
	if journal.SchemaVersion != 1 || journal.ID == "" || journal.StatePath != statePath {
		return nil, fmt.Errorf("invalid transaction journal: %s", journalPath)
	}
	return &journal, nil
}

func removeJournal(statePath string) error {
	path := journalPathForState(statePath)
	if err := os.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return syncDirectory(filepath.Dir(path))
}

func fault(point string) error {
	if transactionFaultHook == nil {
		return nil
	}
	return transactionFaultHook(point)
}

func newTransactionID() (string, error) {
	nonce, err := randomNonce()
	if err != nil {
		return "", err
	}
	return time.Now().UTC().Format("20060102T150405.000000000Z") + "-" + nonce, nil
}

func startJournal(kind, statePath string, before, after stateFile, operations []journalOperation) (*transactionJournal, error) {
	if _, err := os.Stat(journalPathForState(statePath)); err == nil {
		return nil, fmt.Errorf("unfinished transaction exists; run anywork recover")
	} else if !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}
	id, err := newTransactionID()
	if err != nil {
		return nil, err
	}
	for index := range operations {
		operation := &operations[index]
		if operation.Action == "adopt" {
			continue
		}
		root := filepath.Dir(operation.Target)
		operation.Staging, operation.Backup, operation.Recovery = transactionPaths(root, id, operation.Agent, operation.Skill)
		if operation.Action == "uninstall" {
			operation.Staging = ""
		}
		if err := validateOperationContainment(*operation); err != nil {
			return nil, err
		}
	}
	after.LastTransaction = id
	journal := &transactionJournal{
		SchemaVersion: 1,
		ID:            id,
		Kind:          kind,
		Phase:         "applying",
		StatePath:     statePath,
		CreatedAt:     time.Now().UTC().Format(time.RFC3339Nano),
		BeforeState:   before,
		AfterState:    after,
		Operations:    operations,
	}
	if err := writeJournal(journal); err != nil {
		return nil, err
	}
	return journal, fault("after-journal-begin")
}

func ensureDirectoryAbsent(path string) error {
	if _, err := os.Lstat(path); err == nil {
		return fmt.Errorf("transaction path already exists: %s", path)
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return nil
}

func movePreserving(path, destination string) error {
	if _, err := os.Lstat(path); errors.Is(err, os.ErrNotExist) {
		return nil
	} else if err != nil {
		return err
	}
	selected := destination
	for suffix := 1; ; suffix++ {
		if _, err := os.Lstat(selected); errors.Is(err, os.ErrNotExist) {
			break
		} else if err != nil {
			return err
		}
		selected = destination + "." + strconv.Itoa(suffix)
	}
	if err := os.MkdirAll(filepath.Dir(selected), 0o755); err != nil {
		return err
	}
	return os.Rename(path, selected)
}

func rollbackJournal(journal *transactionJournal) error {
	journal.Phase = "rolling_back"
	if err := writeJournal(journal); err != nil {
		return err
	}
	for index := len(journal.Operations) - 1; index >= 0; index-- {
		operation := &journal.Operations[index]
		if err := validateOperationContainment(*operation); err != nil {
			return err
		}
		if operation.Action == "adopt" {
			operation.Status = "rolled_back"
			continue
		}
		backupExists := false
		if operation.Backup != "" {
			if _, err := os.Lstat(operation.Backup); err == nil {
				backupExists = true
			} else if !errors.Is(err, os.ErrNotExist) {
				return err
			}
		}
		if operation.OriginalExists && backupExists {
			if err := movePreserving(operation.Target, operation.Recovery); err != nil {
				return err
			}
			if err := os.MkdirAll(filepath.Dir(operation.Target), 0o755); err != nil {
				return err
			}
			if err := os.Rename(operation.Backup, operation.Target); err != nil {
				return err
			}
		} else if !operation.OriginalExists {
			if err := movePreserving(operation.Target, operation.Recovery); err != nil {
				return err
			}
		}
		if operation.Staging != "" {
			if err := os.RemoveAll(operation.Staging); err != nil {
				return err
			}
		}
		operation.Status = "rolled_back"
		if err := writeJournal(journal); err != nil {
			return err
		}
	}
	if err := writeState(journal.StatePath, journal.BeforeState); err != nil {
		return err
	}
	return removeJournal(journal.StatePath)
}

func finishCommittedJournal(journal *transactionJournal) error {
	for index := range journal.Operations {
		if err := validateOperationContainment(journal.Operations[index]); err != nil {
			return err
		}
		if staging := journal.Operations[index].Staging; staging != "" {
			if err := os.RemoveAll(staging); err != nil {
				return err
			}
		}
		journal.Operations[index].Status = "committed"
	}
	journal.Phase = "committed"
	if err := writeJournal(journal); err != nil {
		return err
	}
	return removeJournal(journal.StatePath)
}

func failTransaction(journal *transactionJournal, cause error) error {
	if errors.Is(cause, errSimulatedCrash) {
		return cause
	}
	if rollbackErr := rollbackJournal(journal); rollbackErr != nil {
		return fmt.Errorf("transaction failed: %v; automatic rollback failed: %w; run anywork recover", cause, rollbackErr)
	}
	return cause
}

func applyJournal(c catalogIndex, journal *transactionJournal, emit eventEmitter) error {
	for index := range journal.Operations {
		operation := &journal.Operations[index]
		if err := validateOperationContainment(*operation); err != nil {
			return failTransaction(journal, err)
		}
		switch operation.Action {
		case "adopt":
			operation.Status = "adopted"
		case "install":
			skill := c.Skills[operation.Skill]
			if err := ensureDirectoryAbsent(operation.Staging); err != nil {
				return failTransaction(journal, err)
			}
			if err := materializeSkill(skill, operation.Agent, operation.Staging); err != nil {
				return failTransaction(journal, err)
			}
			operation.Status = "staged"
			if err := writeJournal(journal); err != nil {
				return failTransaction(journal, err)
			}
			if err := fault("after-stage"); err != nil {
				return failTransaction(journal, err)
			}
			if operation.OriginalExists {
				if err := ensureDirectoryAbsent(operation.Backup); err != nil {
					return failTransaction(journal, err)
				}
				if err := os.MkdirAll(filepath.Dir(operation.Backup), 0o755); err != nil {
					return failTransaction(journal, err)
				}
				if err := os.Rename(operation.Target, operation.Backup); err != nil {
					return failTransaction(journal, err)
				}
				if err := fault("after-backup"); err != nil {
					return failTransaction(journal, err)
				}
				emit("backup", operation.Backup)
			}
			if err := os.MkdirAll(filepath.Dir(operation.Target), 0o755); err != nil {
				return failTransaction(journal, err)
			}
			if err := os.Rename(operation.Staging, operation.Target); err != nil {
				return failTransaction(journal, err)
			}
			operation.Status = "applied"
			if err := writeJournal(journal); err != nil {
				return failTransaction(journal, err)
			}
			if err := fault("after-target-rename"); err != nil {
				return failTransaction(journal, err)
			}
			if operation.OriginalExists {
				emit("updated", operation.Skill, operation.Agent, operation.Target)
			} else {
				emit("installed", operation.Skill, operation.Agent, operation.Target)
			}
		case "uninstall":
			if err := ensureDirectoryAbsent(operation.Backup); err != nil {
				return failTransaction(journal, err)
			}
			if err := os.MkdirAll(filepath.Dir(operation.Backup), 0o755); err != nil {
				return failTransaction(journal, err)
			}
			if err := os.Rename(operation.Target, operation.Backup); err != nil {
				return failTransaction(journal, err)
			}
			if err := fault("after-backup"); err != nil {
				return failTransaction(journal, err)
			}
			operation.Status = "applied"
			if err := writeJournal(journal); err != nil {
				return failTransaction(journal, err)
			}
			emit("removed", operation.Skill, operation.Agent, operation.Backup)
		default:
			return failTransaction(journal, fmt.Errorf("unknown journal operation: %s", operation.Action))
		}
		if err := writeJournal(journal); err != nil {
			return failTransaction(journal, err)
		}
	}
	journal.Phase = "committing_state"
	if err := writeJournal(journal); err != nil {
		return failTransaction(journal, err)
	}
	if err := writeState(journal.StatePath, journal.AfterState); err != nil {
		return failTransaction(journal, err)
	}
	if err := fault("after-state-write"); err != nil {
		return failTransaction(journal, err)
	}
	return finishCommittedJournal(journal)
}

func ensureTargetDirectory(target string) (bool, error) {
	info, err := os.Lstat(target)
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return false, fmt.Errorf("skill target must not be a symlink: %s", target)
	}
	if !info.IsDir() {
		return false, fmt.Errorf("skill target is not a directory: %s", target)
	}
	return true, nil
}

func installNative(c catalogIndex, skillIDs, agents []string, scope, projectDirectory string, dryRun, force bool, emit eventEmitter) error {
	statePath, err := anyworkStatePath(scope, projectDirectory)
	if err != nil {
		return err
	}
	var release func()
	if !dryRun {
		release, err = acquireMutationLock(statePath)
		if err != nil {
			return err
		}
		defer release()
	}
	state, err := loadState(statePath)
	if err != nil {
		return err
	}
	after := cloneState(state)
	operations := []journalOperation{}
	operationByKey := map[string]int{}
	for _, rawAgent := range agents {
		agent := normalizeAgent(rawAgent)
		root, err := agentTargetRoot(agent, scope, projectDirectory)
		if err != nil {
			return err
		}
		for _, skillID := range skillIDs {
			skill, ok := c.Skills[skillID]
			if !ok {
				return fmt.Errorf("unknown skill: %s", skillID)
			}
			target := filepath.Join(root, skillID)
			if !pathWithin(root, target) {
				return fmt.Errorf("skill target escapes agent root: %s", target)
			}
			expectedHash, err := digestEmbeddedSkill(skill, agent)
			if err != nil {
				return err
			}
			key := installationKey(target, skillID)
			record, managed := state.Installations[key]
			if plannedIndex, planned := operationByKey[key]; planned {
				plannedRecord := addConsumer(after.Installations[key], agent)
				after.Installations[key] = plannedRecord
				if operations[plannedIndex].Action == "adopt" {
					emit("unchanged", skillID, agent)
				}
				continue
			}
			targetExists, err := ensureTargetDirectory(target)
			if err != nil {
				return err
			}
			if targetExists {
				currentHash, err := hashDirectory(target)
				if err != nil {
					return err
				}
				if currentHash == expectedHash {
					emit("unchanged", skillID, agent)
					if !managed || !containsString(record.Consumers, agent) {
						if !managed {
							record = installationRecord{Skill: skillID, Target: target, Hash: expectedHash, Version: skill.Version, InstalledAt: time.Now().UTC().Format(time.RFC3339Nano)}
						}
						record = addConsumer(record, agent)
						after.Installations[key] = record
						operationByKey[key] = len(operations)
						operations = append(operations, journalOperation{Action: "adopt", Agent: agent, Skill: skillID, Target: target, ExpectedHash: expectedHash, OriginalExists: true, Status: "planned"})
					}
					continue
				}
				if (!managed || record.Hash != currentHash) && !force {
					return fmt.Errorf(msg("en", "conflict", target))
				}
				if dryRun {
					emit("plan_replace", skillID, target)
					continue
				}
			} else if dryRun {
				emit("plan_install", skillID, target)
				continue
			}
			replacement := installationRecord{Skill: skillID, Target: target, Hash: expectedHash, Version: skill.Version, InstalledAt: time.Now().UTC().Format(time.RFC3339Nano)}
			if managed {
				replacement.Consumers = append(replacement.Consumers, record.Consumers...)
			}
			replacement = addConsumer(replacement, agent)
			after.Installations[key] = replacement
			operationByKey[key] = len(operations)
			operations = append(operations, journalOperation{Action: "install", Agent: agent, Skill: skillID, Target: target, ExpectedHash: expectedHash, OriginalExists: targetExists, Status: "planned"})
		}
	}
	if dryRun || len(operations) == 0 {
		return nil
	}
	journal, err := startJournal("install", statePath, state, after, operations)
	if err != nil {
		if errors.Is(err, errSimulatedCrash) {
			return err
		}
		return err
	}
	return applyJournal(c, journal, emit)
}

func uninstallNative(skillIDs, agents []string, scope, projectDirectory string, dryRun, force bool, emit eventEmitter) error {
	statePath, err := anyworkStatePath(scope, projectDirectory)
	if err != nil {
		return err
	}
	var release func()
	if !dryRun {
		release, err = acquireMutationLock(statePath)
		if err != nil {
			return err
		}
		defer release()
	}
	state, err := loadState(statePath)
	if err != nil {
		return err
	}
	after := cloneState(state)
	operations := []journalOperation{}
	for _, rawAgent := range agents {
		agent := normalizeAgent(rawAgent)
		root, err := agentTargetRoot(agent, scope, projectDirectory)
		if err != nil {
			return err
		}
		for _, skillID := range skillIDs {
			target := filepath.Join(root, skillID)
			key := installationKey(target, skillID)
			record, managed := after.Installations[key]
			if !managed || !containsString(record.Consumers, agent) {
				continue
			}
			record = removeConsumer(record, agent)
			if len(record.Consumers) > 0 {
				after.Installations[key] = record
				operations = append(operations, journalOperation{Action: "adopt", Agent: agent, Skill: skillID, Target: target, OriginalExists: true, Status: "planned"})
				continue
			}
			targetExists, err := ensureTargetDirectory(target)
			if err != nil {
				return err
			}
			if !targetExists {
				delete(after.Installations, key)
				operations = append(operations, journalOperation{Action: "adopt", Agent: agent, Skill: skillID, Target: target, Status: "planned"})
				continue
			}
			currentHash, err := hashDirectory(target)
			if err != nil {
				return err
			}
			if currentHash != record.Hash && !force {
				return fmt.Errorf(msg("en", "not_owned", target))
			}
			if dryRun {
				emit("plan_remove", skillID, target)
				continue
			}
			delete(after.Installations, key)
			operations = append(operations, journalOperation{Action: "uninstall", Agent: agent, Skill: skillID, Target: target, OriginalExists: true, Status: "planned"})
		}
	}
	if dryRun || len(operations) == 0 {
		return nil
	}
	journal, err := startJournal("uninstall", statePath, state, after, operations)
	if err != nil {
		return err
	}
	return applyJournal(catalogIndex{}, journal, emit)
}

func recoverNative(scope, projectDirectory string, emit eventEmitter) error {
	statePath, err := anyworkStatePath(scope, projectDirectory)
	if err != nil {
		return err
	}
	release, err := acquireMutationLock(statePath)
	if err != nil {
		return err
	}
	defer release()
	journal, err := loadJournal(statePath)
	if errors.Is(err, os.ErrNotExist) {
		emit("recover_none")
		return nil
	}
	if err != nil {
		return err
	}
	state, err := loadState(statePath)
	if err != nil {
		return err
	}
	if state.LastTransaction == journal.ID {
		if err := finishCommittedJournal(journal); err != nil {
			return err
		}
		emit("recover_committed", journal.ID)
		return nil
	}
	if err := rollbackJournal(journal); err != nil {
		return err
	}
	emit("recover_rolled_back", journal.ID)
	return nil
}
