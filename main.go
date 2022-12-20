package builder

import (
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
}

func (stripeQuery *StripeQuery) ToString() string {

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

	return strings.Join(queries, " AND ")

}

type Option func(StripeQuery) StripeQuery

var defaultOptions []Option = []Option{
	WithActive(true),
}

func ResetDefaultQueryOptions() {
	defaultOptions = []Option{
		WithActive(true),
	}
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

	return query.ToString()

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

func Bool(b bool) *bool {
	return &b
}

func String(s string) *string {
	return &s
}
