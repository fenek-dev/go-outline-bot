package notifier

type NotifyDTO struct {
	UserID         uint64
	SubscriptionID uint64
	TariffName     string
}

type NotifyBandwidthDTO struct {
	NotifyDTO
	BandwidthSpent uint64
	BandwidthTotal uint64
}

type NotifyPartnerDTO struct {
	UserID            uint64
	RecipientUsername string
	Amount            uint64
	RewardAmount      uint64
}
