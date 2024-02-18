package store

type Asset struct {
	ID       string `gorm:"primarykey"`
	MimeType string
	Blob     []byte
}

func (db *Database) CreateAsset(ID string, mimeType string, blob []byte) {
	add := Asset{ID: ID, MimeType: mimeType, Blob: blob}
	db.Client.Create(&add)
}

func (db *Database) GetAsset(ID string) *Asset {
	var asset Asset
	db.Client.Where("id = ?", ID).First(&asset)
	return &asset
}
