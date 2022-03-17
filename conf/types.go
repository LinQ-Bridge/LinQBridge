package conf

type Config struct {
	RunMode    string
	DBConfig   *DBConfig
	Chains     []*ChainListenConfig
	LinQConfig *LinQConfig
}

type DBConfig struct {
	URL      string
	User     string
	Password string
	Scheme   string
	Debug    bool
}

type RedisConfig struct {
	Addr     string
	Port     int
	Password string
	Db       int
}

type Restful struct {
	URL string
	//Key string
}

type ChainListenConfig struct {
	ChainName          string
	ChainID            uint64
	ListenSlot         uint64
	BatchSize          uint64
	Defer              uint64
	Nodes              []*Restful
	NFTWrapperContract string
	NFTProxyContract   string
	NFTQueryContract   string
	CCMContract        string
}

type ContractAddrs struct {
	ChainID   uint64
	URL       string
	ProxyAddr string
	QueryAddr string
}
type ValidatorConfig struct {
	Addr string
}

type QueryConfig struct {
	ChainID   uint64
	URL       string
	QueryAddr string
	LockAddr  string
}

type LinQConfig struct {
	DefaultBootNodes []string
	Addr             string
	Port             uint
}
