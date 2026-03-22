# script-launcher
Very simple script launcher

<img width="549" height="532" alt="image" src="https://github.com/user-attachments/assets/194b031a-1b96-49ac-8fd6-e442f9c50e62" />

1. put all "*.sh" script for launch into $HOME/system folder
2. launch your shell script by this app

- build

### for Mac M1, M2...
- *This will install many dependency packages into your Mac.*
- Please review make_mac_arm64.sh before you build.

You should have installed "go" in your mac.

```
$./make_mac_arm64.sh
```

### for Linux (arm64, x86)

- *This will install many dependency packages into your PC.*
- If you don't want it, make a another build PC(VM) and then compile it.

```
make_linux_amd64.sh
```

```
make_linux_arm64.sh
```

