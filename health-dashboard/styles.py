# Modern Dark Glassmorphism Stylesheet
MAIN_STYLE = """
QMainWindow {
    background-color: #0f172a;
}

QWidget#CentralWidget {
    background-color: #0f172a;
}

QFrame#HeaderFrame {
    background-color: qlineargradient(x1:0, y1:0, x2:1, y2:0, stop:0 #1e293b, stop:1 #334155);
    border-radius: 12px;
    margin: 10px;
}

QLabel#TitleLabel {
    color: #f8fafc;
    font-size: 24px;
    font-weight: bold;
    font-family: 'Inter', sans-serif;
}

QFrame#Card {
    background-color: #1e293b;
    border: 1px solid #334155;
    border-radius: 15px;
    padding: 10px;
}

QLabel#CardTitle {
    color: #94a3b8;
    font-size: 14px;
    text-transform: uppercase;
    font-weight: 600;
}

QLabel#CardValue {
    color: #f8fafc;
    font-size: 28px;
    font-weight: bold;
}

QLabel#StatusActive {
    color: #4ade80;
    font-weight: bold;
}

QLabel#StatusInactive {
    color: #f87171;
    font-weight: bold;
}

QPushButton#ActionButton {
    background-color: #3b82f6;
    color: white;
    border-radius: 8px;
    padding: 10px 20px;
    font-weight: bold;
}

QPushButton#ActionButton:hover {
    background-color: #2563eb;
}
"""
