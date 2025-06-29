package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/fatih/color"
	"github.com/google/shlex"
	"google.golang.org/genai"
)

func fatalErr(err error, msg string) {
	if msg != "" {
		color.Red(msg)
	}

	if err != nil {
		log.Fatal(err)
	}

	os.Exit(1)
}

func checkGeminiAPIKey() {
	if os.Getenv("GEMINI_API_KEY") == "" {
		log.Fatal("GEMINI_API_KEY is not set.")
	}
}

func formatPrompt(prompt string) string {
	prompt = fmt.Sprintf("You are a cli specialist. You only output the resulting cli command requested without format or markdown style, only pure text, nothing more. Get the CLI command for this: %s", prompt)
	return prompt
}

func getPromptFromArgs() string {
	if len(os.Args) < 2 {
		fatalErr(nil, "Please provide a prompt as a command-line argument.")
	}

	prompt := ""
	prompt = fmt.Sprint(os.Args[1:])
	prompt = prompt[1 : len(prompt)-1]

	return prompt
}

func removeMarkdownCommandQuotes(command string) string {
	if len(command) >= 2 && command[0] == '`' && command[len(command)-1] == '`' {
		return command[1 : len(command)-1]
	}
	return command
}

func confirmRunCommand(command string) bool {
	color.Cyan(command)

	fmt.Print("Run this command? (")
	color.New(color.FgYellow).Print("y or enter to run")
	fmt.Print(" | ")
	color.New(color.FgRed).Print("any other key to cancel")
	fmt.Print("): ")

	b := make([]byte, 1)
	os.Stdin.Read(b)
	if b[0] == 'y' || b[0] == 'Y' || b[0] == '\n' {
		return true
	}
	return false
}

func main() {
	checkGeminiAPIKey()

	ctx := context.Background()
	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		fatalErr(err, "")
	}

	prompt := getPromptFromArgs()
	prompt = formatPrompt(prompt)
	result, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash",
		genai.Text(prompt),
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	resultText := result.Text()

	runCommand := confirmRunCommand(resultText)
	if !runCommand {
		color.Red("cancelled")
		return
	}

	fmt.Println()

	commandStr := resultText
	commandStr = removeMarkdownCommandQuotes(commandStr)

	args, err := shlex.Split(commandStr)
	if err != nil || len(args) == 0 {
		log.Fatalf("Failed to parse command: %v", err)
	}
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		log.Fatalf("Command failed: %v", err)
	}
}
