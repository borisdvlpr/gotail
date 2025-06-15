# gotail

Bootstrap Tailscale into your Raspberry Pi and join it to your tailnet automatically from the very first boot. This project is a Go implementation of the original Tailscale grafter for Raspberry Pi, [tailgraft](https://github.com/tailscale-dev/tailgraft/blob/main/README.md).

gotail is intended to be used after you've flashed Ubuntu onto an SD card, but before you've booted it in a Raspberry Pi for the first time. It assumes you are using Linux or macOS.

## Development Requirements

For development and source builds, the following tools are needed:

- **Go**: Version 1.24.0 or later. Download from [golang.org](https://go.dev/dl/).
- **Task**: A task runner / build tool. Installation instructions and documentation can be found at [taskfile.dev](https://taskfile.dev/).

## Installation

Choose one of these methods to install `gotail`:

### Option 1: Pre-built Binary (Recommended)

1. Download the latest release from the [GitHub releases page](https://github.com/borisdvlpr/gotail/releases)
2. Make the binary executable:

   ```sh
   chmod +x gotail
   ```
3. Optionally, move it to your PATH:

   ```sh
   sudo mv gotail /usr/local/bin/
   ```

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
sudo mv gotail /usr/local/bin/
```

## Usage

### Interactive Setup

Run the setup command to configure a new device interactively:

**Using installed binary:**
```sh
sudo gotail setup
```

**Using Go (development):**
```sh
sudo go run main.go setup
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

Then run `gotail` with the configuration file:

**Using installed binary:**
```sh
sudo gotail setup --file /path/to/config.yaml
```

**Using Go (development):**
```sh
sudo go run main.go setup --file /path/to/config.yaml
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