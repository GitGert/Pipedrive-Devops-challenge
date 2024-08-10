curl -X PUT http://localhost:8080/put_deals?dealId=20 \
-H "Content-Type: application/json" \
-d '{
  "title": "Deal Title",
  "value": 1000.00,
  "currency": "EUR",
  "is_deleted": false,
  "status": "won"
}'