# script-launcher
Very simple script launcher

<img width="549" height="532" alt="image" src="https://github.com/user-attachments/assets/194b031a-1b96-49ac-8fd6-e442f9c50e62" />

1. Put all your `*.sh` scripts into the `$HOME/system` folder.
2. Launch your shell scripts easily using this app.

## Build

While this app is simple to use, building it requires some preparation due to its UI dependencies.

### For Mac (M1, M2, Apple Silicon)

> [!IMPORTANT]
> **This build process will install several dependency packages on your Mac.**

* **Prerequisite:** You must have `go` installed on your Mac.

```bash
./make_mac_arm64.sh
```

### for Linux (arm64, x86)

> [!IMPORTANT]
> **This build process will install several dependency packages on your PC.**

- If you don't want it, make a another build PC(VM) and then compile it.

```
make_linux_amd64.sh
```

```
make_linux_arm64.sh
```

