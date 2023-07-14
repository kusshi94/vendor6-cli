# vendor6-cli
An Interactive CLI tool to identify vendors by IPv6 address

## Installation

### go install

```
go install github.com/kusshi94/vendor6-cli@latest
```

### Homebrew

```
brew install kusshi94/tap/vendor6-cli
```

### Manual

Download the latest release from the [releases page](https://github.com/kusshi94/vendor6-cli/releases) and extract the archive file.

## Usage

Launch the tool with the following command:

```
vendor6-cli
```

vendor6-cli interactively accepts IPv6 addresses from user input and returns the vendor name.
If you want to exit the tool, type `exit`.

oui.txt file is automatically downloaded from [IEEE](https://standards-oui.ieee.org/oui/oui.txt) to run the tool. If you want to use your own file, use `-f` / `--oui-file` option.

### Example

```
$ vendor6-cli
>: 2001:db8::0a00:7ff:fe12:3456
Apple, Inc.
>: 2001:db8::6666:b3ff:fe11:1111
TP-LINK TECHNOLOGIES CO.,LTD.
>: exit
$
```
