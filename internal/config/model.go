package config

type Config struct {
	Server   Server
	Jwt      Jwt
	Database Database
	Redis    Redis
	Email    Email
	Midtrans Midtrans
}

type Server struct {
	Host string
	Port string
}

type Jwt struct {
	AccessTokenSecret  string
	RefreshTokenSecret string
}

type Redis struct {
	Addr string
	Pass string
}

type Database struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type Email struct {
	Host     string
	Port     string
	User     string
	Password string
}

type Midtrans struct {
	Key    string
	IsProd bool
}
