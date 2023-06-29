package client

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Prompt struct {
	history []string
}

func (p *Prompt) Write(instruction string) (string, error) {
	array := make([]string, 0)
	scanner := bufio.NewScanner(os.Stdin)

	input := "Enter"
	if instruction != "" {
		input = fmt.Sprintf("%s(%s)", input, instruction)
	}
	fmt.Printf("%s:", input)

	for {
		scanner.Scan()
		text := scanner.Text()
		if len(text) == 0 {
			break
		}

		array = append(array, text)
	}

	text := strings.Join(array, "\n")
	p.history = append(p.history, text)
	return text, nil
}

func NewPrompt() *Prompt {
	return &Prompt{}
}
