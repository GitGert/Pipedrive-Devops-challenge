curl -X POST http://localhost:8080/deals \
     -H "Content-Type: application/json" \
     -d '{
           "title": "Example Deal",
           "value": "1000",
           "currency": "EUR",
           "status": "open"
         }'
