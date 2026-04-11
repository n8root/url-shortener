package models

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

// Default format 2006-01-02
type Date struct {
	Time   *time.Time
	Format string
}

func (d Date) Value() (driver.Value, error) {
	if d.Time == nil {
		return nil, nil
	}
	return *d.Time, nil
}

func (d *Date) MarshalJSON() ([]byte, error) {
	if d.Time == nil {
		return []byte("null"), nil
	}
	format := d.Format
	if format == "" {
		format = "2006-01-02"
	}

	return []byte(fmt.Sprintf("\"%s\"", d.Time.Format(format))), nil
}

func (d *Date) UnmarshalJSON(data []byte) error {
	s := string(data)
	if s == "null" || s == "" {
		return nil
	}
	// Убираем лишние кавычки из JSON-строки
	s = strings.Trim(s, "\"")

	format := d.Format
	if format == "" {
		format = "2006-01-02"
	}

	t, err := time.Parse(format, s)
	if err != nil {
		return err
	}

	d.Time = &t

	return nil
}
