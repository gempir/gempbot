module github.com/gempir/gempbot

go 1.18

require (
	github.com/carlmjohnson/requests v0.23.5
	github.com/gempir/go-twitch-irc/v4 v4.0.0
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/google/uuid v1.6.0
	github.com/gorilla/websocket v1.5.1
	github.com/nicklaw5/helix/v2 v2.26.0
	github.com/puzpuzpuz/xsync v1.5.2
	github.com/rs/cors v1.10.1
	github.com/sirupsen/logrus v1.9.3
	github.com/stretchr/testify v1.8.4
	github.com/teris-io/shortid v0.0.0-20220617161101-71ec9f2aa569
	gorm.io/driver/postgres v1.5.6
	gorm.io/gorm v1.25.7
)

replace github.com/nicklaw5/helix/v2 v2.12.0 => github.com/gempir/helix/v2 v2.0.2-0.20221223221449-fe5671ac8ea7

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/golang-jwt/jwt/v4 v4.5.0 // indirect
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgconn v1.14.1 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.3.2 // indirect
	github.com/jackc/pgservicefile v0.0.0-20231201235250-de7065d80cb9 // indirect
	github.com/jackc/pgtype v1.14.2 // indirect
	github.com/jackc/pgx/v4 v4.18.1 // indirect
	github.com/jackc/pgx/v5 v5.4.3 // indirect
	github.com/jellydator/ttlcache/v2 v2.11.1
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/objx v0.5.0 // indirect
	golang.org/x/crypto v0.19.0 // indirect
	golang.org/x/net v0.21.0 // indirect
	golang.org/x/sync v0.6.0 // indirect
	golang.org/x/sys v0.17.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
