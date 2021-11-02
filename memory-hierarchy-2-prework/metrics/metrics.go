package metrics

import (
	"encoding/csv"
	"log"
	"math"
	"os"
	"strconv"
	"time"
)

type UserId int

// type UserMap map[UserId]*User

type UserMap map[string][]int

type Address struct {
	fullAddress string
	zip         int
}

type DollarAmount struct {
	dollars, cents uint64
}

type Payment struct {
	amount DollarAmount
	time   time.Time
}

type User struct {
	id       UserId
	name     string
	age      int
	address  Address
	payments []Payment
}

func AverageAge(users UserMap) float64 {
	count := 0
	sum := 0
	for _, age := range users["age"] {
		count += 1
		sum += age
	}
	average := float64(sum) / float64(count)
	return average
}

// users[UserId(userId)].payments = append(users[UserId(userId)].payments, Payment{
// 	DollarAmount{uint64(paymentCents / 100), uint64(paymentCents % 100)},
// 	datetime

func AveragePaymentAmount(users UserMap) float64 {
	count := 0
	total := 0
	for _, payment := range users["payments"] {
		count += 1
		total += payment
	}

	dollars := float64(total / 100)
	cents := float64(total%100) / 100
	average := (dollars + cents) / float64(count)
	return average
}

// Compute the standard deviation of payment amounts
func StdDevPaymentAmount(users UserMap) float64 {
	mean := AveragePaymentAmount(users)
	squaredDiffs := 0.0
	count := 0.0
	// total := 0
	for _, p := range users["payments"] {
		count += 1
		dollars := float64(p / 100)
		cents := float64(p%100) / 100
		amount := dollars + cents
		diff := amount - mean
		squaredDiffs += diff * diff
	}
	return math.Sqrt(squaredDiffs / count)
}

func LoadData() UserMap {
	f, err := os.Open("users.csv")
	if err != nil {
		log.Fatalln("Unable to read users.csv", err)
	}
	reader := csv.NewReader(f)
	userLines, err := reader.ReadAll()
	if err != nil {
		log.Fatalln("Unable to parse users.csv as csv", err)
	}

	users := make(UserMap, len(userLines))
	ageArray := []int{}

	for _, line := range userLines {
		age, _ := strconv.Atoi(line[2])
		ageArray = append(ageArray, age)
	}

	f, err = os.Open("payments.csv")
	if err != nil {
		log.Fatalln("Unable to read payments.csv", err)
	}
	reader = csv.NewReader(f)
	paymentLines, err := reader.ReadAll()
	if err != nil {
		log.Fatalln("Unable to parse payments.csv as csv", err)
	}

	paymentsArray := []int{}
	for _, line := range paymentLines {
		paymentCents, _ := strconv.Atoi(line[0])
		paymentsArray = append(paymentsArray, paymentCents)
	}

	users["age"] = ageArray
	users["payments"] = paymentsArray

	return users
}
