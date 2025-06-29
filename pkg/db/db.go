package db

import "time"

const (
	timeFormat = "2006-01-02 15:04:05"
)

type Time time.Time

type BaseModel struct {
	Id        int64
	CreatedAt Time
	UpdatedAt Time
}

func (t *Time) UnmarshalJSON(data []byte) (err error) {
	now, err := time.ParseInLocation(`"`+timeFormat+`"`, string(data), time.Local)
	*t = Time(now)
	return
}

func (t Time) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(timeFormat)+2)
	b = append(b, '"')
	b = time.Time(t).AppendFormat(b, timeFormat)
	b = append(b, '"')
	return b, nil
}

func (t Time) String() string { return time.Time(t).Format(timeFormat) }
