# FactCheck-AI: Real-Time RAG Fact-Checker

FactCheck-AI is a full-stack, real-time fact-checking application designed to combat misinformation. It leverages a sophisticated **Retrieval-Augmented Generation (RAG)** pipeline, using live web search results to provide Google's Gemini AI with up-to-the-minute context for accurate, trustworthy analysis.

## 📸 Application Preview

![App Screenshot](./images/image%202.png)

## ✨ Core Features

* **Real-Time Analysis**: Solves the "knowledge cut-off" problem by using the Tavily Search API to fetch live information from the web for every query.
* **RAG Pipeline**: Augments prompts to the Gemini AI with real-time context, ensuring analysis is based on the latest information.
* **Transparent Sourcing**: Displays the web sources used for each fact-check, allowing users to verify the information themselves.
* **Scalable Go Backend**: A high-performance backend written in Go handles concurrent requests efficiently.
* **Interactive Web UI**: A clean and simple user interface built with Streamlit.
* **Persistent History**: Saves all fact-checks to a SQLite database and displays them on a dedicated history page.
* **Containerized**: The entire full-stack application is orchestrated with Docker and Docker Compose for easy setup and deployment.

## ⚙️ How It Works: The RAG Pipeline

When a user submits a statement, the application follows a three-step process:

1. **Retrieve**: The Go backend sends the user's statement as a query to the **Tavily Search API**.
2. **Augment**: The backend takes the top search results and constructs a new, "augmented" prompt. This prompt includes both the live search context and the user's original statement.
3. **Generate**: This augmented prompt is sent to the **Google Gemini API**. The AI then generates its fact-check analysis based primarily on the fresh, real-time context provided, not its own internal (and potentially outdated) knowledge.

## 🛠️ Technology Stack

* **Frontend**: Streamlit, Python
* **Backend**: Go, Gin Web Framework
* **AI & Search**: Google Gemini API, Tavily Search API
* **Database**: SQLite
* **Containerization**: Docker, Docker Compose

## 🚀 Project Status: Complete

* [x] **Phase 1: Prototype & Core Inference**
* [x] **Phase 2: Golang-Based Model Serving**
* [x] **Phase 3: Full Integration, Database & RAG Pipeline**

## ⚙️ Setup and Installation

**Prerequisites:**

* Docker and Docker Compose
* A `.env` file in the project root with your API keys:

    ```
    GEMINI_API_KEY="YOUR_GEMINI_KEY"
    TAVILY_API_KEY="YOUR_TAVILY_KEY"
    ```

**1. Clone the repository:**

```bash
git clone https://github.com/josephed37/FactCheck-AI.git
cd FactCheck-AI
```

**2. Run the application:**
The entire application stack can be launched with a single command:

```bash
docker compose up --build
```

* The Streamlit frontend will be available at `http://localhost:8501`.
* The Go backend API will be available at `http://localhost:8080`.

To stop the application, press `Ctrl+C` in the terminal where it's running, and then run:

```bash
docker compose down
