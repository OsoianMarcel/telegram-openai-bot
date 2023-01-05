# Telegram OpenAI Bot

Telegram Bot made for communication with OpenAI.

*The conversation mode is not implemented.

## Build

To build the executable run:

```
go build ./cmd/tgbot
```

To build the executable for different platforms, check the Go documentation.

## Environment variables

To run this project, you will need to set the following environment variables:

**Required**

`TG_APITOKEN` - Telegram bot API token. 

`GPT_AUTH_TOKEN` - OpenAI API key.

**Optional**

`TG_ADMIN_CHATID` - Bot owner chat id (used for feedback command).

## Example

A bot example can be found here https://t.me/ask_openai_bot

Warning: This bot is served from my home Raspberry Pi, so I cannot guarantee 100% uptime. 

## Give a Star ⭐
If you find this repository useful, please give it a star to show your support for this project. 😊

## License
All contents of this repository are licensed under the [MIT license].

[MIT license]: https://github.com/OsoianMarcel/telegram-openai-bot/blob/main/LICENSE
