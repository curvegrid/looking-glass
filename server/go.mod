module github.com/curvegrid/looking-glass/server

go 1.15

require (
	github.com/curvegrid/gofig v1.1.0
	github.com/ethereum/go-ethereum v1.9.24
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/websocket v1.4.2
	github.com/labstack/echo/v4 v4.1.17
	github.com/sirupsen/logrus v1.7.0
)
// until https://github.com/labstack/echo/pull/1651 is merged
replace github.com/labstack/echo/v4 => github.com/curvegrid/echo/v4 v4.1.18-0.20201124072549-e6f24aa8b1cb
