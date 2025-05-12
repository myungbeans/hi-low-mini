import { createClient } from "@connectrpc/connect";
import { createConnectTransport } from "@connectrpc/connect-web";
import { GameEngineService } from "$lib/gen/game_engine/v1/game_engine_pb";

// Create a transport for the game engine service
const transport = createConnectTransport({
  baseUrl: "http://localhost:8080",
  useBinaryFormat: false, // Use JSON format for easier debugging
});

// Create the client for type-safe RPC calls
export const gameEngineClient = createClient(GameEngineService, transport); 