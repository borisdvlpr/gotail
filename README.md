# gotail

[![Release](https://img.shields.io/github/v/release/borisdvlpr/gotail)](https://github.com/borisdvlpr/gotail/releases)
[![CI](https://github.com/borisdvlpr/gotail/actions/workflows/pull-request.yaml/badge.svg)](https://github.com/borisdvlpr/gotail/actions/workflows/pull-request.yaml)
[![Codecov](https://codecov.io/github/borisdvlpr/gotail/graph/badge.svg?token=QJTVO5NX58)](https://codecov.io/github/borisdvlpr/gotail)
[![License](https://img.shields.io/github/license/borisdvlpr/gotail)](LICENSE)
[![OSS hosting by Cloudsmith](https://img.shields.io/badge/OSS%20hosting%20by-cloudsmith-blue?logo=cloudsmith)](https://cloudsmith.com)

Bootstrap Tailscale into your Raspberry Pi and join it to your tailnet automatically from the very first boot.

Available for Linux, macOS, and Windows, gotail is meant to be run after flashing an SD card but before the device's first boot. It's a Go implementation of the original Tailscale grafter for Raspberry Pi, [tailgraft](https://github.com/tailscale-dev/tailgraft/blob/main/README.md).

## Table of Contents

- [Installation](#installation)
  - [Package Managers](#package-managers-recommended)
  - [Pre-built Binary](#pre-built-binary)
  - [Build from Source](#build-from-source)
- [Usage](#usage)
- [Development](#development)
- [Contributing](#contributing)
- [License](#license)

## Installation

### Package Managers (recommended)

#### Homebrew (macOS / Linux)

```sh
brew install borisdvlpr/gotail/gotail
```

#### Debian / Ubuntu (APT)

```sh
curl -1sLf 'https://dl.cloudsmith.io/public/borisdvlpr/gotail/setup.deb.sh' | sudo -E bash
sudo apt install gotail
```

#### Fedora / RHEL (DNF)

```sh
curl -1sLf 'https://dl.cloudsmith.io/public/borisdvlpr/gotail/setup.rpm.sh' | sudo -E bash
sudo dnf install gotail
```

#### Alpine (APK)

```sh
curl -1sLf 'https://dl.cloudsmith.io/public/borisdvlpr/gotail/setup.alpine.sh' | sudo -E bash
sudo apk add gotail
```

#### Windows (winget)

```powershell
winget install borisdvlpr.gotail
```

> Linux package repository hosting is graciously provided by [Cloudsmith](https://cloudsmith.com). Cloudsmith is the only fully hosted, cloud-native, universal 
package management solution, that enables your organization to create, store and share packages in any format, to any place, with total confidence.

### Pre-built Binary

Download the latest release from the [GitHub releases page](https://github.com/borisdvlpr/gotail/releases) and follow the steps for your platform.

#### Linux / macOS

1. Make the binary executable:

   ```sh
   chmod +x gotail
   ```

2. Optionally, move it to your PATH:

   ```sh
   sudo mv gotail /usr/local/bin/
   ```

#### Windows

Optionally, move the binary to a folder in your PATH, or add its location to your `PATH` environment variable via **System Properties → Environment Variables**.

### Build from Source

Building from source requires Go and Task (see [Development](#development)).

#### Using Task

```sh
git clone https://github.com/borisdvlpr/gotail.git
cd gotail
task all
```

#### Using Go directly

```sh
git clone https://github.com/borisdvlpr/gotail.git
cd gotail
go build -o gotail main.go
```

Then move the binary to your PATH as described above.

You can confirm the installation with:

```sh
gotail --version
```

## Usage

### Interactive Setup

Run the setup command to configure a new device interactively.

**Linux / macOS:**

```sh
sudo gotail setup
```

**Windows** (run as Administrator):

```powershell
.\gotail.exe setup
```

**Using Go (development):**

```sh
go run main.go setup
```

Follow the prompts to configure your Tailscale settings.

### Configuration File Setup

Create a YAML configuration file with your settings:

```yaml
exit_node: n                            # 'y' to offer this device as an exit node
subnet_router: y                        # 'y' to advertise the subnet(s) below
subnets: "192.0.2.1/24,192.168.2.1/24"  # comma-separated subnet(s) to advertise (required if subnet_router is 'y')
hostname: raspberrypi                   # hostname for your device
tags: tag-1,tag-2                       # comma-separated ACL tags to apply (must already exist under tagOwners in your ACL policy)
auth_key: tskey-auth-xxxxxxxxxxxxxxxxx  # your Tailscale auth key
```

> Keep your auth key secret — treat the configuration file like any other credential and avoid committing it to version control.

Then run `gotail` with the configuration file.

**Linux / macOS:**

```sh
sudo gotail setup --file /path/to/config.yaml
```

**Windows** (run as Administrator):

```powershell
.\gotail.exe setup --file C:\path\to\config.yaml
```

**Using Go (development):**

```sh
go run main.go setup --file /path/to/config.yaml
```

### Getting Your Auth Key

You'll need a Tailscale auth key, which you can generate from your [Tailscale admin console](https://login.tailscale.com/admin/settings/keys).

### After Setup

When your Raspberry Pi boots up, you should see it in your admin console's [machines](https://login.tailscale.com/admin/machines) page and you should be able to use [Tailscale SSH](https://tailscale.com/tailscale-ssh/) to connect to it:

```sh
tailscale ssh ubuntu@<hostname>
```

Depending on your ACL configuration, you may be prompted to authenticate with Tailscale.

## Development

For development and source builds, the following tools are needed:

- **Go**: Version 1.26 or later. Download from [golang.org](https://go.dev/dl/).
- **Task**: A task runner / build tool. Installation instructions and documentation can be found at [taskfile.dev](https://taskfile.dev/).

Common tasks:

```sh
task test     # run the unit tests
task build    # build the binary for your current platform
task install  # build and install into your user bin directory
task clean    # remove build artifacts and the installed binary
```

## Contributing

Feel free to open any issues or pull requests!

## License

gotail is released under the [BSD-3-Clause License](LICENSE).
