from pydantic import BaseModel, Field
from typing import Literal

class FactCheckResponse(BaseModel):
    """
    A Pydantic model to validate the structured response for a fact-check.
    """
    verdict: Literal["True", "False", "Uncertain"] = Field(
        ..., description="The final verdict on the statement."
    )
    confidence: Literal["Low", "Medium", "High"] = Field(
        ..., description="The confidence level of the verdict."
    )
    reason: str = Field(
        ..., description="A brief explanation for the verdict."
    )
    additional_context: str = Field(
        ..., description="Any additional context to clarify the facts."
    )

class FactCheckHistoryItem(BaseModel):
    """
    Defines the structure for a single record retrieved from our database's history.
    """
    id: int
    statement: str
    verdict: str
    confidence: str
    reason: str
    additional_context: str
    created_at: str

