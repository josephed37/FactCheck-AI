import requests
from typing import List
from pydantic import ValidationError
from models import FactCheckResponse, FactCheckHistoryItem 

# The address of our Go backend service inside the Docker network.
API_BASE_URL = "http://backend:8080/v1"

def get_fact_check_from_backend(statement: str) -> FactCheckResponse | None:
    """
    Sends a statement to our Go backend API and gets the fact-check result.
    """
    try:
        endpoint = f"{API_BASE_URL}/fact-check"
        payload = {"statement": statement}
        response = requests.post(endpoint, json=payload, timeout=20)
        response.raise_for_status()
        response_data = response.json()
        validated_response = FactCheckResponse(**response_data)
        return validated_response
    except (requests.exceptions.RequestException, ValidationError) as e:
        print(f"Error in get_fact_check_from_backend: {e}")
        return None

def get_history_from_backend() -> List[FactCheckHistoryItem] | None:
    """
    Fetches all saved fact-check records from the Go backend API.
    """
    try:
        endpoint = f"{API_BASE_URL}/history"
        response = requests.get(endpoint, timeout=10)
        response.raise_for_status()
        history_data = response.json()
        validated_history = [FactCheckHistoryItem(**item) for item in history_data]
        return validated_history
    except (requests.exceptions.RequestException, ValidationError) as e:
        print(f"Error in get_history_from_backend: {e}")
        return None
