package model

type Player struct {
	Id          int            `json:"id"`
	FirstName   string         `json:"first_name"`
	LastName    string         `json:"last_name"`
	Country     string         `json:"country"`
	Age         uint8          `json:"age"`
	Position    PlayerPosition `json:"position"`
	MarketValue int            `json:"market_value"`
	TeamId      *int
}

type Team struct {
	Id            int      `json:"id"`
	Name          string   `json:"name"`
	Country       string   `json:"country"`
	AvailableCash int      `json:"available_cash"`
	Players       []Player `json:"players"`
	AccountId     int
}

type Account struct {
	Id                int    `json:"id"`
	Username          string `json:"email"`
	Password          string
	FirstName         string `json:"first_name"`
	LastName          string `json:"last_name"`
	Team              *Team  `json:"team"`
	VerificationToken string
	LoginAttempts     uint8
	Locked            bool
	Confirmed         bool
}

type Transfer struct {
	Id              int    `json:"id"`
	Player          Player `json:"player"`
	MarketValue     int    `json:"market_value"`
	AskedPrice      int    `json:"asked_price"`
	TransferredFrom Team   `json:"transferred_from"`
	TransferredTo   Team   `json:"transferred_to"`
	Transferred     bool   `json:"transferred"`
}

type PlayerPosition string

const (
	GoalKeeper PlayerPosition = "GK"
	Defender                  = "DF"
	Midfielder                = "MF"
	Forward                   = "FW"
)
