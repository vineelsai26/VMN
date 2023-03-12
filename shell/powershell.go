package shell

import "fmt"

func SetEnvForPowershell() {
	fmt.Println(`
function changedir($argList) {
    Set-Location $argList
    if (Test-Path .vmnrc) {
        Write-Output "Found .vmnrc file"
        Set-Item -Path Env:PATH -Value "C:\Users\Vineel\.vmn\node\$(Get-Content .vmnrc);$Env:PATH"
        Write-Output "Using node version $(node --version)"
    }

	if (Test-Path .nvmrc) {
		Write-Output "Found .nvmrc file"
		Set-Item -Path Env:PATH -Value "C:\Users\Vineel\.vmn\node\$(Get-Content .nvmrc);$Env:PATH"
		Write-Output "Using node version $(node --version)"
	}

	if (Test-Path .node-version) {
		Write-Output "Found .node-version file"
		Set-Item -Path Env:PATH -Value "C:\Users\Vineel\.vmn\node\$(Get-Content .node-version);$Env:PATH"
		Write-Output "Using node version $(node --version)"
	}
}

if (Test-Path .vmnrc) {
    Write-Output "Found .vmnrc file"
    Set-Item -Path Env:PATH -Value "C:\Users\Vineel\.vmn\node\$(Get-Content .vmnrc);$Env:PATH"
    Write-Output "Using node version $(node --version)"
}

if (Test-Path .nvmrc) {
    Write-Output "Found .nvmrc file"
    Set-Item -Path Env:PATH -Value "C:\Users\Vineel\.vmn\node\$(Get-Content .nvmrc);$Env:PATH"
    Write-Output "Using node version $(node --version)"
}

if (Test-Path .node-version) {
    Write-Output "Found .node-version file"
    Set-Item -Path Env:PATH -Value "C:\Users\Vineel\.vmn\node\$(Get-Content .node-version);$Env:PATH"
    Write-Output "Using node version $(node --version)"
}

Set-Alias -Name cd -Option AllScope -Value changedir
	`)
}
