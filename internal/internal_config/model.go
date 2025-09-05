package internal_config

type Config struct {
	Server      Server      `json:"server"`
	Application Application `json:"application"`
	Something   string      `json:"something"`
}

type Server struct {
	ServerName string `json:"server_name"`
	PortRun    int    `json:"port_run"`
	Realm      string `json:"realm"`
	IsDev      bool   `json:"is_dev"`
}

type Application struct {
	SecretKey         string `json:"secret_key"` //user for JWT & passHash
	TokenExpHours     int    `json:"token_exp_hours"`
	AccessTknTimeout  int64  `json:"accessTknTimeout"`
	RefreshTknTimeout int64  `json:"refreshTknTimeout"`
	AccessKey         string `json:"access_key"`
	RefreshKey        string `json:"refresh_key"`
}
