module github.com/gempir/gempbot

go 1.16

require (
	github.com/gempir/go-twitch-irc/v2 v2.7.0
	github.com/go-sql-driver/mysql v1.6.0
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/golang-jwt/jwt/v4 v4.1.0 // indirect
	github.com/nicklaw5/helix/v2 v2.1.0
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.7.0 // indirect
	golang.org/x/sys v0.0.0-20211007075335-d3039528d8ac // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
	gorm.io/driver/mysql v1.1.2
	gorm.io/gorm v1.21.15
)

replace github.com/nicklaw5/helix/v2 v2.0.1 => github.com/gempir/helix/v2 v2.0.2
