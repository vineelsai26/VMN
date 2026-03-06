//go:build windows

package setup

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/sys/windows/registry"
	"vineelsai.com/vmn/src/utils"
)

func Install() {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	vmnHome := filepath.Join(utils.GetHome(), ".vmn")
	if err := os.MkdirAll(vmnHome, 0755); err != nil {
		panic(err)
	}

	srcFile, err := os.Open(filepath.Join(dir, "vmn.exe"))
	if err != nil {
		panic(err)
	}
	defer srcFile.Close()

	destPath := filepath.Join(vmnHome, "vmn.exe")
	destFile, err := os.Create(destPath)
	if err != nil {
		panic(err)
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, srcFile); err != nil {
		panic(err)
	}

	if err := ensureUserPathEntry(vmnHome); err != nil {
		panic(err)
	}

	if err := ensurePowerShellProfileHook(); err != nil {
		panic(err)
	}

	fmt.Println("VMN installed successfully!")
}

func SetPath(path string) {
	if _, err := os.Stat(path); err != nil {
		panic(err)
	}

	if err := os.Setenv("VMN_VERSION", path); err != nil {
		panic(err)
	}
}

func ensureUserPathEntry(entry string) error {
	key, err := registry.OpenKey(registry.CURRENT_USER, `Environment`, registry.QUERY_VALUE|registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer key.Close()

	current, _, err := key.GetStringValue("Path")
	if err != nil && err != registry.ErrNotExist {
		return err
	}

	parts := splitAndDedupPath(current)
	if containsPath(parts, entry) {
		return nil
	}

	parts = append(parts, entry)
	return key.SetStringValue("Path", strings.Join(parts, ";"))
}

func ensurePowerShellProfileHook() error {
	hookLine := "vmn env powershell | Out-String | Invoke-Expression"
	paths := []string{
		filepath.Join(utils.GetHome(), "Documents", "WindowsPowerShell", "Microsoft.PowerShell_profile.ps1"),
		filepath.Join(utils.GetHome(), "Documents", "PowerShell", "Microsoft.PowerShell_profile.ps1"),
	}

	for _, profilePath := range paths {
		if err := os.MkdirAll(filepath.Dir(profilePath), 0755); err != nil {
			return err
		}

		content, err := os.ReadFile(profilePath)
		if err != nil && !os.IsNotExist(err) {
			return err
		}

		if strings.Contains(string(content), hookLine) {
			continue
		}

		appendContent := hookLine + "\r\n"
		if len(content) > 0 && !strings.HasSuffix(string(content), "\n") {
			appendContent = "\r\n" + appendContent
		}

		file, err := os.OpenFile(profilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}

		if _, err := file.WriteString(appendContent); err != nil {
			file.Close()
			return err
		}
		file.Close()
	}

	return nil
}

func splitAndDedupPath(pathValue string) []string {
	parts := strings.Split(pathValue, ";")
	result := make([]string, 0, len(parts))
	seen := map[string]struct{}{}

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		normalized := strings.ToLower(strings.TrimRight(part, `\`))
		if _, ok := seen[normalized]; ok {
			continue
		}

		seen[normalized] = struct{}{}
		result = append(result, part)
	}

	return result
}

func containsPath(parts []string, target string) bool {
	normalizedTarget := strings.ToLower(strings.TrimRight(target, `\`))
	for _, part := range parts {
		if strings.ToLower(strings.TrimRight(part, `\`)) == normalizedTarget {
			return true
		}
	}

	return false
}
