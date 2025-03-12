package user

// func {
// 1 - Fetch Character informations from the User Blizzard Account
// -- Profile Summary to get class, spec, the ids, media
// -- Equipment to get the item level
// -- Mythic Plus to get the scores
// Will not update users char at every request
// Will store user characters data in the database
// 2 - Store those informations into the users characters table in database
// }

// func {
// 3 - Allow players to refresh his Battle.Net blizzard account, so it can launch 1 and 2 func on demand
// Might need to allow a limit so ppl dont abuse it, like rate limiting this.
// }

/*

The auto data refresh happens because it is sending the GetBattleNetProfile on the way after the oauth2.0 flow
As showcasing here

and because of the GetUserInfo func in the blizzard/auth/service.go file

// RegisterRoutes registers all Battle.net authentication routes

		// Routes requiring a linked Battle.net account
		bnetProtected := authed.Group("")
		bnetProtected.Use(h.middleware.RequireBattleNetAccount())
		bnetProtected.Use(h.middleware.RequireValidToken())
		{
			bnetProtected.GET("/profile", h.GetBattleNetProfile)
		}
	}
}

*/
