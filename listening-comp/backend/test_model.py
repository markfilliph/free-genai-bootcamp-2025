from transformers import AutoTokenizer, T5ForConditionalGeneration
import torch

def test_model():
    # Initialize model and tokenizer
    model_name = "google/flan-t5-base"
    tokenizer = AutoTokenizer.from_pretrained(model_name)
    model = T5ForConditionalGeneration.from_pretrained(model_name)
    
    # Set model to evaluation mode
    model.eval()
    if torch.cuda.is_available():
        model = model.cuda()
    
    # Test prompts
    prompts = [
        "Translate to Japanese: hello",
        "Translate to Japanese: thank you",
        "Translate to Japanese: goodbye"
    ]
    
    print("Testing direct model responses:")
    print("-" * 50)
    
    for prompt in prompts:
        # Tokenize input
        inputs = tokenizer(prompt, return_tensors="pt", max_length=128, truncation=True)
        if torch.cuda.is_available():
            inputs = inputs.to('cuda')
        
        # Generate response
        with torch.no_grad():
            outputs = model.generate(
                **inputs,
                max_length=128,
                num_beams=4,
                temperature=0.7,
                do_sample=True
            )
        
        # Decode and print response
        response = tokenizer.decode(outputs[0], skip_special_tokens=True)
        print(f"\nPrompt: {prompt}")
        print(f"Response: {response}")

if __name__ == "__main__":
    test_model()
