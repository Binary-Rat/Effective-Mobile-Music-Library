package models

import (
	"fmt"
	"time"
)

type Song struct {
	Group   string      `json:"group" example:"group"`
	Song    string      `json:"song" example:"song"`
	Details SongDetails `json:"details"`
}

type SongDetails struct {
	ReleaseDate CustomTime `json:"releaseDate"`
	Lyrics      string     `json:"text" example:"lyrics"`
	Link        string     `json:"link" example:"link"`
}

type Verse string

type CustomTime time.Time

// Формат времени
const customTimeFormat = "02.01.2006"

// Реализация UnmarshalJSON
func (ct *CustomTime) UnmarshalJSON(b []byte) error {
	// Удаляем кавычки из JSON строки
	str := string(b)
	str = str[1 : len(str)-1]

	// Парсим время в заданном формате
	parsedTime, err := time.Parse(customTimeFormat, str)
	if err != nil {
		return err
	}
	*ct = CustomTime(parsedTime)
	return nil
}

func (ct CustomTime) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, time.Time(ct).Format(customTimeFormat))), nil
}
