A plugin for [go-bot](https://github.com/go-chat-bot/bot) that allows defining custom responses for given matches

Usage: 
```
!responses set match response
!responses unset match
!responses list
```

Examples:
```
!responses set "why did the chicken cross the road?" "to get to the other side"
!responses set "Error processing request of user fernando.almeida" "Hey @fernando, take a look"
```

To-do:
- [x] Create project basics
- [x] Define methods structure
- [x] Create some tests
- [x] Connect with redis
- [x] Create command to set patterns/responses
- [x] Apply regex to find responses from patterns
- [x] Create and configure heroku redis app
- [x] Deploy a bot instance and test with slack -> [repo](https://github.com/fernand-o/got-bot-heroku)
- [x] Create command to list defined responses
- [x] Create command to delete defined responses
- [ ] Create command to delete all responses
- [ ] Allow defining prefixes for conditions (to avoid processing all received messages)
