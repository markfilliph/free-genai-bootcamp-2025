```sh
curl -X POST http://localhost:8008/v1/example-service \
  -H "Content-Type: application/json" \
  -d '{
    "model": "llama3.2:1b",
    "messages": "Hello, this is just a test, what do you think?"
  }'
```