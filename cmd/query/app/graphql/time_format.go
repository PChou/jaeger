package graphql

import (
	"time"
)

func ParseSkyWalkingTimeFormat(timeLiteral string) (time.Time, error) {
	return time.ParseInLocation("2006-1-2 15:4", timeLiteral, time.Now().Location())
}
