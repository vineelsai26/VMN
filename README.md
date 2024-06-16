# VMN

VMN is a simple tool for managing Node.js versions. It's written in Go and is available for Linux and macOS.

## Linux

### Manual Install

```bash
wget https://github.com/vineelsai26/VMN/releases/latest/download/vmn-linux-amd64.tar.gz -O vmn-linux-amd64.tar.gz
tar -xvf vmn-linux-amd64.tar.gz
sudo mv vmn /usr/local/bin
```

### Arch Linux Install

Add the following to `/etc/pacman.conf`:

```bash
[vineelsai-arch-repo]
Server = https://repo.vineelsai.com/linux/arch/$arch
```

Then run:

```bash
sudo pacman-key --lsign-key 4431E64723B4ADDE
sudo pacman -Syu vmn
```

### Debian/Ubuntu Install

```bash
curl -fsSL https://repo.vineelsai.com/linux/debian/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/vineelsai.gpg

echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/vineelsai.gpg] https://repo.vineelsai.com/linux/debian stable main" | sudo tee /etc/apt/sources.list.d/vineelsai.list > /dev/null

sudo apt update

sudo apt install vmn
```

## macOS

```bash
wget https://github.com/vineelsai26/VMN/releases/latest/download/vmn-macos-arm64.tar.gz -O vmn-macos-arm64.tar.gz
tar -xvf vmn-macos-arm64.tar.gz
sudo mv vmn /usr/local/bin
```

## Usage Node.js

### Install a Node.js version

```bash
vmn node install 20
```

### Use a Node.js version

```bash
vmn node use 20
```

### List installed Node.js versions

```bash
vmn node list installed
```

### Remove a Node.js version

```bash
vmn node uninstall 20
```

## Usage Python

### Install a Python version

```bash
vmn python install 3.11
```

### Use a Python version

```bash
vmn python use 3.11
```

### List installed Python versions

```bash
vmn python list installed
```

### Remove a Python version

```bash
vmn python uninstall 3.11
```
