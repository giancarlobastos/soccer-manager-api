package domain

type Player struct {
	Id          int            `json:"id"`
	FirstName   string         `json:"firstName"`
	LastName    string         `json:"lastName"`
	Country     string         `json:"country"`
	Age         uint8          `json:"age"`
	Position    PlayerPosition `json:"position"`
	MarketValue int            `json:"marketValue"`
	TeamId      *int           `json:"-"`
}

type Team struct {
	Id            int      `json:"id"`
	Name          string   `json:"name"`
	Country       string   `json:"country"`
	AvailableCash int      `json:"availableCash"`
	Players       []Player `json:"players"`
	AccountId     int      `json:"-"`
}

type Account struct {
	Id                int    `json:"id"`
	Username          string `json:"email"`
	Password          string `json:"-"`
	FirstName         string `json:"firstName"`
	LastName          string `json:"lastName"`
	Team              *Team  `json:"team"`
	VerificationToken string `json:"-"`
	LoginAttempts     uint8  `json:"-"`
	Locked            bool   `json:"-"`
	Confirmed         bool   `json:"-"`
}

type User struct {
	Username string
	TeamId   int
}

type Transfer struct {
	Id              int    `json:"id"`
	Player          Player `json:"player"`
	MarketValue     int    `json:"marketValue"`
	AskedPrice      int    `json:"askedPrice"`
	TransferredFrom *Team  `json:"-"`
	TransferredTo   *Team  `json:"-"`
	Transferred     bool   `json:"-"`
}

type PlayerPosition string

const (
	GoalKeeper PlayerPosition = "GK"
	Defender                  = "DF"
	Midfielder                = "MF"
	Forward                   = "FW"
)
