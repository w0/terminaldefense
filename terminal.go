package main

import (
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/creack/pty"
)

type Terminal struct {
	pty     *os.File
	hub     *Hub
	pending map[string]*PendingCommand
}

type PendingCommand struct {
	Command  []byte
	Executed bool
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
	go t.handleAction()

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

			cmdId := hashIt(string(cmd))

			t.pending[cmdId] = &PendingCommand{
				Command:  cmd,
				Executed: false,
			}

			t.notifyAdmin(cmd, cmdId)
			continue
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

func (t *Terminal) notifyAdmin(cmd []byte, id string) {
	log.Println("notifying admins")
	res := make(map[string]string)

	res["command"] = string(cmd)
	res["dangerous"] = "true"
	res["id"] = id

	for client := range t.hub.clients {
		if client.role == "sysadmin" {
			client.conn.WriteJSON(&res)
		}
	}
}

func (t *Terminal) handleAction() {
	for {
		a := <-t.hub.action

		if a.Action == "allow" {
			t.pending[a.Id].Executed = true
			t.writeToShell(t.pending[a.Id].Command)
		} else {
			log.Printf("blocked %v", a)
		}

	}
}
