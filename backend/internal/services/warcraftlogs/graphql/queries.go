package warcraftlogs

const (
	DungeonLeaderboardQuery = `
query getDungeonLeaderboard($encounterId: Int!, $region: String!, $page: Int!) {
	worldData {
		encounter(id: $encounterId) {
			name
			fightRankings(serverRegion: $region, page: $page) {
				page
				hasMorePages
				count
				rankings {
					server {
						id
						name
						region
					}
					duration
					startTime
					deaths
					tanks
					healers
					melee
					ranged
					bracketData
					affixes
					team {
						id
						name
						class
						spec
						role
					}
					medal
					score
					leaderboard
				}
			}
		}
	}
}`

	DungeonLogsQuery = `
	query getDungeonLogs($encounterId: Int!, $metric: String!, $className: String!) {
		worldData {
			encounter(id: $encounterId) {
				name
				characterRankings(
					metric: $metric
					includeCombatantInfo: false
					className: $className
				) {
					page
					hasMorePages
					count
					rankings {
						name
						class
						spec
						amount
						hardModeLevel
						duration
						startTime
						report {
							code
							fightID
							startTime
						}
						server {
							id
							name
							region
						}
						bracketData
						faction
						affixes
						medal
						score
					}
				}
			}
		}
	}`
)
