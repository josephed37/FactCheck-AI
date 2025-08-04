import streamlit as st
from models import FactCheckResponse # Import from our models file
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
    "This app uses a RAG pipeline with live search results for up-to-date analysis."
)

statement = st.text_area(
    "Enter the statement to fact-check:",
    height=100,
    placeholder="e.g., Who won the last FIFA World Cup?",
)

analyze_button = st.button("Analyze Statement", type="primary")


# --- Logic and Response Handling ---
if analyze_button:
    if not statement.strip():
        st.warning("Please enter a statement to analyze.")
    else:
        with st.spinner("Performing live web search and AI analysis... 👨‍💻"):
            result: FactCheckResponse | None = get_fact_check_from_backend(statement)

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

            # --- NEW: Display the sources ---
            # Check if the sources list is not empty.
            if result.sources:
                st.subheader("Sources Used for Analysis")
                # Loop through each source and display it as a clickable link.
                for source in result.sources:
                    st.markdown(f"- [{source.title}]({source.url})")

        else:
            st.error(
                "Could not get a valid response from the backend API. "
                "Please ensure the Go server is running and accessible."
            )
