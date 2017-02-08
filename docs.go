package qb

import (
	"math"
	"strconv"
	"time"

	"github.com/google/uuid"
)

const (
	sizeOverhead = 481
	chars        = "9CjASFTWkKgHrNl8eJXzfphmyb6ncvR2IDU3P1qiL0s4xYotuEQGB7dwaZ5VOM"
	numChars     = int64(len(chars))
)

func newString(i int64, s string, size int) string {
	newString := make([]byte, size)
	bytes := []byte(s)
	numShifts := len(bytes)

	for j := 0; j < size; j++ {
		shift := bytes[numShifts-1]
		idx := (i + int64(j) + int64(shift)) % numChars
		newString[j] = chars[idx]
		numShifts--
		if numShifts == 0 {
			numShifts = len(bytes)
		}
	}
	return string(newString)
}

func newAlphabet(i int64, key string) string {
	return newString(i, key, 64)
}

func newFirstName(alphabet string) string {
	return alphabet[:8]
}

func newLastName(alphabet string) string {
	return alphabet[8:16]
}

func newEmail(alphabet string) string {
	return alphabet[16:24] + "@" + alphabet[24:30] + ".com"
}

func newCity(alphabet string) string {
	return alphabet[30:38]
}

func newCounty(alphabet string) string {
	return alphabet[38:45]
}

func newCountry(alphabet string) string {
	return alphabet[45:54]
}

func newCompany(alphabet string, i int64) string {
	idx := i % 4
	return alphabet[54:64] + " " + corporateType[idx]
}

func newStreet(i int64) string {
	building := strconv.FormatInt(i%5000, 10)
	idx := i % numSuffixes
	capped_small := newGroup(i, 10)
	capped_large := newGroup(i, 1000*(1+i%3))

	return building + " " + capped_small + " " + capped_large + " " + streetSuffixes[idx]
}

func newZipCode(i int64) int64 {
	return 70000 + i%20000
}

func newBalance(alphabet string) float64 {
	var balance, _ = strconv.ParseInt(alphabet[:3], 36, 0)
	return math.Max(0.1, float64(balance)/100.0)
}

func newCategory(i int64) int64 {
	return i % 5
}

func newAge(alphabet string) int64 {
	var age, _ = strconv.ParseInt(string(alphabet[5]), 36, 0)
	return age
}

func newState(i int64) string {
	var idx = i % numStates
	return unitedStates[idx][0]
}

func newFullState(i int64) string {
	var idx = i % numStates
	return unitedStates[idx][1]
}

func newDateOfBirth(alphabet string) string {
	var id, _ = strconv.ParseInt(string(alphabet[:2]), 36, 0)
	seconds := 30 * 24 * 3600 * id
	d := time.Duration(seconds) * time.Second
	t := time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC).Add(d)
	return t.Format(time.RFC3339)
}

func newUUID() string {
	return uuid.New().String()
}

func newAvatar(alphabet string) string {
	return "https://www.gravatar.com/avatar/" + alphabet[32:]
}

func newGroup(i, capacity int64) string {
	return strconv.FormatInt(i/capacity, 16)
}

type address struct {
	City      string `json:"city"`
	County    string `json:"county"`
	Country   string `json:"country"`
	FullState string `json:"fullState"`
	State     string `json:"state"`
	Street    string `json:"street"`
	Zip       int64  `json:"zip"`
}

type doc struct {
	FirstName   string  `json:"firstName"`
	LastName    string  `json:"lastName"`
	Email       string  `json:"email"`
	Address     address `json:"address"`
	Category    int64   `json:"category"`
	Balance     float64 `json:"balance"`
	DateOfBirth string  `json:"dob"`
	UUID        string  `json:"uuid"`
	Notes       string  `json:"notes"`
	Avatar      string  `json:"avatar"`
	Age         int64   `json:"age"`
	Company     string  `json:"company"`
	LocalGroup  string  `json:"localGroup"`
}

func newKey(prefix string, i int64) string {
	return prefix + "-" + strconv.FormatInt(i*i, 16)
}

func newDoc(i int64, key string, size int) doc {
	alphabet := newAlphabet(i, key)

	var notes string
	if size-sizeOverhead > 0 {
		notes = newString(i<<1, alphabet, size-sizeOverhead)
	}

	return doc{
		FirstName: newFirstName(alphabet),
		LastName:  newLastName(alphabet),
		Email:     newEmail(alphabet),
		Address: address{
			City:      newCity(alphabet),
			County:    newCounty(alphabet),
			Country:   newCountry(alphabet),
			FullState: newFullState(i),
			State:     newState(i),
			Street:    newStreet(i),
			Zip:       newZipCode(i),
		},
		Category:    newCategory(i),
		Balance:     newBalance(alphabet),
		DateOfBirth: newDateOfBirth(alphabet),
		UUID:        newUUID(),
		Notes:       notes,
		Avatar:      newAvatar(alphabet),
		Company:     newCompany(alphabet, i),
		Age:         newAge(alphabet),
		LocalGroup:  newGroup(i, 100),
	}
}
