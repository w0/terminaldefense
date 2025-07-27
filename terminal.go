package main

import (
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/creack/pty"
)

type Terminal struct {
	pty *os.File
	hub *Hub
}

var dangerousCommands = []string{
	"rm -rf",
	"sudo rm",
	"rm",
	"format",
	"dd",
}

func (t *Terminal) start() {
	pty, err := pty.Start(exec.Command("bash"))
	if err != nil {
		log.Fatalf("bash start %s", err)
	}

	t.pty = pty

	go t.readFromShell()
	go t.handleCommand()

}

func (t *Terminal) readFromShell() {
	buf := make([]byte, 1024)

	for {
		n, err := t.pty.Read(buf)
		if err != nil {
			log.Printf("terminal read failed %s", err)
		}

		t.hub.broadcast <- buf[:n]
	}

}

func (t *Terminal) writeToShell(cmd []byte) {
	cmdCR := append(cmd, '\n')

	_, err := t.pty.Write(cmdCR)
	if err != nil {
		log.Panicln(err)
	}
}

func (t *Terminal) handleCommand() {
	for {
		cmd := <-t.hub.cmd

		if checkDanger(cmd) {
			log.Println("im in danger")
			//notify ppl
			// give chance to block
			// idk
		}

		t.writeToShell(cmd)
	}
}

func checkDanger(cmd []byte) bool {
	for _, v := range dangerousCommands {
		if strings.Contains(string(cmd), v) {
			return true
		}
	}

	return false
}
