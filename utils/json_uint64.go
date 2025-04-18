package utils

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type FlexibleJsonUInt64 uint64

func (f *FlexibleJsonUInt64) UnmarshalJSON(b []byte) error {
	if len(b) == 0 {
		// If the input data is empty, let's do nothing and throw an error
		return fmt.Errorf("cannot unmarshal empty int")
	} else if b[0] == '"' {
		// If the input data starts with a quote, let's process it like it were a string

		// First, we unmarshal it just like it were a regular string
		var targetStr string
		err := json.Unmarshal(b, &targetStr)
		if err != nil {
			return err
		}

		// Now, let's parse the string and extract the int that's in it (or not)
		parsedInt, err := strconv.ParseUint(targetStr, 10, 64)
		if err != nil {
			return err
		}

		// Everything worked and we have an int
		*f = FlexibleJsonUInt64(parsedInt)
		return nil
	} else {
		// If the input data does not starts with a quote, let's process it like it were a regular int
		var targetInt int64
		err := json.Unmarshal(b, &targetInt)
		if err != nil {
			return err
		}

		// Everything worked and we have an int
		*f = FlexibleJsonUInt64(targetInt)
		return nil
	}
}

func (f *FlexibleJsonUInt64) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var targetStr string
	if err := unmarshal(&targetStr); err != nil {
		return err
	}

	parsedInt, err := strconv.ParseUint(targetStr, 10, 64)
	if err != nil {
		return err
	}

	*f = FlexibleJsonUInt64(parsedInt)
	return nil
}
