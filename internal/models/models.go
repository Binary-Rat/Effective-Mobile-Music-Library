package models

import (
	"fmt"
	"time"
)

type Song struct {
	Group   string      `json:"group" example:"Nirvana"`
	Song    string      `json:"song" example:"Lithium"`
	Details SongDetails `json:"details"`
}

type SongDetails struct {
	ReleaseDate CustomTime `json:"releaseDate" example:"13.07.1993"`
	Lyrics      string     `json:"text" example:"I am so happy\ncause today..."`
	Link        string     `json:"link" example:"https://www.youtube.com/watch?v=Xsp3_a-PMTw"`
}

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
