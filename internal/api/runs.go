package api

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"os"

	"github.com/gorilla/websocket"
)

type RunStatusResponse struct {
	State  string  `json:"state"`
	Step   Section `json:"step"`
	MsLeft int64   `json:"ms_left"`
	// errors
	Err string `json:"error,omitempty"`
}

type RunInput struct {
	Action string `json:"action"`
	Step   *int32 `json:"step,omitempty"`
}

func ConnectToRun(
	ctx context.Context,
	client APIClient,
	logger *slog.Logger,
	presentationID int,
) (<-chan RunStatusResponse, chan<- RunInput, *websocket.Conn, error) {
	u, err := url.Parse(client.Url(os.LookupEnv))
	if err != nil {
		return nil, nil, nil, err
	}

	u.Scheme = "ws"
	u.Path = fmt.Sprintf("/run/%d", presentationID)

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return nil, nil, nil, err
	}

	res := make(chan RunStatusResponse)
	go func() {
		defer func() {
			logger.Info("connet to run", "exit", "read loop")
		}()
		defer close(res)

		var response RunStatusResponse

		for {
			err := c.ReadJSON(&response)
			logger.Info("received message", "msg", response, "err", err)
			if err != nil {
				return
			}

			res <- response
		}
	}()

	in := make(chan RunInput)
	go func() {
		defer func() {
			logger.Info("connet to run", "exit", "write loop")
		}()

		for {
			logger.Info("write loop")
			select {
			case <-ctx.Done():
				return
			case input, ok := <-in:
				logger.Info("api", "send input", input)
				if !ok {
					logger.Info("api close input ch")
					return
				}

				if err := c.WriteJSON(&input); err != nil {
					logger.Error("error writting to ws", "err", err)
				}
			}
		}
	}()

	return res, in, c, nil
}
