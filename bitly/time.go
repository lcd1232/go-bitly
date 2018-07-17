package bitly

import (
	"time"
)

// JSONDate is simple time.Time object with custom json formatter
type JSONDate time.Time

const timeFormat = "2006-01-02T15:04:05-0700"

func (jd JSONDate) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(timeFormat)+2)
	b = append(b, '"')
	b = time.Time(jd).AppendFormat(b, timeFormat)
	b = append(b, '"')
	return b, nil
}

func (jd *JSONDate) UnmarshalJSON(p []byte) error {
	t, err := time.Parse(`"`+timeFormat+`"`, string(p))
	if err != nil {
		return err
	}
	*jd = JSONDate(t)
	return nil
}
