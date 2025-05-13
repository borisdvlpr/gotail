# gotail

Bootstrap Tailscale into your Raspberry Pi and join it to your tailnet automatically from the very first boot. This project is a `go` implementation of the original Tailscale grafter for Raspberry Pi, [tailgraft](https://github.com/tailscale-dev/tailgraft/blob/main/README.md).

## Get Started

`gotail` is intended to be used after you've flashed Ubuntu onto an SD card, but before you've booted it in a Raspberry Pi for the first time. It assumes you are using Linux or macOS. Choose one of these methods to use `gotail`:

### Using the pre-built binary (recommended)

- Download the latest release from the [GitHub releases page](https://github.com/borisdvlpr/gotail/releases)
- Make the binary executable: `chmod +x gotail`
- Run it with superuser permissions: `sudo ./gotail`

### Building and installing with Make

- Clone this repository: `git clone https://github.com/borisdvlpr/gotail.git`
- Navigate to the directory: `cd gotail`
- Install using Make: `make`
- Run it: `sudo gotail`

### Running directly with Go

- Clone this repository: `git clone https://github.com/borisdvlpr/gotail.git`
- Navigate to the directory: `cd gotail`
- Run with Go: `sudo go run main.go`

## Configuring a new device

After installing `gotail`, you have two options to configure a new device:

### Interactive Setup

Run the setup command to configure a new device interactively:

```sh
sudo gotail setup
```

If you're running directly with Go:

```sh
sudo go run main.go setup
```

Once you run the setup command, answer the prompts to configure Tailscale.

### Configuration File

Alternatively, you can create a YAML configuration file with your settings:

```yaml
exit_node: n          # 'y' to enable exit node functionality
subnet_router: n      # 'y' to enable subnet router functionality
subnets: ""           # comma-separated list of subnets (required if subnet_router is 'y')
hostname: raspberrypi # hostname for your device
auth_key: tskey_1234  # your Tailscale auth key
```

Save this file and provide its path when running `gotail`:

```sh
sudo gotail setup --config /path/to/config.yaml
```

When using gotail with a configuration file or on interactivr mode, an auth key will be required, which you can generate from your [Tailscale admin console](https://login.tailscale.com/admin/settings/keys).

When your Raspberry Pi boots up, you should see it in your admin console's [machines](https://login.tailscale.com/admin/machines) page and you should be able to use [Tailscale SSH](https://tailscale.com/tailscale-ssh/) to connect to it:

```sh
tailscale ssh ubuntu@<hostname>
```

Depending on your ACL configuration, you may be prompted to authenticate with Tailscale.

## Contributing

Feel free to open any issues or pull requests!
