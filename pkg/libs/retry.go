package libs

import (
	"time"
)

var retrySleepTimes = []int{0, 0, 500, 1000, 2000, 5000, 10000}

func Retry(times int, fun func(retryTimes int) ([]byte, error)) (by []byte, err error) {
	length := len(retrySleepTimes)
	for i := 0; times == 0 || i <= times; i++ {
		var sleepTime int
		if i >= length {
			sleepTime = retrySleepTimes[length-1]
		} else {
			sleepTime = retrySleepTimes[i]
		}
		time.Sleep(time.Millisecond * time.Duration(sleepTime))

		by, err = fun(i)
		if err == nil {
			return by, err
		}
	}
	return nil, err
}
