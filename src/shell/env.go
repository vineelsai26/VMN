package shell

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"vineelsai.com/vmn/src/node"
	"vineelsai.com/vmn/src/utils"
)

func setEnvForPosixShell() {
	fmt.Println(`
export PATH="$(cat $HOME/.vmn/node-version):$PATH"

function vmn {
	$(whereis vmn | cut -d" " -f2) $@
	if [[ "$1" == "node" ]]
	then
		export PATH="$(cat $HOME/.vmn/node-version):$PATH"
	fi
}

function setNodeVersion {
	if [ -f .vmnrc ]
	then
		echo "Found .vmnrc file"
		if [ -f $HOME/.vmn/node/v$(ls "$HOME/.vmn/node" 2> /dev/null | grep "$(cat .vmnrc)" | tail -1 | cut -f2 -d"v")/bin/node ]
		then
			export PATH="$HOME/.vmn/node/v$(ls "$HOME/.vmn/node" 2> /dev/null | grep "$(cat .vmnrc)" | tail -1 | cut -f2 -d"v")/bin:$PATH"
		else
			vmn install $(cat .vmnrc)
			export PATH="$HOME/.vmn/node/v$(ls "$HOME/.vmn/node" 2> /dev/null | grep "$(cat .vmnrc)" | tail -1 | cut -f2 -d"v")/bin:$PATH"
		fi
		echo "Using node version $(node --version)"
	elif [ -f .nvmrc ]
	then
		echo "Found .nvmrc file"
		if [ -f $HOME/.vmn/node/v$(ls "$HOME/.vmn/node" 2> /dev/null | grep "$(cat .nvmrc)" | tail -1 | cut -f2 -d"v")/bin/node ]
		then
			export PATH="$HOME/.vmn/node/v$(ls "$HOME/.vmn/node" 2> /dev/null | grep "$(cat .nvmrc)" | tail -1 | cut -f2 -d"v")/bin:$PATH"
		else
			vmn install $(cat .nvmrc)
			export PATH="$HOME/.vmn/node/v$(ls "$HOME/.vmn/node" 2> /dev/null | grep "$(cat .nvmrc)" | tail -1 | cut -f2 -d"v")/bin:$PATH"
		fi
		echo "Using node version $(node --version)"
	elif [ -f .node-version ]
	then
		echo "Found .node-version file"
		if [ -f $HOME/.vmn/node/v$(ls "$HOME/.vmn/node" 2> /dev/null | grep "$(cat .node-version)" | tail -1 | cut -f2 -d"v")/bin/node ]
		then
			export PATH="$HOME/.vmn/node/v$(ls "$HOME/.vmn/node" 2> /dev/null | grep "$(cat .node-version)" | tail -1 | cut -f2 -d"v")/bin:$PATH"
		else
			vmn install $(cat .node-version)
			export PATH="$HOME/.vmn/node/v$(ls "$HOME/.vmn/node" 2> /dev/null | grep "$(cat .node-version)" | tail -1 | cut -f2 -d"v")/bin:$PATH"
		fi
		echo "Using node version $(node --version)"
	fi
}

function cd {
	builtin cd "$@"
	setNodeVersion
}

setNodeVersion
	`)
}

func setEnvForPowershell() {
	fmt.Println(`
function Get-VmnExecutable {
	$candidates = @(
		(Join-Path $HOME ".vmn\vmn.exe"),
		"vmn.exe",
		"vmn"
	)

	foreach ($candidate in $candidates) {
		if (Test-Path $candidate) {
			return $candidate
		}
	}

	return "vmn"
}

function Get-VmnManagedPath {
	$pathEntries = $Env:PATH -split ';' | Where-Object { $_ -and $_.Trim() -ne "" }
	$filtered = foreach ($entry in $pathEntries) {
		$trimmed = $entry.Trim()
		if ($trimmed -like (Join-Path $HOME ".vmn\node\*")) {
			continue
		}
		if ($trimmed -eq "%VMN_VERSION%") {
			continue
		}
		$trimmed
	}

	$seen = @{}
	foreach ($entry in $filtered) {
		$key = $entry.TrimEnd('\').ToLowerInvariant()
		if (-not $seen.ContainsKey($key)) {
			$seen[$key] = $true
			$entry
		}
	}
}

function Find-VmnInstalledVersionPath([string]$versionSpec) {
	$nodeRoot = Join-Path $HOME ".vmn\node"
	if (-not (Test-Path $nodeRoot)) {
		return $null
	}

	$match = Get-ChildItem $nodeRoot -Directory -ErrorAction SilentlyContinue |
		Where-Object { $_.Name -like ("v" + $versionSpec.Trim() + "*") } |
		Sort-Object Name |
		Select-Object -Last 1

	if ($match) {
		return $match.FullName
	}

	return $null
}

function Resolve-VmnVersionFromProject {
	$versionFiles = @(".vmnrc", ".nvmrc", ".node-version")
	$nodeVersionFile = Join-Path $HOME ".vmn\node-version"

	foreach ($versionFile in $versionFiles) {
		if (-not (Test-Path $versionFile)) {
			continue
		}

		$requestedVersion = (Get-Content $versionFile -Raw).Trim()
		if (-not $requestedVersion) {
			continue
		}

		Write-Output ("Found " + $versionFile)
		$resolvedPath = Find-VmnInstalledVersionPath $requestedVersion
		if (-not $resolvedPath) {
			& (Get-VmnExecutable) use $requestedVersion | Out-Host
			if (Test-Path $nodeVersionFile) {
				$resolvedPath = (Get-Content $nodeVersionFile -Raw).Trim()
			}
		}

		if ($resolvedPath) {
			return $resolvedPath
		}
	}

	return $null
}

function Update-VmnPath {
	$pathEntries = @(Get-VmnManagedPath)
	$selectedPath = Resolve-VmnVersionFromProject

	if (-not $selectedPath) {
		$versionFile = Join-Path $HOME ".vmn\node-version"
		if (Test-Path $versionFile) {
			$selectedPath = (Get-Content $versionFile -Raw).Trim()
		}
	}

	if ($selectedPath) {
		$pathEntries = @($selectedPath) + $pathEntries
	}

	$Env:PATH = ($pathEntries -join ';')
}

function vmn {
	& (Get-VmnExecutable) @Args
	if ($LASTEXITCODE -eq 0) {
		Update-VmnPath
	}
}

function Set-VmnLocation {
	param([Parameter(ValueFromRemainingArguments = $true)] $PathArgs)
	Microsoft.PowerShell.Management\Set-Location @PathArgs
	Update-VmnPath
}

Set-Alias -Name cd -Value Set-VmnLocation -Option AllScope
Update-VmnPath
	`)
}

func PrintEnv() {
	if runtime.GOOS == "linux" {
		fmt.Println("eval \"SHELL=`ps -p $$ -o comm=`; `vmn env $SHELL`\"")
	} else if runtime.GOOS == "darwin" {
		fmt.Println("eval \"`vmn env zsh`\"")
	} else if runtime.GOOS == "windows" {
		fmt.Println("vmn env powershell | Out-String | Invoke-Expression")
	}
}

func RunShellSpecificCommands(args []string) {
	if _, err := os.Stat(filepath.Join(utils.GetHome(), ".vmn")); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Join(utils.GetHome(), ".vmn"), 0755); err != nil {
			panic(err)
		}
	}

	if _, err := os.Stat(filepath.Join(utils.GetHome(), ".vmn", "node-version")); os.IsNotExist(err) {
		f, err := os.OpenFile(filepath.Join(utils.GetHome(), ".vmn", "node-version"), os.O_RDWR|os.O_CREATE, 0755)
		if err != nil {
			panic(err)
		}

		defer f.Close()
		installedNodeVersions := node.GetInstalledVersions()
		if len(installedNodeVersions) > 0 {
			if _, err := f.Stat(); err == nil {
				f.Truncate(0)
				f.Seek(0, 0)
				f.WriteString(utils.GetDestination(installedNodeVersions[0], "node"))
			}
		}

	}

	if runtime.GOOS == "windows" && len(args) > 1 && args[1] == "powershell" {
		setEnvForPowershell()
		return
	}

	setEnvForPosixShell()
}
