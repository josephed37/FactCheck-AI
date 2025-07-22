import streamlit as st
# We still need FactCheckResponse for type hinting
from gemini_api import FactCheckResponse 
# UPDATED: Import the new function that calls our backend
from api_client import get_fact_check_from_backend
from logger_config import setup_logging

# --- Initialize Logging ---
setup_logging()

# --- Page Configuration ---
st.set_page_config(
    page_title="Fact-Checker AI",
    page_icon="🔎",
    layout="centered",
)

# --- UI Elements ---
st.title("🔎 Real-Time AI Fact-Checker")
st.write(
    "Enter a statement to check for its factual accuracy. "
    "This app is powered by a Go backend for speed and scalability."
)

statement = st.text_area(
    "Enter the statement to fact-check:",
    height=100,
    placeholder="e.g., The Eiffel Tower is in London.",
)

analyze_button = st.button("Analyze Statement", type="primary")


# --- Logic and Response Handling ---
if analyze_button:
    if not statement.strip():
        st.warning("Please enter a statement to analyze.")
    else:
        with st.spinner("Contacting the Go backend... The AI is thinking 🤔"):
            # UPDATED: Call the new function that talks to our Go backend
            result: FactCheckResponse | None = get_fact_check_from_backend(statement)

        # The rest of the result display logic is IDENTICAL, because
        # our new function returns the same data structure.
        st.divider()
        if result:
            st.subheader("Analysis Complete")

            if result.verdict == "True":
                st.success(f"✅ Verdict: **{result.verdict}**")
            elif result.verdict == "False":
                st.error(f"❌ Verdict: **{result.verdict}**")
            else:
                st.warning(f"🤔 Verdict: **{result.verdict}**")

            col1, col2 = st.columns(2)
            with col1:
                st.metric(label="Confidence Level", value=result.confidence)
            with col2:
                st.info(f"**Reasoning:**\n\n{result.reason}")
            
            st.info(f"**Additional Context:**\n\n{result.additional_context}")

        else:
            st.error(
                "Could not get a valid response from the backend API. "
                "Please ensure the Go server is running and accessible."
            )
