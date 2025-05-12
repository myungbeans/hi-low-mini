import { json } from '@sveltejs/kit';
import type { RequestHandler } from '@sveltejs/kit';
import { gameEngineClient } from '$lib/client';
import { create } from '@bufbuild/protobuf';
import { HandSchema, PlayHandRequestSchema } from '$lib/gen/game_engine/v1/game_engine_pb';

export const POST: RequestHandler = async ({ request }) => {
  try {
    const body = await request.json();
    
    // Validate request structure
    if (!body?.value?.hand?.cards) {
      throw new Error('Invalid request structure: missing required fields');
    }

    // Create the Hand message
    const hand = create(HandSchema, {
      cards: body.value.hand.cards.map((card: { value: string; type: number; }) => ({
        value: card.value.toString(),
        type: card.type
      })),
      gameId: body.value.hand.game_id
    });

    // Create the PlayHandRequest message
    const playHandRequest = create(PlayHandRequestSchema, {
      gameId: body.value.game_id,
      timestamp: {
        seconds: BigInt(body.value.timestamp.seconds),
        nanos: body.value.timestamp.nanos
      },
      hand: hand,
      secsElapsed: body.value.secs_elapsed
    });

    // Make the request using the Connect client
    const response = await gameEngineClient.playHand(playHandRequest);
    
    return json({
      playOrder: response.playOrder,
      scoreCounter: response.scoreCounter
    });
  } catch (error: any) {
    return new Response(JSON.stringify({ 
      error: error?.message || 'Failed to play hand'
    }), {
      status: 500,
      headers: {
        'Content-Type': 'application/json',
      },
    });
  }
}; 