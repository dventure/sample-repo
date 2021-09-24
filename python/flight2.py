import requests

HEADERS = {
    'User-Agent': 'Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.87 Safari/537.36',
}
# url for COK to MAA on 13/10/21 
url = "https://flight.yatra.com/lowest-fare-service/dom2/get-fare?origin=COK&destination=MAA&from=13-10-2021&to=13-10-2021&tripType=O&airlines=all&src=srp"

r = requests.get(url, timeout=5, headers=HEADERS)
data = r.json()
flights = data['day']['2021-10-13']['af']
lowest_fare=0
#lowest_fare_flight=""
print("***** Flight List *****")
for flight in flights:
   flight_data = flights[flight]
   code = flight_data['ac']
   fare = flight_data['tf']
   name = flight_data['ow'][0]['an']
   print(name+" : "+str(fare))
   if(lowest_fare == 0 or fare < lowest_fare):
      lowest_fare = fare
      lowest_fare_flight = name 
   ##stop = flight_data['ow'].len
##print(flights)
print("*****************")      
print("Lowest fare is "+ str(lowest_fare)+" for " + lowest_fare_flight)    
