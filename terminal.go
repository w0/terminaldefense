package main

import (
	"log"
	"os"
	"os/exec"

	"github.com/creack/pty"
)

type Terminal struct {
	pty *os.File
	hub *Hub
}

func (t *Terminal) start() {
	pty, err := pty.Start(exec.Command("bash"))
	if err != nil {
		log.Fatalf("bash start %s", err)
	}

	t.pty = pty

	go t.readFromShell()
	go t.writeToShell()

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

func (t *Terminal) writeToShell() {
	for {
		msg := <-t.hub.cmd

		msg = append(msg, '\n')

		_, err := t.pty.Write(msg)
		if err != nil {
			log.Panicln(err)
		}

	}
}
