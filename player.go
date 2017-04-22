package main

import (
	"errors"
	"io"
	"os/exec"
)

var commandList = map[string]string{
	"play":     "p",
	"pause":    "p",
	"subs":     "m",
	"stop":     "q",
	"backward": "\x1b[D",
	"forward":  "\x1b[C",
}

// Player is the struct that controls the playback on the omxplayer
// using an exec.Cmd command
type Player struct {
	Command *exec.Cmd
	PipeIn  io.WriteCloser
	Playing bool
}

// Starts playback for a given filename
// it will stop any playback that is currently running
func (p *Player) Start(filename string) error {
	var err error
	if p.Playing == true {
		p.SendCommand("stop")
	}
	p.Command = exec.Command("omxplayer", "-o", "hdmi", filename)
	p.PipeIn, err = p.Command.StdinPipe()
	if err == nil {
		//p.Command.Stdout = os.Stdout
		err = p.Command.Start()
	}

	if err == nil {
		p.Playing = true
	} else {
		p.Playing = false
	}

	return err
}

// Sends command to player. Command can be one of:
// "play", "pause", "subs", "stop", "backward", "forward"
func (p *Player) SendCommand(command string) error {
	if _, ok := commandList[command]; ok {
		_, err := p.PipeIn.Write([]byte(commandList[command]))
		return err
	}
	return errors.New("player.sendcommand: unknown command " + command)
}
