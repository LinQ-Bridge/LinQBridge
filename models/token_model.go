package models

import "math/big"

const (
	TokenTypeErc20 uint8 = iota
	TokenTypeErc721
)

type TokenBasic struct {
	ID              int64          `gorm:"primaryKey;autoIncrement"`
	Name            string         `gorm:"uniqueIndex;size:64;not null"`
	Precision       uint64         `gorm:"type:bigint(20);not null"`
	Price           int64          `gorm:"size:64;not null"`
	ChainID         uint64         `gorm:"type:bigint(20);not null"`
	Ind             uint64         `gorm:"type:bigint(20);not null"`
	Time            int64          `gorm:"type:bigint(20);not null"`
	Property        int64          `gorm:"type:bigint(20);not null"`
	Standard        uint8          `gorm:"type:int(8);not null"` // 0: erc20， 1: erc721
	Meta            string         `gorm:"type:varchar(128)"`
	TotalAmount     *BigInt        `gorm:"type:varchar(64)"`
	TotalCount      uint64         `gorm:"type:bigint(20)"`
	StatsUpdateTime int64          `gorm:"type:bigint(20)"`
	SocialTwitter   string         `gorm:"type:varchar(256)"`
	SocialTelegram  string         `gorm:"type:varchar(256)"`
	SocialWebsite   string         `gorm:"type:varchar(256)"`
	SocialOther     string         `gorm:"type:varchar(256)"`
	MetaFetcherType int            `gorm:"type:int(8);not null"` // nft meta profile fetcher type, e.g: unknown 0, opensea: 1, standard: 2,
	PriceMarkets    []*PriceMarket `gorm:"foreignKey:TokenBasicName;references:Name"`
	Tokens          []*Token       `gorm:"foreignKey:TokenBasicName;references:Name"`
}

type PriceMarket struct {
	ID             int64       `gorm:"primaryKey;autoIncrement"`
	TokenBasicName string      `gorm:"uniqueIndex:idx_tokenmarket;size:64;not null"`
	MarketName     string      `gorm:"uniqueIndex:idx_tokenmarket;size:64;not null"`
	Name           string      `gorm:"size:64;not null"`
	Price          int64       `gorm:"type:bigint(20);not null"`
	Ind            uint64      `gorm:"type:bigint(20);not null"`
	Time           int64       `gorm:"type:bigint(20);not null"`
	TokenBasic     *TokenBasic `gorm:"foreignKey:TokenBasicName;references:Name"`
}

type ChainFee struct {
	ID             int64       `gorm:"primaryKey;autoIncrement"`
	ChainID        uint64      `gorm:"uniqueIndex;type:bigint(20);not null"`
	TokenBasicName string      `gorm:"size:64;not null"`
	TokenBasic     *TokenBasic `gorm:"foreignKey:TokenBasicName;references:Name"`
	MaxFee         *BigInt     `gorm:"type:varchar(64);not null"`
	MinFee         *BigInt     `gorm:"type:varchar(64);not null"`
	ProxyFee       *BigInt     `gorm:"type:varchar(64);not null"`
	Ind            uint64      `gorm:"type:bigint(20);not null"`
	Time           int64       `gorm:"type:bigint(20);not null"`
}

type Token struct {
	ID              int64       `gorm:"primaryKey;autoIncrement"`
	Hash            string      `gorm:"uniqueIndex:idx_token;size:66;not null"`
	ChainID         uint64      `gorm:"uniqueIndex:idx_token;type:bigint(20);not null"`
	Name            string      `gorm:"size:64;not null"`
	Precision       uint64      `gorm:"type:bigint(20);not null"`
	TokenBasicName  string      `gorm:"size:64;not null"`
	Property        int64       `gorm:"type:bigint(20);not null"`
	Standard        uint8       `gorm:"type:int(8);not null"`
	TokenType       string      `gorm:"type:varchar(32)"`
	AvailableAmount *BigInt     `gorm:"type:varchar(64)"`
	TokenBasic      *TokenBasic `gorm:"foreignKey:TokenBasicName;references:Name"`
	TokenMaps       []*TokenMap `gorm:"foreignKey:SrcTokenHash,SrcChainID;references:Hash,ChainID"`
}

type TokenStatistic struct {
	ID             int64   `gorm:"primaryKey;autoIncrement"`
	Hash           string  `gorm:"uniqueIndex:idx_token;size:66;not null"`
	ChainID        uint64  `gorm:"uniqueIndex:idx_token;type:bigint(20);not null"`
	InCounter      int64   `gorm:"type:bigint(20)"`
	InAmount       *BigInt `gorm:"type:varchar(64)"`
	InAmountBtc    *BigInt `gorm:"type:varchar(64)"`
	InAmountUsd    *BigInt `gorm:"type:varchar(64)"`
	OutCounter     int64   `gorm:"type:bigint(20)"`
	OutAmount      *BigInt `gorm:"type:varchar(64)"`
	OutAmountBtc   *BigInt `gorm:"type:varchar(64)"`
	OutAmountUsd   *BigInt `gorm:"type:varchar(64)"`
	LastInCheckID  int64   `gorm:"type:int;not null"`
	LastOutCheckID int64   `gorm:"type:int;not null"`
	Token          *Token  `gorm:"foreignKey:Hash,ChainID;references:Hash,ChainID"`
}

type TokenMap struct {
	ID           int64  `gorm:"primaryKey;autoIncrement"`
	SrcChainID   uint64 `gorm:"uniqueIndex:idx_token_map;type:bigint(20);not null"`
	SrcTokenHash string `gorm:"uniqueIndex:idx_token_map;size:66;not null"`
	DstChainID   uint64 `gorm:"uniqueIndex:idx_token_map;type:bigint(20);not null"`
	DstTokenHash string `gorm:"uniqueIndex:idx_token_map;size:66;not null"`
	SrcToken     *Token `gorm:"foreignKey:SrcTokenHash,SrcChainID;references:Hash,ChainID"`
	DstToken     *Token `gorm:"foreignKey:DstTokenHash,DstChainID;references:Hash,ChainID"`
	Standard     uint8  `gorm:"type:int(8);not null"`
	Property     int64  `gorm:"type:bigint(20);not null"`
}

type WrapperTransactionWithToken struct {
	ID           int64   `gorm:"primaryKey;autoIncrement"`
	Hash         string  `gorm:"uniqueIndex;size:66;not null"`
	User         string  `gorm:"size:64"`
	SrcChainID   uint64  `gorm:"type:bigint(20);not null"`
	BlockHeight  uint64  `gorm:"type:bigint(20);not null"`
	Time         uint64  `gorm:"type:bigint(20);not null"`
	DstChainID   uint64  `gorm:"type:bigint(20);not null"`
	DstUser      string  `gorm:"type:varchar(66);not null"`
	ServerID     uint64  `gorm:"type:bigint(20);not null"`
	FeeTokenHash string  `gorm:"size:66;not null"`
	FeeToken     *Token  `gorm:"foreignKey:FeeTokenHash,SrcChainID;references:Hash,ChainID"`
	FeeAmount    *BigInt `gorm:"type:varchar(64);not null"`
	Status       uint64  `gorm:"type:bigint(20);not null"`
}

type CheckFee struct {
	ChainID     uint64
	Hash        string
	PayState    int
	Amount      *big.Float
	MinProxyFee *big.Float
}

type TimeStatistic struct {
	ID         int64  `gorm:"primaryKey;autoIncrement"`
	SrcChainID uint64 `gorm:"uniqueIndex:idx_chains;type:bigint(20);not null"`
	DstChainID uint64 `gorm:"uniqueIndex:idx_chains;type:bigint(20);not null"`
	Time       uint64 `gorm:"type:bigint(20);not null"`
}
