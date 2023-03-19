package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

func main() {
	// Prompt the user for a few things before firing off the conversation.
	reader := bufio.NewReader(os.Stdin)

	// Ask for the background context for each discussion
	fmt.Println("Enter the background/context about this ensuing conversation:")
	convo, _ := reader.ReadString('\n')
	convo = strings.TrimSuffix(convo, "\n")

	// Ask for personal and/or stylistic information about each bot.
	var bots []*bot
	for i := 0; i < 2; i++ {
		fmt.Printf("Tell me Bot#%d's name: ", i+1)
		botName, _ := reader.ReadString('\n')
		botName = strings.TrimSuffix(botName, "\n")
		fmt.Printf("Tell me key personal, background, or stylistic information about %s:\n", botName)
		botContext, _ := reader.ReadString('\n')
		botContext = strings.TrimSuffix(botContext, "\n")

		// Create a new bot and add it to the list of bots.
		bots = append(bots, newBot(convo, botName, botContext))
	}

	// Now that we have the bots, we can start the conversation. Ask for the starting
	// point for the conversation, which Bot#1 will say to Bot#2:
	fmt.Println("--------------------------------------------------")
	fmt.Printf("Let's start the conversation!\n%s: ", bots[0].Name)
	lastMessage, _ := reader.ReadString('\n')
	lastMessage = strings.TrimSuffix(lastMessage, "\n")

	// Now loop and keep the conversation going; if the user wants to escape, they can
	// hit ^C, otherwise hitting <ENTER> will keep the conversation going.
	turn := 0
	for {
		// We alternate turns between the bots.
		asker := bots[turn]
		replier := bots[1-turn]
		fmt.Println("--------------------------------------------------")

		// The asker now asks a question of the replier.
		reply, err := replier.Chat(asker.Name, lastMessage)
		if err != nil {
			panic(err)
		}
		reply = strings.TrimSuffix(reply, "\n")
		fmt.Printf("%s: %s\n", replier.Name, reply)

		// See if the user wants to keep going.
		fmt.Printf("\n[<ENTER> to continue; ^C to quit]\n")
		fmt.Printf("[Feel free to inject new conversation context before <ENTER>: ")
		addedContext, _ := reader.ReadString('\n')
		addedContext = strings.TrimSuffix(addedContext, "\n")
		if addedContext != "" {
			// If there is new conversation context, inject it into both bots.
			asker.InjectContext(addedContext)
			replier.InjectContext(addedContext)
		}

		// If we're continuing onwards, swap the asker and replier, and keep going!
		turn = 1 - turn
		lastMessage = reply
	}
}

type bot struct {
	// Convo is the background/context for the conversation.
	Convo string
	// Name is the bot's name.
	Name string
	// Context is the personal, background, or stylistic information about the bot.
	Context string
	// Client is the ChatGPT client to use for this bot.
	Client *openai.Client
	// History includes the full history of the current conversation, including
	// system messages used to inject the bot's context.
	History []openai.ChatCompletionMessage
}

func newBot(convo, name, context string) *bot {
	// Inject a "system" message generated based on the conversation and bot context.
	history := []openai.ChatCompletionMessage{
		{
			Role: "system",
			Content: fmt.Sprintf(
				"You are about to have a conversation. To prepare you, here is an overview "+
					"of what that conversation is expected to entail: %s. ",
				convo,
			),
		},
		{
			Role: "system",
			Content: fmt.Sprintf(
				"You have a personality. Your name is %s, and you have the following key personal, "+
					"background, and stylistic traits which youre responses should be consistent with: %s",
				name, context,
			),
		},
		{
			Role: "system",
			Content: "All of your replies should be from your perspective and should be a single person's " +
				"response as though you are actually having a conversation with another individual.",
		},
	}
	return &bot{
		Convo:   convo,
		Name:    name,
		Context: context,
		History: history,
	}
}

func (bot *bot) Chat(from, message string) (string, error) {
	// Create a ChatGPT client configured with the API key.
	apiKey := os.Getenv("OPENAI_API_KEY")
	client := openai.NewClient(apiKey)

	// Create a sequence of chat messages to submit to ChatGPT.
	var messages []openai.ChatCompletionMessage
	messages = append(messages, bot.History...)
	messages = append(messages, openai.ChatCompletionMessage{
		Name:    strings.ReplaceAll(from, " ", ""),
		Role:    "user",
		Content: message,
	})

	// Submit the messages to ChatGPT and get the response.
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    openai.GPT3Dot5Turbo,
			Messages: messages,
		},
	)
	if err != nil {
		return "", fmt.Errorf("creating chat completion: %w", err)
	}
	bot.History = append(bot.History, resp.Choices[0].Message)

	// Fetch and return the reply.
	return resp.Choices[0].Message.Content, nil
}

func (bot *bot) InjectContext(context string) {
	// Inject a "system" message generated based on the conversation and bot context.
	bot.History = append(bot.History, openai.ChatCompletionMessage{
		Role: "system",
		Content: fmt.Sprintf(
			"From this point onwards in the conversation, please keep this information in mind: %s",
			context,
		),
	})
}
