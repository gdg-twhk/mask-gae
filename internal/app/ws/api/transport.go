package api

import (
	"context"
	"errors"
	"net/http"

	stdjwt "github.com/dgrijalva/jwt-go"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/go-zoo/bone"
	"github.com/gorilla/websocket"

	"github.com/cage1016/mask/internal/app/ws"

)

var (
	errUnauthorizedAccess = errors.New("missing or invalid credentials provided")
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	logger     log.Logger
	jwtKeyFunc func(token *stdjwt.Token) (interface{}, error)
)

// MakeHandler returns http handler with handshake endpoint.
func MakeHandler(svc ws.Service, l log.Logger) http.Handler {
	logger = l

	mux := bone.New()
	mux.GetFunc("/ws/:topic", handshake(svc))

	return mux
}

func handshake(svc ws.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sub, err := authorize(r)
		if err != nil {
			switch err {
			case errUnauthorizedAccess:
				w.WriteHeader(http.StatusForbidden)
				return
			default:
				w.WriteHeader(http.StatusServiceUnavailable)
				return
			}
		}

		// Create new ws connection.
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			level.Warn(logger).Log("method", "upgrader.Upgrade", "err", err)
			return
		}
		sub.conn = conn

		sub.channel = ws.NewChannel()
		if err := svc.Subscribe(sub.topic, sub.channel); err != nil {
			level.Warn(logger).Log("method", "svc.Subscribe", "err", err)
			conn.Close()
			return
		}

		go sub.listen()

		// Start listening for messages from NATS.
		go sub.broadcast(svc)
	}
}

func authorize(r *http.Request) (subscription, error) {
	//authKey := r.Header.Get("Authorization")
	//if authKey == "" {
	//	authKeys := bone.GetQuery(r, "authorization")
	//	if len(authKeys) == 0 {
	//		return subscription{}, errUnauthorizedAccess
	//	}
	//	authKey = authKeys[0]
	//}
	//
	//token, err := stdjwt.ParseWithClaims(authKey, stdjwt.MapClaims{}, jwtKeyFunc)
	//if err != nil || !token.Valid {
	//	return subscription{}, errUnauthorizedAccess
	//}
	//
	//claim, ok := token.Claims.(stdjwt.MapClaims)
	//if !ok {
	//	return subscription{}, errUnauthorizedAccess
	//}

	topic := bone.GetValue(r, "topic")

	sub := subscription{
		//pubID: claim[jwt.ClaimsUserId].(string),
		topic: topic,
	}

	return sub, nil
}

type subscription struct {
	pubID   string
	topic   string
	conn    *websocket.Conn
	channel *ws.Channel
}

func (sub subscription) broadcast(svc ws.Service) {
	for {
		_, payload, err := sub.conn.ReadMessage()
		if websocket.IsUnexpectedCloseError(err) {
			level.Warn(logger).Log("IsUnexpectedCloseError", err)
			sub.channel.Close()
			return
		}
		if err != nil {
			level.Warn(logger).Log("method", "sub.conn.ReadMessage", "err", err)
			return
		}

		if err := svc.Publish(context.Background(), sub.topic, string(payload)); err != nil {
			level.Warn(logger).Log("method", "svc.Publish", "topic", sub.topic, "err", err)
			if err == ws.ErrFailedConnection {
				sub.conn.Close()
				sub.channel.Closed <- true
				return
			}
		}
	}
}

func (sub subscription) listen() {
	for msg := range sub.channel.Messages {
		if err := sub.conn.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
			level.Warn(logger).Log("method", "sub.conn.WriteMessage", "err", err)
		}
	}
}
