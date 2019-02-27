# Gomegagen

Generate [gomega](https://github.com/onsi/gomega) matchers for any struct


## Installation

```bash
go get -u github.com/andcan/gomegagen/cmd/gomegagen
```

## Usage

```bash
gomegagen \
    -i github.com/example/pkg \
    -i github.com/example1/pkg \
    -p github.com/out/pkg/path \
    --whitelist-struct github.com/example/pkg.StructName \
    --whitelist-struct github.com/example1/pkg.StructName
```

* `-i`, `--input-dirs` list of import paths to get input types from
* `-p`, `--output-package` out package path
* `--whitelist-struct` generates matchers only for whitelisted structs (has precedence over blacklist)
* `--blacklist-struct` generates matchers only for non blacklisted structs (TODO make blacklist smarter, generates matchers for every imported package)

----

Based on [gengo](https://github.com/kubernetes/gengo)