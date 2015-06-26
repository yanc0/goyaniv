# Goyaniv
Online, multiplayer, Yaniv card game

# API
## Action
### Change name
```
{
    "name":"nickname",
    "putcards":[],
    "takecard":0,
    "option":"Alice"
}
```
### Play
Takecard to 0 to take from middledeck
```
{
    "name":"play",
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
    "option":""
}
```

## State

```
{
    "playdeck":[
        {
            "id":25,
            "value":12,
            "symbol":"diamond"
        },
        {
            "id":39,
            "value":12,
            "symbol":"club"
        }
    ],
    "lastaction":{
        "player":"JUmTxQbumZtVEYIbXasf",
        "takecard":23
    },
    "round":0,
    "started":true,
    "terminated":false,
    "players":
    [
      {
        "name":"Alice",
        "id":"Vkz0Gl0tUlPwFsQTNXGh",
        "me": true,
        "playing": false,
        "connected":true,
        "yaniver":false,
        "loss":false,
        "score":36,
        "deckweight":27,
        "deck":[
        {
            "id":33,
            "value":7,
            "symbol":"heart"
        },
        {
            "id":10,
            "value":10,
            "symbol":"spade"
        },
        {
            "id":49,
            "value":10,
            "symbol":"club"
        }
      },
      {
        "name":"Bob",
        "id":"OndJnsnkapmdHgsvzHsk",
        "me": false,
        "playing": true,
        "connected":true,
        "yaniver":false,
        "score":36,
        "deckweight":27,
        "deck":[
        {
            "id":0,
            "value":0,
            "symbol":"nil"
        },
        {
            "id":0,
            "value":0,
            "symbol":"nil"
        },
        {
            "id":0,
            "value":0,
            "symbol":"nil"
        }
      },
    ]
}
```
