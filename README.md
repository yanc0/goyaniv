# Goyaniv
Online, multiplayer, Yaniv card game

# Docker
Launch goyaniv docker container and bind 5000 port to host IP
```
docker run -p 5000:5000 yanc0/goyaniv:latest
```
# API
## Action
### Ready
```
{
    "name":"ready",
    "putcards":[],
    "takecard":0,
    "option":"yes"
}
```
### Change name
```
{
    "name":"name",
    "putcards":[],
    "takecard":0,
    "option":"Alice"
}
```
### Put
Takecard to 0 to take from middledeck
```
{
    "name":"put",
    "putcards":[1,2],
    "takecard":0,
    "option":""
}
```
### Yaniv
```
{
    "name":"yaniv",
    "putcards":[],
    "takecard":0,
    "option":""
}
```
### Asaf
```
{
    "name":"asaf",
    "putcards":[],
    "takecard":0,
    "option":"yes" (or no)
}
```

## State

```
{
    "playdeck":[
        {
            "id":17,
            "value":4,
            "symbol":"heart"
        }
    ],
    "lastlog":{
        "playername":"",
        "action":"",
        "takecard":{
            "id":0,
            "value":0,
            "symbol":""
        },
        "putcards":null,
        "option":""
    },
    "round":0,
    "started":false,
    "terminated":false,
    "players":[
        {
            "name":"WvLdjRIGihpqWlRqaIHH",
            "id":"PZvYDrDRuiArcdsETXNP",
            "me":false,
            "playing":true,
            "connected":true,
            "yaniver":false,
            "asafer":false,
            "ready":false,
            "lost":false,
            "score":0,
            "deckweight":0,
            "deck":[
                {
                    "id":0,
                    "value":0,
                    "symbol":""
                },
                {
                    "id":0,
                    "value":0,
                    "symbol":""
                },
                {
                    "id":0,
                    "value":0,
                    "symbol":""
                },
                {
                    "id":0,
                    "value":0,
                    "symbol":""
                },
                {
                    "id":0,
                    "value":0,
                    "symbol":""
                }
            ]
        },
        {
            "name":"eAFvYsjtKFDxmQtRiXku",
            "id":"LiZkibsHsquOZFbAXuxl",
            "me":true,
            "playing":false,
            "connected":true,
            "yaniver":false,
            "asafer":false,
            "ready":false,
            "lost":false,
            "score":0,
            "deckweight":28,
            "deck":[
                {
                    "id":28,
                    "value":2,
                    "symbol":"diam"
                },
                {
                    "id":16,
                    "value":3,
                    "symbol":"heart"
                },
                {
                    "id":29,
                    "value":3,
                    "symbol":"diam"
                },
                {
                    "id":39,
                    "value":13,
                    "symbol":"diam"
                },
                {
                    "id":26,
                    "value":13,
                    "symbol":"heart"
                }
            ]
        }
    ],
    "error":""
}
```
