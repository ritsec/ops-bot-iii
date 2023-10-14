# Ops Bot III (OBIII)

## Overview

Welcome to Ops Bot III, or as the cool kids like to call it, OBIII! 😎

OBIII is the third iteration of Ops Bot. You can find the previous version, Ops Bot II (OBII), on the RITSEC Gitlab [here](https://gitlab.ritsec.cloud/operations-program/ops-bot-ii). The original is closed-source for your own safety 😉.

OBIII is a Discord bot specifically built for RITSEC and its Discord server. As such, much of the code is purpose-built for RITSEC. If you plan to implement OBIII outside of RITSEC, significant code refactoring may be required.

### Goals

OBIII shares many of the same goals as OBII:

- Modularity
- Meaningful Documentation
- Open Source Community and Contributions

We want OBIII to grow as the needs of its users do. We strive to make it a project that is not only easy to contribute to but also one that people want to contribute to!

### License

OBIII is licensed under the [Apache License 2.0](https://www.apache.org/licenses/LICENSE-2.0).

## Running OBIII for the First Time

### Set up `config.yml`

Make a copy of [config_example.yml](./config_example.yml) and name it `config.yml`:

```
cp config_example.yml config.yml
```

### Discord Bot Token, Application ID, and Guild ID

Instructions for setting up a general bot can be found everywhere, so we won't rehash it here. You can refer to the [Discord.py documentation](https://discordpy.readthedocs.io/en/stable/discord.html) for guidance.

1. Add the token of the bot to `config.yml` under `token`.
2. Add the application ID of the bot to `config.yml` under `app_id`.
3. Add the guild ID (ID of your server) to `config.yml` under `guild_id`.

### Logging Fields

You need to add the IDs of the channels where you want the bot to log.

1. Add the channel IDs to `config.yml` under `logging.[level]_channel` respectively.
2. Set the default logging level in `config.yml` under `logging.level`.
3. Set the default logging file location (ensure OBIII has access to it) in `config.yml` under `logging.log_file`.

### Other Configurations

All other configurations will depend on what you wish to run within your bot.

In the file [commands/enabled.go](./commands/enabled.go), you can enable and disable all the functionality of OBIII by commenting out the respective lines.

## Modular Structure

OBIII is divided into three types of events it can handle: `slash`, `handlers`, and `scheduled`.

Once one of these events is created, it can be enabled in [commands/enabled.go](commands/enabled.go).

These events correspond to different ways a function can be triggered:

### Slash

Slashes are application commands triggered when you press `/` in a Discord message box.

To configure a new slash command, create a new file under [commands/slash/](./commands/slash/). Use the following template:

```go
package slash

import (
	"github.com/bwmarrin/discordgo"
)

func [Name]() *structs.SlashCommand {
	return &structs.SlashCommand{
		Command: &discordgo.ApplicationCommand{
			Name:                     [Name],
			Description:              [Description],
			DefaultMemberPermissions: [Permission],
		},
		Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			[code]
		},
	}
}
```

### Handlers

Handlers act as hooks. They execute based on specific actions, defined by the parameters that the function takes.

To configure a new handler, create a new file under [commands/handlers/](./commands/handler/)

Handlers can respond to various events such as a user joining or leaving, a message being sent, deleted, or edited, and more.

A full set of handlers can be found [here](https://github.com/bwmarrin/discordgo/blob/v0.27.1/events.go)

```go
package handlers

import (
	"github.com/bwmarrin/discordgo"
)

func [Name](s *discordgo.Session, m *discordgo.[event]) {
	[code]
}
```

### Scheduled

Scheduled events run once when the bot starts and manage scheduling tasks on their own.

To configure a new scheduled event, create a new file under [commands/scheduled/](./commands/scheduled/)

There are two main examples that can be used:

#### Cron Events

Cron events run on a schedule or at specific times defined by a cron expression.

```go
package scheduled

import (
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/robfig/cron"
)

func [Name](s *discordgo.Session, quit chan interface{}) error {
	est, _ := time.LoadLocation([location])

	c := cron.NewWithLocation(est)

	c.AddFunc([cron expression], func() { [function] })

	c.Start()
	<-quit
	c.Stop()

	return nil
}
```

#### Continous Events

Continuous events are scheduled on a ticker and execute at regular intervals.

```go
package scheduled

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

func Heartbeat() *structs.ScheduledEvent {
	return structs.NewScheduledTask(
		func(s *discordgo.Session, quit chan interface{}) error {
			ticker := time.NewTicker([interval])
			for {
				select {
				case <-quit:
					return nil
				case <-ticker.C:
					[code]
				}
			}
		},
	)
}
```

## Conclusion

With this information, you have the context needed to understand, contribute to, and deploy OBIII.

Good luck and godspeed!
