import pytube
import requests
from bs4 import BeautifulSoup

#Find video
def find_video():
    page_url = 'https://firstsiteguide.com/what-is-blog/'
    req = requests.get(page_url)
    soup = BeautifulSoup(req.content, "html.parser")
    iframes = soup.find_all('iframe', {'class':"lazy-iframe"})
    return iframes

# Download Youtube video
def dwld_video(url):
    url = pytube.YouTube(url)
    vdo = url.streams.first()
    vdo.download()
    print("Download successfull!!!")

iframes = find_video()
for iframe in iframes:
    url = iframe['data-src']
    if('www.youtube.com' in url):
        dwld_video(url)
    
    

