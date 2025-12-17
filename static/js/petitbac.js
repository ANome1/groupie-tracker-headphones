document.addEventListener('DOMContentLoaded', () => {
    const timerElement = document.getElementById('timer');
    const letterElement = document.getElementById('letter');
    const form = document.getElementById('answers-form');
    
    const initialLetter = letterElement ? letterElement.innerText : '?';
    if (initialLetter !== '?' && initialLetter.length === 1) {
        animateLetter(initialLetter);
    }

    const duration = (window.gameConfig && window.gameConfig.timePerRound) ? window.gameConfig.timePerRound : 60;
    startTimer(duration);

    if (form) {
        form.addEventListener('submit', (e) => {
            e.preventDefault();
            submitAnswers();
        });
    }

    const submitBtn = document.querySelector('.submit-btn');
    if (submitBtn) {
        submitBtn.addEventListener('click', (e) => {
            e.preventDefault();
            submitAnswers();
        });
        submitBtn.addEventListener('touchstart', (e) => {
            e.preventDefault();
            submitAnswers();
        });
    }
    
    const newGameBtn = document.querySelector('.final-results-buttons button') || document.getElementById('new-game-btn');
    if (newGameBtn) {
        newGameBtn.addEventListener('click', () => {
            window.location.href = `/room/lobby?code=${window.roomCode}`;
        });
    }
});

window.submitAnswers = submitAnswers;

let timerInterval;

// Synchronisation du timer Petit Bac avec révélation progressive et soumission auto
function startTimer(duration) {
    let timer = duration;
    const timerElement = document.getElementById('timer');
    
    if (timerInterval) clearInterval(timerInterval);
    
    timerElement.classList.remove('warning');
    
    timerElement.textContent = timer;

    timerInterval = setInterval(() => {
        timer--;
        timerElement.textContent = timer;
        
        if (timer <= 10) {
            timerElement.classList.add('warning');
        }
        
        if (timer <= 0) {
            clearInterval(timerInterval);
            submitAnswers();
        }
    }, 1000);
}

function animateLetter(targetLetter) {
    const letterElement = document.getElementById('letter');
    const alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ";
    let iterations = 0;
    const maxIterations = 20;
    const speed = 50;

    const interval = setInterval(() => {
        letterElement.innerText = alphabet[Math.floor(Math.random() * alphabet.length)];
        iterations++;

        if (iterations >= maxIterations) {
            clearInterval(interval);
            letterElement.innerText = targetLetter;
            letterElement.style.transform = "scale(1.5)";
            setTimeout(() => {
                letterElement.style.transform = "scale(1)";
            }, 200);
        }
    }, speed);
}

function submitAnswers() {
    const form = document.getElementById('answers-form');
    if (!form) return;

    const inputs = form.querySelectorAll('input');
    if (inputs.length > 0 && inputs[0].disabled) {
        return;
    }

    const formData = new FormData(form);
    const answers = {};
    
    for (let [key, value] of formData.entries()) {
        answers[key] = value;
    }

    if (typeof sendMessage === 'function') {
        sendMessage("SUBMIT_ANSWERS", {
            Answers: answers
        });
    }

    const button = form.querySelector('button');
    inputs.forEach(input => input.disabled = true);
    if (button) {
        button.disabled = true;
        button.textContent = "Réponses envoyées !";
    }
}

// Phase de validation: synchronise les réponses et le vote entre joueurs
window.handleValidationPhase = function(data) {
    if (timerInterval) clearInterval(timerInterval);
    const timerElement = document.getElementById('timer');
    if (timerElement) {
        timerElement.classList.remove('warning');
        timerElement.textContent = "Vote";
        timerElement.style.fontSize = "1.5rem";
    }

    let roundData = data;
    let playerNames = {};
    let hostID = null;

    if (data.Round) {
        roundData = data.Round;
        playerNames = data.PlayerNames || {};
        hostID = data.HostID;
        window.playerNames = playerNames;
    }

    document.getElementById('answers-form').style.display = 'none';
    const validationDiv = document.getElementById('validation-phase');
    validationDiv.style.display = 'block';
    
    const container = document.getElementById('all-answers');
    container.innerHTML = '';

    const hostControls = document.getElementById('host-controls');
    if (hostControls) {
        hostControls.innerHTML = '';
        
        if (hostID && (String(window.currentPlayerID) === String(hostID) || String(window.currentPlayerID).includes(hostID))) {
            const button = document.createElement('button');
            button.onclick = () => nextRound();
            button.className = 'btn-primary';
            button.style.cssText = 'background-color: #4CAF50; color: white; padding: 12px 30px; border: none; border-radius: 5px; cursor: pointer; font-size: 16px; font-weight: bold;';
            button.textContent = '➡️ Manche Suivante';
            hostControls.appendChild(button);
        }
    }

    if (!roundData.Responses) {
        container.innerHTML = '<p>Aucune réponse reçue.</p>';
        return;
    }

    for (const [playerID, response] of Object.entries(roundData.Responses)) {
        const playerName = playerNames[playerID] || `Joueur ${playerID}`;
        const playerDiv = document.createElement('div');
        playerDiv.className = 'validation-card';
        playerDiv.innerHTML = `<h4>${playerName}</h4>`;
        
        const list = document.createElement('ul');
        if (response.Answers && Object.keys(response.Answers).length > 0) {
            for (const [category, answer] of Object.entries(response.Answers)) {
                if (category === 'code') continue;

                const li = document.createElement('li');
                
                let voteButtons = '';
                if (window.currentPlayerID && String(window.currentPlayerID) !== String(playerID)) {
                    voteButtons = `
                        <div class="vote-buttons">
                            <button class="vote-btn invalid" onclick="vote('${playerID}', '${category}', false, this)" title="Invalider">❌</button>
                            <button class="vote-btn valid" onclick="vote('${playerID}', '${category}', true, this)" title="Valider">✅</button>
                        </div>
                    `;
                }

                li.innerHTML = `
                    <div>
                        <strong>${category}</strong>
                        <span class="answer-text">${answer || '<em style="opacity:0.5">Pas de réponse</em>'}</span>
                    </div>
                    ${voteButtons}
                `;
                list.appendChild(li);
            }
        } else {
            list.innerHTML = '<li><em>Aucune réponse soumise</em></li>';
        }
        playerDiv.appendChild(list);
        container.appendChild(playerDiv);
    }
};

window.nextRound = function() {
    if (typeof sendMessage === 'function') {
        sendMessage("NEXT_ROUND", {});
    }
};

window.handleNewRound = function(roundUpdate) {
    document.getElementById('validation-phase').style.display = 'none';
    document.getElementById('answers-form').style.display = 'block';
    
    const form = document.getElementById('answers-form');
    form.reset();
    const inputs = form.querySelectorAll('input');
    inputs.forEach(input => input.disabled = false);
    const button = form.querySelector('button');
    if (button) {
        button.disabled = false;
        button.textContent = "Valider mes réponses";
    }

    const roundElement = document.getElementById('current-round');
    if (roundElement && roundUpdate.RoundNumber) {
        roundElement.textContent = roundUpdate.RoundNumber;
    }

    const letterElement = document.getElementById('letter');
    if (letterElement) {
        letterElement.innerText = roundUpdate.Letter;
        animateLetter(roundUpdate.Letter);
    }
    
    const timerElement = document.getElementById('timer');
    if (timerElement) {
        timerElement.style.fontSize = "";
    }
    
    startTimer(roundUpdate.Duration);
};

window.handleGameOver = function(data) {
    clearInterval(timerInterval);
    showFinalResults(data);
};

window.handleRoundResults = function(roundResults) {
    clearInterval(timerInterval);
    showRoundResults(roundResults);
};

function showRoundResults(roundResults) {
    document.getElementById('answers-form').style.display = 'none';
    const votingSection = document.querySelector('[data-phase="voting"]');
    if (votingSection) votingSection.style.display = 'none';
    
    const resultsDiv = document.getElementById('round-results') || createRoundResultsDisplay();
    resultsDiv.innerHTML = `<h2>Résultats de la manche ${roundResults.RoundNumber}</h2>`;
    
    for (const [playerID, answers] of Object.entries(roundResults.Answers)) {
        const playerName = (window.playerNames && window.playerNames[playerID]) || `Joueur ${playerID}`;
        const roundScore = roundResults.RoundScores[playerID] || 0;
        
        let playerResultHTML = `<div class="player-round-result"><h3>${playerName}: +${roundScore} pts</h3>`;
        
        for (const [category, answerData] of Object.entries(answers)) {
            if (!answerData.Answer) {
                playerResultHTML += `<p class="no-answer">${category}: (pas de réponse)</p>`;
            } else {
                const icon = answerData.Valid ? (answerData.Unique ? '✅ 2pts' : '✅ 1pt') : '❌ 0pts';
                playerResultHTML += `<p class="answer-line ${answerData.Valid ? 'valid' : 'invalid'}">${category}: ${answerData.Answer} ${icon}</p>`;
            }
        }
        playerResultHTML += '</div>';
        resultsDiv.innerHTML += playerResultHTML;
    }
    
    if (window.currentUserID == window.hostID) {
        const btnDiv = document.createElement('div');
        btnDiv.innerHTML = '<button onclick="nextRound()" class="submit-btn">Manche suivante</button>';
        resultsDiv.appendChild(btnDiv);
    }
    
    resultsDiv.style.display = 'block';
}

function showFinalResults(data) {
    const answersForm = document.getElementById('answers-form');
    if (answersForm) answersForm.style.display = 'none';
    
    const validationPhase = document.getElementById('validation-phase');
    if (validationPhase) validationPhase.style.display = 'none';
    
    const letterDisplay = document.querySelector('.letter-display');
    if (letterDisplay) letterDisplay.style.display = 'none';
    
    const roundInfo = document.querySelector('.round-info');
    if (roundInfo) roundInfo.style.display = 'none';
    
    const timer = document.querySelector('.timer');
    if (timer) timer.style.display = 'none';
    
    const finalDiv = document.getElementById('final-results') || createFinalResultsDisplay();
    finalDiv.innerHTML = '';
    
    const title = document.createElement('h2');
    title.innerHTML = '🎮 Partie terminée! 🎮';
    finalDiv.appendChild(title);
    
    const subtitle = document.createElement('h3');
    subtitle.textContent = 'Classement final';
    finalDiv.appendChild(subtitle);
    
    const scoreboardDiv = document.createElement('div');
    scoreboardDiv.id = 'final-scoreboard';
    
    const scores = (data && data.scores) ? data.scores : data;
    const sortedScores = Object.entries(scores || {}).sort((a, b) => b[1] - a[1]);
    const medals = ['🥇', '🥈', '🥉'];
    
    sortedScores.forEach((entry, index) => {
        const [playerID, score] = entry;
        const playerName = (window.playerNames && window.playerNames[playerID]) || `Joueur ${playerID}`;
        const medal = medals[index] || '⭐';
        const rank = index + 1;
        
        let rankClass = '';
        if (rank === 1) rankClass = 'first';
        else if (rank === 2) rankClass = 'second';
        else if (rank === 3) rankClass = 'third';
        
        const rankDiv = document.createElement('div');
        rankDiv.className = `final-rank ${rankClass}`;
        rankDiv.innerHTML = `
            <span class="rank-medal">${medal}</span>
            <div class="rank-info">
                <div>
                    <span class="rank-number">#${rank}</span>
                    <span class="rank-username">${playerName}</span>
                </div>
            </div>
            <div class="rank-points">${score} pts</div>
        `;
        scoreboardDiv.appendChild(rankDiv);
    });
    
    finalDiv.appendChild(scoreboardDiv);
    
    const btnDiv = document.createElement('div');
    btnDiv.className = 'final-results-buttons';
    const newGameBtn = document.createElement('button');
    newGameBtn.className = 'submit-btn';
    newGameBtn.id = 'new-game-btn';
    newGameBtn.textContent = 'Nouvelle partie';
    newGameBtn.addEventListener('click', () => {
        window.location.href = `/room/lobby?code=${window.roomCode}`;
    });
    btnDiv.appendChild(newGameBtn);
    finalDiv.appendChild(btnDiv);
    
    finalDiv.style.display = 'block';
}

function createRoundResultsDisplay() {
    const div = document.createElement('div');
    div.id = 'round-results';
    div.className = 'results-phase';
    document.body.appendChild(div);
    return div;
}

function createFinalResultsDisplay() {
    const div = document.createElement('div');
    div.id = 'final-results';
    div.className = 'final-results-phase';
    document.body.appendChild(div);
    return div;
}

// Routage du vote sur les réponses: vérifie la validité et synchronise avec serveur
window.vote = function(targetPlayerID, category, isValid, btnElement) {
    if (typeof sendMessage === 'function') {
        sendMessage("SUBMIT_VOTE", {
            TargetPlayer: targetPlayerID,
            Category: category,
            IsValid: isValid
        });
        
        const parent = btnElement.parentElement;
        const buttons = parent.querySelectorAll('.vote-btn');
        buttons.forEach(btn => btn.classList.remove('active'));
        btnElement.classList.add('active');
    }
};

window.updateScores = function(scores) {
    const scoresList = document.getElementById('scores-list');
    if (!scoresList) return;
    
    scoresList.innerHTML = '';
    for (const [playerID, score] of Object.entries(scores)) {
        const playerName = (window.playerNames && window.playerNames[playerID]) || `Joueur ${playerID}`;
        const li = document.createElement('li');
        li.textContent = `${playerName}: ${score} pts`;
        scoresList.appendChild(li);
    }
};

// Gestion des catégories personnalisées
let currentCategories = ["Artiste", "Groupe de musique", "Album", "Instrument", "Featuring"];

function renderCategoriesConfig() {
    const container = document.getElementById('categories-list-config');
    if (!container) return;
    
    container.innerHTML = '';
    currentCategories.forEach((cat, index) => {
        const tag = document.createElement('div');
        tag.className = 'category-tag';
        tag.style.cssText = 'background: var(--primary); padding: 5px 10px; border-radius: 15px; display: flex; align-items: center; gap: 5px; font-size: 0.9rem;';
        tag.innerHTML = `
            ${cat}
            <span onclick="removeCategory(${index})" style="cursor: pointer; font-weight: bold; margin-left: 5px;">&times;</span>
        `;
        container.appendChild(tag);
    });
}

function addCategory() {
    if (currentCategories.length >= 5) {
        alert("Maximum 5 catégories ! Supprimez-en une pour en ajouter.");
        return;
    }
    const input = document.getElementById('new-category-input');
    const val = input.value.trim();
    if (val) {
        currentCategories.push(val);
        input.value = '';
        renderCategoriesConfig();
    }
}

function removeCategory(index) {
    currentCategories.splice(index, 1);
    renderCategoriesConfig();
}

window.removeCategory = removeCategory;

document.addEventListener('DOMContentLoaded', () => {
    renderCategoriesConfig();
    
    const addBtn = document.getElementById('add-category-btn');
    if (addBtn) {
        addBtn.addEventListener('click', addCategory);
    }
    
    const startBtn = document.getElementById('start-game-btn');
    if (startBtn) {
        startBtn.addEventListener('click', () => {
            const numRounds = parseInt(document.getElementById('config-num-rounds').value) || 5;
            const timePerRound = parseInt(document.getElementById('config-time-round').value) || 60;
            
            if (currentCategories.length === 0) {
                alert("Il faut au moins une catégorie !");
                return;
            }

            if (window.socket && window.socket.readyState === WebSocket.OPEN) {
                window.socket.send(JSON.stringify({
                    type: 'START_GAME',
                    config: {
                        categories: currentCategories,
                        num_rounds: numRounds,
                        time_per_round: timePerRound
                    }
                }));
            }
        });
    }
});
