package bitly

import "time"

// JSONDate is simple time.Time object with custom json formatter
type JSONDate time.Time

func (jd *JSONDate) MarshalJSON() ([]byte, error) {
	panic("implement me")
}

func (jd *JSONDate) UnmarshalJSON([]byte) error {
	panic("implement me")
}
