package types

// ConnectStatus - represents a possible resource status
type ConnectStatus string

const (
	// Offline status
	Offline ConnectStatus = "Offline"
	// Connecting status
	Connecting ConnectStatus = "Connecting"
	// Connected status
	Connected ConnectStatus = "Connected"
	// Failed status
	Failed ConnectStatus = "Failed"
)

// UserCredentials - represents SlashDBs autoload request user credentials
type UserCredentials struct {
	Name     string `json:"dbuser" xml:"dbuser"`
	Password string `json:"dbpass" xml:"dbpass"`
}

// ResourceConfig - represents SlashDBs request definition
type ResourceConfig struct {
	ID              string              `json:"db_id,omitempty" xml:"db_id,omitempty"`
	Type            string              `json:"db_type,omitempty" xml:"db_type,omitempty"`
	Encoding        string              `json:"db_encoding,omitempty" xml:"db_encoding,omitempty"`
	Desc            string              `json:"desc,omitempty" xml:"desc,omitempty"`
	Autoload        bool                `json:"autoload,omitempty" xml:"autoload,omitempty"`
	Autoconnect     bool                `json:"autoconnect,omitempty" xml:"autoconnect,omitempty"`
	Viewable        bool                `json:"viewable,omitempty" xml:"viewable,omitempty"`
	Editable        bool                `json:"editable,omitempty" xml:"editable,omitempty"`
	Executable      bool                `json:"executable,omitempty" xml:"executable,omitempty"`
	Creator         string              `json:"creator,omitempty" xml:"creator,omitempty"`
	Owner           []string            `json:"owners,omitempty" xml:"owners,omitempty"`
	Read            []string            `json:"read,omitempty" xml:"read,omitempty"`
	Write           []string            `json:"write,omitempty" xml:"write,omitempty"`
	Execute         []string            `json:"execute,omitempty" xml:"execute,omitempty"`
	Connection      string              `json:"connection,omitempty" xml:"connection,omitempty"`
	Schema          string              `json:"db_schema,omitempty" xml:"db_schema,omitempty"`
	UserCredentials UserCredentials     `json:"autoload_user,omitempty" xml:"autoload_user,omitempty"`
	AlternateKeys   map[string][]string `json:"alternate_key,omitempty" xml:"alternate_key,omitempty"`
	ExcludedKeys    map[string][]string `json:"excluded_columns,omitempty" xml:"excluded_columns,omitempty"`
	ForeignKeys     map[string][]string `json:"foreign_keys,omitempty" xml:"foreign_keys,omitempty"`
	ConnectStatus   ConnectStatus       `json:"connect_status,omitempty" xml:"connect_status,omitempty"`
}

// QueryConfig - represents SlashDBs custom QueryConfig definition
type QueryConfig struct {
	ID          string          `json:"query_id,omitempty" xml:"query_id,omitempty"`
	DatabaseID  string          `json:"database,omitempty" xml:"database,omitempty"`
	Desc        string          `json:"desc,omitempty" xml:"desc,omitempty"`
	SQLStr      string          `json:"sqlstr,omitempty" xml:"sqlstr,omitempty"`
	HTTPMethods map[string]bool `json:"http_methods,omitempty" xml:"http_methods,omitempty"`
	Creator     string          `json:"creator,omitempty" xml:"creator,omitempty"`
	Read        []string        `json:"read,omitempty" xml:"read,omitempty"`
	Write       []string        `json:"write,omitempty" xml:"write,omitempty"`
	Execute     []string        `json:"execute,omitempty" xml:"execute,omitempty"`
}

// UserConfig - represents SlashDBs custom UserConfig definition
type UserConfig struct {
	ID              string                     `json:"user_id,omitempty" xml:"user_id,omitempty"`
	Password        string                     `json:"password,omitempty" xml:"password,omitempty"`
	Name            string                     `json:"name,omitempty" xml:"name,omitempty"`
	Email           string                     `json:"email,omitempty" xml:"email,omitempty"`
	DSCredentials   map[string]UserCredentials `json:"databases,omitempty" xml:"databases,omitempty"`
	ResourceConfigs []string                   `json:"dbdef,omitempty" xml:"dbdef,omitempty"`
	QueryConfigs    []string                   `json:"querydef,omitempty" xml:"querydef,omitempty"`
	UserConfigs     []string                   `json:"userdef,omitempty" xml:"userdef,omitempty"`
	Creator         string                     `json:"creator,omitempty" xml:"creator,omitempty"`
	View            []string                   `json:"view,omitempty" xml:"view,omitempty"`
	Edit            []string                   `json:"edit,omitempty" xml:"edit,omitempty"`
	APIKey          string                     `json:"api_key,omitempty" xml:"api_key,omitempty"`
}
