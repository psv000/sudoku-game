document.addEventListener('DOMContentLoaded', () => {
    // ----- Get DOM elements -----
    const elements = {
        grid: document.getElementById('sudoku-grid'),
        newGameBtn: document.getElementById('new-game-btn'),
        checkBtn: document.getElementById('check-btn'),
        difficultySelect: document.getElementById('difficulty'),
        message: document.getElementById('message'),
        playerNameInput: document.getElementById('player-name'),
        timer: document.getElementById('timer'),
        leaderboard: document.getElementById('leaderboard'),
        pauseBtn: document.getElementById('pause-btn'),
        pauseOverlay: document.getElementById('pause-overlay'),
    };

    // ----- Global variables  -----
    const state = {
        solution: [],
        timerInterval: null,
        secondsElapsed: 0,
        activeCell: null,
        isPaused: false,
    };
    
    const PLAYER_NAME_KEY = 'sudokuPlayerName';

    // ----- Saving/Loaging logic -----
    const handlePlayerName = {
        load: () => {
            const savedName = localStorage.getItem(PLAYER_NAME_KEY);
            if (savedName) elements.playerNameInput.value = savedName;
        },
        save: () => {
            localStorage.setItem(PLAYER_NAME_KEY, elements.playerNameInput.value);
        }
    };

    // ----- Timer logic -----
    const timer = {
        formatTime: (seconds) => {
            const mins = Math.floor(seconds / 60).toString().padStart(2, '0');
            const secs = (seconds % 60).toString().padStart(2, '0');
            return `${mins}:${secs}`;
        },
        start: () => {
            if (state.timerInterval) clearInterval(state.timerInterval);
            state.timerInterval = setInterval(() => {
                if (!state.isPaused) {
                    state.secondsElapsed++;
                    elements.timer.textContent = timer.formatTime(state.secondsElapsed);
                }
            }, 1000);
        },
        stop: () => {
            clearInterval(state.timerInterval);
        },
        reset: () => {
            state.secondsElapsed = 0;
            elements.timer.textContent = timer.formatTime(0);
        }
    };

    // ----- Resume/Pause logic -----
    function togglePause() {
        state.isPaused = !state.isPaused;
        elements.pauseOverlay.classList.toggle('hidden', !state.isPaused);
        elements.grid.classList.toggle('paused', state.isPaused);
        elements.pauseBtn.textContent = state.isPaused ? '▶️ Resume' : '❚❚ Pause';
    }

    // ----- Handlers -----
    function handleCellClick(e) {
        if (state.isPaused) return;
        const clickedCell = e.target.closest('.cell');
        if (!clickedCell) return;
        if (state.activeCell) state.activeCell.classList.remove('active');
        state.activeCell = clickedCell;
        state.activeCell.classList.add('active');
    }

    function handleKeyPress(e) {
        if (state.isPaused || !state.activeCell || state.activeCell.classList.contains('pre-filled')) return;
        const key = e.key;
        const finalValueEl = state.activeCell.querySelector('.final-value');
        const existingFinalValue = finalValueEl.textContent;

        if (key >= '1' && key <= '9') {
            const hasPencilMarks = state.activeCell.querySelectorAll('.pencil-mark.visible').length > 0;
            
            if (hasPencilMarks) {
                const mark = state.activeCell.querySelector(`.pencil-mark[data-mark="${key}"]`);
                if (mark) mark.classList.toggle('visible');
            } else if (existingFinalValue) {
                if (existingFinalValue !== key) {
                    finalValueEl.textContent = '';
                    state.activeCell.querySelector(`.pencil-mark[data-mark="${existingFinalValue}"]`).classList.add('visible');
                    state.activeCell.querySelector(`.pencil-mark[data-mark="${key}"]`).classList.add('visible');
                }
            } else {
                finalValueEl.textContent = key;
            }
        } else if (key === 'Backspace' || key === 'Delete') {
            finalValueEl.textContent = '';
            clearPencilMarks(state.activeCell);
        }
    }

    // ----- Utilities -----
    function clearPencilMarks(cell) {
        cell.querySelectorAll('.pencil-mark').forEach(m => m.classList.remove('visible'));
    }

    function createGrid() {
        elements.grid.innerHTML = '';
        for (let i = 0; i < 81; i++) {
            const cell = document.createElement('div');
            cell.classList.add('cell');
            cell.dataset.index = i;
            const finalValue = document.createElement('div');
            finalValue.classList.add('final-value');
            const pencilGrid = document.createElement('div');
            pencilGrid.classList.add('pencil-grid');
            for (let j = 1; j <= 9; j++) {
                const pencilMark = document.createElement('div');
                pencilMark.classList.add('pencil-mark');
                pencilMark.dataset.mark = j;
                pencilMark.textContent = j;
                pencilGrid.appendChild(pencilMark);
            }
            cell.appendChild(finalValue);
            cell.appendChild(pencilGrid);
            elements.grid.appendChild(cell);
        }
    }

    function populateGrid(puzzle) {
        const cells = elements.grid.children;
        for (let i = 0; i < 9; i++) {
            for (let j = 0; j < 9; j++) {
                const index = i * 9 + j;
                const cell = cells[index];
                const value = puzzle[i][j];
                const finalValueEl = cell.querySelector('.final-value');
                finalValueEl.textContent = value !== 0 ? value : '';
                clearPencilMarks(cell);
                cell.classList.toggle('pre-filled', value !== 0);
            }
        }
    }

    // ----- Game logic -----
    async function fetchNewGame() {
        state.isPaused = true;
        togglePause(); // Ensure game is unpaused
        timer.reset();
        timer.start();
        const difficulty = elements.difficultySelect.value;
        elements.message.textContent = 'Generate a new game...';
        elements.message.className = '';
        try {
            const response = await fetch(`/generate?difficulty=${difficulty}`);
            if (!response.ok) throw new Error('Network error');
            const data = await response.json();
            state.solution = data.solution;
            populateGrid(data.puzzle);
            elements.message.textContent = '';
        } catch (error) {
            elements.message.textContent = 'Error while game loading.';
            elements.message.className = 'incorrect';
        }
    }

    function checkSolution() {
        if (state.isPaused) return;
        const cells = elements.grid.children;
        let isComplete = true;
        let isCorrect = true;
        
        for (let i = 0; i < 9; i++) {
            for (let j = 0; j < 9; j++) {
                const userValueStr = cells[i * 9 + j].querySelector('.final-value').textContent;
                const userValue = parseInt(userValueStr, 10) || 0;
                
                if (userValue === 0) isComplete = false;
                if (userValue !== state.solution[i][j] && state.solution[i][j] !== 0) isCorrect = false;
            }
        }
        
        if (!isComplete) {
            elements.message.textContent = 'Puzzle not completed!';
            elements.message.className = 'incorrect';
        } else if (isCorrect) {
            timer.stop();
            elements.message.textContent = `Congratulations! You solved for ${timer.formatTime(state.secondsElapsed)}!`;
            elements.message.className = 'correct';
            saveStats();
        } else {
            elements.message.textContent = 'Error. Try one more time.';
            elements.message.className = 'incorrect';
        }
    }

    // ----- Statistics -----
    async function saveStats() {
        const playerName = elements.playerNameInput.value.trim();
        if (!playerName) {
            console.log("There no name, statistic won't be saved.");
            return;
        }
        const stats = {
            playerName: playerName,
            difficulty: elements.difficultySelect.value,
            timeTakenSeconds: state.secondsElapsed
        };
        try {
            const response = await fetch('/stats', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(stats)
            });
            if (response.ok) fetchAndDisplayStats();
        } catch (error) {
            console.error('Error while stats uploading:', error);
        }
    }

    async function fetchAndDisplayStats() {
        try {
            const response = await fetch('/stats');
            if (!response.ok) throw new Error('Network error');
            const stats = await response.json();
            let html = `<table><tr><th>Player</th><th>Difficulty</th><th>Time</th><th>Date</th></tr>`;
            if (stats) {
                stats.forEach(s => {
                    html += `<tr><td>${s.playerName}</td><td>${s.difficulty}</td><td>${timer.formatTime(s.timeTakenSeconds)}</td><td>${s.solvedAt}</td></tr>`;
                });
            }
            html += '</table>';
            elements.leaderboard.innerHTML = html;
        } catch (error) {
            elements.leaderboard.innerHTML = '<p>Failed to statistic loading.</p>';
        }
    }

    // ----- Initialization -----
    createGrid();
    handlePlayerName.load();
    fetchNewGame();
    fetchAndDisplayStats();

    // ----- Event Listeners -----
    elements.newGameBtn.addEventListener('click', fetchNewGame);
    elements.checkBtn.addEventListener('click', checkSolution);
    elements.grid.addEventListener('click', handleCellClick);
    elements.pauseBtn.addEventListener('click', togglePause);
    window.addEventListener('keydown', handleKeyPress);
    elements.playerNameInput.addEventListener('input', handlePlayerName.save);
});