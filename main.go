package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var nativeVersion = "0.3.0-alpha.1"

type mutationOptions struct {
	Pack             string
	Agents           []string
	Scope            string
	ProjectDirectory string
	DryRun           bool
	Force            bool
	Help             bool
}

type recoveryOptions struct {
	Scope            string
	ProjectDirectory string
	Help             bool
}

func main() {
	os.Exit(runNative(os.Args[1:]))
}

func runNative(arguments []string) int {
	language, cleaned, err := extractLanguage(arguments)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 2
	}
	if len(cleaned) == 0 || cleaned[0] == "--help" || cleaned[0] == "-h" {
		fmt.Println(msg(language, "usage"))
		return 0
	}
	if cleaned[0] == "--version" || cleaned[0] == "version" {
		fmt.Printf("anywork %s\n", nativeVersion)
		return 0
	}
	catalog, err := loadCatalog()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	switch cleaned[0] {
	case "list":
		fmt.Println(msg(language, "catalog_title"))
		fmt.Println(msg(language, "quality_notice"))
		for _, pack := range catalog.File.Packs {
			skills, err := catalog.resolvePack(pack.ID)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				return 1
			}
			switch language {
			case "zh-CN":
				fmt.Printf("%s：%s（%d 个技能）\n", localize(pack.DisplayName, language), localize(pack.Description, language), len(skills))
			case "ja":
				fmt.Printf("%s：%s（%d スキル）\n", localize(pack.DisplayName, language), localize(pack.Description, language), len(skills))
			default:
				fmt.Printf("%s: %s (%d skills)\n", localize(pack.DisplayName, language), localize(pack.Description, language), len(skills))
			}
		}
		return 0
	case "demo":
		if len(cleaned) > 2 {
			fmt.Fprintln(os.Stderr, demoUsage(language))
			return 2
		}
		demoID := ""
		if len(cleaned) == 2 {
			demoID = cleaned[1]
		}
		output, err := renderDemo(language, demoID)
		if err != nil {
			fmt.Fprintln(os.Stderr, msg(language, "unknown_demo", demoID))
			return 2
		}
		fmt.Println(output)
		return 0
	case "doctor":
		problems := validateCatalog(catalog)
		if len(problems) > 0 {
			fmt.Println(msg(language, "doctor_fail", len(problems)))
			for _, problem := range problems {
				fmt.Printf("- %s\n", problem)
			}
			return 1
		}
		fmt.Println(msg(language, "doctor_ok", len(catalog.Skills), len(catalog.Packs), len(catalog.File.Agents)))
		return 0
	case "recover":
		options, err := parseRecoveryOptions(cleaned[1:])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 2
		}
		if options.Help {
			fmt.Println(recoveryUsage(language))
			return 0
		}
		emit := func(key string, values ...any) { fmt.Println(msg(language, key, values...)) }
		if err := recoverNative(options.Scope, options.ProjectDirectory, emit); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 2
		}
		return 0
	case "setup", "install", "update", "uninstall":
		optionArguments := cleaned[1:]
		if cleaned[0] == "setup" {
			optionArguments = append([]string{"complete"}, optionArguments...)
		}
		options, err := parseMutationOptions(optionArguments)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 2
		}
		if options.Help {
			if cleaned[0] == "setup" {
				fmt.Printf("%s\n", setupUsage(language))
			} else {
				fmt.Printf("%s\n", mutationUsage(language, cleaned[0]))
			}
			return 0
		}
		skillIDs, err := catalog.resolvePack(options.Pack)
		if err != nil {
			fmt.Fprintln(os.Stderr, msg(language, "unknown_pack", options.Pack))
			return 2
		}
		emit := func(key string, values ...any) { fmt.Println(msg(language, key, values...)) }
		if cleaned[0] == "setup" && !options.DryRun {
			emit = func(key string, values ...any) {
				if key == "backup" {
					fmt.Println(msg(language, key, values...))
				}
			}
		}
		if cleaned[0] == "uninstall" {
			err = uninstallNative(skillIDs, options.Agents, options.Scope, options.ProjectDirectory, options.DryRun, options.Force, emit)
		} else {
			err = installNative(catalog, skillIDs, options.Agents, options.Scope, options.ProjectDirectory, options.DryRun, options.Force, emit)
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 2
		}
		if cleaned[0] == "setup" && !options.DryRun {
			fmt.Println(msg(language, "setup_complete", len(skillIDs), len(options.Agents)))
			fmt.Println(msg(language, "setup_next"))
		}
		return 0
	default:
		fmt.Fprintln(os.Stderr, msg(language, "usage"))
		return 2
	}
}

func setupUsage(language string) string {
	agentValues := "all|" + strings.Join(supportedAgents, "|")
	switch language {
	case "zh-CN":
		return fmt.Sprintf("用法：anywork setup [--agent %s] [--scope user|project] [--project-dir 路径] [--dry-run] [--force]", agentValues)
	case "ja":
		return fmt.Sprintf("使い方：anywork setup [--agent %s] [--scope user|project] [--project-dir パス] [--dry-run] [--force]", agentValues)
	default:
		return fmt.Sprintf("Usage: anywork setup [--agent %s] [--scope user|project] [--project-dir path] [--dry-run] [--force]", agentValues)
	}
}

func demoUsage(language string) string {
	switch language {
	case "zh-CN":
		return "用法：anywork demo [research|meeting|automation]"
	case "ja":
		return "使い方：anywork demo [research|meeting|automation]"
	default:
		return "Usage: anywork demo [research|meeting|automation]"
	}
}

func parseRecoveryOptions(arguments []string) (recoveryOptions, error) {
	workingDirectory, err := os.Getwd()
	if err != nil {
		return recoveryOptions{}, err
	}
	result := recoveryOptions{Scope: "user", ProjectDirectory: workingDirectory}
	for index := 0; index < len(arguments); index++ {
		current := arguments[index]
		nextValue := func() (string, error) {
			if index+1 >= len(arguments) {
				return "", fmt.Errorf("%s requires a value", current)
			}
			index++
			return arguments[index], nil
		}
		switch {
		case current == "--help" || current == "-h":
			result.Help = true
		case current == "--scope":
			value, err := nextValue()
			if err != nil {
				return result, err
			}
			result.Scope = value
		case strings.HasPrefix(current, "--scope="):
			result.Scope = strings.TrimPrefix(current, "--scope=")
		case current == "--project-dir":
			value, err := nextValue()
			if err != nil {
				return result, err
			}
			result.ProjectDirectory = value
		case strings.HasPrefix(current, "--project-dir="):
			result.ProjectDirectory = strings.TrimPrefix(current, "--project-dir=")
		default:
			return result, fmt.Errorf("unknown recover option: %s", current)
		}
	}
	if result.Scope != "user" && result.Scope != "project" {
		return result, fmt.Errorf("scope must be user or project")
	}
	if result.ProjectDirectory, err = filepath.Abs(result.ProjectDirectory); err != nil {
		return result, err
	}
	return result, nil
}

func parseMutationOptions(arguments []string) (mutationOptions, error) {
	workingDirectory, err := os.Getwd()
	if err != nil {
		return mutationOptions{}, err
	}
	result := mutationOptions{Scope: "user", ProjectDirectory: workingDirectory}
	for index := 0; index < len(arguments); index++ {
		current := arguments[index]
		nextValue := func() (string, error) {
			if index+1 >= len(arguments) {
				return "", fmt.Errorf("%s requires a value", current)
			}
			index++
			return arguments[index], nil
		}
		switch {
		case current == "--help" || current == "-h":
			result.Help = true
		case current == "--dry-run":
			result.DryRun = true
		case current == "--force":
			result.Force = true
		case current == "--agent":
			value, err := nextValue()
			if err != nil {
				return result, err
			}
			result.Agents = append(result.Agents, value)
		case strings.HasPrefix(current, "--agent="):
			result.Agents = append(result.Agents, strings.TrimPrefix(current, "--agent="))
		case current == "--scope":
			value, err := nextValue()
			if err != nil {
				return result, err
			}
			result.Scope = value
		case strings.HasPrefix(current, "--scope="):
			result.Scope = strings.TrimPrefix(current, "--scope=")
		case current == "--project-dir":
			value, err := nextValue()
			if err != nil {
				return result, err
			}
			result.ProjectDirectory = value
		case strings.HasPrefix(current, "--project-dir="):
			result.ProjectDirectory = strings.TrimPrefix(current, "--project-dir=")
		case strings.HasPrefix(current, "-"):
			return result, fmt.Errorf("unknown option: %s", current)
		case result.Pack == "":
			result.Pack = current
		default:
			return result, fmt.Errorf("unexpected argument: %s", current)
		}
	}
	if result.Help {
		return result, nil
	}
	if result.Pack == "" {
		return result, fmt.Errorf("pack id is required")
	}
	if result.Scope != "user" && result.Scope != "project" {
		return result, fmt.Errorf("scope must be user or project")
	}
	if result.ProjectDirectory, err = filepath.Abs(result.ProjectDirectory); err != nil {
		return result, err
	}
	if len(result.Agents) == 0 {
		result.Agents = append([]string(nil), supportedAgents...)
	} else {
		normalized := []string{}
		seen := map[string]bool{}
		for _, raw := range result.Agents {
			if raw == "all" {
				normalized = append([]string(nil), supportedAgents...)
				for _, agent := range normalized {
					seen[agent] = true
				}
				break
			}
			agent := normalizeAgent(raw)
			if !isSupportedAgent(agent) {
				return result, fmt.Errorf("unsupported agent: %s", raw)
			}
			if !seen[agent] {
				seen[agent] = true
				normalized = append(normalized, agent)
			}
		}
		result.Agents = normalized
	}
	return result, nil
}

func mutationUsage(language, command string) string {
	agentValues := "all|" + strings.Join(supportedAgents, "|")
	switch language {
	case "zh-CN":
		return fmt.Sprintf("用法：anywork %s <工作包> [--agent %s] [--scope user|project] [--project-dir 路径] [--dry-run] [--force]", command, agentValues)
	case "ja":
		return fmt.Sprintf("使い方：anywork %s <パック> [--agent %s] [--scope user|project] [--project-dir パス] [--dry-run] [--force]", command, agentValues)
	default:
		return fmt.Sprintf("Usage: anywork %s <pack> [--agent %s] [--scope user|project] [--project-dir path] [--dry-run] [--force]", command, agentValues)
	}
}

func recoveryUsage(language string) string {
	switch language {
	case "zh-CN":
		return "用法：anywork recover [--scope user|project] [--project-dir 路径]"
	case "ja":
		return "使い方：anywork recover [--scope user|project] [--project-dir パス]"
	default:
		return "Usage: anywork recover [--scope user|project] [--project-dir path]"
	}
}
