from fastapi import FastAPI, HTTPException
from comps.cores.proto.api_protocol import (
    ChatCompletionRequest,
    ChatCompletionResponse,
    ChatCompletionResponseChoice,
    ChatMessage,
    UsageInfo
)
from comps.cores.mega.constants import ServiceType, ServiceRoleType
from comps import MicroService, ServiceOrchestrator
import uvicorn
import os

EMBEDDING_SERVICE_HOST_IP = os.getenv("EMBEDDING_SERVICE_HOST_IP", "0.0.0.0")
EMBEDDING_SERVICE_PORT = os.getenv("EMBEDDING_SERVICE_PORT", 6000)
LLM_SERVICE_HOST_IP = os.getenv("LLM_SERVICE_HOST_IP", "localhost")
LLM_SERVICE_PORT = int(os.getenv("LLM_SERVICE_PORT", 9000))

app = FastAPI()

SERVICE_PORT = int(os.getenv("SERVICE_PORT", 8008))

class ExampleService:
    def __init__(self, host="0.0.0.0", port=SERVICE_PORT):
        print('hello')
        os.environ["TELEMETRY_ENDPOINT"] = ""
        self.host = host
        self.port = port
        self.endpoint = "/v1/example-service"
        self.megaservice = ServiceOrchestrator()

    def add_remote_service(self):
        llm = MicroService(
            name="llm",
            host=LLM_SERVICE_HOST_IP,
            port=LLM_SERVICE_PORT,
            endpoint="/api/generate",
            use_remote_service=True,
            service_type=ServiceType.LLM,
        )
        self.megaservice.add(llm)

    async def handle_request(self, request: ChatCompletionRequest) -> ChatCompletionResponse:
        try:
            # Format the request for Ollama
            ollama_request = {
                "model": request.model,
                "prompt": request.messages,
                "stream": False
            }
            
            # Schedule the request through the orchestrator
            result = await self.megaservice.schedule(ollama_request)
            
            # Extract the actual content from the response
            if isinstance(result, list) and len(result) > 0:
                llm_response = result[0].get('llm/MicroService')
                content = llm_response.get('body', "No response content available")
            else:
                content = "Invalid response format"

            # Create the response
            response = ChatCompletionResponse(
                model=request.model or "example-model",
                choices=[
                    ChatCompletionResponseChoice(
                        index=0,
                        message=ChatMessage(
                            role="assistant",
                            content=content
                        ),
                        finish_reason="stop"
                    )
                ],
                usage=UsageInfo(
                    prompt_tokens=0,
                    completion_tokens=0,
                    total_tokens=0
                )
            )
            
            return response
            
        except Exception as e:
            # Handle any errors
            raise HTTPException(status_code=500, detail=str(e))

example = ExampleService(host="0.0.0.0", port=LLM_SERVICE_PORT)
example.add_remote_service()

@app.post("/v1/example-service")
async def handle_chat_request(request: ChatCompletionRequest):
    return await example.handle_request(request)

if __name__ == "__main__":
    uvicorn.run(app, host="0.0.0.0", port=SERVICE_PORT)