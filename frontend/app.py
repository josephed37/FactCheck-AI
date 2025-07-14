import streamlit as st
from gemini_api import get_fact_check, FactCheckResponse

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
    "The AI will provide a verdict, confidence level, and context."
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
        # Handle empty input
        st.warning("Please enter a statement to analyze.")
    else:
        # Show a spinner while processing
        with st.spinner("Analyzing... The AI is thinking 🤔"):
            result: FactCheckResponse | None = get_fact_check(statement)

        # Display results once processing is complete
        st.divider()
        if result:
            st.subheader("Analysis Complete")

            # Display verdict with a colored box and icon
            if result.verdict == "True":
                st.success(f"✅ Verdict: **{result.verdict}**")
            elif result.verdict == "False":
                st.error(f"❌ Verdict: **{result.verdict}**")
            else:
                st.warning(f"🤔 Verdict: **{result.verdict}**")

            # Display confidence and reason in columns
            col1, col2 = st.columns(2)
            with col1:
                st.metric(label="Confidence Level", value=result.confidence)
            with col2:
                st.info(f"**Reasoning:**\n\n{result.reason}")
            
            # Display additional context
            st.info(f"**Additional Context:**\n\n{result.additional_context}")

        else:
            # Handle cases where the API call failed
            st.error(
                "Could not get a valid response from the AI. "
                "This might be due to a network issue or an unclear API response. "
                "Please try again or rephrase your statement."
            )