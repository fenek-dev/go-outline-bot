package models

import "time"

type SubscriptionStatus string

const (
	SubscriptionStatusActive           SubscriptionStatus = "active"
	SubscriptionStatusExpired          SubscriptionStatus = "expired"
	SubscriptionStatusBandwidthReached SubscriptionStatus = "bandwidth_reached"
)

type Subscription struct {
	ID             uint64    `json:"id" db:"id"`
	UserID         uint64    `json:"user_id" db:"user_id"`
	ServerID       uint64    `json:"server_id" db:"server_id"`
	TariffID       uint64    `json:"tariff_id" db:"tariff_id"`
	InitialPrice   uint32    `json:"initial_price" db:"initial_price"`
	BandwidthSpent uint64    `json:"bandwidth_spent" db:"bandwidth_spent"`
	BandwidthTotal *uint64   `json:"bandwidth_total" db:"bandwidth_total"` // Optional field from subscription join tariff
	KeyUUID        string    `json:"-" db:"key_uucid"`
	AccessUrl      string    `json:"access_url" db:"access_url"`
	ServerIP       string    `json:"server_ip" db:"server_ip"`
	ServerPort     uint16    `json:"server_port" db:"server_port"`
	Password       string    `json:"password" db:"password"`
	Method         string    `json:"method" db:"method"`
	IsTrial        bool      `json:"is_trial" db:"is_trial"`
	AutoProlong    bool      `json:"auto_prolong" db:"auto_prolong"`
	ExpiredAt      time.Time `json:"expired_at" db:"expired_at"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`

	Status SubscriptionStatus `json:"status" db:"status"`
}
