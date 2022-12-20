package main

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

func NewQuery(options ...Option) StripeQuery {

	query := StripeQuery{
		Active: Bool(true),
		Custom: &map[string]interface{}{
			"active": true,
		},
	}

	for _, option := range options {
		query = option(query)
	}

	return query

}

func WithActive(active bool) Option {
	return func(stripeQuery StripeQuery) StripeQuery {
		stripeQuery.Active = &active
		(*stripeQuery.Custom)["active"] = active
		return stripeQuery
	}
}

func WithDeleted(deleted bool) Option {
	return func(stripeQuery StripeQuery) StripeQuery {
		stripeQuery.Deleted = &deleted
		(*stripeQuery.Custom)["deleted"] = deleted
		return stripeQuery
	}
}

func WithShippable(shippable bool) Option {
	return func(stripeQuery StripeQuery) StripeQuery {
		stripeQuery.Shippable = &shippable
		(*stripeQuery.Custom)["shippable"] = shippable
		return stripeQuery
	}
}

func WithId(id string) Option {
	return func(stripeQuery StripeQuery) StripeQuery {
		stripeQuery.ID = &id
		(*stripeQuery.Custom)["id"] = id
		return stripeQuery
	}
}

func WithPriceId(priceId string) Option {
	return func(stripeQuery StripeQuery) StripeQuery {
		stripeQuery.PriceId = &priceId
		(*stripeQuery.Custom)["default_price.id"] = priceId
		return stripeQuery
	}
}

func WithDescription(description string) Option {
	return func(stripeQuery StripeQuery) StripeQuery {
		stripeQuery.Description = &description
		(*stripeQuery.Custom)["description"] = description
		return stripeQuery
	}
}

func WithCreated(created int) Option {
	return func(stripeQuery StripeQuery) StripeQuery {
		stripeQuery.Created = &created
		(*stripeQuery.Custom)["created"] = created
		return stripeQuery
	}
}

func WithType(t string) Option {
	return func(stripeQuery StripeQuery) StripeQuery {
		stripeQuery.Type = &t
		(*stripeQuery.Custom)["type"] = t
		return stripeQuery
	}
}

func WithCurrency(currency string) Option {
	return func(stripeQuery StripeQuery) StripeQuery {
		stripeQuery.Currency = &currency
		(*stripeQuery.Custom)["currency"] = currency
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
