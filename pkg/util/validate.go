package util

import (
	"encoding/json"
	"errors"
	"strconv"
)

func CheckBool(async string) (bool, error) {
	if async == "" {
		return true, nil
	}
	isAysnc, err := strconv.ParseBool(async)
	if err != nil {
		return true, err
	}
	return isAysnc, nil
}


func CheckInt64(num string) (int64, error) {
	if num == "" {
		return 0, errors.New("it is empty")
	}
	n, err := strconv.ParseInt(num, 10, 64)
	if err != nil {
		return 0, err
	}
	return n, nil
}

func CheckUint64(num string) (uint64, error) {
	if num == "" {
		return 0, errors.New("it is empty")
	}
	n, err := strconv.ParseUint(num, 10, 64)
	if err != nil {
		return 0, err
	}
	return n, nil
}
//check length of string

func CheckStringLength(min, max int, str string) error {
	if str == "" {
		return errors.New("empty string")
	}

	length := len(str)
	if max == -1 {
		if length < min {
			return errors.New("unexpected length")
		}else {
			return nil
		}
	}else {
		if length < min || length > max {
			return errors.New("unexpected length")
		}
		return nil
	}
}
//check json string

func CheckJsonArgs(str string, param interface{}) (bool, error) {
	if str == "" {
		return false, errors.New("empty string")
	}
	b := []byte(str)
	err := json.Unmarshal(b, &param)
	if err != nil {
		return false, errors.New("error byte")
	}
	return true, nil
}