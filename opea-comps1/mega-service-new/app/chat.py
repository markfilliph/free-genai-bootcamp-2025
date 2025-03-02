from typing import List, Optional, Dict, Any
from fastapi import FastAPI, Request, HTTPException
from pydantic import BaseModel, Field
from comps import ServiceOrchestrator, ServiceRoleType, MicroService
from comps.cores.proto.api_protocol import (
    ChatCompletionRequest,
    ChatCompletionResponse,
    Message
)
import asyncio
import logging

# Configure logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

class Chat:
    def __init__(self):
        logger.info('Initializing Chat Service')
        self.megaservice = ServiceOrchestrator()
        self.endpoint = '/james-is-great'
        self.host = '0.0.0.0'
        self.port = 8888
        self.remote_services = {}

    async def add_remote_services(self):
        """Initialize connections to remote services required for chat completion"""
        logger.info('Setting up remote services')
        try:
            # Add connection to VLLM service
            self.remote_services['vllm'] = await self.megaservice.connect_service(
                'vllm-service',
                'http://vllm-service:80/v1/completions'
            )
            
            # Add connection to embedding service
            self.remote_services['embedding'] = await self.megaservice.connect_service(
                'tei-embedding-service',
                'http://tei-embedding-service:8080/embed'
            )
            
            # Add connection to retriever service
            self.remote_services['retriever'] = await self.megaservice.connect_service(
                'retriever-service',
                'http://retriever:8081/retrieve'
            )
            
            # Add connection to reranking service
            self.remote_services['reranking'] = await self.megaservice.connect_service(
                'tei-reranking-service',
                'http://tei-reranking-service:8082/rerank'
            )
            
            logger.info('Successfully connected to all remote services')
        except Exception as e:
            logger.error(f'Failed to connect to remote services: {str(e)}')
            raise

    def start(self):
        """Start the chat service"""
        logger.info('Starting Chat Service')
        try:
            self.service = MicroService(
                self.__class__.__name__,
                service_role=ServiceRoleType.MEGASERVICE,
                host=self.host,
                port=self.port,
                endpoint=self.endpoint,
                input_datatype=ChatCompletionRequest,
                output_datatype=ChatCompletionResponse,
            )

            self.service.add_route(self.endpoint, self.handle_request, methods=["POST"])
            logger.info(f'Added route: {self.endpoint}')

            self.service.start()
            logger.info(f'Service started on {self.host}:{self.port}')
        except Exception as e:
            logger.error(f'Failed to start service: {str(e)}')
            raise

    async def handle_request(self, request: Request) -> ChatCompletionResponse:
        """Handle incoming chat completion requests"""
        try:
            # Parse request body
            body = await request.json()
            chat_request = ChatCompletionRequest(**body)
            logger.info(f'Received chat request: {chat_request}')

            # Validate request
            if not chat_request.messages:
                raise HTTPException(status_code=400, detail="No messages provided")

            # Process messages through the pipeline
            response = await self._process_chat_request(chat_request)
            
            logger.info('Successfully processed chat request')
            return response

        except Exception as e:
            logger.error(f'Error processing request: {str(e)}')
            raise HTTPException(status_code=500, detail=str(e))

    async def _process_chat_request(self, request: ChatCompletionRequest) -> ChatCompletionResponse:
        """Process a chat completion request through the service pipeline"""
        try:
            # 1. Get embeddings for the latest message
            latest_message = request.messages[-1].content
            embeddings = await self.remote_services['embedding'].embed_text(latest_message)

            # 2. Retrieve relevant context
            context = await self.remote_services['retriever'].get_relevant_context(embeddings)

            # 3. Rerank retrieved context
            reranked_context = await self.remote_services['reranking'].rerank_context(
                query=latest_message,
                contexts=context
            )

            # 4. Prepare prompt with context
            augmented_prompt = self._prepare_prompt(request.messages, reranked_context)

            # 5. Get completion from LLM
            completion = await self.remote_services['vllm'].complete(
                messages=augmented_prompt,
                temperature=request.temperature or 0.7,
                max_tokens=request.max_tokens or 150
            )

            # 6. Prepare response
            return ChatCompletionResponse(
                id=request.id,
                choices=[{
                    'message': Message(
                        role="assistant",
                        content=completion.text
                    ),
                    'finish_reason': completion.finish_reason
                }],
                usage=completion.usage
            )

        except Exception as e:
            logger.error(f'Error in chat pipeline: {str(e)}')
            raise

    def _prepare_prompt(self, messages: List[Message], context: List[str]) -> List[Message]:
        """Prepare the prompt by augmenting it with relevant context"""
        # Add context as system message
        context_message = Message(
            role="system",
            content=f"Context for answering: {' '.join(context)}"
        )
        return [context_message] + messages

if __name__ == '__main__':
    logger.info('Starting main application')
    chat = Chat()
    asyncio.run(chat.add_remote_services())
    chat.start()