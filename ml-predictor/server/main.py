from fastapi import FastAPI, Request
from pydantic import BaseModel
import random

app = FastAPI()

class PredictRequest(BaseModel):
    key: str

@app.post("/predict")
def predict_ttl(req: PredictRequest):
    # Simulate ML prediction with random TTL (10–60 seconds)
    ttl_seconds = random.randint(10, 60)
    return {"ttl": ttl_seconds}
