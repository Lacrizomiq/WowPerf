# WowPerf

WowPerf is a web application that allows users to track their performance in World of Warcraft.

It also uses the [World of Warcraft API](https://worldofwarcraft.blizzard.com/en-gb/) to fetch character data.

In the future, the goals is to includes other data sources such as WarcraftLogs in order to provide a more comprehensive view of your performance or to compare your performance with other players.

## Features to date

- Fetch character data from Raider.io API (profile, mythic plus scores, raid progression, talents, gear and more)
- Fetch character data from World of Warcraft API (profile, mythic keystone profile, equipment, specializations, media)
- Get all information for a character
- Get all equipment for a character
- Calculate mythic plus scores
- Calculate raid progression
- Get character media from Blizzard API

## Todo

- Setting up the interaction with WarcraftLogs API
- Implementing Wowhead tooltip for the item
- Implementing link to the item on Wowhead
- Implementing link to Wow calculator talent tree
- Make the frontend to display the data and let the user interact with it
- Make the whole backend to be more like a wrapper around the APIs

## Tech stack

- Golang
- Gin
- Next.js
- Tailwind CSS
- React Query
- Postgres
- Docker
- Docker Compose
- OAuth 2.0
- NextAuth.js

## Contributing

Contributions are welcome! Feel free to open an issue or submit a pull request.
