# VMN

VMN is a simple tool for managing Node.js versions. It's written in Go and is available for Linux, macOS and Windows.

## Installation

### Linux

```bash
wget https://github.com/vineelsai26/VMN/releases/download/v0.0.2/vmn-linux-amd64.tar.gz -O vmn-linux-amd64.tar.gz
tar -xvf vmn-linux-amd64.tar.gz
sudo mv vmn /usr/local/bin
```

### macOS

```bash
wget https://github.com/vineelsai26/VMN/releases/download/v0.0.2/vmn-darwin-amd64.tar.gz -O vmn-linux-amd64.tar.gz
tar -xvf vmn-linux-amd64.tar.gz
sudo mv vmn /usr/local/bin
```

### Windows

- Download the [latest release](https://github.com/vineelsai26/VMN/releases/latest) for Windows

- Extract the zip file

- Add the extracted folder to your PATH

## Usage

### Install a Node.js version

```bash
vmn install lts
```

### Use a Node.js version

```bash
vmn use lts
```

### List installed Node.js versions

```bash
vmn list installed
```

### Remove a Node.js version

```bash
vmn uninstall lts
```

### See all available commands

```bash
vmn help
```
