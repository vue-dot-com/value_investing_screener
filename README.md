# value_investing_screener
This repository contains code and files for the creation of a value investing screener using web scraping.

### Input
The code takes two .csv files containing tickers traded in the NASDAQ and the NYSE. You can download and put them in the working directory; here the links:
- NYSE http://www.nasdaq.com/screening/companies-by-industry.aspx?exchange=NYSE&render=download
- NASDAQ http://www.nasdaq.com/screening/companies-by-industry.aspx?exchange=NASDAQ&render=download

The .csv files should be replaced with updated ones whenever you want, as new IPOs take place or firms are delisted during time.

### Code
The code scrapes data from the https://www.gurufocus.com/ website for every ticker found. The process uses regular expressions to locate the numbers. Before running the code you should install the packages required. Download the `requirements.txt` in the repo in your local directory and run this in your terminal
```
pip install -r /path/to/requirements.txt
```
For each ticker, the code downloads (if available):
- Stock price
- ROIC (in %)	
- Owner Earnings per Share	
- 10y, 5y Revenue Growth	
- 10y, 5y EPS Growth	
- 10y, 5y Ebit Growth
- 10y, 5y Ebitda Growth
- 10y, 5y Free Cash Flow Growth
- 10y, 5y Dividend Growth	
- 10y, 5y BV Growth
- 10y, 5y Stock Price Growth	
- Ten Cap Valuation (computed)

### Output
The output is a pandas dataframe which is exported as a .csv file in your working directory. Here's a simple screenshot:
![image](https://user-images.githubusercontent.com/104139268/165929628-83de152f-2c82-4eff-8a53-ba5880e823de.png)

### Suggestions and improvements
I have made the work available to everyone, feel free to share, improve and commit to it. Inside the code you will find some `#TODO` which I wanted to implement.
Hope you will enjoy it and find it useful as much as I did.
