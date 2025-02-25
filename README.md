# gotail

Bootstrap Tailscale into your Raspberry Pi and join it to your tailnet automatically from the very first boot. This project is a `go` implementation of the original Tailscale grafter for Raspberry Pi, [tailgraft](https://github.com/tailscale-dev/tailgraft/blob/main/README.md).

## Get Started

This script is intended to be used after you've flashed Ubuntu onto an SD card, but before you've booted it in a Raspberry Pi for the first time. It assumes you are using Linux or macOS and have `go` installed but otherwise has no external dependencies

- Clone this repository
- Once the operating system is flashed, run the script with `sudo go run main.go` under the `cmd` folder or `sudo go run ./...` on root
- Answer the prompts to configure Tailscale. One of the prompts will request an auth key, which you can generate from your [Tailscale admin console](https://login.tailscale.com/admin/settings/keys)

When your Rasbperry Pi boots up, you should see it in your admin console's [machines](https://login.tailscale.com/admin/machines) page and you should be able to use to [Tailscale SSH](https://tailscale.com/tailscale-ssh/) to connect to it

```
tailscale ssh ubuntu@<hostname>
```

Depending on your ACL configuration, you may be prompted to authenticate with Tailscale.

## Contributing

Feel free to open any issues or pull requests!

## TEST!