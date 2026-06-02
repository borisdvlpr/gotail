# gotail

Bootstrap Tailscale into your Raspberry Pi and join it to your tailnet automatically from the very first boot. This project is a Go implementation of the original Tailscale grafter for Raspberry Pi, [tailgraft](https://github.com/tailscale-dev/tailgraft/blob/main/README.md).

gotail is intended to be used after you've flashed Ubuntu onto an SD card, but before you've booted it in a Raspberry Pi for the first time. It supports Linux, macOS, and Windows.

## Development Requirements

For development and source builds, the following tools are needed:

- **Go**: Version 1.26 or later. Download from [golang.org](https://go.dev/dl/).
- **Task**: A task runner / build tool. Installation instructions and documentation can be found at [taskfile.dev](https://taskfile.dev/).

## Installation

Choose one of these methods to install `gotail`:

### Option 1: Pre-built Binary (Recommended)

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

### Option 2: Build from Source

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
exit_node: n          # 'y' to enable exit node functionality
subnet_router: n      # 'y' to enable subnet router functionality
subnets: ""           # comma-separated list of subnets (required if subnet_router is 'y')
hostname: raspberrypi # hostname for your device
auth_key: tskey_1234  # your Tailscale auth key
```

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

## Contributing

Feel free to open any issues or pull requests!
