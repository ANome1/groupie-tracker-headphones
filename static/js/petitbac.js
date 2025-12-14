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

    const formData = new FormData(form);
    const answers = {};
    
    for (let [key, value] of formData.entries()) {
        answers[key] = value;
    }

    // Envoyer via WebSocket
    // Note: sendMessage est défini dans websocket.js
    if (typeof sendMessage === 'function') {
        sendMessage("SUBMIT_ANSWERS", {
            answers: answers
        });
    } else {
        console.error("sendMessage function not found");
    }

    // Désactiver le formulaire
    const inputs = form.querySelectorAll('input');
    const button = form.querySelector('button');
    inputs.forEach(input => input.disabled = true);
    if (button) {
        button.disabled = true;
        button.textContent = "Réponses envoyées !";
    }
}

// Fonction appelée par websocket.js lors de la réception de VALIDATION_PHASE
window.handleValidationPhase = function(roundData) {
    document.getElementById('answers-form').style.display = 'none';
    const validationDiv = document.getElementById('validation-phase');
    validationDiv.style.display = 'block';
    
    const container = document.getElementById('all-answers');
    container.innerHTML = ''; // Clear previous

    // Afficher les réponses pour validation
    for (const [playerID, response] of Object.entries(roundData.Responses)) {
        const playerDiv = document.createElement('div');
        playerDiv.className = 'validation-card';
        playerDiv.innerHTML = `<h4>Joueur ${playerID}</h4>`;
        
        const list = document.createElement('ul');
        for (const [category, answer] of Object.entries(response.Answers)) {
            const li = document.createElement('li');
            
            let voteButtons = '';
            // Ne pas afficher les boutons de vote pour ses propres réponses
            if (window.currentPlayerID && window.currentPlayerID !== playerID) {
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
        playerDiv.appendChild(list);
        container.appendChild(playerDiv);
    }
};

window.vote = function(targetPlayerID, category, isValid, btnElement) {
    if (typeof sendMessage === 'function') {
        sendMessage("SUBMIT_VOTE", {
            targetPlayerID: targetPlayerID,
            category: category,
            isValid: isValid
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
        const li = document.createElement('li');
        li.textContent = `Joueur ${playerID}: ${score} pts`;
        scoresList.appendChild(li);
    }
};
