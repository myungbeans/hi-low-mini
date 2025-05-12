import { json } from '@sveltejs/kit';
import type { RequestHandler } from '@sveltejs/kit';
import { GetGameRequestSchema, GetGameResponseSchema, GameSchema } from '$lib/gen/game_engine/v1/game_engine_pb';
import { create } from '@bufbuild/protobuf';

export const GET: RequestHandler = async () => {
  try {
    console.log('Handling GET request to /api/game');
    
    // Create the request message using protobuf
    const request = create(GetGameRequestSchema, {
      timestamp: {
        seconds: BigInt(Math.floor(Date.now() / 1000)),
        nanos: 0
      },
      userId: ''
    });

    // Convert to JSON in the format the server expects (Connect protocol format)
    const requestJson = {
      "@type": "type.googleapis.com/game_engine.v1.GetGameRequest",
      "value": {
        timestamp: {
          seconds: request.timestamp?.seconds.toString(),
          nanos: request.timestamp?.nanos || 0
        },
        userId: request.userId
      }
    };

    console.log('Sending request to game engine:', requestJson);
    const response = await fetch('http://localhost:8080/game_engine.v1.GameEngineService/GetGame', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Connect-Protocol-Version': '1'
      },
      body: JSON.stringify(requestJson),
    });

    if (!response.ok) {
      const errorText = await response.text();
      console.error('Game engine error:', response.status, errorText);
      throw new Error(`Game engine error: ${response.status} ${errorText}`);
    }

    const data = await response.json();
    console.log('Raw response from game engine:', data);
    
    // If we got an error response with a code, throw it
    if (data.code) {
      throw new Error(`Game engine error: ${data.code} - ${data.message}`);
    }

    // Process the game data to ensure proper type handling
    const processedGame = {
      ...data.game,
      pool: {
        ...data.game.pool,
        cards: data.game.pool.cards.map((card: any) => {
          // Check if the value looks like an operator
          const isOperator = ['+', '-', '*', '/', '^'].includes(card.value);
          return {
            value: card.value,
            // Use the numeric enum values directly
            type: isOperator ? 2 : 1 // 2 for CARD_TYPE_OPERATOR, 1 for CARD_TYPE_NUMBER
          };
        })
      }
    };

    console.log('Processed game data:', processedGame);

    // Create proper protobuf messages to ensure correct type handling
    const game = create(GameSchema, processedGame);
    console.log('Game after protobuf creation:', game);
    
    const gameResponse = create(GetGameResponseSchema, {
      game
    });
    console.log('Final game response:', gameResponse);
    
    return json(gameResponse);
  } catch (error) {
    console.error('Error in GET /api/game:', error);
    return new Response(
      JSON.stringify({ 
        error: error instanceof Error ? error.message : 'Failed to fetch game',
        details: error instanceof Error ? error.stack : undefined
      }), 
      {
        status: 500,
        headers: {
          'Content-Type': 'application/json',
        },
      }
    );
  }
}; 