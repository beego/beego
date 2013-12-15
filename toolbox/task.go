package toolbox

import (
	"log"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"
)

// bounds provides a range of acceptable values (plus a map of name to value).
type bounds struct {
	min, max uint
	names    map[string]uint
}

// The bounds for each field.
var (
	AdminTaskList map[string]Tasker
	stop          chan bool
	seconds       = bounds{0, 59, nil}
	minutes       = bounds{0, 59, nil}
	hours         = bounds{0, 23, nil}
	days          = bounds{1, 31, nil}
	months        = bounds{1, 12, map[string]uint{
		"jan": 1,
		"feb": 2,
		"mar": 3,
		"apr": 4,
		"may": 5,
		"jun": 6,
		"jul": 7,
		"aug": 8,
		"sep": 9,
		"oct": 10,
		"nov": 11,
		"dec": 12,
	}}
	weeks = bounds{0, 6, map[string]uint{
		"sun": 0,
		"mon": 1,
		"tue": 2,
		"wed": 3,
		"thu": 4,
		"fri": 5,
		"sat": 6,
	}}
)

const (
	// Set the top bit if a star was included in the expression.
	starBit = 1 << 63
)

type Schedule struct {
	Second uint64
	Minute uint64
	Hour   uint64
	Day    uint64
	Month  uint64
	Week   uint64
}

type TaskFunc func() error

type Tasker interface {
	GetStatus() string
	Run() error
	SetNext(time.Time)
	GetNext() time.Time
	SetPrev(time.Time)
	GetPrev() time.Time
}

type taskerr struct {
	t       time.Time
	errinfo string
}

type Task struct {
	Taskname string
	Spec     *Schedule
	DoFunc   TaskFunc
	Prev     time.Time
	Next     time.Time
	Errlist  []*taskerr //errtime:errinfo
	ErrLimit int        //max length for the errlist 0 stand for there' no limit
}

func NewTask(tname string, spec string, f TaskFunc) *Task {

	task := &Task{
		Taskname: tname,
		DoFunc:   f,
		ErrLimit: 100,
	}
	task.SetCron(spec)
	return task
}

func (tk *Task) GetStatus() string {
	var str string
	for _, v := range tk.Errlist {
		str += v.t.String() + ":" + v.errinfo + "\n"
	}
	return str
}

func (tk *Task) Run() error {
	err := tk.DoFunc()
	if err != nil {
		if tk.ErrLimit > 0 && tk.ErrLimit > len(tk.Errlist) {
			tk.Errlist = append(tk.Errlist, &taskerr{t: tk.Next, errinfo: err.Error()})
		}
	}
	return err
}

func (tk *Task) SetNext(now time.Time) {
	tk.Next = tk.Spec.Next(now)
}

func (tk *Task) GetNext() time.Time {
	return tk.Next
}
func (tk *Task) SetPrev(now time.Time) {
	tk.Prev = now
}

func (tk *Task) GetPrev() time.Time {
	return tk.Prev
}

//前6个字段分别表示：
//       秒钟：0-59
//       分钟：0-59
//       小时：1-23
//       日期：1-31
//       月份：1-12
//       星期：0-6（0表示周日）

//还可以用一些特殊符号：
//       *： 表示任何时刻
//       ,：　表示分割，如第三段里：2,4，表示2点和4点执行
//　　    －：表示一个段，如第三端里： 1-5，就表示1到5点
//       /n : 表示每个n的单位执行一次，如第三段里，*/1, 就表示每隔1个小时执行一次命令。也可以写成1-23/1.
/////////////////////////////////////////////////////////
//	0/30 * * * * *                        每30秒 执行
//	0 43 21 * * *                         21:43 执行
//	0 15 05 * * * 　　                     05:15 执行
//	0 0 17 * * *                          17:00 执行
//	0 0 17 * * 1                          每周一的 17:00 执行
//	0 0,10 17 * * 0,2,3                   每周日,周二,周三的 17:00和 17:10 执行
//	0 0-10 17 1 * *                       毎月1日从 17:00到7:10 毎隔1分钟 执行
//	0 0 0 1,15 * 1                        毎月1日和 15日和 一日的 0:00 执行
//	0 42 4 1 * * 　 　                     毎月1日的 4:42分 执行
//	0 0 21 * * 1-6　　                     周一到周六 21:00 执行
//	0 0,10,20,30,40,50 * * * *　           每隔10分 执行
//	0 */10 * * * * 　　　　　　              每隔10分 执行
//	0 * 1 * * *　　　　　　　　               从1:0到1:59 每隔1分钟 执行
//	0 0 1 * * *　　　　　　　　               1:00 执行
//	0 0 */1 * * *　　　　　　　               毎时0分 每隔1小时 执行
//	0 0 * * * *　　　　　　　　               毎时0分 每隔1小时 执行
//	0 2 8-20/3 * * *　　　　　　             8:02,11:02,14:02,17:02,20:02 执行
//	0 30 5 1,15 * *　　　　　　              1日 和 15日的 5:30 执行
func (t *Task) SetCron(spec string) {
	t.Spec = t.parse(spec)
}

func (t *Task) parse(spec string) *Schedule {
	if len(spec) > 0 && spec[0] == '@' {
		return t.parseSpec(spec)
	}
	// Split on whitespace.  We require 5 or 6 fields.
	// (second) (minute) (hour) (day of month) (month) (day of week, optional)
	fields := strings.Fields(spec)
	if len(fields) != 5 && len(fields) != 6 {
		log.Panicf("Expected 5 or 6 fields, found %d: %s", len(fields), spec)
	}

	// If a sixth field is not provided (DayOfWeek), then it is equivalent to star.
	if len(fields) == 5 {
		fields = append(fields, "*")
	}

	schedule := &Schedule{
		Second: getField(fields[0], seconds),
		Minute: getField(fields[1], minutes),
		Hour:   getField(fields[2], hours),
		Day:    getField(fields[3], days),
		Month:  getField(fields[4], months),
		Week:   getField(fields[5], weeks),
	}

	return schedule
}

func (t *Task) parseSpec(spec string) *Schedule {
	switch spec {
	case "@yearly", "@annually":
		return &Schedule{
			Second: 1 << seconds.min,
			Minute: 1 << minutes.min,
			Hour:   1 << hours.min,
			Day:    1 << days.min,
			Month:  1 << months.min,
			Week:   all(weeks),
		}

	case "@monthly":
		return &Schedule{
			Second: 1 << seconds.min,
			Minute: 1 << minutes.min,
			Hour:   1 << hours.min,
			Day:    1 << days.min,
			Month:  all(months),
			Week:   all(weeks),
		}

	case "@weekly":
		return &Schedule{
			Second: 1 << seconds.min,
			Minute: 1 << minutes.min,
			Hour:   1 << hours.min,
			Day:    all(days),
			Month:  all(months),
			Week:   1 << weeks.min,
		}

	case "@daily", "@midnight":
		return &Schedule{
			Second: 1 << seconds.min,
			Minute: 1 << minutes.min,
			Hour:   1 << hours.min,
			Day:    all(days),
			Month:  all(months),
			Week:   all(weeks),
		}

	case "@hourly":
		return &Schedule{
			Second: 1 << seconds.min,
			Minute: 1 << minutes.min,
			Hour:   all(hours),
			Day:    all(days),
			Month:  all(months),
			Week:   all(weeks),
		}
	}
	log.Panicf("Unrecognized descriptor: %s", spec)
	return nil
}

func (s *Schedule) Next(t time.Time) time.Time {

	// Start at the earliest possible time (the upcoming second).
	t = t.Add(1*time.Second - time.Duration(t.Nanosecond())*time.Nanosecond)

	// This flag indicates whether a field has been incremented.
	added := false

	// If no time is found within five years, return zero.
	yearLimit := t.Year() + 5

WRAP:
	if t.Year() > yearLimit {
		return time.Time{}
	}

	// Find the first applicable month.
	// If it's this month, then do nothing.
	for 1<<uint(t.Month())&s.Month == 0 {
		// If we have to add a month, reset the other parts to 0.
		if !added {
			added = true
			// Otherwise, set the date at the beginning (since the current time is irrelevant).
			t = time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
		}
		t = t.AddDate(0, 1, 0)

		// Wrapped around.
		if t.Month() == time.January {
			goto WRAP
		}
	}

	// Now get a day in that month.
	for !dayMatches(s, t) {
		if !added {
			added = true
			t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
		}
		t = t.AddDate(0, 0, 1)

		if t.Day() == 1 {
			goto WRAP
		}
	}

	for 1<<uint(t.Hour())&s.Hour == 0 {
		if !added {
			added = true
			t = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 0, 0, 0, t.Location())
		}
		t = t.Add(1 * time.Hour)

		if t.Hour() == 0 {
			goto WRAP
		}
	}

	for 1<<uint(t.Minute())&s.Minute == 0 {
		if !added {
			added = true
			t = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, t.Location())
		}
		t = t.Add(1 * time.Minute)

		if t.Minute() == 0 {
			goto WRAP
		}
	}

	for 1<<uint(t.Second())&s.Second == 0 {
		if !added {
			added = true
			t = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), 0, t.Location())
		}
		t = t.Add(1 * time.Second)

		if t.Second() == 0 {
			goto WRAP
		}
	}

	return t
}

func dayMatches(s *Schedule, t time.Time) bool {
	var (
		domMatch bool = 1<<uint(t.Day())&s.Day > 0
		dowMatch bool = 1<<uint(t.Weekday())&s.Week > 0
	)

	if s.Day&starBit > 0 || s.Week&starBit > 0 {
		return domMatch && dowMatch
	}
	return domMatch || dowMatch
}

func StartTask() {
	go run()
}

func run() {
	now := time.Now().Local()
	for _, t := range AdminTaskList {
		t.SetNext(now)
	}

	for {
		sortList := NewMapSorter(AdminTaskList)
		sortList.Sort()
		var effective time.Time
		if len(AdminTaskList) == 0 || sortList.Vals[0].GetNext().IsZero() {
			// If there are no entries yet, just sleep - it still handles new entries
			// and stop requests.
			effective = now.AddDate(10, 0, 0)
		} else {
			effective = sortList.Vals[0].GetNext()
		}
		select {
		case now = <-time.After(effective.Sub(now)):
			// Run every entry whose next time was this effective time.
			for _, e := range sortList.Vals {
				if e.GetNext() != effective {
					break
				}
				go e.Run()
				e.SetPrev(e.GetNext())
				e.SetNext(effective)
			}
			continue
		case <-stop:
			return
		}
	}
}

func StopTask() {
	stop <- true
}

func AddTask(taskname string, t Tasker) {
	AdminTaskList[taskname] = t
}

//sort map for tasker
type MapSorter struct {
	Keys []string
	Vals []Tasker
}

func NewMapSorter(m map[string]Tasker) *MapSorter {
	ms := &MapSorter{
		Keys: make([]string, 0, len(m)),
		Vals: make([]Tasker, 0, len(m)),
	}
	for k, v := range m {
		ms.Keys = append(ms.Keys, k)
		ms.Vals = append(ms.Vals, v)
	}
	return ms
}

func (ms *MapSorter) Sort() {
	sort.Sort(ms)
}

func (ms *MapSorter) Len() int { return len(ms.Keys) }
func (ms *MapSorter) Less(i, j int) bool {
	if ms.Vals[i].GetNext().IsZero() {
		return false
	}
	if ms.Vals[j].GetNext().IsZero() {
		return true
	}
	return ms.Vals[i].GetNext().Before(ms.Vals[j].GetNext())
}
func (ms *MapSorter) Swap(i, j int) {
	ms.Vals[i], ms.Vals[j] = ms.Vals[j], ms.Vals[i]
	ms.Keys[i], ms.Keys[j] = ms.Keys[j], ms.Keys[i]
}

func getField(field string, r bounds) uint64 {
	// list = range {"," range}
	var bits uint64
	ranges := strings.FieldsFunc(field, func(r rune) bool { return r == ',' })
	for _, expr := range ranges {
		bits |= getRange(expr, r)
	}
	return bits
}

// getRange returns the bits indicated by the given expression:
//   number | number "-" number [ "/" number ]
func getRange(expr string, r bounds) uint64 {

	var (
		start, end, step uint
		rangeAndStep     = strings.Split(expr, "/")
		lowAndHigh       = strings.Split(rangeAndStep[0], "-")
		singleDigit      = len(lowAndHigh) == 1
	)

	var extra_star uint64
	if lowAndHigh[0] == "*" || lowAndHigh[0] == "?" {
		start = r.min
		end = r.max
		extra_star = starBit
	} else {
		start = parseIntOrName(lowAndHigh[0], r.names)
		switch len(lowAndHigh) {
		case 1:
			end = start
		case 2:
			end = parseIntOrName(lowAndHigh[1], r.names)
		default:
			log.Panicf("Too many hyphens: %s", expr)
		}
	}

	switch len(rangeAndStep) {
	case 1:
		step = 1
	case 2:
		step = mustParseInt(rangeAndStep[1])

		// Special handling: "N/step" means "N-max/step".
		if singleDigit {
			end = r.max
		}
	default:
		log.Panicf("Too many slashes: %s", expr)
	}

	if start < r.min {
		log.Panicf("Beginning of range (%d) below minimum (%d): %s", start, r.min, expr)
	}
	if end > r.max {
		log.Panicf("End of range (%d) above maximum (%d): %s", end, r.max, expr)
	}
	if start > end {
		log.Panicf("Beginning of range (%d) beyond end of range (%d): %s", start, end, expr)
	}

	return getBits(start, end, step) | extra_star
}

// parseIntOrName returns the (possibly-named) integer contained in expr.
func parseIntOrName(expr string, names map[string]uint) uint {
	if names != nil {
		if namedInt, ok := names[strings.ToLower(expr)]; ok {
			return namedInt
		}
	}
	return mustParseInt(expr)
}

// mustParseInt parses the given expression as an int or panics.
func mustParseInt(expr string) uint {
	num, err := strconv.Atoi(expr)
	if err != nil {
		log.Panicf("Failed to parse int from %s: %s", expr, err)
	}
	if num < 0 {
		log.Panicf("Negative number (%d) not allowed: %s", num, expr)
	}

	return uint(num)
}

// getBits sets all bits in the range [min, max], modulo the given step size.
func getBits(min, max, step uint) uint64 {
	var bits uint64

	// If step is 1, use shifts.
	if step == 1 {
		return ^(math.MaxUint64 << (max + 1)) & (math.MaxUint64 << min)
	}

	// Else, use a simple loop.
	for i := min; i <= max; i += step {
		bits |= 1 << i
	}
	return bits
}

// all returns all bits within the given bounds.  (plus the star bit)
func all(r bounds) uint64 {
	return getBits(r.min, r.max, 1) | starBit
}

func init() {
	AdminTaskList = make(map[string]Tasker)
	stop = make(chan bool)
}
