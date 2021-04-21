# Elevators
This is a simple application that represent a simulation of an office building. You can set the number of floors and elevators for it.
## The application has the following entities:
* storage - responsible for providing elevators;
* elevator - the most common elevator, but we can set speed for it;
* worker - building worker, has his own schedule.
# How to start
First of all we need to run the server:  
`go run main.go server [number of floors] [number of elevators]`  
Then we need to run the `client`:  
`go run main.go client [worker name] [worker schedule]`  
Worker can visit the different floors per day, so his `schedule` introduced in following format:  
`[floor number]:[and residence time]_[floor number]:[and residence time]...`
