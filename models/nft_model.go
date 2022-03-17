package models

type NFTProfile struct {
	ID             int64  `gorm:"primaryKey;autoIncrement"`
	TokenBasicName string `gorm:"uniqueIndex:idx_name_token;size:64;not null"`
	NftTokenID     string `gorm:"uniqueIndex:idx_name_token;type:varchar(64);not null"`
	Name           string `gorm:"size:64;not null"`
	URL            string `gorm:"size:64;not null"`
	Image          string `gorm:"size:64;not null"`
	Description    string `gorm:"type:varchar(256)"`
	Text           string `gorm:"type:text"`
}
