# Changelog

## [0.1.0](https://github.com/borisdvlpr/gotail/compare/v0.1.0...v0.1.0) (2025-02-25)


### Features

* add StatusError type and optimize error handling for file module ([7fdbab9](https://github.com/borisdvlpr/gotail/commit/7fdbab9e7215de34f673d3e3e7d37558a50f7269))
* add user input for Tailscale configuration ([7352ea8](https://github.com/borisdvlpr/gotail/commit/7352ea8a9dc27aec9e581dbb72d025576dbc25e1))
* add validation to subnet addresses input ([c2ff163](https://github.com/borisdvlpr/gotail/commit/c2ff1633470e93d8e8ec2f05e8ed12c1a88ede7f))
* added promptUser function ([079a583](https://github.com/borisdvlpr/gotail/commit/079a583ff3b5288543e4e3f90c0ae7fe80380158))
* implement checkRoot function tu ensure root execution ([6dc0663](https://github.com/borisdvlpr/gotail/commit/6dc0663cae27d0c9f5a6e2fb42e70200877009ee))
* implement findUserData function to find user-data file ([d86f9c2](https://github.com/borisdvlpr/gotail/commit/d86f9c246e8a6d22fc2017a87b0375f88ca12b87))
* implement lsblk command execution for linux systems ([7d7b310](https://github.com/borisdvlpr/gotail/commit/7d7b310648d308db57774c6abf184583b32e96e4))
* implement writer to add Tailscale configurations to cloud-init file ([471d896](https://github.com/borisdvlpr/gotail/commit/471d8969011943e140bf77e992db662a37eb7c59))
* optimize error handling for input module ([fd22601](https://github.com/borisdvlpr/gotail/commit/fd2260132d3fecbdf74e8640b6d3c035be58c63b))
* parse lsbls json data intro struct ([8aae3e5](https://github.com/borisdvlpr/gotail/commit/8aae3e5ef50287031deefeec0e4dd354ba59fb1f))


### Bug Fixes

* **build:** add missing go installation check to build step on Makefile ([995aa84](https://github.com/borisdvlpr/gotail/commit/995aa847f30a8c9012416482c9e30a5c562a8676))
* change file permission ([61da6d6](https://github.com/borisdvlpr/gotail/commit/61da6d66d962fe3c3f8a91cf96d2eb1526fa54c4))
* **ci:** fix format and lint actions ([9068f11](https://github.com/borisdvlpr/gotail/commit/9068f11a3d5bdf1b79ce0bc3dfe6f70cc6b38ef1))
* **ci:** fix typo on build steps of release pipeline ([1b58c0b](https://github.com/borisdvlpr/gotail/commit/1b58c0b8d7435e864fd925c3862eb610aed10154))
* **ci:** fix typo on pull requests pipeline ([9b47804](https://github.com/borisdvlpr/gotail/commit/9b47804e99dee9ece59a1eebf3f94b9f02c94945))
* **ci:** fix typo on pull requests pipeline ([f62d049](https://github.com/borisdvlpr/gotail/commit/f62d049ecbafe3e1617ad2c486076b46a15c72cf))
* **ci:** remove duplicate test action ([2b7fdd8](https://github.com/borisdvlpr/gotail/commit/2b7fdd8bf115e22b1b3bd6e2e68b6a1d04dbeb71))
* fix configuration added to cloud-init file ([7cc23ff](https://github.com/borisdvlpr/gotail/commit/7cc23ffa1e8efae5dc5a0e2e984f6cb4a95cc697))
* fix issue for macOS not returning the path for the file once found ([d4a7944](https://github.com/borisdvlpr/gotail/commit/d4a794440350b9e2392d0efb9a210d60bf23dca1))
* fix prompt message format for non-existent allowed replies ([abe0eab](https://github.com/borisdvlpr/gotail/commit/abe0eab3b11a6074e69ac41eab6826efdca4d689))
* formatted input prompt of promptUser function ([cd35d9d](https://github.com/borisdvlpr/gotail/commit/cd35d9da5872de3f41c49eda86f833eeac51e633))
* ignore all hidden files and folders, to avoid execution errors ([577d98b](https://github.com/borisdvlpr/gotail/commit/577d98b8b578cd29432bf4a33aea39b51ce80b0e))
* **input:** handle invalid user input by returning an error ([7f21778](https://github.com/borisdvlpr/gotail/commit/7f217784602fcf7ce47b4582a524aa8985c2e333))
* remove unecessary for loop from PrompUser function ([09a0d80](https://github.com/borisdvlpr/gotail/commit/09a0d80dcad2caa8853a5f61432a8609c2f98e8d))
* split subnet router from default config ([44dc5eb](https://github.com/borisdvlpr/gotail/commit/44dc5eb3ea39d07858e05c96ad60e4cba6fa1e0f))
* wrap file closing error handle in closure ([d2ee068](https://github.com/borisdvlpr/gotail/commit/d2ee06833e2e9f33063589826e3a7eb8fed35618))


### Performance Improvements

* add validation against paths to ignore in linux ([040a8c2](https://github.com/borisdvlpr/gotail/commit/040a8c21ebb39e328ea2b5f36f1d4881ca173225))


### Continuous Integration

* fix artifacts name on build ([b2f67b9](https://github.com/borisdvlpr/gotail/commit/b2f67b9322d41261633c05d84c0fd51884a64de6))
* fix artifacts path ([845d137](https://github.com/borisdvlpr/gotail/commit/845d137e81f5e060e5c7c23c0778da3a4173c6d2))
* fix pipeline execution ([b22ab94](https://github.com/borisdvlpr/gotail/commit/b22ab94df7127568b9f5eaeeda6edcd8e58ef832))
