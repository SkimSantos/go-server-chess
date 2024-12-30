const loginBtn = document.getElementById('login-btn');
const signupBtn = document.getElementById('signup-btn');
const createGameBtn = document.getElementById('create-game-btn');
const joinGameBtn = document.getElementById('join-game-btn');
const resignBtn = document.getElementById('resign-btn');
const userInfoDiv = document.getElementById('user-info');
const chessBoard = document.getElementById('chess-board');

const serverHost = "http://localhost:8080/"

// Placeholder for user and game state
let user = null;
let currentGame = null;
let vsPlayer = "";
let whitePlayer = true;

let selectedSquare = null; // Track the currently selected square
let myTurn = false;

function addMoveListeners() {
  const squares = chessBoard.querySelectorAll('div');
  squares.forEach((square, index) => {
    square.addEventListener('click', () => {
      handleSquareClick(square, index);
    });
  });
}

function handleSquareClick(square, index) {
  let row = Math.floor(index / 8);
  let col = index % 8;

  if(!whitePlayer) {
    row = 7 - row;
    col = 7 - col;
  }

  if (!selectedSquare && square.textContent.trim()) {
    // Select the square with the piece
    selectedSquare = { square, row, col };
    square.classList.add('selected');
  } else if (selectedSquare) {
    // Attempt to make a move
    let from = `${String.fromCharCode(97 + selectedSquare.col)}${8 - selectedSquare.row}`;
    let to = `${String.fromCharCode(97 + col)}${8 - row}`;
    makeMove(from, to);
    selectedSquare.square.classList.remove('selected');
    selectedSquare = null;
  }
}

async function makeMove(from, to) {
  if (!currentGame) {
    alert('No active game!');
    return;
  }

  try {
    const response = await fetch(serverHost + 'game/move', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${user.token}`,
      },
      body: JSON.stringify({ game_id: currentGame.ID, move: `${from}${to}` }),
    });

    if (response.ok) {
      const game = await response.json();
      currentGame = game;
      myTurn = game.CurrentTurn == user.id
      renderChessBoard(game.State, game.Player1ID === user.id); // Update the board with the new state
      checkGameState()
    } else {
      const errorMessage = await response.text();
      alert(`Invalid move: ${errorMessage}`);
    }
  } catch (error) {
    console.error(error);
    alert('Error making move.');
  }
}


function renderChessBoard(fen = "8/8/8/8/8/8/8/8", isPlayer1 = true, otherPlayer = "") {
  chessBoard.innerHTML = ''; // Clear the board

  const [position] = fen.split(' '); // Extract FEN piece placement
  const pieceMap = {
    r: '♜', n: '♞', b: '♝', q: '♛', k: '♚', p: '♟', // Black pieces
    R: '♖', N: '♘', B: '♗', Q: '♕', K: '♔', P: '♙', // White pieces
  };

  if (user != null) {
    // Set usernames
    const playerTop = document.getElementById('player-top');
    const playerBottom = document.getElementById('player-bottom');
  
    playerTop.textContent = vsPlayer; // Opponent on top
    playerBottom.textContent = user.username; // Current player on bottom
  }

  // Split the position into rows
  let rows = position.split('/');
  if (!isPlayer1) {
    // Reverse rows for Player 2
    rows = rows.reverse();
  }

  // Render the board row by row
  rows.forEach((row, rowIndex) => {
    let colIndex = 0;
    let coladd = 1;
    if (!isPlayer1) {
      const sepRow = row.split("");
      const reverseRow = sepRow.reverse();
      row = reverseRow.join("")
    }


    // Iterate over each character in the row
    for (const char of row) {
      if (isNaN(char)) {
        // If it's a piece, add it to the board
        const square = document.createElement('div');
        square.className = (rowIndex + colIndex) % 2 === 0 ? 'light' : 'dark';
        square.textContent = pieceMap[char] || '';
        chessBoard.appendChild(square);
        colIndex += coladd;
      } else {
        // If it's a number, add that many empty squares
        for (let i = 0; i < parseInt(char); i++) {
          const square = document.createElement('div');
          square.className = (rowIndex + colIndex) % 2 === 0 ? 'light' : 'dark';
          chessBoard.appendChild(square);
          colIndex += coladd;
        }
      }
    }
  });

  addMoveListeners();
}

async function checkGameState() {
  if(currentGame != null && !myTurn) {
    try {
      const gameID = currentGame.ID
      const response = await fetch(serverHost + 'game/get/' + gameID, {
        method: 'GET',
        headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${user.token}` },
      });
  
      if (response.ok) {
        const game = await response.json();
        if(game.Status == "ongoing"){
          currentGame = game;
          if(vsPlayer == "") {
            getVsPlayer();
          }
          createGameBtn.style.display = 'none';
          joinGameBtn.style.display = 'none';
          resignBtn.style.display = 'inline-block';
          if(game.CurrentTurn == user.id) {
            myTurn = true;
          }
          renderChessBoard(game.State, game.Player1ID === user.id, "something"); // Reset or initialize the chessboard
        } else if(game.Status != "pending") {
          currentGame = game;
          createGameBtn.style.display = 'inline-block';
          joinGameBtn.style.display = 'inline-block';
          resignBtn.style.display = 'none';
          myTurn = true;
          if(game.WinnerID == user.id) {
            alert("You Won!")
          } else if(game.WinnerID != 0) {
            alert("You Lost!")
          } else {
            alert("Draw")
          }
          renderChessBoard(game.State, game.Player1ID === user.id, "something"); // Reset or initialize the chessboard
        }
      } else {

      }
    } catch (error) {
      console.error(error);
    }
    setTimeout(checkGameState, 5000); // check again in 5 seconds
  }
}

async function getVsPlayer() {
  try {
    const vsUserID = currentGame.Player1ID === user.id ? currentGame.Player2ID : currentGame.Player1ID 
    const response = await fetch(serverHost + 'user/get/' + vsUserID, {
      method: 'GET',
      headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${user.token}` },
    });
    if(response.ok) {
      const vsUser = await response.json();
      vsPlayer = vsUser.username;
      playerTop.textContent = vsPlayer;
    } else {
      
    }
  } catch (error) {
    alert("Error getting other player")
  }
}

// Login functionality
loginBtn.addEventListener('click', async () => {
  const username = prompt('Enter your username:');
  const password = prompt('Enter your password:');
  try {
    const response = await fetch(serverHost + 'login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ username, password }),
    });
    if (response.ok) {
      const data = await response.json();
      user = { username, token: data.token, id: data.id };
      userInfoDiv.innerHTML = `<span>${username}</span> <button id="logout-btn">Logout</button>`;
      createGameBtn.style.display = 'inline-block';
      joinGameBtn.style.display = 'inline-block';
      resignBtn.style.display = 'none';
      document.getElementById('logout-btn').addEventListener('click', logout);
    } else {
      alert('Login failed!');
    }
  } catch (error) {
    console.error(error);
    alert('Error logging in');
  }
});

// Login functionality
signupBtn.addEventListener('click', async () => {
  const username = prompt('Enter new username:');
  const password = prompt('Enter new password:');
  const cpassword = prompt('Enter password again:');
  if(password != cpassword) {
    alert("Passwords don't match")
  } else {
    try {
      const response = await fetch(serverHost + 'register', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username, password }),
      });
      if (response.ok) {
        const data = await response.json();
        alert("User Created, you can now login.")
      } else {
        alert('Sing in failed!');
      }
    } catch (error) {
      console.error(error);
      alert('Error sing in');
    }
  }
});

// Logout functionality
function logout() {
  user = null;
  userInfoDiv.innerHTML = '<button id="login-btn">Login</button><button id="signup-btn">Sign Up</button>';
  document.getElementById('login-btn').addEventListener('click', () => loginBtn.click());
  document.getElementById('signup-btn').addEventListener('click', () => signupBtn.click());
  currentGame = null
  vsPlayer = "";
  createGameBtn.style.display = 'none';
  joinGameBtn.style.display = 'none';
  resignBtn.style.display = 'none';
  renderChessBoard()
}

// Create game functionality
createGameBtn.addEventListener('click', async () => {
  if (!user) {
    alert('Please log in to create a game!');
    return;
  }

  try {
    const response = await fetch(serverHost + 'game/create', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${user.token}` },
      body: JSON.stringify({}),
    });

    if (response.ok) {
      const game = await response.json();
      currentGame = game;
      createGameBtn.style.display = 'none';
      joinGameBtn.style.display = 'none';
      resignBtn.style.display = 'inline-block';
      myTurn = false
      whitePlayer = true;
      renderChessBoard(game.State, whitePlayer)
      alert('Game created, waiting for player!');
      checkGameState();
    } else {
      alert('Failed to create game.');
    }
  } catch (error) {
    console.error(error);
    alert('Error creating game');
  }
});

joinGameBtn.addEventListener('click', async () => {
  if (!user) {
    alert('Please log in to join a game!');
    return;
  }

  try {
    // Call the join endpoint directly
    const response = await fetch(serverHost + 'game/join', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${user.token}` },
      body: JSON.stringify({}),
    });

    if (response.ok) {
      const game = await response.json();
      currentGame = game;
      alert(`Joined game with ID: ${game.ID}`);
      createGameBtn.style.display = 'none';
      joinGameBtn.style.display = 'none';
      resignBtn.style.display = 'inline-block';
      myTurn = false
      whitePlayer = false;
      renderChessBoard(game.State, whitePlayer, "something"); // Reset or initialize the chessboard
      checkGameState();
    } else {
      alert('Failed to join a game. No games might be available.');
    }
  } catch (error) {
    console.error(error);
    alert('Error joining game.');
  }
});

resignBtn.addEventListener('click', async () => {
  if (!user) {
    alert('Error Resign without user!');
    return;
  }

  try {
    const response = await fetch(serverHost + 'game/resign', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${user.token}` },
      body: JSON.stringify({game_id: currentGame.ID}),
    });

    if (response.ok) {
      const game = await response.json();
      currentGame = null;
      createGameBtn.style.display = 'inline-block';
      joinGameBtn.style.display = 'inline-block';
      resignBtn.style.display = 'none';
      renderChessBoard()
      alert('Resigned game!');
    } else {
      alert('Failed to resign game.');
    }
  } catch (error) {
    console.error(error);
    alert('Error resigning game');
  }
});

// Initial setup
renderChessBoard();
