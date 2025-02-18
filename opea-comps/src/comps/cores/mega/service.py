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

    def add_route(self, path, handler, methods=None):
        if methods is None:
            methods = ["GET"]
        self.routes[path] = {"handler": handler, "methods": methods}

    def start(self):
        print(f"Starting {self.name} service on {self.host}:{self.port}")

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
        # This is a simplified implementation
        results = []
        for service_name, service in self.services.items():
            if service.service_type == ServiceType.LLM:
                # Simulate LLM service response
                response = {"llm/MicroService": {"body": "Simulated LLM response"}}
                results.append(response)
        return results
