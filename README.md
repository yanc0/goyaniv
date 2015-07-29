# Goyaniv
Online, multiplayer, Yaniv card game

# Docker
Launch goyaniv docker container and bind 5000 port to host IP
```
docker run -p 5000:5000 yanc0/goyaniv:latest
```
# API
State :
-------

```
{
  "action_last": {
    "name": "card_take",
    "options": {
      "discarded": [{
        "id": 17,
        "value": 3
      }, {
        "id": 18,
        "value": 4
      }, {
        "id": 19,
        "value": 5
      }],
      "taken": {
        "id": 0
      }
    },
    "player": {
      "id": 345,
      "name": "Totoro"
    }
  },
  "errors": [{
    "message": "Shit! You got an error! ='("
  }],
  "stack": [{
    "id": 15,
    "value": 2
  }],
  "user": {
    "id": 345,
    "name": "Totoro",
    "score": 0,
    "state": {
      "asaf": false,
      "loser": false,
      "online": true,
      "player": true,
      "playing": false,
      "ready": true,
      "yaniv": false
    }
  },
  "game": {
    "ended": false,
    "round": 1,
    "started": true
  },
  "opponents": [{
    "hand_size": 5,
    "id": 134,
    "name": "Totoro 24",
    "score": 0,
    "state": {
      "asaf": false,
      "online": true,
      "playing": false,
      "ready": true,
      "yaniv": false
    }
  }],
  "spectators": [{
    "id": 10345,
    "name": "Totoro 56",
    "state": {
      "loser": false,
      "online": true
    }
  }]
}
```

Actions :
---------
- Set ready

```
{
  "name": "ready_set",
  "options": {
    "is_ready": true
  }
}
```

- Set user name

```
{
  "name": "user_name_set",
  "options": {
    "name": "New name"
  }
}
```

- Take card

```
{
  "name": "card_take",
  "options": {
    "take": 3,
    "discard": [17, 18, 19]
  }
}

```
- Yaniv

```
{
  "name": "yaniv",
  "options": null
}

```
- Asaf

```
{
  "name": "asaf",
  "options": {
    "try_asaf": true
  }
}
```
