package mongoserve

// conn
const (
	// HostKey is the key of mongo host, if is cluster, use comma to split
	HostKey = "mongo.hostStr"
	// MinPoolSizeKey is the key of min pool size
	MinPoolSizeKey = "mongo.minPoolSizePerHost"
	// MaxPoolSizeKey is the key of max pool size
	MaxPoolSizeKey = "mongo.maxPoolSizePerHost"
	// ServerSelectTimeoutKey is the key of server selection timeout
	ServerSelectionTimeoutKey = "mongo.serverSelectionTimeout"
	// ConnectTimeoutKey is the key of connect timeout
	ConnectTimeoutKey = "mongo.connectTimeout"
	// SocketTimeoutKey is the key of socket timeout
	SocketTimeoutKey = "mongo.socketTimeout"
	// MaxConnIdleTimeKey is the key of max connection idle time
	MaxConnIdleTimeKey = "mongo.maxConnIdleTime"
)

// auth
const (
	// AuthMechanismKey is the key of mongo auth mechanism
	AuthMechanismKey = "mongo.auth.authMechanism"
	// UsernameKey is the key of mongo username
	UsernameKey = "mongo.auth.username"
	// PasswordKey is the key of mongo password
	PasswordKey = "mongo.auth.password"
	// AuthSourceKey is the key of mongo auth source
	AuthSourceKey = "mongo.auth.authSource"
)

// read preference
const (
	// PreferPrimaryKey is the key of mongo read preference
	ReadPreferenceKey = "mongo.preferPrimary"
)
