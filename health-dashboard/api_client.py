import requests
import time

class DaemonClient:
    def __init__(self, host="localhost", port=8081):
        self.base_url = f"http://{host}:{port}"

    def get_health(self):
        try:
            response = requests.get(f"{self.base_url}/health", timeout=2)
            if response.status_code == 200:
                return response.json()
        except Exception as e:
            return {"error": str(e), "daemon": "inactive"}
        return {"daemon": "inactive"}
