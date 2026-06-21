package common

import (
	"encoding/json"
	"strconv"
	"time"
)

type FlexTime string

// UnmarshalJSON
func (f *FlexTime) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || string(data) == "null" {
		*f = ""
		return nil
	}

	if data[0] == '"' {
		var s string
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}
		*f = FlexTime(s)
		return nil
	}

	var ejson struct {
		Date *int64 `json:"$date"`
	}
	if err := json.Unmarshal(data, &ejson); err == nil && ejson.Date != nil {
		t := time.UnixMilli(*ejson.Date).UTC()
		*f = FlexTime(t.Format(time.RFC3339))
		return nil
	}

	if n, err := strconv.ParseInt(string(data), 10, 64); err == nil {
		t := time.UnixMilli(n).UTC()
		*f = FlexTime(t.Format(time.RFC3339))
		return nil
	}

	*f = FlexTime(string(data))
	return nil
}

func (f FlexTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(f))
}

func (f FlexTime) String() string {
	return string(f)
}
