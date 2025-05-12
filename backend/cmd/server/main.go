package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"connectrpc.com/connect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	pb "hi_low_mini/gen/game_engine/v1"
	"hi_low_mini/runtime/engine"
	"hi_low_mini/runtime/models"
)

// CORSMiddleware adds CORS headers to responses
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000") //TODO: update with real domain once implemented
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Connect-Protocol-Version, Accept")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

type GameEngineServer struct{}

// TODO: this should read from db to fetch today's game instead of generating a new one
func (s *GameEngineServer) GetGame(
	ctx context.Context,
	req *connect.Request[pb.GetGameRequest],
) (*connect.Response[pb.GetGameResponse], error) {
	res := connect.NewResponse(&pb.GetGameResponse{
		Game: models.NewGame(engine.NewSet()),
	})
	res.Header().Set("Version", "v1")
	return res, nil
}

func (s *GameEngineServer) GenerateGame(
	ctx context.Context,
	req *connect.Request[pb.GenerateGameRequest],
) (*connect.Response[pb.GenerateGameResponse], error) {
	res := connect.NewResponse(&pb.GenerateGameResponse{
		Game: models.NewGame(engine.NewSet()),
	})
	res.Header().Set("Version", "v1")
	return res, nil
}

func (s *GameEngineServer) PlayHand(
	ctx context.Context,
	req *connect.Request[pb.PlayHandRequest],
) (*connect.Response[pb.PlayHandResponse], error) {
	log.Println("Request headers: ", req.Header())
	log.Printf("Raw request: %+v\n", req)
	log.Printf("Request message: %+v\n", req.Msg)

	if req.Msg == nil {
		log.Println("Error: request message is nil")
		return nil, fmt.Errorf("invalid request: message is nil")
	}

	if req.Msg.Hand == nil {
		log.Printf("Error: hand is nil in message: %+v\n", req.Msg)
		return nil, fmt.Errorf("invalid request: missing hand")
	}

	cards := req.Msg.GetHand().GetCards()
	if len(cards) == 0 {
		log.Println("Error: no cards in hand")
		return nil, fmt.Errorf("invalid request: no cards in hand")
	}

	log.Printf("Processing hand with %d cards: %+v\n", len(cards), cards)
	scores := []float32{}

	scores, err := engine.Calculate(cards, scores)
	if err != nil {
		log.Printf("Error calculating scores: %v\n", err)
		return nil, err
	}

	log.Printf("Calculated scores: %v\n", scores)
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
		h2c.NewHandler(CORSMiddleware(mux), &http2.Server{}),
	)
}
