document.addEventListener('DOMContentLoaded', () => {
    // Initialisation
    const timerElement = document.getElementById('timer');
    const letterElement = document.getElementById('letter');
    const form = document.getElementById('answers-form');
    
    // Récupérer la lettre initiale si elle existe
    const initialLetter = letterElement ? letterElement.innerText : '?';
    if (initialLetter !== '?' && initialLetter.length === 1) {
        animateLetter(initialLetter);
    }

    // Démarrer le timer (60s par défaut)
    // TODO: Récupérer la durée réelle depuis le serveur
    startTimer(60);

    // Gestion de la soumission du formulaire
    if (form) {
        form.addEventListener('submit', (e) => {
            e.preventDefault();
            submitAnswers();
        });
    }

    // Attacher l'événement au bouton de soumission explicitement
    const submitBtn = document.querySelector('.submit-btn');
    if (submitBtn) {
        submitBtn.addEventListener('click', (e) => {
            e.preventDefault();
            console.log("Bouton valider cliqué");
            submitAnswers();
        });
        // Support tactile pour mobile
        submitBtn.addEventListener('touchstart', (e) => {
            e.preventDefault();
            console.log("Bouton valider touché");
            submitAnswers();
        });
    }
    
    // Bouton nouvelle partie - Rediriger vers le lobby
    const newGameBtn = document.querySelector('.final-results-buttons button') || document.getElementById('new-game-btn');
    if (newGameBtn) {
        newGameBtn.addEventListener('click', () => {
            window.location.href = `/room/lobby?code=${window.roomCode}`;
        });
    }
});

// Expose submitAnswers to global scope for button onclick
window.submitAnswers = submitAnswers;

let timerInterval;

function startTimer(duration) {
    let timer = duration;
    const timerElement = document.getElementById('timer');
    
    if (timerInterval) clearInterval(timerInterval);
    
    timerElement.classList.remove('warning');
    
    // Mise à jour immédiate
    timerElement.textContent = timer;

    timerInterval = setInterval(() => {
        timer--; // Décrémenter d'abord
        timerElement.textContent = timer;
        
        if (timer <= 10) {
            timerElement.classList.add('warning');
        }
        
        if (timer <= 0) {
            clearInterval(timerInterval);
            // Auto-submit when time is up
            submitAnswers();
        }
    }, 1000);
}

function animateLetter(targetLetter) {
    const letterElement = document.getElementById('letter');
    const alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ";
    let iterations = 0;
    const maxIterations = 20; // Nombre de changements de lettre
    const speed = 50; // Vitesse initiale en ms

    const interval = setInterval(() => {
        letterElement.innerText = alphabet[Math.floor(Math.random() * alphabet.length)];
        iterations++;

        if (iterations >= maxIterations) {
            clearInterval(interval);
            letterElement.innerText = targetLetter;
            // Effet de "pop" final
            letterElement.style.transform = "scale(1.5)";
            setTimeout(() => {
                letterElement.style.transform = "scale(1)";
            }, 200);
        }
    }, speed);
}

function submitAnswers() {
    const form = document.getElementById('answers-form');
    if (!form) return; // Sécurité si le formulaire est déjà caché

    // Check if already submitted (inputs disabled)
    const inputs = form.querySelectorAll('input');
    if (inputs.length > 0 && inputs[0].disabled) {
        console.log("Already submitted, ignoring.");
        return;
    }

    const formData = new FormData(form);
    const answers = {};
    
    for (let [key, value] of formData.entries()) {
        answers[key] = value;
    }

    // Envoyer via WebSocket
    // Note: sendMessage est défini dans websocket.js
    if (typeof sendMessage === 'function') {
        console.log("Sending answers:", answers);
        sendMessage("SUBMIT_ANSWERS", {
            Answers: answers
        });
    } else {
        console.error("sendMessage function not found");
    }

    // Désactiver le formulaire
    const button = form.querySelector('button');
    inputs.forEach(input => input.disabled = true);
    if (button) {
        button.disabled = true;
        button.textContent = "Réponses envoyées !";
    }
}

// Fonction appelée par websocket.js lors de la réception de VALIDATION_PHASE
window.handleValidationPhase = function(data) {
    console.log("Received VALIDATION_PHASE data:", data);

    // Stop timer
    if (timerInterval) clearInterval(timerInterval);
    const timerElement = document.getElementById('timer');
    if (timerElement) {
        timerElement.classList.remove('warning');
        timerElement.textContent = "Vote";
        timerElement.style.fontSize = "1.5rem"; // Adjust font size for text
    }

    let roundData = data;
    let playerNames = {};
    let hostID = null;

    if (data.Round) {
        roundData = data.Round;
        playerNames = data.PlayerNames || {};
        hostID = data.HostID;
        window.playerNames = playerNames; // Store for scoreboard
        console.log("VALIDATION_PHASE - HostID:", hostID, "CurrentPlayerID:", window.currentPlayerID);
    }

    document.getElementById('answers-form').style.display = 'none';
    const validationDiv = document.getElementById('validation-phase');
    validationDiv.style.display = 'block';
    
    const container = document.getElementById('all-answers');
    container.innerHTML = ''; // Clear previous

    // Show Host Button if applicable
    const hostControls = document.getElementById('host-controls');
    if (hostControls) {
        hostControls.innerHTML = ''; // Clear previous
        
        // Always show button for host
        if (hostID && (String(window.currentPlayerID) === String(hostID) || String(window.currentPlayerID).includes(hostID))) {
            console.log("Showing NEXT_ROUND button for host");
            const button = document.createElement('button');
            button.onclick = () => nextRound();
            button.className = 'btn-primary';
            button.style.cssText = 'background-color: #4CAF50; color: white; padding: 12px 30px; border: none; border-radius: 5px; cursor: pointer; font-size: 16px; font-weight: bold;';
            button.textContent = '➡️ Manche Suivante';
            hostControls.appendChild(button);
        } else {
            console.log("Host check failed - HostID:", hostID, "CurrentPlayerID:", window.currentPlayerID);
        }
    }

    if (!roundData.Responses) {
        console.error("No responses found in round data");
        container.innerHTML = '<p>Aucune réponse reçue.</p>';
        return;
    }

    // Afficher les réponses pour validation
    for (const [playerID, response] of Object.entries(roundData.Responses)) {
        console.log("Processing response for player:", playerID, response);
        const playerName = playerNames[playerID] || `Joueur ${playerID}`;
        const playerDiv = document.createElement('div');
        playerDiv.className = 'validation-card';
        playerDiv.innerHTML = `<h4>${playerName}</h4>`;
        
        const list = document.createElement('ul');
        if (response.Answers && Object.keys(response.Answers).length > 0) {
            for (const [category, answer] of Object.entries(response.Answers)) {
                // Skip hidden fields like 'code'
                if (category === 'code') continue;

                const li = document.createElement('li');
                
                let voteButtons = '';
                // Allow voting for everyone except self
                // Also ensure host can vote if they are playing
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
    // Reset UI
    document.getElementById('validation-phase').style.display = 'none';
    document.getElementById('answers-form').style.display = 'block';
    
    // Reset inputs
    const form = document.getElementById('answers-form');
    form.reset();
    const inputs = form.querySelectorAll('input');
    inputs.forEach(input => input.disabled = false);
    const button = form.querySelector('button');
    if (button) {
        button.disabled = false;
        button.textContent = "Valider mes réponses";
    }

    // Update Round Number
    const roundElement = document.getElementById('current-round');
    if (roundElement && roundUpdate.RoundNumber) {
        roundElement.textContent = roundUpdate.RoundNumber;
    }

    // Update Letter and Timer
    const letterElement = document.getElementById('letter');
    if (letterElement) {
        letterElement.innerText = roundUpdate.Letter;
        animateLetter(roundUpdate.Letter);
    }
    
    const timerElement = document.getElementById('timer');
    if (timerElement) {
        timerElement.style.fontSize = ""; // Reset font size
    }
    
    startTimer(roundUpdate.Duration);
};

window.handleGameOver = function(data) {
    // Afficher les résultats finaux
    clearInterval(timerInterval);
    showFinalResults(data);
};

window.handleRoundResults = function(roundResults) {
    // Afficher les résultats détaillés de la manche
    clearInterval(timerInterval);
    showRoundResults(roundResults);
};

function showRoundResults(roundResults) {
    // Hide voting phase
    document.getElementById('answers-form').style.display = 'none';
    const votingSection = document.querySelector('[data-phase="voting"]');
    if (votingSection) votingSection.style.display = 'none';
    
    // Create results display
    const resultsDiv = document.getElementById('round-results') || createRoundResultsDisplay();
    resultsDiv.innerHTML = `<h2>Résultats de la manche ${roundResults.RoundNumber}</h2>`;
    
    // Show answers and points
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
    
    // Add next round button for host
    if (window.currentUserID == window.hostID) {
        const btnDiv = document.createElement('div');
        btnDiv.innerHTML = '<button onclick="nextRound()" class="submit-btn">Manche suivante</button>';
        resultsDiv.appendChild(btnDiv);
    }
    
    resultsDiv.style.display = 'block';
}

function showFinalResults(data) {
    // Hide game form and validation phase
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
    
    // Show final results
    const finalDiv = document.getElementById('final-results') || createFinalResultsDisplay();
    finalDiv.innerHTML = ''; // Clear previous content
    
    // Add title
    const title = document.createElement('h2');
    title.innerHTML = '🎮 Partie terminée! 🎮';
    finalDiv.appendChild(title);
    
    const subtitle = document.createElement('h3');
    subtitle.textContent = 'Classement final';
    finalDiv.appendChild(subtitle);
    
    // Create scoreboard container
    const scoreboardDiv = document.createElement('div');
    scoreboardDiv.id = 'final-scoreboard';
    
    // Sort and display final scores - handle both wrapped and direct data
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
    
    // Add button to go back to lobby
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

window.nextRound = function() {
    sendMessage("NEXT_ROUND", {});
};

window.endRound = function() {
    sendMessage("END_ROUND", {});
};

window.vote = function(targetPlayerID, category, isValid, btnElement) {
    if (typeof sendMessage === 'function') {
        console.log("Sending vote:", { TargetPlayer: targetPlayerID, Category: category, IsValid: isValid });
        sendMessage("SUBMIT_VOTE", {
            TargetPlayer: targetPlayerID,
            Category: category,
            IsValid: isValid
        });
        
        // Feedback visuel
        const parent = btnElement.parentElement;
        const buttons = parent.querySelectorAll('.vote-btn');
        buttons.forEach(btn => btn.classList.remove('active'));
        btnElement.classList.add('active');
    }
};

// Fonction pour mettre à jour les scores (appelée par websocket.js)
window.updateScores = function(scores) {
    const scoresList = document.getElementById('scores-list');
    if (!scoresList) return;
    
    scoresList.innerHTML = '';
    // scores est une map: PlayerID -> Score
    for (const [playerID, score] of Object.entries(scores)) {
        const playerName = (window.playerNames && window.playerNames[playerID]) || `Joueur ${playerID}`;
        const li = document.createElement('li');
        li.textContent = `${playerName}: ${score} pts`;
        scoresList.appendChild(li);
    }
};
