# WowPerf

WowPerf is a web application that allows users to track their performance in World of Warcraft.

It uses the [Raider.io API](https://raider.io/) to fetch character data, calculates mythic plus scores, raid progression and more.

It also uses the [World of Warcraft API](https://worldofwarcraft.blizzard.com/en-gb/) to fetch character data.

In the future, the goals is to includes other data sources such as WarcraftLogs in order to provide a more comprehensive view of your performance or to compare your performance with other players.

## Features to date

- Fetch character data from Raider.io API
- Fetch character data from World of Warcraft API
- Get all information for a character
- Get all equipment for a character
- Calculate mythic plus scores
- Calculate raid progression

## Todo

- Setting up the interaction with WarcraftLogs API
- Implementing Wowhead tooltip for the item
- Implementing link to the item on Wowhead
- Implementing link to Wow calculator talent tree
- Make the frontend to display the data and let the user interact with it

## Tech stack

- Golang
- Gin
- Next.js
- Tailwind CSS
- React Query
- Postgres
- Docker
- OAuth 2.0

## Contributing

Contributions are welcome! Feel free to open an issue or submit a pull request.
