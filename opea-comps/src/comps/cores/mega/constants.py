from enum import Enum

class ServiceType(str, Enum):
    LLM = "llm"
    EMBEDDING = "embedding"

class ServiceRoleType(str, Enum):
    MEGASERVICE = "megaservice"
