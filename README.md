# apspace-attendance

cli tool for submitting attendance code to APSpace/AttendiX

also features a pretty tui made with [bubbletea](https://github.com/charmbracelet/bubbletea) if that's your thing

## installation
build the binary, then put the binary in your PATH. for example `mv attendance ~/.local/bin/`

### building
```sh
go build ./cmd/attendance
```

### uninstallation
log out from the app to remove the stored password, then delete the binary

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

