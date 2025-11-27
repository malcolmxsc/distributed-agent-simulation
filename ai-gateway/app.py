# Application Logic will go here

import time
import random
import uvicorn
from fastapi import FastAPI, Response
from pydantic import BaseModel
from prometheus_client import Counter, Histogram, generate_latest, CONTENT_TYPE_LATEST

#1. initialize the app
app = FastAPI(title="Apple SFTWR Evaluation Gateway")


# 2. Define metrics (Sensors)
# These allow us to monitor the app's health and safety performance 

# Tracks total requests, broken down by "persona" and "status" (safe/unsafe)
REQUEST_COUNT = Counter(
    'eval_requests_total',
    'Total AI requests processed',
    ['persona','status']

)

# Tracks how long the "AI" takes to generate a response.
LATENCY = Histogram(
    'llm_latency_seconds',
    'Time spent processing LLM request'
)

#Tracks how many times the AI said something forbidden

SAFETY_VIOLATIONS = Counter(
    'safety_violation_total',
    'Total safety violations detected',
)

# 3. Define the Data Model
# This ensures that any data sent to our API is structured correctly.

class ChatRequest(BaseModel):
    prompt: str
    persona: str # E.G. "Expert Software Engineer", "Beginner Software Engineer" 'Hacker'


# 4. define the chat endpoint

@app.post("/chat")
async def chat_endpoint(request: ChatRequest):
    """
    Simulates an LLM interaction and runs an Evaluation check.
    """
    start_time = time.time()



# --- Phase 1: Simulation ---
# simulate the "thinking" time of a real ai model 0.1 to 1 second.

    time.sleep(random.uniform(0.1, 1.0))

    #Generate a dummy response
    response_text = f"Response to {request.persona}: {request.prompt[:20]}..."

    # Phase 2: Evaluation (the judge) ---

    #simulate a safety check.

    #we roll a die: there is a 10% chance of a safety violation.

    is_safe = True
    if random.random() < 0.1:
        is_safe = False
        SAFETY_VIOLATIONS.inc() # increment the safety violation counter
        response_text = "[BLOCKED BY SAFETY CHECK]"


    # Phase 3: Observability metrics collection ---

    #calculate duration.

    process_time = time.time() - start_time
    LATENCY.observe(process_time)
    

    # log the result status
    status = "safe" if is_safe else "unsafe"
    REQUEST_COUNT.labels(persona=request.persona, status=status).inc()

    #return the json response
    return { "response": response_text,
    "evaluation": {
        "is_safe": is_safe,
        "latency_ms": round(process_time * 1000, 2)
    }
    }

    # 5. Define the Metrics Endpoint
    # This is where Prometheus/Grafana will look to gather data.

@app.get("/metrics")
def metrics():
    return Response(generate_latest(), media_type=CONTENT_TYPE_LATEST)

# 6. entry point
# this allows you to run the app.

if __name__ == "__main__":
    uvicorn.run(app, host="0.0.0.0", port=8000)