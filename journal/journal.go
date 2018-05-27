package journal

import (
	"fmt"
	"log"
	"path"
	"runtime"
	"strconv"
	"sync"
	"time"
)

var (
	LEVEL_FLAGS = [...] string {"DEBUG", " INFO", " WARN", "ERROR", "FATAL"}
	recordPool * sync.Pool
)

const (
	DEBUG = iota
	INFO
	WARNING
	ERROR
	FATAL
)

const tunnel_size_default = 1024

type Record struct {
	time  string
	code  string
	info  string
	level int
}

func (r * Record) String() string {
	return fmt.Sprintf("%s [%s] <%s> %s\n", r.time, LEVEL_FLAGS[r.level], r.code, r.info)
}

type Voyager interface {
	Init() error
	Write(* Record) error
}

type Rotater interface {
	Rotate() error
	SetPathPattern(string) error
}

type Flusher interface {
	Flush() error
}

type Journal struct {
	writers     [] Voyager
	tunnel      chan * Record
	level       int
	lastTime    int64
	lastTimeStr string
	c           chan bool
	layout      string
}

func init() {
	journal_default = NewJournal()

	recordPool = & sync.Pool {
		New: func() interface{} {
			return & Record{}
		},
	}
}

func NewJournal() * Journal {
	if journal_default != nil && takeup == false {
		takeup = true
		return journal_default
	}

	j := new(Journal)
	j.writers = make([] Voyager, 0, 2)
	j.tunnel = make(chan * Record, tunnel_size_default)
	j.c = make(chan bool, 1)
	j.level = DEBUG
	j.layout = "2006/01/02 15:04:05"

	go boostrapJournalVoyager(j)

	return j
}

func (j * Journal) Register(w Voyager) {
	if err := w.Init(); err != nil {
		panic(err)
	}

	j.writers = append(j.writers, w)
}

func (j * Journal) SetLevel(level int) {
	j.level = level
}

func (j * Journal) SetLayout(layout string) {
	j.layout = layout
}

func (j * Journal) Debug(fmt string, args ...interface{}) {
	j.deliverRecordToVoyager(DEBUG, fmt, args...)
}

func (j * Journal) Warn(fmt string, args ...interface{}) {
	j.deliverRecordToVoyager(WARNING, fmt, args...)
}

func (j * Journal) Info(fmt string, args ...interface{}) {
	j.deliverRecordToVoyager(INFO, fmt, args...)
}

func (j * Journal) Error(fmt string, args ...interface{}) {
	j.deliverRecordToVoyager(ERROR, fmt, args...)
}

func (j * Journal) Fatal(fmt string, args ...interface{}) {
	j.deliverRecordToVoyager(FATAL, fmt, args...)
}

func (j * Journal) Close() {
	close(j.tunnel)

	<-j.c

	for _, w := range j.writers {
		if f, ok := w.(Flusher); ok {
			if err := f.Flush(); err != nil {
				log.Println(err)
			}
		}
	}
}

func (j * Journal) deliverRecordToVoyager(level int, format string, args ...interface{}) {
	var inf, code string

	if level < j.level {
		return
	}

	if format != "" {
		inf = fmt.Sprintf(format, args...)
	} else {
		inf = fmt.Sprint(args...)
	}

	// source code, file and line num
	_, file, line, ok := runtime.Caller(2)
	if ok {
		code = path.Base(file) + ":" + strconv.Itoa(line)
	}

	// format time
	now := time.Now()
	if now.Unix() != j.lastTime {
		j.lastTime = now.Unix()
		j.lastTimeStr = now.Format(j.layout)
	}

	r := recordPool.Get().(* Record)
	r.info = inf
	r.code = code
	r.time = j.lastTimeStr
	r.level = level

	j.tunnel <- r
}

func boostrapJournalVoyager(journal *Journal) {
	if journal == nil {
		panic("journal is nil")
	}

	var (
		r * Record
		ok bool
	)

	if r, ok = <- journal.tunnel; !ok {
		journal.c <- true
		return
	}

	for _, w := range journal.writers {
		if err := w.Write(r); err != nil {
			log.Println(err)
		}
	}

	flushTimer := time.NewTimer(time.Millisecond * 500)
	rotateTimer := time.NewTimer(time.Minute * 10)

	for {
		select {
			case r, ok = <-journal.tunnel:
				if !ok {
					journal.c <- true
					return
				}

				for _, w := range journal.writers {
					if err := w.Write(r); err != nil {
						log.Println(err)
					}
				}

				recordPool.Put(r)
			case <- flushTimer.C:
				for _, w := range journal.writers {
					if f, ok := w.(Flusher); ok {
						if err := f.Flush(); err != nil {
							log.Println(err)
						}
					}
				}

				flushTimer.Reset(time.Millisecond * 1000)
			case <- rotateTimer.C:
				for _, w := range journal.writers {
					if r, ok := w.(Rotater); ok {
						if err := r.Rotate(); err != nil {
							log.Println(err)
						}
					}
				}

				rotateTimer.Reset(time.Second * 10)
		}
	}
}

// default
var (
	journal_default * Journal
	takeup         = false
)

func SetLevel(level int) {
	journal_default.level = level
}

func SetLayout(layout string) {
	journal_default.layout = layout
}

func Debug(fmt string, args ...interface{}) {
	journal_default.deliverRecordToVoyager(DEBUG, fmt, args...)
}

func Warn(fmt string, args ...interface{}) {
	journal_default.deliverRecordToVoyager(WARNING, fmt, args...)
}

func Info(fmt string, args ...interface{}) {
	journal_default.deliverRecordToVoyager(INFO, fmt, args...)
}

func Error(fmt string, args ...interface{}) {
	journal_default.deliverRecordToVoyager(ERROR, fmt, args...)
}

func Fatal(fmt string, args ...interface{}) {
	journal_default.deliverRecordToVoyager(FATAL, fmt, args...)
}

func Register(w Voyager) {
	journal_default.Register(w)
}

func Close() {
	journal_default.Close()
}
