export function renderHome() {
  return `      <div class="container">
  <div><img src="./public/assets/logo.jpeg" /></div>
  <div class="card-connexion">
    <div><h1>BOMBERMAN ONLINE GAME</h1></div>
    <div id="error"></div>
    <div class="block_input username">
      <label>Username</label>
      <input type="text" required id ="username" name="username" placeholder="give your username"/>
    </div>
    <div class="block_input username">
      <label>Fullname</label>
      <input type="text" required id="name" name="name" placeholder="give your full name"/>
    </div>
    <div>
      <button type="submit" id="participate" class="submit">
        <span>Participate</span>
      </submit>
    </div>
</div>
</div>`;
}

export function renderWaitingPlayers() {
  return `


<div class="left-game">  

</div>
<div class="game" >
<div id="timing-block"><div class="timing-container">
<h1 id="time"></h1>

</div></div>

</div>
<div class="right-game">


</div>



`;
}

export function renderMap(players, player) {
  let render = ""
   players.forEach(player => {
    render += `<li> <span class="player-name">${player.username}</span>: <span id="${player.username}" class="liveValue">
    ${player.livesCount} </span>
    <svg
            xmlns="http://www.w3.org/2000/svg"
            width="20"
            height="32"
            viewBox="0 0 24 24"
            class="heart"
          >
            <path
              fill="white"
              d="m12 21.35l-1.45-1.32C5.4 15.36 2 12.27 2 8.5C2 5.41 4.42 3 7.5 3c1.74 0 3.41.81 4.5 2.08C13.09 3.81 14.76 3 16.5 3C19.58 3 22 5.41 22 8.5c0 3.77-3.4 6.86-8.55 11.53z"/></svg>
</li>`
  });
  return `  
    <div class="left-game">
    <div class="card">
      <div class="card-header">
        <div class="bgImage">
          <img
            src="./public/assets/player-1.png"
            alt="player-1"
            class="imgPlayer"
          />
        </div>
      </div>
      <div class="playerName">
        <h2 id="playerUsername">${player.username}</h2>
        <span class="position-player"> You are positionned on </span>
        <p class="lineHeader"></p>
      </div>
      <div class="life">
        <ul id="live-content">
          ${render}
        </ul>
      </div>
      <div class="dataGame">
        <div class="powerUp">
          <div class="asset">
            <img src="./public/assets/burst.jpeg" alt="flame" class="flame" />
          </div>
          <span class="numberBombs">3</span>
        </div>
        <div class="powerUp">
          <div class="asset">
            <img src="./public/assets/bombup.png" alt="flame" class="flame" />
          </div>
          <span class="numberBombs">1</span>
        </div>
        <div class="powerUp">
          <div class="asset">
            <img src="./public/assets/speed.png" alt="flame" class="flame" />
          </div>
          <span class="numberBombs">2</span>
        </div>
      </div>
    </div>
  </div>  
  
  <div class="game"></div>

  <div class="right-game">
    <div class="chat-container">
      <div class="chat-header">
        <h2>Game Chat</h2>
      </div>
      <div class="chat-messages">
        <!-- Chat messages will be dynamically added here -->
      </div>
      <div class="chat-input">
        <input type="text" placeholder="Type your message..." id="chatInput" />
        <button id="sendBtn">Send</button>
      </div>
    </div>
  </div>
  `;
}
