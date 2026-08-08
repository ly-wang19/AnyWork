package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/fs"
	"path"
	"regexp"
	"sort"
	"strings"
)

var skillIDPattern = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

type localizedText map[string]string

type sourceMetadata struct {
	Type       string `json:"type"`
	Repository string `json:"repository"`
}

type languageReview struct {
	Status       string `json:"status"`
	Revision     string `json:"revision"`
	ReviewedAt   string `json:"reviewed_at,omitempty"`
	ReviewerRole string `json:"reviewer_role,omitempty"`
}

type skillRecord struct {
	ID           string                    `json:"id"`
	Path         string                    `json:"path"`
	Version      string                    `json:"version"`
	Maturity     string                    `json:"maturity"`
	Class        string                    `json:"class"`
	Languages    []string                  `json:"languages"`
	Reviews      map[string]languageReview `json:"language_reviews"`
	Capabilities []string                  `json:"capabilities"`
	License      string                    `json:"license"`
	Source       sourceMetadata            `json:"source"`
	Risk         string                    `json:"risk"`
	SideEffects  []string                  `json:"side_effects"`
}

type packRecord struct {
	ID          string        `json:"id"`
	DisplayName localizedText `json:"display_name"`
	Description localizedText `json:"description"`
	Extends     []string      `json:"extends"`
	Skills      []string      `json:"skills"`
}

type catalogFile struct {
	SchemaVersion int           `json:"schema_version"`
	ID            string        `json:"id"`
	Version       string        `json:"version"`
	Languages     []string      `json:"languages"`
	Agents        []string      `json:"agents"`
	Skills        []skillRecord `json:"skills"`
	Packs         []packRecord  `json:"packs"`
}

type catalogIndex struct {
	File   catalogFile
	Skills map[string]skillRecord
	Packs  map[string]packRecord
}

func loadCatalog() (catalogIndex, error) {
	data, err := officialContent.ReadFile("registry/catalog.json")
	if err != nil {
		return catalogIndex{}, fmt.Errorf("read embedded catalog: %w", err)
	}
	var parsed catalogFile
	if err := json.Unmarshal(data, &parsed); err != nil {
		return catalogIndex{}, fmt.Errorf("decode embedded catalog: %w", err)
	}
	index := catalogIndex{File: parsed, Skills: map[string]skillRecord{}, Packs: map[string]packRecord{}}
	for _, skill := range parsed.Skills {
		if _, exists := index.Skills[skill.ID]; exists {
			return catalogIndex{}, fmt.Errorf("duplicate skill id %q", skill.ID)
		}
		index.Skills[skill.ID] = skill
	}
	for _, pack := range parsed.Packs {
		if _, exists := index.Packs[pack.ID]; exists {
			return catalogIndex{}, fmt.Errorf("duplicate pack id %q", pack.ID)
		}
		index.Packs[pack.ID] = pack
	}
	return index, nil
}

func (c catalogIndex) resolvePack(id string) ([]string, error) {
	seenPacks := map[string]bool{}
	visiting := map[string]bool{}
	seenSkills := map[string]bool{}
	result := []string{}
	var visit func(string) error
	visit = func(current string) error {
		if seenPacks[current] {
			return nil
		}
		if visiting[current] {
			return fmt.Errorf("cyclic pack inheritance at %s", current)
		}
		pack, ok := c.Packs[current]
		if !ok {
			return fmt.Errorf("unknown pack %s", current)
		}
		visiting[current] = true
		for _, parent := range pack.Extends {
			if err := visit(parent); err != nil {
				return err
			}
		}
		for _, skill := range pack.Skills {
			if _, ok := c.Skills[skill]; !ok {
				return fmt.Errorf("pack %s references unknown skill %s", current, skill)
			}
			if !seenSkills[skill] {
				seenSkills[skill] = true
				result = append(result, skill)
			}
		}
		delete(visiting, current)
		seenPacks[current] = true
		return nil
	}
	if err := visit(id); err != nil {
		return nil, err
	}
	return result, nil
}

func validateCatalog(c catalogIndex) []string {
	problems := []string{}
	if c.File.SchemaVersion != 1 {
		problems = append(problems, "catalog schema_version must be 1")
	}
	if !sameStrings(c.File.Languages, []string{"en", "zh-CN", "ja"}) {
		problems = append(problems, "catalog must support en, zh-CN, and ja")
	}
	for id, skill := range c.Skills {
		if !skillIDPattern.MatchString(id) || len(id) > 64 {
			problems = append(problems, fmt.Sprintf("invalid skill id %s", id))
		}
		if skill.Path != "skills/"+id {
			problems = append(problems, fmt.Sprintf("skill %s path must match its id", id))
		}
		if !sameStrings(skill.Languages, []string{"en", "zh-CN", "ja"}) {
			problems = append(problems, fmt.Sprintf("skill %s lacks required languages", id))
		}
		if len(skill.Capabilities) == 0 {
			problems = append(problems, fmt.Sprintf("skill %s lacks declared capabilities", id))
		}
		for _, language := range []string{"en", "zh-CN", "ja"} {
			review, ok := skill.Reviews[language]
			if !ok || (review.Status != "machine-drafted" && review.Status != "human-reviewed") || review.Revision != skill.Version {
				problems = append(problems, fmt.Sprintf("skill %s has invalid %s review metadata", id, language))
				continue
			}
			if skill.Maturity == "stable" && (review.Status != "human-reviewed" || review.ReviewerRole == "" || review.ReviewedAt == "") {
				problems = append(problems, fmt.Sprintf("stable skill %s lacks signed %s review", id, language))
			}
		}
		if skill.License == "" || skill.Source.Type == "" {
			problems = append(problems, fmt.Sprintf("skill %s lacks license or provenance", id))
		}
		skillFile := path.Join(skill.Path, "SKILL.md")
		data, err := officialContent.ReadFile(skillFile)
		if err != nil {
			problems = append(problems, fmt.Sprintf("skill %s lacks SKILL.md", id))
			continue
		}
		text := string(data)
		if !strings.HasPrefix(text, "---\nname: "+id+"\ndescription:") {
			problems = append(problems, fmt.Sprintf("skill %s has invalid frontmatter", id))
		}
		if strings.Count(text, "\n")+1 > 500 {
			problems = append(problems, fmt.Sprintf("skill %s exceeds 500 lines", id))
		}
		if _, err := officialContent.ReadFile(path.Join(skill.Path, "agents/openai.yaml")); err != nil {
			problems = append(problems, fmt.Sprintf("skill %s lacks Codex UI metadata", id))
		}
	}
	for id := range c.Packs {
		if _, err := c.resolvePack(id); err != nil {
			problems = append(problems, err.Error())
		}
	}
	problems = append(problems, validateSupportMatrix(c)...)
	problems = append(problems, validateClaimsLedger()...)
	var suite struct {
		SchemaVersion int      `json:"schema_version"`
		Languages     []string `json:"languages"`
		Difficulties  []string `json:"difficulties"`
		Cases         []struct {
			ID         string            `json:"id"`
			Skill      string            `json:"skill"`
			Difficulty string            `json:"difficulty"`
			Prompts    map[string]string `json:"prompts"`
			Routing    struct {
				ExpectedSkill string              `json:"expected_skill"`
				Neighbor      string              `json:"neighbor"`
				NotFor        []string            `json:"not_for"`
				HardNegative  bool                `json:"hard_negative"`
				Negative      map[string][]string `json:"negative_triggers"`
			} `json:"routing"`
			Assertions struct {
				Must    []string `json:"must"`
				MustNot []string `json:"must_not"`
			} `json:"assertions"`
		} `json:"cases"`
	}
	data, err := officialContent.ReadFile("evals/cases.json")
	if err != nil || json.Unmarshal(data, &suite) != nil {
		problems = append(problems, "invalid multilingual evaluation suite")
	} else {
		if suite.SchemaVersion != 2 || !sameStrings(suite.Difficulties, []string{"basic", "ambiguous", "conflicting"}) {
			problems = append(problems, "evaluation suite must use schema 2 and all three difficulties")
		}
		if len(suite.Cases) != len(c.Skills)*3 {
			problems = append(problems, fmt.Sprintf("evaluation suite has %d cases; expected %d", len(suite.Cases), len(c.Skills)*3))
		}
		covered := map[string]map[string]int{}
		caseIDs := map[string]bool{}
		for _, testCase := range suite.Cases {
			if caseIDs[testCase.ID] {
				problems = append(problems, fmt.Sprintf("duplicate evaluation id %s", testCase.ID))
			}
			caseIDs[testCase.ID] = true
			if covered[testCase.Skill] == nil {
				covered[testCase.Skill] = map[string]int{}
			}
			covered[testCase.Skill][testCase.Difficulty]++
			if len(testCase.Prompts) != 3 || testCase.Prompts["en"] == "" || testCase.Prompts["zh-CN"] == "" || testCase.Prompts["ja"] == "" {
				problems = append(problems, fmt.Sprintf("eval for %s is not trilingual", testCase.Skill))
			}
			if len(testCase.Assertions.Must) == 0 {
				problems = append(problems, fmt.Sprintf("eval for %s lacks assertions", testCase.Skill))
			}
			if testCase.Routing.ExpectedSkill != testCase.Skill {
				problems = append(problems, fmt.Sprintf("eval %s routes to the wrong skill", testCase.ID))
			}
			if _, ok := c.Skills[testCase.Routing.Neighbor]; !ok || testCase.Routing.Neighbor == testCase.Skill {
				problems = append(problems, fmt.Sprintf("eval %s has an invalid neighbor", testCase.ID))
			}
			expectedHardNegative := testCase.Difficulty != "basic"
			if testCase.Routing.HardNegative != expectedHardNegative {
				problems = append(problems, fmt.Sprintf("eval %s has an invalid hard-negative flag", testCase.ID))
			}
		}
		for id := range c.Skills {
			for _, difficulty := range []string{"basic", "ambiguous", "conflicting"} {
				if covered[id][difficulty] != 1 {
					problems = append(problems, fmt.Sprintf("skill %s lacks exactly one %s evaluation", id, difficulty))
				}
			}
		}
	}
	sort.Strings(problems)
	return problems
}

func validateSupportMatrix(c catalogIndex) []string {
	var matrix struct {
		SchemaVersion int `json:"schema_version"`
		Agents        []struct {
			ID             string `json:"id"`
			Status         string `json:"status"`
			TargetFamily   string `json:"target_family"`
			DefaultInstall *bool  `json:"default_install"`
		} `json:"agents"`
	}
	data, err := officialContent.ReadFile("registry/support-matrix.json")
	if err != nil || json.Unmarshal(data, &matrix) != nil || matrix.SchemaVersion != 1 {
		return []string{"invalid host support matrix"}
	}
	problems := []string{}
	stable := []string{}
	seen := map[string]bool{}
	for _, agent := range matrix.Agents {
		if seen[agent.ID] {
			problems = append(problems, fmt.Sprintf("duplicate support adapter %s", agent.ID))
		}
		seen[agent.ID] = true
		if agent.Status == "stable" {
			stable = append(stable, agent.ID)
		}
		if agent.ID == "cursor" && (agent.Status != "experimental" || agent.DefaultInstall == nil || *agent.DefaultInstall) {
			problems = append(problems, "Cursor must remain explicit and experimental until host smoke tests pass")
		}
		descriptorData, readErr := officialContent.ReadFile(path.Join("adapters", agent.ID+".json"))
		var descriptor map[string]any
		if readErr != nil || json.Unmarshal(descriptorData, &descriptor) != nil {
			problems = append(problems, fmt.Sprintf("invalid adapter descriptor %s", agent.ID))
			continue
		}
		if descriptor["id"] != agent.ID || descriptor["status"] != agent.Status || descriptor["target_family"] != agent.TargetFamily {
			problems = append(problems, fmt.Sprintf("adapter descriptor %s disagrees with support matrix", agent.ID))
		}
		if descriptor["install_executes_skill_code"] != false {
			problems = append(problems, fmt.Sprintf("adapter %s may execute skill code at install time", agent.ID))
		}
		if _, exists := descriptor["allowed_tools"]; exists {
			problems = append(problems, fmt.Sprintf("adapter %s pre-approves tools", agent.ID))
		}
	}
	if !sameStrings(stable, c.File.Agents) {
		problems = append(problems, "stable support matrix agents must match the catalog")
	}
	return problems
}

func validateClaimsLedger() []string {
	var ledger struct {
		SchemaVersion int `json:"schema_version"`
		Claims        []struct {
			ID       string   `json:"id"`
			Status   string   `json:"status"`
			Evidence []string `json:"evidence"`
		} `json:"claims"`
	}
	data, err := officialContent.ReadFile("evidence/claims.json")
	if err != nil || json.Unmarshal(data, &ledger) != nil || ledger.SchemaVersion != 1 {
		return []string{"invalid claims ledger"}
	}
	problems := []string{}
	seen := map[string]bool{}
	for _, claim := range ledger.Claims {
		if claim.ID == "" || seen[claim.ID] {
			problems = append(problems, fmt.Sprintf("invalid or duplicate claim %s", claim.ID))
		}
		seen[claim.ID] = true
		if claim.Status != "verified" && claim.Status != "partial" && claim.Status != "not-verified" {
			problems = append(problems, fmt.Sprintf("claim %s has invalid status", claim.ID))
		}
		if claim.Status == "verified" && len(claim.Evidence) == 0 {
			problems = append(problems, fmt.Sprintf("verified claim %s lacks evidence", claim.ID))
		}
	}
	return problems
}

func sameStrings(left, right []string) bool {
	if len(left) != len(right) {
		return false
	}
	l := append([]string(nil), left...)
	r := append([]string(nil), right...)
	sort.Strings(l)
	sort.Strings(r)
	for i := range l {
		if l[i] != r[i] {
			return false
		}
	}
	return true
}

func digestEmbeddedSkill(skill skillRecord, agent string) (string, error) {
	family := targetFamily(agent)
	hash := sha256.New()
	err := fs.WalkDir(officialContent, skill.Path, func(filePath string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			return nil
		}
		relative := strings.TrimPrefix(filePath, skill.Path+"/")
		if family == "portable-skills" && strings.HasPrefix(relative, "agents/") {
			return nil
		}
		data, err := officialContent.ReadFile(filePath)
		if err != nil {
			return err
		}
		hash.Write([]byte(relative))
		hash.Write([]byte{0})
		hash.Write(data)
		hash.Write([]byte{0})
		return nil
	})
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}
