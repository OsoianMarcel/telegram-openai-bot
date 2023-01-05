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

## Give a Star ⭐
If you find this repository useful, please give it a star to show your support for this project. 😊

## License
All contents of this repository are licensed under the [MIT license].

[MIT license]: https://github.com/OsoianMarcel/telegram-openai-bot/blob/main/LICENSE
