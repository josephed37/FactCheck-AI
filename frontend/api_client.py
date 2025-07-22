import requests
from pydantic import ValidationError
from gemini_api import FactCheckResponse

# We use the service name 'backend' as defined in our docker-compose.yaml file.
API_BASE_URL = "http://backend:8080/v1"

def get_fact_check_from_backend(statement: str) -> FactCheckResponse | None:
    """
    Sends a statement to our Go backend API and gets the fact-check result.

    Args:
        statement: The user-provided statement to be fact-checked.

    Returns:
        An instance of FactCheckResponse if the API call is successful and
        the response data is valid.
        None if any error occurs (e.g., network issue, server error, bad data).
    """
    try:
        endpoint = f"{API_BASE_URL}/fact-check"
        payload = {"statement": statement}
        response = requests.post(endpoint, json=payload, timeout=20)
        response.raise_for_status()
        response_data = response.json()
        validated_response = FactCheckResponse(**response_data)
        return validated_response

    except requests.exceptions.RequestException as e:
        print(f"API request error: {e}")
        return None
    except ValidationError as e:
        print(f"Backend response validation error: {e}")
        return None
    except Exception as e:
        print(f"An unexpected error occurred in the api_client: {e}")
        return None
