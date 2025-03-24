package common
import (
	"fmt"
	"strings"
)

func formatBetMessage(bet Bet) string {
	return fmt.Sprintf(
		"%s;%s;%s;%s;%s;%s\n",
		bet.ClientId,
		bet.FirstName,
		bet.LastName,
		bet.Document,
		bet.BirthDate,
		bet.Number,
	)
}

func isOkMsg(msg string) bool{
	return strings.Contains(msg, "OK")
}

func isEndOfMsg(aByte byte) bool{
	return aByte == '\n'
}
