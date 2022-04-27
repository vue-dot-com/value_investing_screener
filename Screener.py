#This script tries to find possible stocks that applies for value investing principles.
##########IMPORT TICKERS#############
#importing CSV files with tickers and names
#You need to have those CSV files in your working directory. The link for the download is:
#NYSE: http://www.nasdaq.com/screening/companies-by-industry.aspx?exchange=NYSE&render=download
#NASDAQ: http://www.nasdaq.com/screening/companies-by-industry.aspx?exchange=NASDAQ&render=download
#They change when new IPOs (new listed companies) happen. Anyway update them every 1 or 2 years.
import csv, os
import requests
import bs4
import re
import time
import pickle
import pandas as pd
import yfinance as yf
import pandas as pd


nyse = pd.read_csv('NYSE.csv')
nasdaq = pd.read_csv('NASDAQ.csv')
stocks = pd.merge(nyse, nasdaq, how="outer")
#TODO: Code for downloading ticker prices. Example yfinance.download(ticker, threads=false, period=1d)
#creating ticker list

ticker = stocks['Symbol'].to_list()
company = stocks['Name'].to_list()
ipo_year = stocks['IPO Year'].to_list()
country = stocks['Country'].to_list()
sector = stocks['Sector'].to_list()
industry =stocks['Industry'].to_list()
quote = stocks['Last Sale'].to_list()
capitalization = stocks['Market Cap'].to_list()

print('Ticker list created with ' + str(len(ticker))+ ' elements')
print('Collecting URLs...')

#Storing URLs
evUrl = list([])
oeUrl = list([])
roicUrl = list([])
growthsUrl = list([])
summaryUrl = list([])
for i in ticker:
    evUrl.append('https://www.gurufocus.com/term/ev/'+str(i)+'/Enterprise-Value')
    oeUrl.append('https://www.gurufocus.com/term/Owner_Earnings/'+str(i)+'/Owner-Earnings-per-Share-(TTM)')
    roicUrl.append('https://www.gurufocus.com/term/ROIC/'+str(i)+'/ROIC-Percentage')
    growthsUrl.append('https://www.gurufocus.com/financials/'+str(i))
    summaryUrl.append('https://www.gurufocus.com/stock/'+str(i)+'/summary')

##########################################DEFINITIONS#########################################################
def adj_prices():
    n = 0
    pass
    while len(quote) < len(ticker):
        try:
            prices = yf.download(ticker[n], period='1d', group_by='ticker', threads=False)
            quote.append(prices.iloc[0, 4])
            n=n+1
        except:
            quote.append('')
            n=n+1
            pass
        print(str(len(quote)) + '/' + str(len(ticker)))

#def adj_prices():
#    while len(quote) < len(ticker):
#        import pandas_datareader.data as pdr
#        yf.pdr_override()
#        prices = yf.download(ticker, period='1d', threads=False)
#        try:
#            quote.append(prices.iloc[0, 4])

#        except:
#            quote.append('')
            #n=n+1
#        print (str(len(quote) )+ '/' + str(len(ticker)))
            #continue

#Retrieve Enterprise Value numbers
def ev_catcher():
    for u in evUrl[len(enterpriseValue):len(evUrl)] :
        res = requests.get(u)
        from bs4 import SoupStrainer

        # Match only font tag in the HTML source
        only_font_tag = SoupStrainer('font', style='font-size: 24px; font-weight: 700; color: #337ab7')
        evSoup = bs4.BeautifulSoup(res.text, 'lxml', parse_only=only_font_tag)
        evText = evSoup.select('font[style="font-size: 24px; font-weight: 700; color: #337ab7"]')
        # find EV number inside the page and attach it to the list roic

        evRegex = re.compile(r'(\$\d+([,\.]\d+)?\d+([,\.]\d+)?k?)')
        # \b(? < !\.)(?!0+(?:\.0+)? % )(?:\d |[1-9]\d | 100)(?:(? < !100)\.\d +)? % used before
        text = str(evText)
        try:
            enterpriseValue.append(evRegex.search(text)[0])
        except:
            enterpriseValue.append('')
        print(str(len(enterpriseValue)) + ' / ' + str(len(ticker)) + ' Done. ' + '( ' + str(
            len(enterpriseValue) / len(ticker) * 100) + '% )')

#Retrieve Owner Earnings Numbers
def oe_catcher():
    for u in oeUrl[len(ownerEarnings):len(oeUrl)]:
        res = requests.get(u)
        from bs4 import SoupStrainer

        # Match only font tag in the HTML source
        only_font_tag = SoupStrainer('font', style='font-size: 24px; font-weight: 700; color: #337ab7')
        oeSoup = bs4.BeautifulSoup(res.text, 'lxml', parse_only=only_font_tag)
        oeText = oeSoup.select('font[style="font-size: 24px; font-weight: 700; color: #337ab7"]')
        # find oe number inside the page and attach it to the list oe
        import re

        oeRegex = re.compile(r'(?:\+|\-|)(?<!\.)(?!0+(?:\.0+)?%)(?:\d|[1-9]\d|100)(?:(?<!100)\.\d+)')
        # \b(? < !\.)(?!0+(?:\.0+)? % )(?:\d |[1-9]\d | 100)(?:(? < !100)\.\d +)? % used before
        text = str(oeText)
        try:
            ownerEarnings.append(oeRegex.findall(text)[0])
        except:
            ownerEarnings.append('')
        print(str(len(ownerEarnings)) + ' / ' + str(len(ticker)) + ' Done. ' + '( ' + str(
            len(ownerEarnings) / len(ticker) * 100) + '% )')
        
#Retrieve Roic numbers
def ro_catcher():
    for u in roicUrl[len(roic):len(roicUrl)]:
        res = requests.get(u)
        from bs4 import SoupStrainer

        # Match only font tag in the HTML source
        only_tr_tag = SoupStrainer('font', style='font-size: 24px; font-weight: 700; color: #337ab7')
        roicSoup = bs4.BeautifulSoup(res.text, 'lxml', parse_only=only_tr_tag)
        # The code down here will select only the tag with that specific parameter, so you will not work with the entire HTML but only a fraction
        roicText = roicSoup.select('font[style="font-size: 24px; font-weight: 700; color: #337ab7"]')
        # find roic number inside the page and attach it to the list roic

        roicRegex = re.compile(r'(?:\+|\-|)(?<!\.)(?!0+(?:\.0+)?%)(?:\d|[1-9]\d|100)(?:(?<!100)\.\d+)?%')
        # \b(? < !\.)(?!0+(?:\.0+)? % )(?:\d |[1-9]\d | 100)(?:(? < !100)\.\d +)? % used before
        text = str(roicText)
        try:
            roic.append(roicRegex.findall(text)[0])
        except:
            roic.append('')
        print(str(len(roic)) + ' / ' + str(len(ticker)) + ' Done. ' + '( ' + str(
            len(roic) / len(ticker) * 100) + '% )')
        
#Retrieve Growth numbers
def gr_catcher():
    for u in growthsUrl[len(revenueGrowth10y):len(growthsUrl)]:
        res = requests.get(u)
        from bs4 import SoupStrainer

        only_tr_tag = SoupStrainer('tr')
        grSoup = bs4.BeautifulSoup(res.text, 'lxml', parse_only=only_tr_tag)
        grText = grSoup.select('tr')
        # find oe number inside the page and attach it to the list oe


        grRegex = re.compile(r'(?:\+|\-|)(?<!\.)(?!0+(?:\.0+)?%)(?:\d|[1-9]\d|100)(?:(?<!100)\.\d+)|(?:N\/A)')

        text = str(grText)

        # 10Y revenue growth
        try:

            if grRegex.findall(text)[4] == 'N/A':
                revenueGrowth10y.append(
                    grRegex.findall(text)[5])  # If it doesn't find 10y growth it looks for the 5y growth
            else:
                revenueGrowth10y.append(grRegex.findall(text)[4])
        except:
            revenueGrowth10y.append('')
        print(str(len(revenueGrowth10y)) + ' / ' + str(len(ticker)) + ' Done. ' + '( ' + str(
            len(revenueGrowth10y) / len(ticker) * 100) + '% )')


        #10Y eps growth
        try:
            if grRegex.findall(text)[14] == 'N/A':
                epsGrowth10y.append(
                    grRegex.findall(text)[14])  # If it doesn't find 10y growth it looks for the 5y growth
            else:
                epsGrowth10y.append(grRegex.findall(text)[13])
        except:
            epsGrowth10y.append('')
        print(str(len(epsGrowth10y)) + ' / ' + str(len(ticker)) + ' Done. ' + '( ' + str(
            len(epsGrowth10y) / len(ticker) * 100) + '% )')


        #10Y fcf growth
        try:
            if grRegex.findall(text)[16] == 'N/A':
                fcfGrowth10y.append(
                    grRegex.findall(text)[17])  # If it doesn't find 10y growth it looks for the 5y growth
            else:
                fcfGrowth10y.append(grRegex.findall(text)[16])
        except:
            fcfGrowth10y.append('')
        print(str(len(fcfGrowth10y)) + ' / ' + str(len(ticker)) + ' Done. ' + '( ' + str(
            len(fcfGrowth10y) / len(ticker) * 100) + '% )')

        #10Y bv growth
        try:
            if grRegex.findall(text)[19] == 'N/A':
                bvGrowth10y.append(
                    grRegex.findall(text)[20])  # If it doesn't find 10y growth it looks for the 5y growth
            else:
                bvGrowth10y.append(grRegex.findall(text)[19])
        except:
            bvGrowth10y.append('')
        print(str(len(bvGrowth10y)) + ' / ' + str(len(ticker)) + ' Done. ' + '( ' + str(
            len(bvGrowth10y) / len(ticker) * 100) + '% )')
        
#TODO: cash to asset, operating and profit margin
def summary_catcher():
    for u in summaryUrl[len(cash_to_debt):len(summaryUrl)]:
        res = requests.get(u)
        from bs4 import SoupStrainer

        only_td_tag = SoupStrainer('td')
        sumSoup = bs4.BeautifulSoup(res.text, 'lxml', parse_only=only_td_tag)
        sumText = sumSoup.select('td')

        sumRegex = re.compile(r'(?:\+|\-|)(?<!\.)(?!0+(?:\.0+)?)(?:\d|[1-9]\d|100)(?:(?<!100)\.\d+)|(?:N\/A)|\d+$|(0\.[1-9])')

        text = str(sumText)

        #Find cash to debt ratio
        try:

            if sumRegex.findall(text)[0] == 'N/A':
                cash_to_debt.append('')
            else:
                cash_to_debt.append(sumRegex.findall(text)[0])
        except:
            cash_to_debt.append('')
        print(str(len(cash_to_debt)) + ' / ' + str(len(ticker)) + ' Done. ' + '( ' + str(
            len(cash_to_debt) / len(ticker) * 100) + '% )')

    #TODO find operating margin
     #   try:

      #      if sumRegex.findall(text)[8] == 'N/A':
       #         op_margin.append('')
        #    else:
         #       op_margin.append(sumRegex.findall(text)[8])
       # except:
        #   op_margin.append('')
       # print(str(len(op_margin)) + ' / ' + str(len(ticker)) + ' Done. ' + '( ' + str(
        #    len(op_margin) / len(ticker) * 100) + '% )')


#threading for running all the functions above together
def download_thread():
    import threading
    import time

    e1 = threading.Thread(target=adj_prices, name='e1')
    e2 = threading.Thread(target=ev_catcher, name='e2')
    e3 = threading.Thread(target=ro_catcher, name='e3')
    e4 = threading.Thread(target=oe_catcher, name='e4')
    e5 = threading.Thread(target=gr_catcher, name='e5')

    origin_time =time.time()
    print('Starting downloading prices in 1 sec...')
    time.sleep(1.0)
    e1.start()
    time.sleep(0.05)
    print('Starting downloading Enterprise Values in 1 sec...')
    time.sleep(1.0)
    e2.start()
    time.sleep(0.05)
    print('Starting downloading ROIC in 1 sec...')
    time.sleep(1.0)
    e3.start()
    time.sleep(0.05)
    print('Starting downloading Owner Earnings in 1 sec...')
    time.sleep(1.0)
    e4.start()
    time.sleep(0.05)
    print('Starting downloading Growth numbers in 1 sec...')
    time.sleep(1.0)
    e5.start()
    time.sleep(0.05)

    #Wait until all finished
    e1.join()
    e2.join()
    e3.join()
    e4.join()
    e5.join()
    print('Download ended, it took' + str(time.time()-origin_time))

##########################################DEFINITIONS END#########################################################


print('Do you want to update the data? It will take a while if yes')

print('Below you can type:',
      '- yes (if you want to update all data)', '- restart (if the download failed)',
      '- no or press ENTER (if you do not want to update)', sep='\n')

update = input()
if update == 'yes':
    #1. ROIC (Return on invested capital). Measures how well a company generates cash flow relative to the capital it
    # has invested in its business. >10%
    roic = []
    #2. 10Y Revenue (per share) Growth Rate.
    revenueGrowth10y =[]
    #3. 10Y EPS Growth Rate.
    epsGrowth10y =[]
    #4. 10Y Free Cash Flow (per share) Growth Rate.
    fcfGrowth10y=[]
    #5. 10Y Book Value (per share) Growth Rate.
    bvGrowth10y =[]
    #a. Enterprise Value.
    enterpriseValue = []
    #b. Owner Earnings.
    ownerEarnings =[]
    #6. Summary financial
    cash_to_debt =[]
    op_margin=[]
    net_margin=[]

    print('starting download...')
    download_thread()

    print('Download ended')

#Save data into pickle file
    with open('Newobjs.pickle', 'wb') as f: # Python 2: open(...,'w')
        pickle.dump([quote, roic, enterpriseValue, ownerEarnings, revenueGrowth10y, epsGrowth10y, fcfGrowth10y, bvGrowth10y], f)

    print('Objects saved in the pickle file')

elif update == 'restart':
    with open('Newobjs.pickle', 'rb') as f:  # Python 3: open(..., 'rb')
        quote, roic, enterpriseValue, ownerEarnings, revenueGrowth10y, epsGrowth10y, fcfGrowth10y, bvGrowth10y = pickle.load(f)

    print('Objects loaded, re-starting the download process where it ended')

    download_thread()

    print('Download ended')

    with open('Newobjs.pickle', 'wb') as f: # Python 2: open(...,'w')
        pickle.dump([quote, roic, enterpriseValue, ownerEarnings, revenueGrowth10y, epsGrowth10y, fcfGrowth10y, bvGrowth10y], f)

    print('Objects save in pickle file')

else:
    import pickle
    with open('Newobjs.pickle', 'rb') as f:  # Python 3: open(..., 'rb')
        quote, roic, enterpriseValue, ownerEarnings,revenueGrowth10y, epsGrowth10y, fcfGrowth10y, bvGrowth10y = pickle.load(f)

    #load those URLs to put in the CSV file
    print('ROIC numbers are loaded and unchanged. Remember it is useful to update at least twice a year.')
    growthsUrl = list([])
    for i in ticker:
        growthsUrl.append('https://www.gurufocus.com/financials/'+str(i))

#create a dataframe where storing the data
data ={'Ticker': ticker,
       'Company': company,
       'Sector': sector,
       'Industry': industry,
       'Price': quote,
       'Enterprise Value (in Mill.)': enterpriseValue,
       'Market Cap': capitalization,
       'ROIC (in %)': roic,
       'Owner Earnings per Share': ownerEarnings,
       '10y Revenue Growth': revenueGrowth10y,
       '10y EPS Growth': epsGrowth10y,
       '10y FCF Growth': fcfGrowth10y,
       '10y BV Growth':bvGrowth10y,
       'Financial & News': growthsUrl}

df = pd.DataFrame(data, columns=['Ticker',
                                 'Company',
                                 'Sector',
                                 'Industry',
                                 'Price',
                                 'Market Cap',
                                 #'Enterprise Value (in Mill.)',
                                 'ROIC (in %)',
                                 'Owner Earnings per Share',
                                 '10y Revenue Growth',
                                 '10y EPS Growth',
                                 '10y FCF Growth',
                                 '10y BV Growth',
                                 'Financial & News'], index=ticker)
#Replace values
df['ROIC (in %)' ]=df['ROIC (in %)'].replace({'\%':''}, regex = True) # or use this code: list(map(lambda x: x[:-1], df['ROIC (in %)'].values))
df['10y Revenue Growth'] = df['10y Revenue Growth'].replace({'N\/A':''}, regex= True)
df['10y EPS Growth'] = df['10y EPS Growth'].replace({'N\/A':''}, regex= True)
df['10y FCF Growth'] = df['10y FCF Growth'].replace({'N\/A':''}, regex= True)
df['10y BV Growth'] = df['10y BV Growth'].replace({'N\/A':''}, regex= True)
#Convert dollar values into numeric
#df['Enterprise Value (in Mill.)'] = df['Enterprise Value (in Mill.)'].replace({'\$':''}, regex = True)
#df['Enterprise Value (in Mill.)'] = df['Enterprise Value (in Mill.)'].replace({'\,':''}, regex = True)

#Convert into numeric
df['Owner Earnings per Share'] = pd.to_numeric(df['Owner Earnings per Share'])
#df['Enterprise Value (in Mill.)'] = pd.to_numeric(df['Enterprise Value (in Mill.)'])
df['Price'] = pd.to_numeric(df['Price'])
df['Market Cap'] = pd.to_numeric(df['Market Cap'])
df['ROIC (in %)'] = pd.to_numeric(df['ROIC (in %)'])
df['10y Revenue Growth'] = pd.to_numeric(df['10y Revenue Growth'])
df['10y EPS Growth'] = pd.to_numeric(df['10y EPS Growth'])
df['10y FCF Growth'] = pd.to_numeric(df['10y FCF Growth'])
df['10y BV Growth'] = pd.to_numeric(df['10y BV Growth'])

tenCap = df['Owner Earnings per Share']*10
df.insert(12, 'Ten Cap  Valuation', tenCap)
#df['Ten Cap Valuation'] = pd.to_numeric(df['Ten Cap Valuation'])


print(df)
df.to_csv('Screener.csv', index=False, sep=';')
print(r'CSV file called ''Screener'' created')


#TODO: Roic 5y average
####5y roic####
#def ro5_catcher():
    #for u in roUrl[len(roic):len(roUrl)]:
        #res = requests.get(u)
        #from bs4 import SoupStrainer

        #only_tr_tag = SoupStrainer('tr')
        #roSoup = bs4.BeautifulSoup(res.text, 'lxml', parse_only=only_tr_tag)
        #roText = roSoup.select('tr')
        # find oe number inside the page and attach it to the list oe

        #roRegex = re.compile(r'(?:\+|\-|)(?<!\.)(?!0+(?:\.0+)?%)(?:\d|[1-9]\d|100)(?:(?<!100)\.\d+)|(?:N\/A)')

        #text = str(roText)


        #try:

            #if roRegex.findall(text)[3] == 'N/A':
                #roic.append(
                    #mean(roRegex.findall(text)))  # If it doesn't find 10y growth it looks for the 5y growth
            #else:
                #ro_catcher()
        #except:
            #roic.append('')
        #print(str(len(revenueGrowth10y)) + ' / ' + str(len(ticker)) + ' Done. ' + '( ' + str(
            #len(revenueGrowth10y) / len(ticker) * 100) + '% )')



#TODO: Payback time. Grow FCF to a windage growth rate and see what is the payback time
#TODO: Try multiprocessing