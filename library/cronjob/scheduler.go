package cronjob

import "time"

const Timer = `2006-01-02 15:04:05`

var (
	timeDur = map[string]time.Duration{
		"H": time.Hour, "m": time.Minute, "s": time.Second,
	}
)

// 自定义scheduler
type Scheduler struct {
	rate     int    // 类型的值
	rateType string // 类型
	ntime    string // 定时时间点
}

func NewScheduler(rate int, rateType, ntime string) Scheduler {
	return Scheduler{rate, rateType, ntime}
}

func (s *Scheduler) Next(t time.Time) time.Time {
	TimeLocation, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return time.Time{}
	}
	firstTime, err := time.ParseInLocation(Timer, s.ntime, TimeLocation)
	if err != nil {
		return time.Time{}
	}
	if firstTime.After(t) {
		// 说明还没开始第一次任务
		return firstTime
	}
	var duration time.Duration
	switch s.rateType {
	case "Y": // 年
		for i := 1; i < 5; i++ {
			temp := firstTime.AddDate(s.rate, 0, 0)
			if temp.After(t) {
				return temp
			}
			firstTime = temp
		}
	case "M": // 月
		for {
			temp := firstTime.AddDate(0, s.rate, 0)
			if temp.After(t) {
				return temp
			}
			firstTime = temp
		}
	case "D":
		for {
			temp := firstTime.AddDate(0, 0, s.rate)
			if temp.After(t) {
				return temp
			}
			firstTime = temp
		}
	case "H", "m", "s":
		duration = timeDur[s.rateType]
		var subs int
		if s.rateType == "s" {
			subs = int(t.Sub(firstTime).Seconds()/float64(s.rate)) + 1
		} else if s.rateType == "m" {
			subs = int(t.Sub(firstTime).Minutes()/float64(s.rate)) + 1
		} else {
			subs = int(t.Sub(firstTime).Hours()/float64(s.rate)) + 1
		}
		return firstTime.Add(time.Duration(s.rate*subs) * duration)
	}
	// 第一次需要运行的时间
	return time.Time{}
}
