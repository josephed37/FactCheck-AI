# Veritas-AI: Real-Time Fact-Checker

![Build Status](https://img.shields.io/badge/build-passing-brightgreen)
![License](https://img.shields.io/badge/license-MIT-blue)
![Go Version](https://img.shields.io/badge/go-1.22-blue.svg)
![Python Version](https://img.shields.io/badge/python-3.11-blue.svg)

A high-performance fake news detector using Google's Gemini AI, served via a Go backend and consumed by a Streamlit web interface. This project is designed to provide rapid, real-time analysis of statements to combat the spread of misinformation.

---

## 📸 Application Preview

![App Screenshot](./images/image%201.png)

---

## ✨ Core Features

- **AI-Powered Analysis**: Leverages the reasoning capabilities of Google's Gemini models to evaluate statements.
- **Structured Output**: Returns a clear verdict (`True`, `False`, `Uncertain`), a confidence level, and detailed reasoning.
- **Contextual Insights**: Provides additional context to help users understand the nuances of a claim.
- **Scalable Architecture**: Built with a high-performance Go backend to handle concurrent users efficiently.
- **Interactive Web UI**: A clean and simple user interface built with Streamlit for ease of use.

---

## 🛠️ Technology Stack

This project uses a modern, decoupled architecture to ensure performance and scalability.

**Frontend (Phase 1):**

- **Framework**: [Streamlit](https://streamlit.io/)
- **Language**: Python 3.11+
- **API Client**: `google-generativeai`
- **Data Validation**: `pydantic`
- **Resilience**: `tenacity` for API retries

**Backend (Phase 2):**

- **Language**: Go 1.22+
- **HTTP Server**: Gin or Fiber
- **Logging**: `logrus` or `zerolog`
- **Security**: Rate-limiting and secure API key management

**Deployment (Phase 3):**

- **Containerization**: Docker & Docker Compose
- **Hosting**: GCP Cloud Run / Vercel (Optional)

---

## 🧪 Testing & Automation

To ensure code quality and reliability, this project utilizes:

- **Frontend Testing**: `pytest` for unit tests on the Python-based Streamlit application.
- **Backend Testing**: Go's built-in `testing` package for unit tests on the API service.
- **CI/CD**: GitHub Actions to automatically run tests on every push and pull request to the `main` branch.

---

## 🚀 Project Roadmap

This project is being developed in three distinct phases:

- [x] **Phase 1: Prototype & Core Inference**
  - Build a working Streamlit prototype.
  - Connect directly to the Gemini API.
  - Implement data validation, logging, and error handling.

- [ ] **Phase 2: Golang-Based Model Serving**
  - Develop a high-performance Go API to act as middleware.
  - Offload all Gemini API interaction to the Go service.
  - Containerize the Go backend with Docker.

- [ ] **Phase 3: Full Integration & Advanced Features**
  - Connect the Streamlit frontend to the Go backend.
  - Integrate a database (SQLite/Firestore) for query history.
  - **Implement Retrieval-Augmented Generation (RAG)** by integrating a live search API to provide the model with real-time information.

---

## ⚙️ Setup and Installation

To run this project locally, follow these steps.

**Prerequisites:**

- Python 3.11+
- Go 1.22+ (for Phase 2)
- An active Google AI Studio API key.
- Docker (optional for Phase 3)

**1. Clone the repository:**

```bash
git clone https://github.com/josephed37/FactCheck-AI.git
cd FactCheck-AI
```

**2. Set up the Python environment (Phase 1):**

```bash
# Create and activate a virtual environment
python3 -m venv venv
source venv/bin/activate

# Install required packages
pip install -r requirements.txt
```

**3. Configure your API Key:**

- Create a file named `.env` in the root of the project.
- Add your Gemini API key to it:

    ```
    GEMINI_API_KEY="YOUR_API_KEY_HERE"
    ```

---

## ▶️ How to Run

To launch the Streamlit web application (Phase 1):

```bash
streamlit run frontend/app.py
```

Navigate to `http://localhost:8501` in your web browser to use the application.
