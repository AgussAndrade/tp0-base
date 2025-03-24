package common
import (
	"os"
)
type Bet struct {
	ClientId  string
	FirstName string
	LastName  string
	Document  string
	BirthDate string
	Number    string
}

func getBets(clientId string) []Bet {
	bet := Bet{
		ClientId:  clientId,
		FirstName: os.Getenv("NOMBRE"),
		LastName:  os.Getenv("APELLIDO"),
		Document:  os.Getenv("DOCUMENTO"),
		BirthDate: os.Getenv("NACIMIENTO"),
		Number:    os.Getenv("NUMERO"),
	}
	return []Bet{bet}
}
