# Terminal Defense

Welcome, it is your first day at BIG_CORP! You have been tasked with keeping the hackers out of the mainframe.

A submission for the [boot.dev](https://boot.dev) July 2025 hackathon.

## Game Features
* Keep the hackers at bay! Sysadmins should block any dangerous commands sent to the terminal
* Multiplayer terminal emulator! Run a command and everyone sees the results.
* All commands get executed over a PTY to a running bash shell (watch out for rm -rf / !)
* Utilizes xterm.js for rendering!

## Installation

⚠️ USE DOCKER IF YOU WANT TO RUN DANGEROUS COMMANDS! ⚠️

### Docker (recommended)
* run `docker image pull w0ct/terminaldefense:lastest`
* run `docker run --rm -p 8080:8080 w0ct/terminaldefense:lastest`

### GO RUN
* Install bash shell
* Install go version 1.24
* Close this repo `git clone https://github.com/w0/terminaldefense.git`
* Run `mv .env.example .env` modify TD_PORT if needed.
* Run `go run .`

The server will open up a bash shell running as the user who started it. If you run any dangerous commands, they will execute on your local system. Data loss can and will happen.

The server will be running on localhost:8080 by default. Export `TD_PORT` env variable to change if needed.

## Gameplay/Instructions

1. Open terminaldefense in two tabs. (localhost:8080)
2. One tab will be the hacker, the other will have the sysadmin role.
  1. You will need two open tabs at a minimum. If you don't have a window with the sysadmin role (not reciving dangerous command approvals), try refreshing the page. Currently the server just hands roles in order of hacker -> sysadmin -> hacker ...
3. Start entering commands into the terminal window. The results will appear for all clients.
4. Try entering dangerous commands `rm /usr/bin/ls`. The sysadmin role will be able to allow or block the command from executing. The approval buttons will appear below the terminal window.


## Demo

![demo](https://i.imgur.com/XnVjndI.gif)
