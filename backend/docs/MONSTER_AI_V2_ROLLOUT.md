# Monster AI V2 Rollout

## Feature Flags

- Backend: `ENABLE_MONSTER_AI_V2=true`
- Frontend: `VITE_ENABLE_MONSTER_AI_V2=true`

## Rollout Steps

1. Deploy backend with flag disabled (default).
2. Deploy frontend with flag disabled (default) to ensure no behavior changes.
3. Enable backend flag in staging and validate:
   - `POST /api/monsters/preview`
   - `POST /api/monsters/{id}/regenerate-section`
   - `POST /api/monsters/{id}/variant`
   - `POST /api/monsters/{id}/duplicate`
4. Enable frontend flag in staging and verify preview + sheet actions.
5. Monitor generation event summary endpoint:
   - `GET /api/user/{userID}/monsters/generation-summary`
6. Enable both flags in production.
