package builder

import (
	"errors"
	"fmt"
	"strings"
)

type StripeQuery struct {
	ConnectionType ConnectionType
	Collections    []QueryCollection
}

type QueryCollection struct {
	ConnectionType ConnectionType
	Custom         map[string]interface{}
	Metadata       *map[string]string
	Entries        *[]QueryEntry[interface{}]
	RawStrings     *[]string
}

type QueryEntry[T interface{}] struct {
	Key      string
	Operator Operator
	Value    T
}

type Option func(QueryCollection) QueryCollection

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

func CreateDefaultCollection(connectionType ConnectionType) QueryCollection {

	var collection QueryCollection
	collection.ConnectionType = connectionType
	collection.Custom = map[string]interface{}{}

	for _, option := range defaultOptions {
		collection = option(collection)
	}

	return collection

}

func Options(options ...Option) []Option {
	return options
}

func NewQuery(connectionType ConnectionType, collections ...QueryCollection) StripeQuery {

	query := StripeQuery{
		ConnectionType: connectionType,
		Collections:    collections,
	}

	return query
}

func NewAndQuery(options ...Option) StripeQuery {

	var query StripeQuery = StripeQuery{}
	query.Collections = append(query.Collections, And(options...))

	return query

}

func NewOrQuery(options ...Option) StripeQuery {

	var query StripeQuery = StripeQuery{}
	query.Collections = append(query.Collections, Or(options...))

	return query

}

func NewStringQuery(connectionType ConnectionType, collections ...QueryCollection) string {
	query := NewQuery(connectionType, collections...)
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

func And(options ...Option) QueryCollection {
	return connect(EnumConnectionType.And, options...)
}

func Or(options ...Option) QueryCollection {
	return connect(EnumConnectionType.Or, options...)
}

func connect(connectionType ConnectionType, options ...Option) QueryCollection {

	collection := CreateDefaultCollection(connectionType)

	for _, option := range options {
		collection = option(collection)
	}

	return collection

}

func WithActive(active bool) Option {
	return func(stripeQuery QueryCollection) QueryCollection {
		stripeQuery.Custom["active"] = active
		return stripeQuery
	}
}

func WithDeleted(deleted bool) Option {
	return func(stripeQuery QueryCollection) QueryCollection {
		stripeQuery.Custom["deleted"] = deleted
		return stripeQuery
	}
}

func WithShippable(shippable bool) Option {
	return func(stripeQuery QueryCollection) QueryCollection {
		stripeQuery.Custom["shippable"] = shippable
		return stripeQuery
	}
}

func WithId(id string) Option {
	return func(stripeQuery QueryCollection) QueryCollection {
		stripeQuery.Custom["id"] = id
		return stripeQuery
	}
}

func WithPriceId(priceId string) Option {
	return func(stripeQuery QueryCollection) QueryCollection {
		stripeQuery.Custom["default_price.id"] = priceId
		return stripeQuery
	}
}

func WithDescription(description string) Option {
	return func(stripeQuery QueryCollection) QueryCollection {
		stripeQuery.Custom["description"] = description
		return stripeQuery
	}
}

func WithType(t string) Option {
	return func(stripeQuery QueryCollection) QueryCollection {
		stripeQuery.Custom["type"] = t
		return stripeQuery
	}
}

func WithCurrency(currency string) Option {
	return func(stripeQuery QueryCollection) QueryCollection {
		stripeQuery.Custom["currency"] = currency
		return stripeQuery
	}
}

func WithMetadata(key string, value string) Option {
	return func(stripeQuery QueryCollection) QueryCollection {

		if stripeQuery.Metadata == nil {
			stripeQuery.Metadata = &map[string]string{}
		}

		(*stripeQuery.Metadata)[key] = value

		return stripeQuery
	}
}

func WithMetadataMap(metadataMap map[string]string) Option {
	return func(stripeQuery QueryCollection) QueryCollection {

		if stripeQuery.Metadata == nil {
			stripeQuery.Metadata = &map[string]string{}
		}

		for key, value := range metadataMap {
			(*stripeQuery.Metadata)[key] = value
		}

		return stripeQuery
	}
}

func With(key string, value interface{}) Option {
	return func(stripeQuery QueryCollection) QueryCollection {

		stripeQuery.Custom[key] = value
		return stripeQuery
	}
}

func WithRawString(raw string) Option {
	return func(stripeQuery QueryCollection) QueryCollection {

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
	return func(stripeQuery QueryCollection) QueryCollection {

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

/* String Methods */

func (stripeQuery *StripeQuery) String() string {

	queryStrings := []string{}

	for _, collection := range *&stripeQuery.Collections {

		var queries []string

		if collection.Metadata != nil {

			for key, value := range *collection.Metadata {
				queries = append(queries, fmt.Sprintf("metadata['%s']:'%s'", key, value))
			}

		}

		for key, value := range collection.Custom {
			queries = append(queries, fmt.Sprintf("%s:'%v'", key, value))
		}

		if collection.RawStrings != nil {
			queries = append(queries, *collection.RawStrings...)
		}

		if collection.Entries != nil {

			for _, entry := range *collection.Entries {
				queries = append(queries, entry.String())
			}

		}

		queryStrings = append(queryStrings, fmt.Sprintf("(%s)", strings.Join(queries, fmt.Sprintf(" %s ", collection.ConnectionType.String()))))

	}

	return strings.Join(queryStrings, fmt.Sprintf(" %s ", stripeQuery.ConnectionType))

}

func (collection *QueryCollection) String() string {

	var queries []string

	if collection.Metadata != nil {

		for key, value := range *collection.Metadata {
			queries = append(queries, fmt.Sprintf("metadata['%s']:'%s'", key, value))
		}

	}

	for key, value := range collection.Custom {
		queries = append(queries, fmt.Sprintf("%s:'%v'", key, value))
	}

	if collection.RawStrings != nil {
		queries = append(queries, *collection.RawStrings...)
	}

	if collection.Entries != nil {

		for _, entry := range *collection.Entries {
			queries = append(queries, entry.String())
		}

	}

	return fmt.Sprintf("(%s)", strings.Join(queries, fmt.Sprintf(" %s ", collection.ConnectionType.String())))

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
