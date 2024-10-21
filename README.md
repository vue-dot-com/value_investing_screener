# Value Investing Screener

This Go-based project is designed to screen value investing metrics for various stock tickers. It retrieves data such as stock prices, enterprise values, ROIC (Return on Invested Capital), owner earnings, and growth numbers from [GuruFocus](https://www.gurufocus.com/). The output is then saved into a CSV file. Stocks in the NASDAQ and NYSE are pre-loaded and used as the base input source if no tickers are provided.

## Prerequisites

1. **Go**: Make sure Go is installed. You can download it [here](https://golang.org/doc/install).
2. **.env file**: You can configure your environment variables through a `.env` file or by exporting them directly into your terminal session.

## Installation

1. Clone this repository:

   ```bash
   git clone https://github.com/vue-dot-com/value_investing_screener.git
   ```

2. Navigate to the project directory:

   ```bash
   cd value_investing_screener
   ```

3. Create a `.env` file in the project root based on the example provided in `env.example`:

   ```bash
   cp env.example .env
   ```

4. Modify `.env` as needed for your specific environment:

   ```bash
   TICKERS=AAPL,GOOG # Example tickers
   MAX_CONCURRENCY=10 # Set the concurrency level
   OUTPUT_FILE=Screener.csv # Output CSV file location
   ```

## Usage

1. **Build the Project**:

   Before running the program, you need to build the Go executable:

   ```bash
   go build -o screener
   ```

2. **Run the Program**:

   You can run the program either by providing environment variables through the `.env` file or by exporting them directly in the terminal.

   - **Using the `.env` file**:

     After setting up your `.env` file, you can simply run the program:

     ```bash
     ./screener
     ```

   - **Without a `.env` file** (by exporting environment variables):

     You can export the environment variables directly in your terminal session:

     ```bash
     export TICKERS="AAPL,GOOG"
     export MAX_CONCURRENCY="10"
     export OUTPUT_FILE="usr/home/Screener.csv"

     ./screener
     ```

## Configuration

The following environment variables are used to configure the screener:

- **`TICKERS`**: A comma-separated list of tickers to process. If not set, the program will default to all tickers in NASDAQ and NYSE.
- **`MAX_CONCURRENCY`**: The maximum number of concurrent processes. Recommended max is 30. Default is 20.
- **`OUTPUT_FILE`**: The file path where the results will be saved in CSV format. Default is `Screener.csv`.

## Output

The program will save the results of the screener to the specified output file in CSV format. If no file path is provided, it defaults to `Screener.csv`.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change or add.
