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
)

type GameEngineServer struct{}

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
	res := connect.NewResponse(&pb.PlayHandResponse{
		PlayOrder:    []int32{},
		ScoreCounter: []float32{},
	})
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
