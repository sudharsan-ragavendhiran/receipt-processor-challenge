package rules

import (
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
	"receipt-processor/models"
)

// Rule interface
type Rule interface {
	Calculate(receipt models.Receipt) int
}

// RulesEngine applies all rules dynamically
type RulesEngine struct {
	rules []Rule
}

// NewRulesEngine initializes the rules engine with all defined rules
func NewRulesEngine() *RulesEngine {
	return &RulesEngine{
		rules: []Rule{
			RetailerNameRule{},
			RoundDollarRule{},
			MultipleOfQuarterRule{},
			ItemPairRule{},
			ItemDescriptionRule{},
			OddPurchaseDayRule{},
			TimeRangeRule{},
		},
	}
}

// CalculatePoints executes all rules and returns the total score
func (re *RulesEngine) CalculatePoints(receipt models.Receipt) int {
	totalPoints := 0
	for _, rule := range re.rules {
		totalPoints += rule.Calculate(receipt)
	}
	return totalPoints
}

// Rule 1: One point per alphanumeric character in retailer name
type RetailerNameRule struct{}

func (r RetailerNameRule) Calculate(receipt models.Receipt) int {
	reg := regexp.MustCompile(`[a-zA-Z0-9]`)
	return len(reg.FindAllString(receipt.Retailer, -1))
}

// Rule 2: 50 points if total is a round dollar amount
type RoundDollarRule struct{}

func (r RoundDollarRule) Calculate(receipt models.Receipt) int {
	total, err := strconv.ParseFloat(receipt.Total, 64)
	if err == nil && math.Mod(total, 1.00) == 0 {
		return 50
	}
	return 0
}

// Rule 3: 25 points if total is a multiple of 0.25
type MultipleOfQuarterRule struct{}

func (r MultipleOfQuarterRule) Calculate(receipt models.Receipt) int {
	total, err := strconv.ParseFloat(receipt.Total, 64)
	if err == nil && math.Mod(total, 0.25) == 0 {
		return 25
	}
	return 0
}

// Rule 4: 5 points for every two items
type ItemPairRule struct{}

func (r ItemPairRule) Calculate(receipt models.Receipt) int {
	return (len(receipt.Items) / 2) * 5
}

// Rule 5: Points based on item description length
type ItemDescriptionRule struct{}

func (r ItemDescriptionRule) Calculate(receipt models.Receipt) int {
	points := 0
	for _, item := range receipt.Items {
		desc := strings.TrimSpace(item.ShortDescription)
		if len(desc)%3 == 0 {
			priceFloat, _ := strconv.ParseFloat(item.Price, 64)
			points += int(math.Ceil(priceFloat * 0.2))
		}
	}
	return points
}

// Rule 6: 6 points if the purchase date is odd
type OddPurchaseDayRule struct{}

func (r OddPurchaseDayRule) Calculate(receipt models.Receipt) int {
	date, err := time.Parse("2006-01-02", receipt.PurchaseDate)
	if err == nil && date.Day()%2 != 0 {
		return 6
	}
	return 0
}

// Rule 7: 10 points if purchase time is between 2:00 PM - 4:00 PM
type TimeRangeRule struct{}

func (r TimeRangeRule) Calculate(receipt models.Receipt) int {
	purchaseTime, err := time.Parse("15:04", receipt.PurchaseTime)
	if err == nil && purchaseTime.Hour() == 14 {
		return 10
	}
	return 0
}
