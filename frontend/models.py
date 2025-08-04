# frontend/models.py

from pydantic import BaseModel, Field
from typing import List, Literal

# NEW: A Pydantic model for a single source.
class Source(BaseModel):
    title: str
    url: str

class FactCheckResponse(BaseModel):
    """
    A Pydantic model to validate the structured response for a fact-check.
    """
    verdict: Literal["True", "False", "Uncertain"]
    confidence: Literal["Low", "Medium", "High"]
    reason: str
    additional_context: str
    # NEW: Add the sources field, which is a list of Source objects.
    # We default it to an empty list to prevent errors if it's missing.
    sources: List[Source] = []

class FactCheckHistoryItem(BaseModel):
    """
    Defines the structure for a single record from the database's history.
    """
    id: int
    statement: str
    verdict: str
    confidence: str
    reason: str
    additional_context: str
    created_at: str
