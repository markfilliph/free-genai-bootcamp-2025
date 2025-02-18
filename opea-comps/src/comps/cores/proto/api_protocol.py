from pydantic import BaseModel
from typing import List, Optional

class ChatMessage(BaseModel):
    role: str
    content: str

class ChatCompletionResponseChoice(BaseModel):
    index: int
    message: ChatMessage
    finish_reason: str

class UsageInfo(BaseModel):
    prompt_tokens: int
    completion_tokens: int
    total_tokens: int

class ChatCompletionRequest(BaseModel):
    model: Optional[str] = None
    messages: str

class ChatCompletionResponse(BaseModel):
    model: str
    choices: List[ChatCompletionResponseChoice]
    usage: UsageInfo
