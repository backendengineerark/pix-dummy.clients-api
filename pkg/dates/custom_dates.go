package dates

import "time"

var birhFormatExample = "2006-01-02"

func GenerateBirthDate(birthDate string) (time.Time, error) {
	date, err := time.Parse(birhFormatExample, birthDate)
	if err != nil {
		return date, err
	}
	return date, nil
}

func DateToString(date time.Time) string {
	return date.Format(birhFormatExample)
}
