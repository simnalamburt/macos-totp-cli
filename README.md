macos-totp-cli
========
macos-totp-cli is a simple TOTP CLI, powered by keychain of macOS.

### Installation
```bash
brew install simnalamburt/x/totp
```

totp is also provided as a single static universal binary. Whether you have an
Mac with Apple silicon or Intel-based Mac, you can install totp by downloading
just one file.

```bash
curl -LO https://github.com/simnalamburt/macos-totp-cli/releases/download/v1.1.1/totp && chmod +x totp
```

### Usage
```console
$ totp
Usage:
  totp [command]

Available Commands:
  add         Manually add a secret to the macOS keychain
  completion  Generate the autocompletion script for the specified shell
  delete      Delete a TOTP code
  get         Get a TOTP code
  help        Help about any command
  list        List all registered TOTP codes
  scan        Scan a QR code image
  temp        Get a TOTP code from a secret without saving it to the keychain

Flags:
  -h, --help      help for totp
  -v, --version   version for totp

Use "totp [command] --help" for more information about a command.

$ totp scan google ./image.jpg
Given QR code successfully registered as "google".

$ totp add github
Type secret: ABCDEFGHIJKLMNOPQRSTUVWXYZ
Given secret successfully registered as "github".

$ totp list
google
github

$ totp get google
123456

$ totp delete google
Successfully deleted "google".

$ totp temp
Type secret: ABCDEFGHIJKLMNOPQRSTUVWXYZ
123456
```

&nbsp;

--------
*macos-totp-cli* is primarily distributed under the terms of both the [Apache
License (Version 2.0)] and the [MIT license]. See [COPYRIGHT] for details.

[MIT license]: LICENSE-MIT
[Apache License (Version 2.0)]: LICENSE-APACHE
[COPYRIGHT]: COPYRIGHT
