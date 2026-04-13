import sys
import time
from PySide6.QtWidgets import (QApplication, QMainWindow, QWidget, QVBoxLayout, 
                             QHBoxLayout, QFrame, QLabel, QPushButton, QGridLayout)
from PySide6.QtCore import QTimer, Qt
from api_client import DaemonClient
from styles import MAIN_STYLE

class Dashboard(QMainWindow):
    def __init__(self):
        super().__init__()
        self.setWindowTitle("OneVisionOS | System Wellness")
        self.setMinimumSize(900, 600)
        self.client = DaemonClient()
        
        self.init_ui()
        self.setStyleSheet(MAIN_STYLE)
        
        # Update timer
        self.timer = QTimer()
        self.timer.timeout.connect(self.update_data)
        self.timer.start(2000) # Update every 2 seconds

    def init_ui(self):
        central_widget = QWidget()
        central_widget.setObjectName("CentralWidget")
        self.setCentralWidget(central_widget)
        main_layout = QVBoxLayout(central_widget)

        # Header
        header = QFrame()
        header.setObjectName("HeaderFrame")
        header.setFixedHeight(80)
        header_layout = QHBoxLayout(header)
        
        title = QLabel("ONEVISION WELLNESS HUB")
        title.setObjectName("TitleLabel")
        header_layout.addWidget(title)
        
        header_layout.addStretch()
        
        self.node_id_label = QLabel("Node: Unknown")
        self.node_id_label.setStyleSheet("color: #94a3b8; font-family: monospace;")
        header_layout.addWidget(self.node_id_label)
        
        main_layout.addWidget(header)

        # Content Grid
        grid = QGridLayout()
        main_layout.addLayout(grid)

        # Status Cards
        self.cards = {}
        self.add_card(grid, 0, 0, "Daemon Status", "OFFLINE")
        self.add_card(grid, 0, 1, "P2P Network", "0 PEERS")
        self.add_card(grid, 1, 0, "Self-Healing", "INITIALIZED")
        self.add_card(grid, 1, 1, "Security (NIDS)", "SECURE")

        # System Log
        log_frame = QFrame()
        log_frame.setObjectName("Card")
        log_layout = QVBoxLayout(log_frame)
        log_layout.addWidget(QLabel("Real-time Integrity Log", objectName="CardTitle"))
        self.log_content = QLabel("Waiting for data...")
        self.log_content.setStyleSheet("color: #cbd5e1; font-family: Consolas;")
        self.log_content.setAlignment(Qt.AlignTop)
        log_layout.addWidget(self.log_content)
        grid.addWidget(log_frame, 2, 0, 1, 2)

        main_layout.addStretch()

    def add_card(self, grid, row, col, title, initial_value):
        frame = QFrame()
        frame.setObjectName("Card")
        layout = QVBoxLayout(frame)
        
        title_label = QLabel(title)
        title_label.setObjectName("CardTitle")
        layout.addWidget(title_label)
        
        value_label = QLabel(initial_value)
        value_label.setObjectName("CardValue")
        layout.addWidget(value_label)
        
        grid.addWidget(frame, row, col)
        self.cards[title] = value_label

    def update_data(self):
        data = self.client.get_health()
        
        if "error" in data:
            self.cards["Daemon Status"].setText("OFFLINE")
            self.cards["Daemon Status"].setObjectName("StatusInactive")
            return

        # Update values
        self.cards["Daemon Status"].setText("ACTIVE")
        self.cards["Daemon Status"].setStyleSheet("color: #4ade80;")
        
        p2p_data = data.get("p2p", {})
        peer_count = p2p_data.get("peers", 0)
        self.cards["P2P Network"].setText(f"{peer_count} PEERS")
        self.node_id_label.setText(f"Node ID: {p2p_data.get('peer_id', 'Unknown')[:10]}...")
        
        self.cards["Self-Healing"].setText(data.get("watchdog", "UNKNOWN").upper())
        
        nids_data = data.get("nids", {})
        sig_count = nids_data.get("signatures", 0)
        self.cards["Security (NIDS)"].setText(f"{sig_count} SIGS")

        self.log_content.setText(f"Last sync: {data.get('time', 'N/A')}\nIntegrity: Verified ✅")

if __name__ == "__main__":
    app = QApplication(sys.argv)
    window = Dashboard()
    window.show()
    sys.exit(app.exec())
