package websocket

import (
	"os"
	"time"

	"github.com/op/go-logging"
)

// Parameters defines adapter behavior.
type Parameters struct {
	AutoReconnect         bool
	ReconnectInterval     time.Duration
	ReconnectAttempts     int
	reconnectTry          int
	ShutdownTimeout       time.Duration
	CapacityPerConnection int
	Logger                *logging.Logger

	ResubscribeOnReconnect bool

	HeartbeatTimeout time.Duration
	LogTransport     bool

	URL             string
	ManageOrderbook bool
}

func NewDefaultParameters() *Parameters {
	format := logging.MustStringFormatter(
		`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
	)

	lg := logging.MustGetLogger("bitfinex-ws")
	be := logging.NewLogBackend(os.Stderr, "", 0)
	bef := logging.NewBackendFormatter(be, format)
	bel := logging.AddModuleLevel(bef)
	bel.SetLevel(logging.ERROR, "")
	logging.SetBackend(bel)

	return &Parameters{
		AutoReconnect:          true,
		CapacityPerConnection:  25,
		ReconnectInterval:      time.Second * 3,
		reconnectTry:           0,
		ReconnectAttempts:      15,
		URL:                    productionBaseURL,
		ManageOrderbook:        false,
		ShutdownTimeout:        time.Second * 5,
		ResubscribeOnReconnect: true,
		HeartbeatTimeout:       time.Second * 30,
		LogTransport:           false, // log transport send/recv
		Logger:                 lg,
	}
}
