# apspace-attendance

![video demo](https://github.com/meow-d/apspace-attendance/raw/refs/heads/main/demo.mp4)

cli tool for submitting attendance code to APSpace/AttendiX

also features a pretty tui made with [bubbletea](https://github.com/charmbracelet/bubbletea) if that's your thing

## installation
build the binary, then put the binary in your PATH. for example `mv attendance ~/.local/bin/`

### installation
download the [latest release](https://github.com/meow-d/apspace-attendance/releases) and put the program anywhere you want

### uninstallation
log out from the app to remove the stored password, then delete the binary

## building
```sh
go get ./...
go build ./cmd/attendance
```

## usage
run `attendance` for the tui

run `attendance 000` to submit directly

## todo
- [X] store passwords in a keyring
- [X] login
- [X] proper readme...
- [X] fix cursor blinking in login after changing focus
- [ ] tests (probably won't... but would be nice...)
- [ ] github actions? why am i even putting effort into this at this point

