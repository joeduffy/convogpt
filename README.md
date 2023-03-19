# ConvoGPT

Bots talking to bots.

When you run ConvoGPT, it will ask you some questions:

* Enter the background/context about this ensuing conversation.
* Tell me Bot#1's name.
* Tell me key personal, background, or stylistic information about (Bot#1 name).
* Tell me Bot#2's name.
* Tell me key personal, background, or stylistic information about (Bot#2 name).

From there, it will prompt you to impersonate Bot#1 and say something to Bot#2. Afterwards, the two bots
will continue to interact until you stop it.

After each interaction, you'll be given the opportunity to hit ENTER to continue, or ^C to stop. And you can
add additional contextual information to the conversation going forward by typing that context before hitting ENTER.

Note that the `OPENAI_API_KEY` envvar must contain your OpenAI API key.
