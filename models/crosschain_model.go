package models

type Chain struct {
	ID                  int64  `gorm:"primaryKey;autoIncrement"`
	ChainID             uint64 `gorm:"uniqueIndex;type:bigint(20);not null"`
	Name                string `gorm:"type:varchar(32)"`
	Height              uint64 `gorm:"type:bigint(20);not null"`
	HeightSwap          uint64 `gorm:"type:bigint(20);not null"`
	BackwardBlockNumber uint64 `gorm:"type:bigint(20);not null"`
}

type ChainStatistic struct {
	ID             int64  `gorm:"primaryKey;autoIncrement"`
	ChainID        uint64 `gorm:"uniqueIndex;type:bigint(20);not null"`
	Addresses      int64  `gorm:"type:bigint(20);not null"`
	In             int64  `gorm:"type:bigint(20);not null"`
	Out            int64  `gorm:"type:bigint(20);not null"`
	LastInCheckID  int64  `gorm:"type:int"`
	LastOutCheckID int64  `gorm:"type:int"`
}

type SrcTransaction struct {
	ID          int64        `gorm:"primaryKey;autoIncrement"`
	Hash        string       `gorm:"uniqueIndex;size:66;not null"`
	ChainID     uint64       `gorm:"type:bigint(20);not null"`
	Standard    uint8        `gorm:"type:int(8);not null"`
	State       uint64       `gorm:"type:bigint(20);not null"`
	Time        uint64       `gorm:"type:bigint(20);not null"`
	Fee         *BigInt      `gorm:"type:varchar(64);not null"`
	Height      uint64       `gorm:"type:bigint(20);not null"`
	User        string       `gorm:"type:varchar(66);not null"`
	DstChainID  uint64       `gorm:"type:bigint(20);not null"`
	Contract    string       `gorm:"type:varchar(66);not null"`
	Key         string       `gorm:"type:text;not null"`
	Param       string       `gorm:"type:text;not null"`
	SrcTransfer *SrcTransfer `gorm:"foreignKey:TxHash;references:Hash"`
	SrcSwap     *SrcSwap     `gorm:"foreignKey:TxHash;references:Hash"`
}

type SrcTransfer struct {
	ID         int64   `gorm:"primaryKey;autoIncrement"`
	TxHash     string  `gorm:"uniqueIndex;size:66;not null"`
	ChainID    uint64  `gorm:"type:bigint(20);not null"`
	Standard   uint8   `gorm:"type:int(8);not null"`
	Time       uint64  `gorm:"type:bigint(20);not null"`
	Asset      string  `gorm:"type:varchar(120);not null"`
	From       string  `gorm:"type:varchar(66);not null"`
	To         string  `gorm:"type:varchar(66);not null"`
	TokenID    *BigInt `gorm:"type:varchar(86);not null"`
	DstChainID uint64  `gorm:"type:bigint(20);not null"`
	DstAsset   string  `gorm:"type:varchar(120);not null"`
	DstUser    string  `gorm:"type:varchar(66);not null"`
}

type SrcSwap struct {
	ID         int64   `gorm:"primaryKey;autoIncrement"`
	TxHash     string  `gorm:"uniqueIndex;size:66;not null"`
	ChainID    uint64  `gorm:"type:bigint(20);not null"`
	Time       uint64  `gorm:"type:bigint(20);not null"`
	Asset      string  `gorm:"type:varchar(120);not null"`
	From       string  `gorm:"type:varchar(66);not null"`
	To         string  `gorm:"type:varchar(66);not null"`
	TokenID    *BigInt `gorm:"type:varchar(86);not null"`
	PoolID     uint64  `gorm:"type:bigint(20);not null"`
	DstChainID uint64  `gorm:"type:bigint(20);not null"`
	DstAsset   string  `gorm:"type:varchar(120);not null"`
	DstUser    string  `gorm:"type:varchar(66);not null"`
	Type       uint64  `gorm:"type:bigint(20);not null"`
}

type PolyTransaction struct {
	ID         int64   `gorm:"primaryKey;autoIncrement"`
	Hash       string  `gorm:"uniqueIndex;size:66;not null"`
	ChainID    uint64  `gorm:"type:bigint(20);not null"`
	State      uint64  `gorm:"type:bigint(20);not null"`
	Time       uint64  `gorm:"type:bigint(20);not null"`
	Fee        *BigInt `gorm:"type:varchar(64);not null"`
	Height     uint64  `gorm:"type:bigint(20);not null"`
	SrcChainID uint64  `gorm:"type:bigint(20);not null"`
	SrcHash    string  `gorm:"index;size:66;not null"`
	DstChainID uint64  `gorm:"type:bigint(20);not null"`
	Key        string  `gorm:"type:varchar(8192);not null"`
}

type PolySrcRelation struct {
	SrcHash         string
	SrcTransaction  *SrcTransaction `gorm:"foreignKey:SrcHash;references:Hash"`
	PolyHash        string
	PolyTransaction *PolyTransaction `gorm:"foreignKey:PolyHash;references:Hash"`
}

type DstTransaction struct {
	ID          int64        `gorm:"primaryKey;autoIncrement"`
	Hash        string       `gorm:"uniqueIndex;size:66;not null"`
	ChainID     uint64       `gorm:"type:bigint(20);not null"`
	Standard    uint8        `gorm:"type:int(8);not null"`
	State       uint64       `gorm:"type:bigint(20);not null"`
	Time        uint64       `gorm:"type:bigint(20);not null"`
	Fee         *BigInt      `gorm:"type:varchar(64);not null"`
	Height      uint64       `gorm:"type:bigint(20);not null"`
	SrcChainID  uint64       `gorm:"type:bigint(20);not null"`
	Contract    string       `gorm:"type:varchar(66);not null"`
	PolyHash    string       `gorm:"index;size:66;not null"`
	DstTransfer *DstTransfer `gorm:"foreignKey:TxHash;references:Hash"`
	DstSwap     *DstSwap     `gorm:"foreignKey:TxHash;references:Hash"`
}

type DstTransfer struct {
	ID       int64   `gorm:"primaryKey;autoIncrement"`
	TxHash   string  `gorm:"uniqueIndex;size:66;not null"`
	ChainID  uint64  `gorm:"type:bigint(20);not null"`
	Standard uint8   `gorm:"type:int(8);not null"`
	Time     uint64  `gorm:"type:bigint(20);not null"`
	Asset    string  `gorm:"type:varchar(120);not null"`
	From     string  `gorm:"type:varchar(66);not null"`
	To       string  `gorm:"type:varchar(66);not null"`
	TokenID  *BigInt `gorm:"type:varchar(86);not null"`
}

type DstSwap struct {
	ID         int64   `gorm:"primaryKey;autoIncrement"`
	TxHash     string  `gorm:"uniqueIndex;size:66;not null"`
	ChainID    uint64  `gorm:"type:bigint(20);not null"`
	Time       uint64  `gorm:"type:bigint(20);not null"`
	PoolID     uint64  `gorm:"type:bigint(20);not null"`
	InAsset    string  `gorm:"type:varchar(66);not null"`
	InTokenID  *BigInt `gorm:"type:varchar(86);not null"`
	OutAsset   string  `gorm:"type:varchar(120);not null"`
	OutTokenID *BigInt `gorm:"type:varchar(86);not null"`
	DstChainID uint64  `gorm:"type:bigint(20);not null"`
	DstAsset   string  `gorm:"type:varchar(120);not null"`
	DstUser    string  `gorm:"type:varchar(66);not null"`
	Type       uint64  `gorm:"type:bigint(20);not null"`
}

type WrapperTransaction struct {
	ID           int64   `gorm:"primaryKey;autoIncrement"`
	Hash         string  `gorm:"uniqueIndex;size:66;not null"`
	User         string  `gorm:"type:varchar(66);not null"`
	SrcChainID   uint64  `gorm:"type:bigint(20);not null"`
	Standard     uint8   `gorm:"type:int(8);not null"`
	BlockHeight  uint64  `gorm:"type:bigint(20);not null"`
	Time         uint64  `gorm:"type:bigint(20);not null"`
	DstChainID   uint64  `gorm:"type:bigint(20);not null"`
	DstUser      string  `gorm:"type:varchar(66);not null"`
	ServerID     uint64  `gorm:"type:bigint(20);not null"`
	FeeTokenHash string  `gorm:"size:66;not null"`
	FeeAmount    *BigInt `gorm:"type:varchar(64);not null"`
	Status       uint64  `gorm:"type:bigint(20);not null"`
}

type SrcPolyDstRelation struct {
	SrcHash            string
	WrapperTransaction *WrapperTransaction `gorm:"foreignKey:SrcHash;references:Hash"`
	SrcTransaction     *SrcTransaction     `gorm:"foreignKey:SrcHash;references:Hash"`
	PolyHash           string
	PolyTransaction    *PolyTransaction `gorm:"foreignKey:PolyHash;references:Hash"`
	DstHash            string
	DstTransaction     *DstTransaction `gorm:"foreignKey:DstHash;references:Hash"`
	ChainID            uint64
	TokenHash          string
	FeeTokenHash       string
	Token              *Token `gorm:"foreignKey:TokenHash,ChainID;references:Hash,ChainID"`
	FeeToken           *Token `gorm:"foreignKey:FeeTokenHash,ChainID;references:Hash,ChainID"`
}

type PolyTxRelation struct {
	SrcHash            string
	WrapperTransaction *WrapperTransaction `gorm:"foreignKey:SrcHash;references:Hash"`
	SrcTransaction     *SrcTransaction     `gorm:"foreignKey:SrcHash;references:Hash"`
	PolyHash           string
	PolyTransaction    *PolyTransaction `gorm:"foreignKey:PolyHash;references:Hash"`
	DstHash            string
	DstTransaction     *DstTransaction `gorm:"foreignKey:DstHash;references:Hash"`
	ChainID            uint64          `gorm:"type:bigint(20);not null"`
	ToChainID          uint64          `gorm:"type:bigint(20);not null"`
	DstChainID         uint64          `gorm:"type:bigint(20);not null"`
	TokenHash          string          `gorm:"type:varchar(66);not null"`
	ToTokenHash        string          `gorm:"type:varchar(66);not null"`
	DstTokenHash       string          `gorm:"type:varchar(66);not null"`
	FeeTokenHash       string          `gorm:"type:varchar(66);not null"`
	Token              *Token          `gorm:"foreignKey:TokenHash,ChainID;references:Hash,ChainID"`
	FeeToken           *Token          `gorm:"foreignKey:FeeTokenHash,ChainID;references:Hash,ChainID"`
	ToToken            *Token          `gorm:"foreignKey:ToTokenHash,ToChainID;references:Hash,ChainID"`
	DstToken           *Token          `gorm:"foreignKey:DstTokenHash,DstChainID;references:Hash,ChainID"`
}

type AssetStatistic struct {
	ID             int64       `gorm:"primaryKey;autoIncrement"`
	Amount         *BigInt     `gorm:"type:varchar(64);not null"`
	Txnum          uint64      `gorm:"type:bigint(20);not null"`
	Addressnum     uint64      `gorm:"type:bigint(20);not null"`
	TokenBasicName string      `gorm:"uniqueIndex;size:64;not null"`
	AmountBtc      *BigInt     `gorm:"type:varchar(64);not null"`
	AmountUsd      *BigInt     `gorm:"type:varchar(64);not null"`
	LastCheckID    int64       `gorm:"type:int"`
	TokenBasic     *TokenBasic `gorm:"foreignKey:TokenBasicName;references:Name"`
}

type AssetInfo struct {
	Amount         *BigInt
	Txnum          uint64
	Price          int64
	TokenBasicName string
	Precision      uint64
}

type TransactionOnToken struct {
	Hash    string
	ChainID uint64
	Time    uint64
	Height  uint64
	From    string
	To      string
	Amount  *BigInt
	Direct  uint32
}

type TransactionOnAddress struct {
	Hash      string
	ChainID   uint64
	Time      uint64
	Height    uint64
	From      string
	To        string
	Amount    *BigInt
	TokenHash string
	TokenName string
	TokenType string
	Direct    uint32
	Precision uint64
}

type ErrorTransaction struct {
	ID           uint    `gorm:"primaryKey;autoIncrement"`
	TxHash       string  `gorm:"type:varchar(150);not null"`
	FromChainID  uint64  `gorm:"type:bigint(20);not null"`
	FromContract string  `gorm:"type:varchar(150);not null"`
	ToChainID    uint64  `gorm:"type:bigint(20);not null"`
	ToContract   string  `gorm:"type:varchar(66);not null"`
	ToAssetHash  string  `gorm:"type:varchar(66);not null"`
	ToAddress    string  `gorm:"type:varchar(66);not null"`
	TokenID      *BigInt `gorm:"type:varchar(86);not null"`
	TokenURI     string  `gorm:"type:varchar(255);not null"`
	Signature    string  `gorm:"not null"`
	State        uint    `gorm:"default:1"`
	ErrorMsg     string
}

type TxHashHistory struct {
	ID      uint   `gorm:"primaryKey;autoIncrement"`
	ChainID uint64 `json:"chain_id"`
	TxHash  string `json:"uniqueIndex;tx_hash"`
}
