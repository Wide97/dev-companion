package utility

import "time"

func Parser(str string) (*time.Time, error) {
	if str == "" {
		return nil, nil
	}
	res, err := time.Parse(time.RFC3339, str)
	if err != nil {
		return nil, err
	}
	var converted time.Time

	converted = res

	return &converted, nil
}
