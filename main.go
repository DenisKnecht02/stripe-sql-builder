package builder

import (
	"errors"
	"fmt"
	"strings"
)

type StripeQuery struct {
	ConnectionType ConnectionType
	Custom         map[string]*[]interface{}
	Metadata       *map[string]string
	Entries        *[]QueryEntry[interface{}]
	RawStrings     *[]string
}

type QueryEntry[T interface{}] struct {
	Key      string
	Operator Operator
	Value    T
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

func CreateDefaultQuery(connectionType ConnectionType) StripeQuery {

	var query StripeQuery
	query.ConnectionType = connectionType
	query.Custom = map[string]*[]interface{}{}

	for _, option := range defaultOptions {
		query = option(query)
	}

	return query

}

func NewQuery(connectionType ConnectionType, options ...Option) StripeQuery {

	var query StripeQuery = CreateDefaultQuery(connectionType)

	for _, option := range options {
		query = option(query)
	}

	return query
}

func NewAndQuery(options ...Option) StripeQuery {
	return NewQuery(EnumConnectionType.And, options...)
}

func NewOrQuery(options ...Option) StripeQuery {
	return NewQuery(EnumConnectionType.Or, options...)
}

func NewStringQuery(connectionType ConnectionType, options ...Option) string {
	query := NewQuery(connectionType, options...)
	return query.String()
}

func NewAndStringQuery(options ...Option) string {
	query := NewAndQuery(options...)
	return query.String()
}

func NewOrStringQuery(options ...Option) string {
	query := NewOrQuery(options...)
	return query.String()
}

func addCustom(customMap *map[string]*[]interface{}, key string, values ...interface{}) {

	if (*customMap)[key] == nil {
		(*customMap)[key] = &values
	} else {
		currentValues := *(*customMap)[key]
		currentValues = append(currentValues, values...)
		(*customMap)[key] = &currentValues
	}

}

func WithActive(active bool) Option {
	return func(stripeQuery StripeQuery) StripeQuery {
		addCustom(&stripeQuery.Custom, "active", active)
		return stripeQuery
	}
}

func WithDeleted(deleted bool) Option {
	return func(stripeQuery StripeQuery) StripeQuery {
		addCustom(&stripeQuery.Custom, "deleted", deleted)
		return stripeQuery
	}
}

func WithShippable(shippable bool) Option {
	return func(stripeQuery StripeQuery) StripeQuery {
		addCustom(&stripeQuery.Custom, "shippable", shippable)
		return stripeQuery
	}
}

func WithId(id string) Option {
	return func(stripeQuery StripeQuery) StripeQuery {
		addCustom(&stripeQuery.Custom, "id", id)
		return stripeQuery
	}
}

func WithCustomer(customerId string) Option {
	return func(stripeQuery StripeQuery) StripeQuery {
		addCustom(&stripeQuery.Custom, "customer", customerId)
		return stripeQuery
	}
}

func WithIds(ids ...string) Option {
	return func(stripeQuery StripeQuery) StripeQuery {
		addCustom(&stripeQuery.Custom, "id", ConvertSlice(&ids)...)
		return stripeQuery
	}
}

func WithPriceId(priceId string) Option {
	return func(stripeQuery StripeQuery) StripeQuery {
		addCustom(&stripeQuery.Custom, "default_price.id", priceId)
		return stripeQuery
	}
}

func WithDescription(description string) Option {
	return func(stripeQuery StripeQuery) StripeQuery {
		addCustom(&stripeQuery.Custom, "description", description)
		return stripeQuery
	}
}

func WithType(t string) Option {
	return func(stripeQuery StripeQuery) StripeQuery {
		addCustom(&stripeQuery.Custom, "type", t)
		return stripeQuery
	}
}

func WithCurrency(currency string) Option {
	return func(stripeQuery StripeQuery) StripeQuery {
		addCustom(&stripeQuery.Custom, "currency", currency)
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

func With(key string, values ...interface{}) Option {
	return func(stripeQuery StripeQuery) StripeQuery {
		addCustom(&stripeQuery.Custom, key, values...)
		return stripeQuery
	}
}

func WithRawString(raw string) Option {
	return func(stripeQuery StripeQuery) StripeQuery {

		rawValues := []string{}
		if stripeQuery.RawStrings != nil {
			rawValues = *stripeQuery.RawStrings
		}

		rawValues = append(rawValues, raw)
		stripeQuery.RawStrings = &rawValues
		return stripeQuery
	}
}

func WithEntry(key string, operator Operator, value interface{}) Option {
	return func(stripeQuery StripeQuery) StripeQuery {

		rawValues := []QueryEntry[interface{}]{}
		if stripeQuery.Entries != nil {
			rawValues = *stripeQuery.Entries
		}

		entry := QueryEntry[interface{}]{
			Key:      key,
			Operator: operator,
			Value:    value,
		}

		rawValues = append(rawValues, entry)
		stripeQuery.Entries = &rawValues
		return stripeQuery
	}
}

func WithIsNull(key string) Option {
	return WithEntry(key, EnumOperator.Equals, "null")
}

func Bool(b bool) *bool {
	return &b
}

func String(s string) *string {
	return &s
}

var ErrorParseOperator error = errors.New("INVALID_OPERATOR")
var ErrorParseConnectionType error = errors.New("INVALID_CONNECTION_TYPE")

func ConvertSlice[T interface{}](slice *[]T) []interface{} {

	var interfaceSlice []interface{}
	for _, s := range *slice {
		interfaceSlice = append(interfaceSlice, s)
	}

	return interfaceSlice

}

/* String Methods */

func (stripeQuery *StripeQuery) String() string {

	var queries []string

	if stripeQuery.Metadata != nil {

		for key, value := range *stripeQuery.Metadata {
			queries = append(queries, fmt.Sprintf("metadata['%s']:'%s'", key, value))
		}

	}

	for key, values := range stripeQuery.Custom {

		valueQueries := []string{}

		for _, value := range *values {
			valueQueries = append(valueQueries, fmt.Sprintf("%s:'%v'", key, value))
		}

		queries = append(queries, strings.Join(valueQueries, fmt.Sprintf(" %s ", stripeQuery.ConnectionType.String())))

	}

	if stripeQuery.RawStrings != nil {
		queries = append(queries, *stripeQuery.RawStrings...)
	}

	if stripeQuery.Entries != nil {

		for _, entry := range *stripeQuery.Entries {
			queries = append(queries, entry.String())
		}

	}

	return strings.Join(queries, fmt.Sprintf(" %s ", stripeQuery.ConnectionType.String()))

}

func (entry *QueryEntry[T]) String() string {
	if entry.Operator == EnumOperator.NotEquals {
		return fmt.Sprintf("%s%s%s'%v'", entry.Operator.Operator(), entry.Key, EnumOperator.Equals.Operator(), entry.Value)
	} else if entry.Operator == EnumOperator.NotLike {
		return fmt.Sprintf("%s%s%s'%v'", entry.Operator.Operator(), entry.Key, EnumOperator.Like.Operator(), entry.Value)
	} else {
		return fmt.Sprintf("%s%s'%v'", entry.Key, entry.Operator.Operator(), entry.Value)
	}
}

/* Enum: Connection Type */

type ConnectionType string

type connectionTypeList struct {
	Unknown ConnectionType
	And     ConnectionType
	Or      ConnectionType
}

var EnumConnectionType = &connectionTypeList{
	Unknown: "UNKNOWN",
	And:     "AND",
	Or:      "OR",
}

var connectionTypeMap = map[string]ConnectionType{
	"UNKNOWN": EnumConnectionType.Unknown,
	"AND":     EnumConnectionType.And,
	"OR":      EnumConnectionType.Or,
}

func ParseStringToConnectionType(str string) (ConnectionType, error) {
	connectionType, ok := connectionTypeMap[strings.ToLower(str)]
	if ok {
		return connectionType, nil
	} else {
		return connectionType, ErrorParseConnectionType
	}
}

func (connectionType ConnectionType) String() string {
	switch connectionType {
	case EnumConnectionType.And:
		return "AND"
	case EnumConnectionType.Or:
		return "OR"
	}
	return "UNKNOWN"
}

/* Enum: Operator */

type Operator string

type operatorList struct {
	Unknown          Operator
	Equals           Operator
	GreaterThan      Operator
	LessThan         Operator
	GreaterEqualThan Operator
	LessEqualThan    Operator
	NotEquals        Operator
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
	NotEquals:        "not_equals",
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
	"not_equals":         EnumOperator.NotEquals,
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
	case EnumOperator.NotEquals:
		return "not_equals"
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
	case EnumOperator.NotEquals:
		return "-"
	case EnumOperator.Like:
		return "~"
	case EnumOperator.NotLike:
		return "-"
	}
	return "unknown"
}
