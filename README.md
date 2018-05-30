## Overview
A plugin for [go-bot](https://github.com/go-chat-bot/bot) that allows defining custom bot responses for given matches

## Available commands & Examples
```
!responses match set "why did the chicken cross the road?" "to get to the other side"
!responses match set "Error processing request of user fernando.almeida" "Hey @fernando, take a look"
!responses match unset
```

### To-do:
- [x] Create project basics
- [x] Define methods structure
- [x] Create some tests
- [x] Connect with redis
- [x] Create command to set patterns/responses
- [x] Apply regex to find responses from patterns
- [x] Create and configure heroku redis app
- [x] Deploy a bot instance and test with slack -> [repo](https://github.com/fernand-o/got-bot-heroku)
- [x] Create command to list defined responses -> (!responses list)
- [x] Create command to delete defined responses -> (!responses unset)
- [x] Create command to delete all responses -> (!responses list)
- [x] List responses formatted
- [x] Define lists with random responses to send combined with the defined response
- [x] Send a random item list with the message
- [x] Fix delete list command
- [ ] Update usage examples
- [ ] Create an easier way of deleting responses
