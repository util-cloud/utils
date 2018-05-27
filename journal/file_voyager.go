package journal

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"path"
	"time"
	"strconv"
)

var journalTitleSection map [byte] func(* time.Time) int

type FileVoyager struct {
	pathFmt string
	file * os.File
	bufFileVoyager * bufio.Writer
	actions [] func(* time.Time) int
	variables [] interface{}
}

func NewFileVoyager() * FileVoyager {
	return & FileVoyager{}
}

func (v * FileVoyager) Init() error {
	return v.Rotate()
}

func (v * FileVoyager) SetPathPattern(pattern string) error {
	n := 0

	for _, c := range pattern {
		if c == '%' {
			n++
		}
	}

	if n == 0 {
		v.pathFmt = pattern
		return nil
	}

	v.actions = make([] func(* time.Time) int, 0, n)
	v.variables = make([] interface{}, n, n)

	tmp := [] byte(pattern)

	variable := 0

	for _, c := range tmp {
		if variable == 1 {
			act, ok := journalTitleSection[c]
			if !ok {
				return errors.New("Invalid rotate pattern (" + pattern + ")")
			}
			v.actions = append(v.actions, act)
			variable = 0
			continue
		}

		if c == '%' {
			variable = 1
		}
	}

	v.pathFmt = convertPatternToFmt(tmp)

	return nil
}

func (v * FileVoyager) Write(r * Record) error {
	if v.bufFileVoyager == nil {
		return errors.New("no opened file")
	}

	if _, err := v.bufFileVoyager.WriteString(r.String()); err != nil {
		return err
	}

	return nil
}

func (v * FileVoyager) Rotate() error {
	now := time.Now()
	vv := 0
	rotate := false

	for i, act := range v.actions {
		vv = act(& now)

		if vv != v.variables[i] {
			v.variables[i] = vv
			rotate = true
		}
	}

	if rotate == false {
		return nil
	}

	if v.bufFileVoyager != nil {
		if err := v.bufFileVoyager.Flush(); err != nil {
			return err
		}
	}

	if v.file != nil {
		if err := v.file.Close(); err != nil {
			return err
		}
	}

	filePath := fmt.Sprintf(v.pathFmt, v.variables...)

	if err := os.MkdirAll(path.Dir(filePath), 0755); err != nil {
		if !os.IsExist(err) {
			return err
		}
	}

	if file, err := os.OpenFile(filePath, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0644); err != nil {
		return err
	} else {
		v.file = file
	}

	if v.bufFileVoyager = bufio.NewWriterSize(v.file, 8192); v.bufFileVoyager == nil {
		return errors.New("new bufFileVoyager failed.")
	}

	return nil
}

func (v * FileVoyager) Flush() error {
	if v.bufFileVoyager != nil {
		return v.bufFileVoyager.Flush()
	}

	return nil
}

func convertPatternToFmt(pattern [] byte) string {
	m := map [byte] int {
		'Y': 4,
		'M': 2,
		'D': 2,
		'H': 2,
		'm': 2,
		's': 2,
	}

	keys := [] byte {
		'Y', 'M', 'D', 'H', 'm', 's',
	}

	for _, k := range keys {
		pattern = bytes.Replace(pattern, [] byte("%" + string(k)), [] byte("%0" + strconv.Itoa(m[k]) + "d"), -1)
	}

	return string(pattern)
}

func init() {
	journalTitleSection = make(map [byte] func(* time.Time) int, 5)

	m := map [byte] func(* time.Time) int {
		'Y': func(now * time.Time) int {
			return now.Year()
		},
		'M': func(now * time.Time) int {
			return int(now.Month())
		},
		'D': func(now * time.Time) int {
			return now.Day()
		},
		'H': func(now * time.Time) int {
			return now.Hour()
		},
		'm': func(now * time.Time) int {
			return now.Minute()
		},
		's': func(now * time.Time) int {
			return now.Second()
		},
	}

	keys := [] byte {
		'Y', 'M', 'D', 'H', 'm', 's',
	}

	for _, k := range keys {
		journalTitleSection[k] = m[k]
	}
}
