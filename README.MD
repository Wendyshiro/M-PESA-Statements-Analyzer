🧾 M-Pesa Statement Analyzer
Overview

The M-Pesa Statement Analyzer is a Go-based tool for parsing, analyzing, and visualizing financial data from M-Pesa statements.
It converts PDF statements into structured data, then generates summaries and visual insights to help users track spending, income, and transaction trends.

✨ Features

PDF Parsing – Extracts transaction data from official M-Pesa statements

Summary Reports – Calculates total inflows, outflows, and category-based spending

Category Classification – Groups transactions into categories like Food, Transport, Utilities, and Savings

Visualization Dashboard – Displays charts and trends using Chart.js

Offline-Ready – Runs locally without needing an internet connection

🧠 Tech Stack

Component	Purpose

Go	Backend logic for parsing and data analysis

pdftotext	Converts PDFs to plain text

HTML / Chart.js	Frontend visualization

JSON / CSV	Data storage and summary export

OPENAI- For improved categorization

🧰 Installation

Clone this repository

→ git clone https://github.com/<your-username>/mpesa-statement-analyzer.git

→ cd mpesa-statement-analyzer

Run the Go program

→ go run main.go

💡 Tip: Make sure pdftotext is installed and added to your system path.

⚙️ Usage

Option 1 – Web Interface

1.Run the app with go run main.go.

2.Open your browser and go to http://localhost:/8000/dashboard

3.Upload your M-Pesa statement via the upload form

4.View your financial summary and interactive charts instantly

Option 2 - API/ Postman

1.Send a POST request to /upload with the PDF file as multipart/form data.

2.The response returns parsed JSON data and computed totals

3.Example endpoint:
-POST http://localhost:8080/upload

📊 Example Output

Sample JSON Response:
json
{
  "total_inflows": 18500,
  "total_outflows": 16500,
  "balance": 2000,
  "categories": {
    "Food & Dining": 8450,
    "Transport": 3200,
    "Bills & Utilities": 5100,
    "Savings": 2000
  }
}

🧩 Roadmap

Docker containerization

Cloud Run deployment

Improved category classifier (AI/NLP)

Persistent user dashboard

WhatsApp or SMS summary notifications

🤝 Contributing

Contributions are welcome!
If you’d like to add new features, fix bugs, or improve performance, open a pull request or issue to start the discussion.

📜 License

Licensed under the MIT License.

💬 Acknowledgements

Thanks to the open-source Go and Chart.js communities — and the M-Pesa ecosystem — for enabling accessible financial transparency tools.








