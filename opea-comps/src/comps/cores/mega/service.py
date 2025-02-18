import httpx
import json
from .constants import ServiceType, ServiceRoleType

class MicroService:
    def __init__(self, name, host="0.0.0.0", port=8000, endpoint="/", service_type=None, service_role=None, use_remote_service=False, input_datatype=None, output_datatype=None):
        self.name = name
        self.host = host
        self.port = port
        self.endpoint = endpoint
        self.service_type = service_type
        self.service_role = service_role
        self.use_remote_service = use_remote_service
        self.input_datatype = input_datatype
        self.output_datatype = output_datatype
        self.routes = {}
        self._client = httpx.AsyncClient()

    def add_route(self, path, handler, methods=None):
        if methods is None:
            methods = ["GET"]
        self.routes[path] = {"handler": handler, "methods": methods}

    def start(self):
        print(f"Starting {self.name} service on {self.host}:{self.port}")

    async def close(self):
        await self._client.aclose()

    @property
    def base_url(self):
        return f"http://{self.host}:{self.port}{self.endpoint}"

    async def call_service(self, request_data):
        try:
            response = await self._client.post(
                self.base_url,
                json=request_data,
                headers={"Content-Type": "application/json"}
            )
            response.raise_for_status()
            response_json = response.json()
            # Handle Ollama's response format
            if 'response' in response_json and not response_json.get('done', False):
                # This is a streaming response, we'll just take the first part
                return {"response": response_json['response']}
            return response_json
        except httpx.HTTPError as e:
            print(f"Error calling service {self.name}: {str(e)}")
            return {"error": str(e)}

class ServiceOrchestrator:
    def __init__(self):
        self.services = {}
        self.flow = {}

    def add(self, service):
        self.services[service.name] = service
        return self

    def flow_to(self, from_service, to_service):
        if from_service.name not in self.flow:
            self.flow[from_service.name] = []
        self.flow[from_service.name].append(to_service.name)
        return self

    async def schedule(self, request):
        results = []
        for service_name, service in self.services.items():
            if service.service_type == ServiceType.LLM:
                try:
                    response = await service.call_service(request)
                    results.append({"llm/MicroService": {"body": response.get("response", "No response content")}})
                except Exception as e:
                    print(f"Error in LLM service: {str(e)}")
                    results.append({"llm/MicroService": {"body": f"Error: {str(e)}"}})
        return results
