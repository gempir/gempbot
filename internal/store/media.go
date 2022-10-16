package store

type MediaPlayer struct {
	ChannelTwitchId string `gorm:"primaryKey"`
	CurrentTime     float32
}
