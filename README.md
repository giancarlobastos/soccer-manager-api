## Soccer online manager game API

* You need to write an API for a simple application where football/soccer fans will create fantasy teams and will be able to sell or buy players.
* User must be able to create an account and log in using the API. 
* Each user can have only one team (user is identified by an email)
* When the user is signed up, they should get a team of 20 players (the system should generate players):
  * 3 goalkeepers
  * 6 defenders
  * 6 midfielders
  * 5 attackers

* Each player has an initial value of $1.000.000.
* Each team has an additional $5.000.000 to buy other players.
* When logged in, a user can see their team and player information

* Team has the following information:
  * Team name, and a country (can be edited)
  * Team value (sum of player values)

* Player has the following information
  * First name, last name, country (can be edited by a team owner)
  * Age (random number from 18 to 40) and market value 

* A team owner can set the player on a transfer list
* When a user places a player on a transfer list, they must set the asking price/value for this player. This value should be listed on a market list. When another user/team buys this player, they must be bought for this price. 
* Each user should be able to see all players on a transfer list and filter them by country, team name, player name, and a value.
* With each transfer, team budgets are updated.
* When a player is transferred to another team, their value should be increased between 10 and 100 per cent. Implement a random factor for this purpose.

### Implement administrator role.
* An administrator who can CRUD users, teams, players, add new players to the market or in the team and change all player/team information, including player’s value
* New users need to verify their account by email. Users should not be able to log in until this verification is complete. 
* When a user fails to log in three times in a row, their account should be blocked automatically, and only admins and managers should be able to unblock it. 
