package main

import (
	"context"
	"log"
	"net/http"

	"connectrpc.com/connect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "hi_low_mini/gen/game_engine/v1"
	"hi_low_mini/runtime/engine"
)

type GameEngineServer struct{}

// TODO: needs to be implemented - returns a stub Game with 1 Card in the Pool for local testing
func (s *GameEngineServer) GetGame(
	ctx context.Context,
	req *connect.Request[pb.GetGameRequest],
) (*connect.Response[pb.GetGameResponse], error) {
	log.Println("Request headers: ", req.Header())
	res := connect.NewResponse(&pb.GetGameResponse{
		Game: &pb.Game{
			Id:        "000",
			Timestamp: timestamppb.Now(),
			Pool: &pb.Pool{
				Cards: []*pb.Card{
					&pb.Card{
						Value: "1",
						Type:  pb.CardType_CARD_TYPE_NUMBER,
					},
				},
			},
		},
	})
	res.Header().Set("Version", "v1")
	return res, nil
}

func (s *GameEngineServer) PlayHand(
	ctx context.Context,
	req *connect.Request[pb.PlayHandRequest],
) (*connect.Response[pb.PlayHandResponse], error) {
	log.Println("Request headers: ", req.Header())

	// Hands are pre-validated:
	// Starts with a number, followed by an operator, end with a number, has 7 cards
	// e.g. ["1","+","2","*","3","/","4"]
	cards := req.Msg.GetHand().GetCards()
	scores := []float32{}

	scores, err := engine.Calculate(cards, scores)
	if err != nil {
		return nil, err
	}

	res := connect.NewResponse(&pb.PlayHandResponse{
		ScoreCounter: scores,
	})

	// TODO: consider additional headers needed
	// (if this were work, I would create a ticket and paste the ticket num here track work)
	res.Header().Set("Version", "v1")

	return res, nil
}

func main() {
	server := &GameEngineServer{}
	mux := http.NewServeMux()
	path, handler := pb.NewGameEngineServiceHandler(server)
	mux.Handle(path, handler)
	mux.Handle("/api/", http.StripPrefix("/api", mux))
	log.Println("Server listening on :8080")
	http.ListenAndServe(
		"localhost:8080",
		// Use h2c so we can serve HTTP/2 without TLS.
		h2c.NewHandler(mux, &http2.Server{}),
	)
}
