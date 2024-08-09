# Pipedrive-Devops-challenge




curl -X GET "http://localhost:8080/deals"
curl -X GET "http://localhost:8080/post_deals"
<!-- curl -X POST "http://localhost:8080/post_deals" -->
curl -X GET "http://localhost:8080/put_deals"
<!-- curl -X PUT "http://localhost:8080/put_deals" -->

snap install golangci-lint --classic


docker build -t sampleapp:latest .
docker run -p 8080:8080 sampleapp