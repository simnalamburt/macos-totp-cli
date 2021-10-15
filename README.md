macos-totp-cli
========
macos-totp-cli is a simple TOTP CLI, powered by keychain of macOS.

```console
$ totp
Usage:
  totp [command]

Available Commands:
  completion  generate the autocompletion script for the specified shell
  delete      Delete a TOTP code
  get         Get a TOTP code
  help        Help about any command
  list        List all registered TOTP codes
  scan        Scan a QR code image

Flags:
  -h, --help   help for totp

Use "totp [command] --help" for more information about a command.

$ totp scan ./image.jpg google
Given QR code successfully registered as "google".

$ totp list
facebook
google
openvpn

$ totp get google
123456

$ totp delete google
Successfully deleted "google".
```

&nbsp;

--------
*macos-totp-cli* is primarily distributed under the terms of both the [Apache
License (Version 2.0)] and the [MIT license]. See [COPYRIGHT] for details.

[MIT license]: LICENSE-MIT
[Apache License (Version 2.0)]: LICENSE-APACHE
[COPYRIGHT]: COPYRIGHT