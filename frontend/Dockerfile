# --- Stage 1: The Base Image ---
# We start with an official Python image. Using a specific version ensures
# that our environment is reproducible. The 'slim' variant is smaller than the
# full image, which is good for production.
FROM python:3.11-slim-buster


# --- Environment Setup ---
# Set the working directory inside the container. All subsequent commands
# will be run from this path.
WORKDIR /app

# --- Dependency Management ---
# Copy the requirements.txt file into the container first.
# This is a caching optimization. Docker will only re-run the pip install
# step if this file changes.
COPY requirements.txt .

# Install the Python dependencies.
# --no-cache-dir reduces the image size by not storing the download cache.
RUN pip install --no-cache-dir -r requirements.txt

# --- Application Code ---
# Copy the rest of the frontend application code into the container.
COPY . .

# --- Port Exposure ---
# Expose port 8501, which is the default port that Streamlit runs on.
# This lets Docker know which port the application will listen on.
EXPOSE 8501

# --- Run Command ---
# This is the command that will be executed when the container starts.
# It runs the streamlit application.
# --server.port=8501 explicitly sets the port.
# --server.address=0.0.0.0 makes the app accessible from outside the container.
CMD ["streamlit", "run", "app.py", "--server.port=8501", "--server.address=0.0.0.0"]
