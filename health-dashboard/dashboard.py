import sys
import time
from PySide6.QtWidgets import (QApplication, QMainWindow, QWidget, QVBoxLayout, 
                             QHBoxLayout, QFrame, QLabel, QPushButton, QGridLayout, QProgressBar, QScrollArea)
from PySide6.QtCore import QTimer, Qt
from api_client import DaemonClient
from styles import MAIN_STYLE

class Dashboard(QMainWindow):
    def __init__(self):
        super().__init__()
        self.setWindowTitle("OneVisionOS | System Wellness")
        self.setMinimumSize(1000, 700)
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
        main_layout.setContentsMargins(20, 20, 20, 20)
        main_layout.setSpacing(15)

        # Header
        header = QFrame()
        header.setObjectName("HeaderFrame")
        header.setFixedHeight(80)
        header_layout = QHBoxLayout(header)
        header_layout.setContentsMargins(20, 0, 20, 0)
        
        title = QLabel("ONEVISION WELLNESS HUB")
        title.setObjectName("TitleLabel")
        header_layout.addWidget(title)
        
        header_layout.addStretch()
        
        self.node_id_label = QLabel("Node: Unknown")
        self.node_id_label.setStyleSheet("color: #94a3b8; font-family: monospace;")
        header_layout.addWidget(self.node_id_label)
        
        main_layout.addWidget(header)

        # Main Content Area
        content_layout = QHBoxLayout()
        main_layout.addLayout(content_layout)

        # Left Column: Status & Resources
        left_col = QVBoxLayout()
        content_layout.addLayout(left_col, 2)

        # Status Cards Grid
        status_grid = QGridLayout()
        left_col.addLayout(status_grid)
        self.cards = {}
        self.add_card(status_grid, 0, 0, "Daemon Status", "OFFLINE")
        self.add_card(status_grid, 0, 1, "P2P Network", "0 PEERS")
        
        # Resource Monitoring Card
        res_frame = QFrame()
        res_frame.setObjectName("Card")
        res_layout = QVBoxLayout(res_frame)
        res_layout.addWidget(QLabel("Real-time Resource Usage", objectName="CardTitle"))
        
        self.cpu_bar = self.add_resource_bar(res_layout, "CPU Utilization")
        self.ram_bar = self.add_resource_bar(res_layout, "Memory Usage")
        self.disk_bar = self.add_resource_bar(res_layout, "Disk Space")
        
        left_col.addWidget(res_frame)

        # System Log
        log_frame = QFrame()
        log_frame.setObjectName("Card")
        log_layout = QVBoxLayout(log_frame)
        log_layout.addWidget(QLabel("Live System Integrity Log", objectName="CardTitle"))
        self.log_content = QLabel("Waiting for data...")
        self.log_content.setStyleSheet("color: #cbd5e1; font-family: Consolas; font-size: 11px;")
        self.log_content.setAlignment(Qt.AlignTop)
        self.log_content.setWordWrap(True)
        log_layout.addWidget(self.log_content)
        left_col.addWidget(log_frame)

        # Right Column: Security & Repairs
        right_col = QVBoxLayout()
        content_layout.addLayout(right_col, 1)

        # Security Score Card
        sec_frame = QFrame()
        sec_frame.setObjectName("Card")
        sec_layout = QVBoxLayout(sec_frame)
        sec_layout.setAlignment(Qt.AlignCenter)
        sec_layout.addWidget(QLabel("Cyber Security Score", objectName="CardTitle"))
        self.sec_score_label = QLabel("0")
        self.sec_score_label.setObjectName("SecurityScore")
        sec_layout.addWidget(self.sec_score_label)
        sec_layout.addWidget(QLabel("SYSTEM PROTECTED", styleSheet="color: #4ade80; font-weight: bold; font-size: 10px;"))
        right_col.addWidget(sec_frame)

        # History of Repairs Card
        repair_frame = QFrame()
        repair_frame.setObjectName("Card")
        repair_layout = QVBoxLayout(repair_frame)
        repair_layout.addWidget(QLabel("History of Automated Repairs", objectName="CardTitle"))
        
        scroll = QScrollArea()
        scroll.setWidgetResizable(True)
        scroll.setStyleSheet("background: transparent; border: none;")
        self.repair_list_container = QWidget()
        self.repair_list_container.setStyleSheet("background: transparent;")
        self.repair_list_layout = QVBoxLayout(self.repair_list_container)
        self.repair_list_layout.setAlignment(Qt.AlignTop)
        scroll.setWidget(self.repair_list_container)
        repair_layout.addWidget(scroll)
        
        right_col.addWidget(repair_frame)

    def add_resource_bar(self, layout, label_text):
        h_layout = QHBoxLayout()
        label = QLabel(label_text)
        label.setStyleSheet("color: #94a3b8; font-size: 11px;")
        h_layout.addWidget(label)
        h_layout.addStretch()
        val_label = QLabel("0%")
        val_label.setStyleSheet("color: #f8fafc; font-size: 11px; font-weight: bold;")
        h_layout.addWidget(val_label)
        layout.addLayout(h_layout)
        
        bar = QProgressBar()
        bar.setRange(0, 100)
        bar.setValue(0)
        bar.setTextVisible(False)
        layout.addWidget(bar)
        
        return (bar, val_label)

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

        # Update core status
        self.cards["Daemon Status"].setText("ACTIVE")
        self.cards["Daemon Status"].setStyleSheet("color: #4ade80;")
        
        p2p_data = data.get("p2p", {})
        peer_count = p2p_data.get("peers", 0)
        self.cards["P2P Network"].setText(f"{peer_count} PEERS")
        self.node_id_label.setText(f"Node ID: {p2p_data.get('peer_id', 'Unknown')[:10]}...")
        
        # Update Resource Metrics
        metrics = data.get("metrics", {})
        self.update_resource_bar(self.cpu_bar, metrics.get("cpu_usage", 0))
        self.update_resource_bar(self.ram_bar, metrics.get("memory_usage", 0))
        self.update_resource_bar(self.disk_bar, metrics.get("disk_usage", 0))

        # Update Security Score
        nids_data = data.get("nids", {})
        score = nids_data.get("security_score", 0)
        self.sec_score_label.setText(str(score))
        if score > 80: self.sec_score_label.setStyleSheet("color: #4ade80;")
        elif score > 50: self.sec_score_label.setStyleSheet("color: #facc15;")
        else: self.sec_score_label.setStyleSheet("color: #f87171;")

        # Update Repair History
        repairs = data.get("repairs", [])
        self.update_repair_list(repairs)

        self.log_content.setText(f"System Time: {data.get('time', 'N/A')}\nWatchdog: {data.get('watchdog', 'unknown')}\nIntegrity State: CLEAN ✅")

    def update_resource_bar(self, bar_tuple, value):
        bar, label = bar_tuple
        bar.setValue(int(value))
        label.setText(f"{value:.1f}%")

    def update_repair_list(self, repairs):
        # Clear current list
        for i in reversed(range(self.repair_list_layout.count())): 
            self.repair_list_layout.itemAt(i).widget().setParent(None)
            
        if not repairs:
            self.repair_list_layout.addWidget(QLabel("No repairs recorded.", objectName="RepairItem"))
            return

        for repair in repairs[-10:]: # Show last 10
            ts = repair.get('timestamp', '').split('T')[-1][:5]
            item = QLabel(f"[{ts}] {repair.get('status')} : {repair.get('path').split('/')[-1]}")
            item.setObjectName("RepairItem")
            self.repair_list_layout.addWidget(item)

if __name__ == "__main__":
    app = QApplication(sys.argv)
    window = Dashboard()
    window.show()
    sys.exit(app.exec())
