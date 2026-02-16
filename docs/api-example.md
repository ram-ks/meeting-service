## Event
### Create Event

```
curl --request POST \
  --url http://localhost:8080/events \
  --header 'Content-Type: application/json' \
  --header 'User-Agent: insomnia/12.3.1' \
  --data '{
    "title": "Team Planning Meeting",
    "description": "Q1 planning and goal setting",
    "duration": "1h30m",
    "proposed_slots": [
      {
        "start_time": "2025-02-20T10:00:00Z",
        "end_time": "2025-02-20T11:30:00Z",
        "timezone": "America/New_York"
      },
      {
        "start_time": "2025-02-20T14:00:00Z",
        "end_time": "2025-02-20T15:30:00Z",
        "timezone": "America/New_York"
      },
      {
        "start_time": "2025-02-21T09:00:00Z",
        "end_time": "2025-02-21T10:30:00Z",
        "timezone": "UTC"
      }
    ],
    "participants": [
      {
        "email": "alice@example.com",
        "name": "Alice Johnson"
      },
      {
        "email": "bob@example.com",
        "name": "Bob Smith"
      },
      {
        "email": "charlie@example.com",
        "name": "Charlie Davis"
      }
    ]
}'
```

### GET EVENT
```
curl --request GET \
  --url http://localhost:8080/events/639d904e-bf33-4876-9213-2f636576a85c \
  --header 'User-Agent: insomnia/11.6.1'
```

## Submit Slots
### POST slots for each user via email

```
curl --request POST \
  --url http://localhost:8080/preferred-slots \
  --header 'Content-Type: application/json' \
  --header 'User-Agent: insomnia/12.3.1' \
  --data '{
    "email": "charlie@example.com",
    "start_time": "2025-02-20T09:00:00Z",
    "end_time": "2025-02-20T12:00:00Z",
    "timezone": "America/New_York",
    "day_of_week": 4
}'
```

## Availability
### Submit Availability for each participant

```
curl --request POST \
  --url http://localhost:8080/events/639d904e-bf33-4876-9213-2f636576a85c/availability \
  --header 'Content-Type: application/json' \
  --header 'User-Agent: insomnia/12.3.1' \
  --data '{
	"participant_id": "57593f7f-3d7f-459b-84df-30f13059e970",
	"slots": [
		{
			"slot_id": "96336ee1-a6bd-49c8-b632-da0acc6eed09",
			"status": "available"
		}
	]
}'
```

## Recommendation
### Get Recommendation

```
curl --request GET \
  --url http://localhost:8080/events/639d904e-bf33-4876-9213-2f636576a85c/recommendations \
  --header 'User-Agent: insomnia/12.3.1'
```

### Health
```
curl --request GET \
  --url http://localhost:8080/health \
  --header 'User-Agent: insomnia/12.3.1'
```