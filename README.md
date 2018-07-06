## Overview
A plugin for [go-bot](https://github.com/go-chat-bot/bot) that allows defining custom bot responses for given matches

## Usage
```
Defining matches and responses:
!responses match set "Is someone there?" "Hello"
!responses match set "Rick Sanchez" "@rick" #fun
!responses match showall
!responses match unset 1

Defining lists to be used within the responses:
!responses list add #fun "Hey %s, looks like you're guilty"
!responses list add #fun "Is that you %s?"
!responses list show #fun
!responses list showall
!responses list remove #fun "Is that you %s?"
!responses list delete #fun
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
- [x] Update usage examples
- [ ] Create an easier way of deleting matches and list items
- [ ] Create a !responses showall command
