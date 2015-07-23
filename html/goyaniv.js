(function($) {
    $.fn.onEnter = function(func) {
        this.bind('keypress', function(e) {
            if (e.keyCode == 13) func.apply(this, [e]);    
        });               
        return this; 
     };
})(jQuery);

function getenv() {
  var env = {
    ws: "ws://"
  }
  env.ws += window.location.host
  env.ws += window.location.pathname + "/ws"
  return env
};

function getcardhtml(id, value, symbol) {
  var div = "<div ";
  var val;
  var sym;
  var color = "red";
  if (symbol == "spade" || symbol == "club")
    color = "black";

  if (value == 11)
    val = "J";
  else if (value == 12)
    val = "Q";
  else if (value == 13)
    val = "K";
  else
    val = value;

  sym = "&"+symbol+"s;";
  if (symbol == "") {
    sym = "&#8718;";
    val = "";
    color = "purple"
  }
  div += "id='"+id+"' ";
  div += "class='"+color+" card'>";
  div += val;
  div += sym + " </div>";
  return div
};

function getdeckhtml(deck, mine) {
  var div = "<div class='deck ";
  if (mine == true) 
    div += "mine";
  div += "'>";
  for (var card of deck) { 
    div += getcardhtml(card.id, card.value, card.symbol);
  };
  div += "</div>";
  return div
}

function getplayershtml(players) {
  var div = "";
  for (var player of players) { 
    if (player.connected)
      var conn = "<span class='green'>&#9210;</span>";
    else
      var conn = "<span class='red'>&#9210;</span>";
    div += "<div class='player'> ";
    if (player.ready)
      div += "<div class='ready'>&#10003;</div>";
    else
      div += "<div class='notready'>&#10162;</div>";
    if (player.playing)
      div += "&#9660;";
    div += "<br/>";
    div += "name: "+player.name.substring(0,8)+" "+conn+"<br/>";
    div += "score: "+player.score+"<br/>";
    div += "deck: "+player.deckweight+"<br/>";
      div += getdeckhtml(player.deck, player.me);
    div += "</div>";
  };
  return div 
}

function resetview() {
  $('.head').empty();
  $('.middledeck').empty();
  $('.playdeck').empty();
  $('.players').empty();
  $('.log').empty();
  $('.debug').empty();
}

function cardhighlighter() {
  $('.mine > .card').each(function () {
    if (cardselected.indexOf(parseInt(this.id)) > -1) {
      $(this).addClass('highlight');
    } else {
      $(this).removeClass('highlight');
    };
  });
};

function cardselector() {
  $('.mine > .card').click(function (event) {
    var card = $(event.target)
    var id = parseInt(this.id);
    var index = cardselected.indexOf(id);
    if (index > -1) {
        cardselected.splice(index, 1);
    } else {
      cardselected.push(id);
    }
    cardhighlighter();
  });

  $('.playdeck > .deck > .card').click(function (event){
    action("put",cardselected, parseInt(this.id));
  });
  
  $('.middledeck > .card').click(function (event){
    action("put",cardselected, parseInt(this.id));
  });
}

function readyselector() {
  $('.ready').click(function(event) {
    action("ready", [], 0, "no");
  });
  $('.notready').click(function(event) {
    action("ready", [], 0, "yes");
  });
}

function me(state) {
  for (var player of state.players) {
    if (player.me)
      return player;
  }
}

function actionbar(me) {
  $('.actionbar').empty();
  if(me.deckweight <= 5) {
    $(".actionbar").append("<div class='yaniv'>Yaniv !</div>");
    $(".actionbar").append("<div class='asaf'>  Asaf !</div>");
    $(".actionbar").append("<div class='noasaf'> No Asaf !</div>");
  }
}

function action(name, putcards, takecard, option) {
  var action = {"name":"","putcards":[],"takecard":0,"option":""}
  action.name = name;
  action.putcards = putcards;
  action.takecard = takecard;
  action.option = option;
  ws.send(JSON.stringify(action)); 
  cardselected = Array();
}

function getloghtml(log) {
  var div = "";
  if (log.action == "") {
    return div
  }
  div += "Player ";
  div += log.playername.substring(0,8);
  if (log.action == "put") {
    div += " takes ";
    div += getcardhtml( log.takecard.id, 
                        log.takecard.value, 
                        log.takecard.symbol);
    div += " and puts ";
    if (log.putcards != null) {
      for (var card of log.putcards) { 
        div += getcardhtml( card.id, 
                            card.value, 
                            card.symbol);
      } 
    };
  }

  if (log.action == "yaniv") {
    div += " Yaniv !";
  }

  if (log.action == "asaf" && log.option == "yes") {
    div += " Asaf !";
  }
  
  if (log.action == "asaf" && log.option == "no") {
    div += " doesn't want Asaf !";
  }

  if (log.action == "name") {
    div = log.playername.substring(0,8) + " becomes " + log.option;
  }
  
  return div;
}


$(document).on("click", ".yaniv", function() {
  action('yaniv');
});

$(document).on("click", ".asaf", function() {
  action('asaf', [], 0, "yes");
});

$(document).on("click", ".noasaf", function() {
  action('asaf', [], 0, "no");
});

$( function () {
    console.log($(".input-name"));
    $(".input-name").onEnter( function() {
        action("name", [], 0, $(this).val());                        
    });
});

var env = getenv();
var ws = new WebSocket(env.ws);
var cardselected = Array();

ws.onmessage = function (msg) {
  resetview();
  var state = jQuery.parseJSON(msg.data);
  $(".log").append(getloghtml(state.lastlog));
  $(".middledeck").append(getcardhtml("0","",""));
  $(".playdeck").append(getdeckhtml(state.playdeck));
  $(".players").append(getplayershtml(state.players));
  $(".head").append(state.error + '<br/>');
  $(".debug").append(msg.data);
  actionbar(me(state));
  cardselector();
  cardhighlighter();
  readyselector();
};
