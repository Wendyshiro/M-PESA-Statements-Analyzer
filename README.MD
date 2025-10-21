🧾 M-Pesa Statement Analyzer
Overview

The M-Pesa Statement Analyzer is a lightweight tool for parsing, analyzing, and visualizing financial data from M-Pesa statements.
It converts PDF statements into structured data, then generates summaries and visual insights to help users track spending, income, and transaction trends.

✨ Features

📂 PDF Parsing: Extracts transaction data from official M-Pesa statements.

📊 Summary Reports: Calculates total inflows, outflows, and category-based spending.

💰 Category Classification: Groups transactions into categories such as Food, Transport, Utilities, and Savings.

📈 Visualization Dashboard: Presents trends through charts (powered by Chart.js).

⚡ Fast & Offline-Ready: Runs locally without needing an internet connection or external APIs.

🧠 Tech Stack
Component	Purpose
Go	Backend logic for parsing and data analysis
pdftotext	Converts PDFs to plain text
HTML / Chart.js	Frontend visualization
JSON / CSV	Data storage and summary export