import { DOM } from "./my_framework/dom.js";
import { Router } from "./my_framework/router.js";
import {
  renderHome,
  renderMap,
  renderWaitingPlayers,
} from "./renders/components.templates.js";
let then = Date.now();
let fpsInterval = 400;
const router = new Router();
let player = null;
let players = null
let state = null;
let ws = null;
let mapGame = null;
let direction = null;
let bombPosition = null;
let lost = false;
let gameOver = false;
let canPain = false;
let mapNodes = []
router.addRoute("/", () => {
  let render = renderHome();
  DOM.setHTML("#app", render);
  DOM.addEventListener("#participate", "click", validateUsername);
});
router.addRoute("/running", () => {
  let render = renderMap(players, player);
  DOM.setHTML("#app", render);
  updateMapView()
  mapNodes = getMapNodes(mapGame)
  console.log("---- player ---- ", player);
  let posMessage = null;
  if (player.number === "player1") {
    posMessage = "You are positionned on left top";
  } else if (player.number === "player2") {
    posMessage = "You are positionned on right top";
  } else if (player.number === "player3") {
    posMessage = "You are positionned on left bottom";
  } else if (player.number === "player4") {
    posMessage = "You are positionned on right bottom";
  }
  if (posMessage !== null) {
    DOM.setHTML(".position-player", posMessage);
  }

  requestAnimationFrame(refreshGameMap);




  const chatInput = DOM.getById("chatInput");
  const sendBtn = DOM.getById("sendBtn");
  sendBtn.addEventListener("click", () => {
    const message = chatInput.value.trim();
    if (message !== "") {
      const data = {
        type: "message",
        username: player.username,
        message: message,
      };
      ws.send(JSON.stringify(data));

      chatInput.value = "";
    }
  });

  chatInput.addEventListener("keydown", (event) => {
    if (event.key === "Enter") {
      const message = chatInput.value.trim();
      if (message !== "") {
        const data = {
          type: "message",
          username: player.username,
          message: message,
        };
        ws.send(JSON.stringify(data));

        chatInput.value = "";
      }
    }
  });

});
router.addRoute("/waiting", () => {
  if (player === null) {
    return;
  }
  let render = renderWaitingPlayers();
  DOM.setHTML("#app", render);
  ws = new WebSocket(`ws://localhost:4044/ws?username=${player.username}`);
  ws.onopen = () => {
    console.log("Connecté au serveur WebSocket");
  };
  ws.onmessage = (event) => {
    const message = JSON.parse(event.data);
    console.log("message server:", message);
    if (message.type === "waiting_players") {
      DOM.setHTML("#time", message.data);
    } else if (message.type === "starting_game") {
      DOM.setHTML("#time", message.data);
    } else if (message.type === "note_complete_game") {
      let render = renderHome();
      DOM.setHTML("#app", render);
      DOM.addEventListener("#participate", "click", validateUsername);
      DOM.setHTML("#error", message.data);
      DOM.setAttribute("#error", "class", "error");
    } else if (message.type === "running_game") {
      mapGame = message.data.map;
      players = message.data.listPlayer
      canPain = true;
      message.data.forEach
      if (player.username === message.data.player.username) {
        player = message.data.player;
       // DOM.setHTML("#nbrlive", message.data.player.livesCount);

      }
      router.navigate("/running");
    } else if (message.type === "move") {
      mapGame = message.data.map;
      canPain = true;
      //console.log("from server:", message.data.player);
   //   console.log("from client:", player);
      clearPlayer(message.data.player)
     placePlayer(message.data.player)
      if (player.username === message.data.player.username) {
        player = message.data.player;
        fpsInterval = player.speed;
       // DOM.setHTML("#nbrlive", message.data.player.livesCount);
      }
      //updateMapView();
    } else if (message.type === "pose-bomb") {
      mapGame = message.data.map;
      poseBomb(message.data.player)
      canPain = true
      if (player.username === message.data.player.username) {
        player = message.data.player;
        setTimeout(function () {
          let dataAction = {
            type: "explose-bomb",
            sense: "",
            username: player.username,
          };
          ws.send(JSON.stringify(dataAction));
        }, 1500);
        //DOM.setHTML("#nbrlive", message.data.player.livesCount);

      }
      //updateMapView();
    } else if (message.type === "explose-bomb") {
      mapGame = message.data.map;
      exploseBomb(mapGame)
      canPain = true
      if (player.username === message.data.player.username) {
        player = message.data.player;
        //DOM.setHTML("#nbrlive", message.data.player.livesCount);

      }
      // updateMapView();
    } else if (message.type === "after-explosion") {
      mapGame = message.data.map;
      clearExplosion(mapGame)
      players = message.data.listPlayer
      players.forEach((player)=>{
        DOM.setHTML(`#${player.username}`, player.livesCount)
      })
      canPain = true
      if (
        message.data.listPlayer.filter((player) => player.livesCount == 0)
          .length
      ) {
        // gameOver = true
      }
      if (player.username === message.data.player.username) {
        bombPosition = null;
        if (player.livesCount == 0) {
          lost = true;
        }
        player = message.data.player;
        //DOM.setHTML("#nbrlive", message.data.player.livesCount);

      }
      //updateMapView();
    }

    if (message.type === "message") {
      const chatMessages = DOM.querySelector(".chat-messages");
      const newMessage =  DOM.createElement("div");
      newMessage.classList.add("chat-message");

      if (message.data.username === player.username) {
        DOM.addOneClass(newMessage,"sent-message")
      } else {
        DOM.addOneClass(newMessage,"received-message")
      }

      const usernameElement = DOM.createElement("span");
      DOM.addOneClass(usernameElement,"username");
      usernameElement.innerText = message.data.username + ": ";

      const messageContent = DOM.createElement("span");
      DOM.addOneClass(messageContent,"message-content");
      messageContent.classList.add("message-content");
      messageContent.innerText = message.data.message;

      newMessage.appendChild(usernameElement);
      newMessage.appendChild(messageContent);

      chatMessages.appendChild(newMessage);
      chatMessages.scrollTop = chatMessages.scrollHeight;
    }

  };
  ws.onerror = (error) => {
    console.error("Erreur WebSocket:", error);
  };
  ws.onclose = () => {
    console.log("Connexion WebSocket fermée");
    ws = null
  };
});
// Load the home page by default
router.navigate("/");
function handleLink(event) {
  var link = event.target.closest("a");
  if (link) {
    event.preventDefault();
    const newRoute = link.getAttribute("href");
    history.pushState(null, null, newRoute);
    router.navigate(newRoute);
  }
}
async function validateUsername() {
  let username = DOM.getValue("#username");
  let name = DOM.getValue("#name");
  console.log("username:", username);
  let success = "";
  try {
    const response = await fetch(`/validate?username=${username}&name=${name}`);
    success = await response.text();
  } catch (error) {
    console.error("Erreur lors de la validation du nom d'utilisateur:", error);
    success = "Erreur fetch";
  }
  //affichage erreur
  console.log("sucess:", success);
  if (success.trim() !== "ok") {
    DOM.setHTML("#error", success);
    DOM.setAttribute("#error", "class", "error");
    ws = null;
    return;
  }
  player = { username: username, name: name };
  router.navigate("/waiting");
}

//gestion des deplacements
document.addEventListener("keyup", function (event) {
  if (
    event.key === "a" &&
    event.target != DOM.getById("username") &&
    bombPosition === null &&
    !lost &&
    !gameOver
  ) {
    bombPosition = player.position;
    console.log("pose bomb");
    let dataAction = {
      type: "pose-bomb",
      sense: "",
      username: player.username,
    };
    ws.send(JSON.stringify(dataAction));
    return;
  }
  switch (event.key) {
    case "ArrowUp":
      direction = "up";
      break;
    case "ArrowDown":
      direction = "down";
      break;
    case "ArrowLeft":
      direction = "left";
      break;
    case "ArrowRight":
      direction = "right";
      break;
  }
  if (direction != null && !lost && !gameOver) {
    let dataAction = {
      type: "move",
      sense: direction,
      username: player.username,
    };
    console.log("ws:", ws)
    ws.send(JSON.stringify(dataAction));
    direction = null;
    return;
  }
});

function updateMapView() {
  DOM.setHTML(".game", "");
  for (let i = 0; i < mapGame.length; i++) {
    const element = mapGame[i];
    let line = DOM.createOneElement("div", "class", `line row-${i} line${i}`);
    DOM.append(".game", line);
    for (let j = 0; j < element.length; j++) {
      const tuile = element[j];
      switch (tuile) {
        case "#":
          DOM.append(
            `.line${i}`,
            DOM.createOneElement("div", "class", `line${i} col col-${i} wall`, tuile)
          );
          break;
        case "bomb":
          DOM.append(
            `.line${i}`,
            DOM.createOneElement("div", "class", `line${i} col col-${i} bomb`, tuile)
          );
          break;
        case "explosion":
          DOM.append(
            `.line${i}`,
            DOM.createOneElement("div", "class", `line${i} col col-${i} explosion`, tuile)
          );
          break;
        case "eb":
          DOM.append(
            `.line${i}`,
            DOM.createOneElement("div", "class", `line${i} col col-${i} eb`, tuile)
          );
          break;
        case "ef":
          DOM.append(
            `.line${i}`,
            DOM.createOneElement("div", "class", `line${i} col col-${i} ef`, tuile)
          );
          break;
        case "es":
          DOM.append(
            `.line${i}`,
            DOM.createOneElement("div", "class", `line${i} col col-${i} es`, tuile)
          );
          break;
        case "gb":
          DOM.append(
            `.line${i}`,
            DOM.createOneElement("div", "class", `line${i} col col-${i} gb`, tuile)
          );
          break;
        case "gf":
          DOM.append(
            `.line${i}`,
            DOM.createOneElement("div", "class", `line${i} col col-${i} gf`, tuile)
          );
          break;
        case "gs":
          DOM.append(
            `.line${i}`,
            DOM.createOneElement("div", "class", `line${i} col col-${i} gs`, tuile)
          );
          break;

        case " ":
          DOM.append(
            `.line${i}`,
            DOM.createOneElement("div", "class", `line${i} col col-${i} empty`, tuile)
          );
          break;
        case "player1":
          DOM.append(
            `.line${i}`,
            DOM.createOneElement("div", "class", `line${i} col col-${i} player1`, tuile)
          );
          break;
        case "player2":
          DOM.append(
            `.line${i}`,
            DOM.createOneElement("div", "class", `line${i} col col-${i} player2`, tuile)
          );
          break;
        case "player3":
          DOM.append(
            `.line${i}`,
            DOM.createOneElement("div", "class", `line${i} col col-${i} player3`, tuile)
          );
          break;
        case "player4":
          DOM.append(
            `.line${i}`,
            DOM.createOneElement("div", "class", `line${i} col col-${i} player4`, tuile)
          );
          break;
        default:
          DOM.append(
            `.line${i}`,
            DOM.createOneElement("div", "class", `line${i} col col-${i} brick`, tuile)
          );
          break;
      }
    }
  }
}
function refreshGameMap() {
  requestAnimationFrame(refreshGameMap);
  let now = Date.now();
  let elapsed = now - then;
  if (elapsed > fpsInterval) {
    then = Date.now();
    if (canPain){
     // updateMapView();
    }
    canPain = false
  }
}

function getMapNodes(map){
  let mapNodes = []
for(let i = 0; i < map.length; i++){
  let classCol = `.col-${i}`
  let lines = document.querySelectorAll(classCol)
  mapNodes.push(lines)
}
return mapNodes
}

function clearPlayer(player){
 //console.log("player:", player)
  let x = player.position.x
  let y = player.position.y
  let before_x = player.beforeMove.x
  let before_y = player.beforeMove.y
  if (x === before_x && y === before_y) return
  //console.log("position avant:", before_y, before_x)
  //console.log("clear player:", mapNodes[before_y][before_x])
  //console.log("class list:", mapNodes[before_y][before_x].classList)
  //console.log("player id:", `${player.number}`)
  mapNodes[before_y][before_x].classList.remove(`${player.number}`);
  mapNodes[before_y][before_x].classList.add(`empty`);
}
function placePlayer(player){
  let x = player.position.x
  let y = player.position.y
  let before_x = player.beforeMove.x
  let before_y = player.beforeMove.y
  if (x === before_x && y === before_y) return
  //console.log("position apres:", y, x)

  mapNodes[y][x].className=player.number
  // mapNodes[y][x].classList.add(`${player.number}`);
  // mapNodes[y][x].classList.remove(`empty`);
}

function poseBomb(player){
  console.log("bombs:", player.bombPosition)
  player.bombPosition.forEach((pos)=>{
  mapNodes[pos.y][pos.x].className= "bomb"
})
}

function exploseBomb(map){
  for(let i = 0; i < map.length; i++){
    for(let j = 0; j < map[0].length; j++){
      if (map[i][j] === "explosion"){
        console.log("explosion:",  i, j)
        mapNodes[i][j].className= "explosion"
      }else if(map[i][j] === "ef"){
        mapNodes[i][j].className= "ef"
        console.log("explosion:ef", i, j)
      }else if(map[i][j] === "eb"){
        mapNodes[i][j].className= "eb"
        console.log("explosion:eb", i, j)
      }
      else if(map[i][j] === "es"){
        mapNodes[i][j].className= "es"
        console.log("explosion:es", i, j)
      }
    }
  }
}

function clearExplosion(map){
  for(let i = 0; i < map.length; i++){
    for(let j = 0; j < map[0].length; j++){
      if (map[i][j] === " "){
        mapNodes[i][j].className= "empty"
      }else if(map[i][j] === "gf"){
        mapNodes[i][j].className= "gf"
      }else if(map[i][j] === "gb"){
        mapNodes[i][j].className= "gb"
      }
      else if(map[i][j] === "player1"){
        mapNodes[i][j].className= "player1"
      }else if(map[i][j] === "player2"){
        mapNodes[i][j].className= "player2"
      }else if(map[i][j] === "player3"){
        mapNodes[i][j].className= "player3"
      }else if(map[i][j] === "player4"){
        mapNodes[i][j].className= "player4"
      }else if(map[i][j] === "gs"){
        mapNodes[i][j].className= "gs"
      }
    }
  }
}