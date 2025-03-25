package common

type Bet struct {
	ClientId  string
	FirstName string
	LastName  string
	Document  string
	BirthDate string
	Number    string
}

func ParseBet(line []string, clientID string) (Bet, bool) {
	if len(line) != 5 {
		return Bet{}, false
	}
	bet := Bet{
		ClientId:  clientID,
		FirstName: line[0],
		LastName:  line[1],
		Document:  line[2],
		BirthDate: line[3],
		Number:    line[4],
	}
	return bet, true
}
