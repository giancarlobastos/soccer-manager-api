package main

type Player struct {
	Id          int            `json:"id"`
	FirstName   string         `json:"firstName"`
	LastName    string         `json:"lastName"`
	Country     string         `json:"country"`
	Age         uint8          `json:"age"`
	Position    PlayerPosition `json:"position"`
	MarketValue int            `json:"marketValue"`
	teamId      *int
}

type Team struct {
	Id            int      `json:"id"`
	Name          string   `json:"name"`
	Country       string   `json:"country"`
	AvailableCash int      `json:"availableCash"`
	Players       []Player `json:"players"`
	accountId     int
}

type Account struct {
	Id                int    `json:"id"`
	Username          string `json:"email"`
	password          string
	FirstName         string `json:"firstName"`
	LastName          string `json:"lastName"`
	Team              *Team  `json:"team"`
	verificationToken string
	loginAttempts     uint8
	locked            bool
	confirmed         bool
}

type Transfer struct {
	Id              int    `json:"id"`
	Player          Player `json:"player"`
	MarketValue     int    `json:"marketValue"`
	AskedPrice      int    `json:"askedPrice"`
	TransferredFrom Team   `json:"transferredFrom"`
	TransferredTo   Team   `json:"transferredTo"`
	Transferred     bool   `json:"transferred"`
}

type PlayerPosition string

const (
	GoalKeeper PlayerPosition = "GK"
	Defender                  = "DF"
	Midfielder                = "MF"
	Forward                   = "FW"
)
