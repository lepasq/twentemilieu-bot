# WhatsApp bot for Twente-Milieu

WhatsApp Handler written in Go.   
Twente Milieu is a waste disposal company in the Twente region.  
On their website, the company provides a waster calendar, accessible via an API.  
This bot fetches the calendar and sends reminders to bring out trash.   

## Setup
```sh
git clone https://github.com/lepasq/twentemilieu-bot.git
go get -d ./...
```

Next, open `config.yml` and adapt the values accordingly.


## Packages
* `api`: fetches the address_id and calendar from the twente milieu 
* `config`: parses and creates config struct from `config.yml`
* `handlers`: handle received messages
* `scheduler`: schedules a calendar check once per day


## Contribute
Feel free to contribute to this project. There are multiple things to improve, such as:  
* Saving the Calendar in a file, instead of requesting it daily from the API
* Saving the AddressId in `config.yml` instead of requesting it daily from the API
* Extending the handler functionality