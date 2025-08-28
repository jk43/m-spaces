ls -al /
# For docker compose
# Install dependencies
pip install -r requirements.txt

# Start the FastAPI application
#uvicorn main:app --reload --port=80 --host=0.0.0.0
python main.py