# frontend/gemini_api.py

import os
import json
import logging
from typing import Literal

from dotenv import load_dotenv
import google.generativeai as genai
from pydantic import BaseModel, Field, ValidationError
from tenacity import retry, stop_after_attempt, wait_fixed
from pathlib import Path

load_dotenv()
# Get a logger instance for this module
logger = logging.getLogger(__name__) 

# --- Configuration ---

# Load API key from .env file.
genai.configure(api_key=os.getenv("GEMINI_API_KEY"))

# Define the safety settings for the generative model.
SAFETY_SETTINGS = [
    {"category": "HARM_CATEGORY_HARASSMENT", "threshold": "BLOCK_NONE"},
    {"category": "HARM_CATEGORY_HATE_SPEECH", "threshold": "BLOCK_NONE"},
    {"category": "HARM_CATEGORY_SEXUALLY_EXPLICIT", "threshold": "BLOCK_NONE"},
    {"category": "HARM_CATEGORY_DANGEROUS_CONTENT", "threshold": "BLOCK_NONE"},
]


# --- Pydantic Data Models (Our Data Contract) ---

class FactCheckResponse(BaseModel):
    """
    A Pydantic model to validate the structured response from the Gemini API.
    
    This acts as a data contract, ensuring the AI's output is predictable.
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


# --- Core API Function ---
@retry(stop=stop_after_attempt(3), wait=wait_fixed(2))
def get_fact_check(statement: str) -> FactCheckResponse | None:
    """
    Analyzes a statement using the Gemini API and returns a structured fact-check.

    This function sends a carefully crafted prompt to the AI and validates
    the response using the FactCheckResponse Pydantic model.

    Args:
        statement: The user-provided statement to be fact-checked.

    Returns:
        An instance of FactCheckResponse if the API call and validation succeed.
        None if the API fails, the response is malformed, or validation fails.
    """
    prompt_template = f"""
    You are an expert AI fact-checker. Your sole purpose is to analyze a statement
    and return a structured JSON response. Do not add any conversational text or
    markdown formatting.

    Analyze the following statement: "{statement}"

    Respond ONLY with a JSON object in this exact format:
    {{
      "verdict": "True", "False", or "Uncertain",
      "confidence": "Low", "Medium", or "High",
      "reason": "Your brief reasoning here.",
      "additional_context": "Additional clarifying facts here."
    }}
    """

    try:
        prompt_path = Path(__file__).parent.parent / "prompts" / "fact_check_prompt.txt"

        with open(prompt_path, 'r') as f:
            prompt_template = f.read()
        
        # Inject the user's statement into the loaded template
        final_prompt = prompt_template.format(statement=statement)
        
        model = genai.GenerativeModel('models/gemini-2.0-flash')
        response = model.generate_content(
            final_prompt,
            generation_config=genai.GenerationConfig(
                response_mime_type="application/json"
            ),
            safety_settings=SAFETY_SETTINGS,
        )

        # The response text is parsed from a string into a Python dictionary.
        response_data = json.loads(response.text)

        # Pydantic validates the dictionary against our model.
        validated_response = FactCheckResponse(**response_data)
        logger.info(f"Successfully fact-checked statement: '{statement}'")

        return validated_response

    except (ValidationError, json.JSONDecodeError) as e:
        # Catches errors if the AI's response is not valid JSON
        # or doesn't match our Pydantic model.
        print(f"Data validation error: {e}")
        return None
    except Exception as e:
        # Catches any other potential errors (e.g., API connection issues).
        print(f"An unexpected error occurred: {e}")
        return None