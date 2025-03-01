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

After installing `gotail`, run the setup command to configure a new device:

```sh
sudo gotail setup
```

If you're running directly with Go:

```sh
sudo go run main.go setup
```

Once you run the setup command, answer the prompts to configure Tailscale. One of the prompts will request an auth key, which you can generate from your [Tailscale admin console](https://login.tailscale.com/admin/settings/keys).

When your Rasbperry Pi boots up, you should see it in your admin console's [machines](https://login.tailscale.com/admin/machines) page and you should be able to use to [Tailscale SSH](https://tailscale.com/tailscale-ssh/) to connect to it

```sh
tailscale ssh ubuntu@<hostname>
```

Depending on your ACL configuration, you may be prompted to authenticate with Tailscale.

## Contributing

Feel free to open any issues or pull requests!
