package plex

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

type dynamicBool struct {
	bool
}

// this function attempts to parse either a plain boolean or a binary into a boolean (binary can be an int or string)
func (b *dynamicBool) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		b.bool = false
		return nil
	}
	var isInt int

	if err := json.Unmarshal(data, &isInt); err == nil {
		if isInt == 0 || isInt == 1 {

			if isInt != 0 && isInt != 1 {
				return fmt.Errorf("invalid boolOrInt: %d", isInt)
			}

			b.bool = isInt == 1

			return nil
		}
	}

	var isBool bool

	if err := json.Unmarshal(data, &isBool); err == nil {
		b.bool = isBool
		return nil
	}

	var isString string

	if err := json.Unmarshal(data, &isString); err == nil {
		parsed, err := strconv.ParseBool(isString)
		if err != nil {
			return err
		}
		b.bool = parsed
		return nil
	}

	return nil
}

type dynamicTime struct {
	time.Time
}

const (
	layoutISO     = "2006-01-02"
	layoutUS      = "January 2, 2006"
	layoutRFC3339 = "2006-01-02T15:04:05Z07:00"
)

// this function first attempts to parse time as a unix timestamp, then as RFC3339 string
func (d *dynamicTime) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		d.Time = time.Now()
		return nil
	}
	var unixTime json.Number
	err := json.Unmarshal(data, &unixTime)
	if err != nil {
		var toParse string
		err := json.Unmarshal(data, &toParse)
		if err != nil {
			return err
		}

		layouts := []string{layoutISO, layoutUS, layoutRFC3339}
		var t time.Time
		for _, layout := range layouts {
			t, err = time.Parse(layout, toParse)
			if err == nil {
				break
			}
		}

		d.Time = t
		return nil

	}
	parsed, err := unixTime.Int64()
	if err != nil {
		return err
	}
	d.Time = time.Unix(parsed, 0)
	return nil
}

// type dynamicStringArray struct {
// 	 []string{}
// }

// func (d *dynamicStringArray) UnmarshalJSON(data []byte) error {
// 	var array []string
// 	var single string
// 	err := json.Unmarshal(data, array)
// 	if err != nil {
// 		err := json.Unmarshal(data, single)
// 		if err != nil {
// 			return err
// 		}

// 		dynamicStringArray.Str = [1]string{single}

// 	}
// 	return nil
// }
