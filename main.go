package builder

import (
	"errors"
	"fmt"
	"strings"
)

type StripeQuery struct {
	ID          *string
	Active      *bool
	Deleted     *bool
	Shippable   *bool
	Metadata    *map[string]string
	Created     *int
	Description *string
	Type        *string
	Currency    *string
	PriceId     *string
	Custom      *map[string]interface{}
	Raw         *[]string
	RawEntries  *[]QueryEntry[interface{}]
}

type QueryEntry[T interface{}] struct {
	Key      string
	Operator Operator
	Value    T
}

func (entry *QueryEntry[T]) String() string {
	if entry.Operator == EnumOperator.NotEqual {
		return fmt.Sprintf("%s%s%s'%v'", entry.Operator.Operator(), entry.Key, EnumOperator.Equals.Operator(), entry.Value)
	} else if entry.Operator == EnumOperator.NotLike {
		return fmt.Sprintf("%s%s%s'%v'", entry.Operator.Operator(), entry.Key, EnumOperator.Like.Operator(), entry.Value)
	} else {
		return fmt.Sprintf("%s%s'%v'", entry.Key, entry.Operator.Operator(), entry.Value)
	}
}

func (stripeQuery *StripeQuery) String() string {

	var queries []string

	if stripeQuery.ID != nil {
		queries = append(queries, fmt.Sprintf("%s:'%s'", "id", *stripeQuery.ID))
	}

	if stripeQuery.Active != nil {
		queries = append(queries, fmt.Sprintf("%s:'%t'", "active", *stripeQuery.Active))
	}

	if stripeQuery.Deleted != nil {
		queries = append(queries, fmt.Sprintf("%s:'%t'", "deleted", *stripeQuery.Deleted))
	}

	if stripeQuery.Shippable != nil {
		queries = append(queries, fmt.Sprintf("%s:'%t'", "shippable", *stripeQuery.Shippable))
	}

	if stripeQuery.Created != nil {
		queries = append(queries, fmt.Sprintf("%s:'%d'", "created", *stripeQuery.Created))
	}

	if stripeQuery.Description != nil {
		queries = append(queries, fmt.Sprintf("%s:'%s'", "description", *stripeQuery.Description))
	}

	if stripeQuery.Type != nil {
		queries = append(queries, fmt.Sprintf("%s:'%s'", "type", *stripeQuery.Type))
	}

	if stripeQuery.Currency != nil {
		queries = append(queries, fmt.Sprintf("%s:'%s'", "currency", *stripeQuery.Currency))
	}

	if stripeQuery.PriceId != nil {
		queries = append(queries, fmt.Sprintf("%s:'%s'", "default_price.id", *stripeQuery.PriceId))
	}

	if stripeQuery.Metadata != nil {

		for key, value := range *stripeQuery.Metadata {
			queries = append(queries, fmt.Sprintf("metadata['%s']:'%s'", key, value))
		}

	}

	if stripeQuery.Custom != nil {

		for key, value := range *stripeQuery.Custom {
			queries = append(queries, fmt.Sprintf("%s:'%v'", key, value))
		}

	}

	if stripeQuery.Raw != nil {
		queries = append(queries, *stripeQuery.Raw...)
	}

	if stripeQuery.RawEntries != nil {

		for _, entry := range *stripeQuery.RawEntries {
			queries = append(queries, entry.String())
		}

	}

	return strings.Join(queries, " AND ")

}

type Option func(StripeQuery) StripeQuery

var defaultOptions []Option = []Option{}

func ResetDefaultQueryOptions() {
	defaultOptions = []Option{}
}

func GetDefaultQueryOptions() []Option {
	return defaultOptions
}

func SetDefaultQueryOptions(options ...Option) {
	defaultOptions = options
}

func CreateDefaultQuery() StripeQuery {

	var query StripeQuery

	for _, option := range defaultOptions {
		query = option(query)
	}

	return query

}

func NewQuery(options ...Option) StripeQuery {

	query := CreateDefaultQuery()

	for _, option := range options {
		query = option(query)
	}

	return query

}

func NewStringQuery(options ...Option) string {

	query := CreateDefaultQuery()

	for _, option := range options {
		query = option(query)
	}

	return query.String()

}

func WithActive(active bool) Option {
	return func(stripeQuery StripeQuery) StripeQuery {
		stripeQuery.Active = &active
		return stripeQuery
	}
}

func WithDeleted(deleted bool) Option {
	return func(stripeQuery StripeQuery) StripeQuery {
		stripeQuery.Deleted = &deleted
		return stripeQuery
	}
}

func WithShippable(shippable bool) Option {
	return func(stripeQuery StripeQuery) StripeQuery {
		stripeQuery.Shippable = &shippable
		return stripeQuery
	}
}

func WithId(id string) Option {
	return func(stripeQuery StripeQuery) StripeQuery {
		stripeQuery.ID = &id
		return stripeQuery
	}
}

func WithPriceId(priceId string) Option {
	return func(stripeQuery StripeQuery) StripeQuery {
		stripeQuery.PriceId = &priceId
		return stripeQuery
	}
}

func WithDescription(description string) Option {
	return func(stripeQuery StripeQuery) StripeQuery {
		stripeQuery.Description = &description
		return stripeQuery
	}
}

func WithCreated(created int) Option {
	return func(stripeQuery StripeQuery) StripeQuery {
		stripeQuery.Created = &created
		return stripeQuery
	}
}

func WithType(t string) Option {
	return func(stripeQuery StripeQuery) StripeQuery {
		stripeQuery.Type = &t
		return stripeQuery
	}
}

func WithCurrency(currency string) Option {
	return func(stripeQuery StripeQuery) StripeQuery {
		stripeQuery.Currency = &currency
		return stripeQuery
	}
}

func WithMetadata(key string, value string) Option {
	return func(stripeQuery StripeQuery) StripeQuery {

		if stripeQuery.Metadata == nil {
			stripeQuery.Metadata = &map[string]string{}
		}

		(*stripeQuery.Metadata)[key] = value

		return stripeQuery
	}
}

func WithMetadataMap(metadataMap map[string]string) Option {
	return func(stripeQuery StripeQuery) StripeQuery {

		if stripeQuery.Metadata == nil {
			stripeQuery.Metadata = &map[string]string{}
		}

		for key, value := range metadataMap {
			(*stripeQuery.Metadata)[key] = value
		}

		return stripeQuery
	}
}

func WithCustom(key string, value string) Option {
	return func(stripeQuery StripeQuery) StripeQuery {

		if stripeQuery.Custom == nil {
			stripeQuery.Custom = &map[string]interface{}{}
		}

		(*stripeQuery.Custom)[key] = value
		return stripeQuery
	}
}

func WithRawString(raw string) Option {
	return func(stripeQuery StripeQuery) StripeQuery {

		rawValues := []string{}
		if stripeQuery.Raw != nil {
			rawValues = *stripeQuery.Raw
		}

		rawValues = append(rawValues, raw)
		stripeQuery.Raw = &rawValues
		return stripeQuery
	}
}

func WithRawEntry(key string, operator Operator, value interface{}) Option {
	return func(stripeQuery StripeQuery) StripeQuery {

		rawValues := []QueryEntry[interface{}]{}
		if stripeQuery.Raw != nil {
			rawValues = *stripeQuery.RawEntries
		}

		entry := QueryEntry[interface{}]{
			Key:      key,
			Operator: operator,
			Value:    value,
		}

		rawValues = append(rawValues, entry)
		stripeQuery.RawEntries = &rawValues
		return stripeQuery
	}
}

func Bool(b bool) *bool {
	return &b
}

func String(s string) *string {
	return &s
}

var ErrorParseOperator error = errors.New("INVALID_OPERATOR")

/* Enum: Operator */

type Operator string

type operatorList struct {
	Unknown          Operator
	Equals           Operator
	GreaterThan      Operator
	LessThan         Operator
	GreaterEqualThan Operator
	LessEqualThan    Operator
	NotEqual         Operator
	Like             Operator
	NotLike          Operator
}

var EnumOperator = &operatorList{
	Unknown:          "unknown",
	Equals:           "equals",
	GreaterThan:      "greater_than",
	LessThan:         "less_than",
	GreaterEqualThan: "greater_equal_than",
	LessEqualThan:    "less_equal_than",
	NotEqual:         "not_equal",
	Like:             "like",
	NotLike:          "not_like",
}

var operatorMap = map[string]Operator{
	"unknown":            EnumOperator.Unknown,
	"equals":             EnumOperator.Equals,
	"greater_than":       EnumOperator.GreaterThan,
	"less_than":          EnumOperator.LessThan,
	"greater_equal_than": EnumOperator.GreaterEqualThan,
	"less_equal_than":    EnumOperator.LessEqualThan,
	"not_equal":          EnumOperator.NotEqual,
	"like":               EnumOperator.Like,
	"not_like":           EnumOperator.NotLike,
}

func ParseStringToOperator(str string) (Operator, error) {
	operator, ok := operatorMap[strings.ToLower(str)]
	if ok {
		return operator, nil
	} else {
		return operator, ErrorParseOperator
	}
}

func (operator Operator) String() string {
	switch operator {
	case EnumOperator.Equals:
		return "equals"
	case EnumOperator.GreaterThan:
		return "greater_than"
	case EnumOperator.LessThan:
		return "less_than"
	case EnumOperator.GreaterEqualThan:
		return "greater_equal_than"
	case EnumOperator.LessEqualThan:
		return "less_equal_than"
	case EnumOperator.NotEqual:
		return "not_equal"
	case EnumOperator.Like:
		return "like"
	case EnumOperator.NotLike:
		return "not_like"
	}
	return "unknown"
}

func (operator Operator) Operator() string {
	switch operator {
	case EnumOperator.Equals:
		return ":"
	case EnumOperator.GreaterThan:
		return ">"
	case EnumOperator.LessThan:
		return "<"
	case EnumOperator.GreaterEqualThan:
		return ">="
	case EnumOperator.LessEqualThan:
		return "<="
	case EnumOperator.NotEqual:
		return "-"
	case EnumOperator.Like:
		return "~"
	case EnumOperator.NotLike:
		return "-"
	}
	return "unknown"
}
