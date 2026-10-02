package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

type userMessage struct {
	Type    string      `json:"type"`
	Message messageBody `json:"message"`
}

type messageBody struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type event struct {
	Type              string            `json:"type"`
	Result            string            `json:"result"`
	PermissionDenials []json.RawMessage `json:"permission_denials"`
}

func main() {
	claudePath, err := exec.LookPath("claude")
	if err != nil {
		fmt.Println("claude not found on PATH:", err)
		os.Exit(1)
	}

	cmd := exec.Command(claudePath, "-p",
		"--input-format", "stream-json",
		"--output-format", "stream-json",
		"--verbose",
		//"--dangerously-skip-permissions",
		"--allowedTools", "Read,Write,Edit,Glob,Grep,WebSearch,WebFetch,Bash(*),PowerShell(*)")
	cmd.Stderr = io.Discard

	stdin, err := cmd.StdinPipe()
	if err != nil {
		fmt.Println("stdin pipe:", err)
		os.Exit(1)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("stdout pipe:", err)
		os.Exit(1)
	}
	if err := cmd.Start(); err != nil {
		fmt.Println("start claude:", err)
		os.Exit(1)
	}
	defer cmd.Process.Kill()

	out := bufio.NewReaderSize(stdout, 1<<20)
	in := bufio.NewReader(os.Stdin)

	fmt.Println("Claude started. Type messages (type quit to exit):")
	fmt.Println()

	for {
		fmt.Print("You: ")
		text, err := in.ReadString('\n')
		if err != nil {
			return
		}
		text = strings.TrimSpace(text)
		if text == "" {
			continue
		}
		if text == "quit" {
			return
		}

		payload, _ := json.Marshal(userMessage{
			Type:    "user",
			Message: messageBody{Role: "user", Content: text},
		})
		if _, err := stdin.Write(append(payload, '\n')); err != nil {
			fmt.Println("Claude is no longer running:", err)
			return
		}

		fmt.Println("Thinking...")
		ev, err := waitForResult(out)
		if err != nil {
			fmt.Println("Claude exited:", err)
			return
		}
		fmt.Printf("Claude: %s\n\n", ev.Result)
		for _, d := range ev.PermissionDenials {
			fmt.Printf("[denied] %s\n", d)
		}
		if len(ev.PermissionDenials) > 0 {
			fmt.Println()
		}
	}
}

func waitForResult(out *bufio.Reader) (event, error) {
	for {
		line, err := out.ReadBytes('\n')
		if len(line) > 0 {
			var ev event
			if json.Unmarshal(line, &ev) == nil && ev.Type == "result" {
				return ev, nil
			}
		}
		if err != nil {
			return event{}, err
		}
	}
}
