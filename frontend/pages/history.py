import streamlit as st
from api_client import get_history_from_backend
from models import FactCheckHistoryItem

# --- Page Configuration ---
st.set_page_config(
    page_title="Fact-Check History",
    page_icon="📜",
    layout="wide", # Use wide layout to better display the history table
)

st.title("📜 Fact-Check History")
st.write("Here you can see all the statements that have been previously checked.")

# --- Fetch and Display History ---

# Add a button to manually refresh the history
if st.button("Refresh History"):
    st.cache_data.clear() # Clear the cache to force a new API call

# Use st.cache_data to cache the API response.
# This makes the app much faster and reduces API calls.
@st.cache_data
def fetch_history():
    return get_history_from_backend()

history_list = fetch_history()

if history_list is not None:
    if not history_list:
        st.info("No fact-checks have been made yet. Go to the main page to analyze a statement!")
    else:
        # Display each history item in a formatted way
        for item in history_list:
            with st.container():
                # Use an expander to keep the UI clean
                with st.expander(f"**{item.statement.strip()}** - Verdict: **{item.verdict}**"):
                    st.metric(label="Confidence", value=item.confidence)
                    st.markdown(f"**Reasoning:** \n*{item.reason}*")
                    st.markdown(f"**Additional Context:** \n*{item.additional_context}*")
                    
                    # Display the timestamp in a smaller, lighter font
                    st.markdown(f"<p style='text-align: right; color: grey; font-size: 0.8em;'>Checked on: {item.created_at}</p>", unsafe_allow_html=True)
                st.divider()
else:
    st.error("Failed to fetch history from the backend. Please ensure the backend server is running.")
