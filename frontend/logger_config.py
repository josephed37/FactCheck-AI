import logging
import os
from logging.handlers import RotatingFileHandler

def setup_logging():
    """
    Configures a centralized logger for the application.

    This sets up logging to both the console and a rotating file.
    """
    # Create logs directory if it doesn't exist
    if not os.path.exists('logs'):
        os.makedirs('logs')

    # Define the log format
    log_formatter = logging.Formatter(
        '%(asctime)s - %(name)s - %(levelname)s - %(message)s'
    )

    # Get the root logger
    logger = logging.getLogger()
    logger.setLevel(logging.INFO)

    # Avoid adding duplicate handlers
    if logger.hasHandlers():
        logger.handlers.clear()

    # Console Handler
    console_handler = logging.StreamHandler()
    console_handler.setFormatter(log_formatter)
    logger.addHandler(console_handler)

    # File Handler (with rotation)
    # Rotates logs when they reach 2MB, keeping 5 backup files.
    file_handler = RotatingFileHandler(
        'logs/app.log', maxBytes=2*1024*1024, backupCount=5
    )
    file_handler.setFormatter(log_formatter)
    logger.addHandler(file_handler)

    logging.info("Logging configured successfully.")