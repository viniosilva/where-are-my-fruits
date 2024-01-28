package models

import "github.com/shopspring/decimal"

type BucketFruits struct {
	ID          int64
	Name        string
	Capacity    int
	TotalFruits int64
	TotalPrice  decimal.Decimal
	Percent     decimal.Decimal
}

// Refers: https://martinfowler.com/bliki/DDD_Aggregate.html
