FROM python:3.10-slim

WORKDIR /app

COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

COPY main.py .
COPY suppression.py .
COPY match.json .
COPY config.ini .

CMD ["python", "main.py"]
