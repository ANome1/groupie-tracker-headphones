// blindtest.js - Client-side UI logic for Blind Test (WebSocket communication ONLY)

let timerInterval = null;
let currentRound = 0;
let totalRounds = 5;

// Exposer la fonction pour websocket.js
window.handleBlindTestMessage = handleBlindTestMessage;

// Sélection de playlist
document.addEventListener('DOMContentLoaded', () => {
    const playlistBtns = document.querySelectorAll('.playlist-btn');
    playlistBtns.forEach(btn => {
        btn.addEventListener('click', () => {
            const playlistId = btn.getAttribute('data-playlist-id');
            const playlistName = btn.getAttribute('data-playlist-name');
            selectPlaylist(playlistId, playlistName);
        });
    });

    // Bouton soumission réponse
    const submitBtn = document.getElementById('submit-answer-btn');
    if (submitBtn) {
        submitBtn.addEventListener('click', submitAnswer);
    }

    // Bouton manche suivante
    const nextRoundBtn = document.getElementById('next-round-btn');
    if (nextRoundBtn) {
        nextRoundBtn.addEventListener('click', () => {
            sendWebSocketMessage({
                type: 'next_round'
            });
        });
    }

    // Bouton nouvelle partie
    const newGameBtn = document.getElementById('new-game-btn');
    if (newGameBtn) {
        newGameBtn.addEventListener('click', () => {
            location.reload();
        });
    }
});

// Envoyer la sélection de playlist au serveur
function selectPlaylist(playlistId, playlistName) {
    sendWebSocketMessage({
        type: 'select_playlist',
        playlist_id: playlistId,
        playlist_name: playlistName
    });
    
    // Cacher la sélection
    document.getElementById('playlist-selection').style.display = 'none';
    document.getElementById('game-phase').style.display = 'block';
}

// Soumettre la réponse du joueur
function submitAnswer() {
    const trackAnswer = document.getElementById('track-answer').value.trim();
    const artistAnswer = document.getElementById('artist-answer').value.trim();
    
    if (!trackAnswer && !artistAnswer) {
        alert('Veuillez entrer au moins le titre ou l\'artiste');
        return;
    }
    
    sendWebSocketMessage({
        type: 'submit_answer',
        track: trackAnswer,
        artist: artistAnswer,
        username: window.currentUser
    });
    
    // Désactiver le formulaire
    document.getElementById('answer-form').style.display = 'none';
    document.getElementById('answer-sent').style.display = 'block';
}

// Démarrer une nouvelle manche (message reçu du serveur)
function startRound(data) {
    console.log('Starting round with data:', data);
    currentRound = data.round_number;
    document.getElementById('current-round').textContent = currentRound;
    document.getElementById('total-rounds').textContent = data.total_rounds || totalRounds;
    
    // Afficher le lecteur audio
    const audioPlayer = document.getElementById('track-audio');
    if (data.preview_url) {
        console.log('Setting audio source:', data.preview_url);
        audioPlayer.src = data.preview_url;
        audioPlayer.play().then(() => {
            console.log('Audio started playing');
        }).catch(err => {
            console.error('Audio play error:', err);
        });
    } else {
        console.error('No preview_url in data');
    }
    
    // Afficher la cover
    const coverDiv = document.getElementById('track-cover');
    const coverImg = document.getElementById('cover-image');
    if (data.cover_image) {
        coverImg.src = data.cover_image;
        coverDiv.style.display = 'block';
    } else {
        coverDiv.style.display = 'none';
    }
    
    // Réinitialiser le formulaire
    document.getElementById('answer-form').style.display = 'block';
    document.getElementById('answer-sent').style.display = 'none';
    document.getElementById('track-answer').value = '';
    document.getElementById('artist-answer').value = '';
    
    // Cacher les résultats
    document.getElementById('results-phase').style.display = 'none';
    
    // Démarrer le timer
    startTimer(data.duration || 30);
}

// Timer visuel
function startTimer(duration) {
    let timeLeft = duration;
    const timerEl = document.getElementById('timer');
    timerEl.textContent = timeLeft;
    
    if (timerInterval) clearInterval(timerInterval);
    
    timerInterval = setInterval(() => {
        timeLeft--;
        timerEl.textContent = timeLeft;
        
        if (timeLeft <= 0) {
            clearInterval(timerInterval);
            // Auto-submit si pas encore répondu
            if (document.getElementById('answer-form').style.display !== 'none') {
                submitAnswer();
            }
        }
    }, 1000);
}

// Afficher les résultats de la manche
function showRoundResults(data) {
    clearInterval(timerInterval);
    
    // Arrêter l'audio
    const audioPlayer = document.getElementById('track-audio');
    audioPlayer.pause();
    
    // Afficher la bonne réponse
    document.getElementById('correct-track').textContent = data.correct_track;
    document.getElementById('correct-artist').textContent = data.correct_artist;
    
    // Afficher les scores de la manche
    const roundScoresEl = document.getElementById('round-scores');
    roundScoresEl.innerHTML = '<h4>Scores de cette manche:</h4>';
    if (data.round_scores) {
        for (const [username, score] of Object.entries(data.round_scores)) {
            roundScoresEl.innerHTML += `<p>${username}: +${score} points</p>`;
        }
    }
    
    // Afficher la section résultats
    document.getElementById('game-phase').style.display = 'none';
    document.getElementById('results-phase').style.display = 'block';
    
    // Afficher le bouton pour l'hôte
    if (data.is_host) {
        document.getElementById('next-round-btn').style.display = 'block';
    }
}

// Mettre à jour le scoreboard
function updateScoreboard(scores) {
    const scoreboardEl = document.getElementById('scoreboard-content');
    scoreboardEl.innerHTML = '';
    
    // Trier par score décroissant
    const sortedScores = Object.entries(scores).sort((a, b) => b[1] - a[1]);
    
    sortedScores.forEach(([username, score], index) => {
        const rank = index + 1;
        const playerDiv = document.createElement('div');
        playerDiv.className = 'scoreboard-item';
        playerDiv.innerHTML = `<span>${rank}. ${username}</span><span>${score} pts</span>`;
        scoreboardEl.appendChild(playerDiv);
    });
}

// Afficher les résultats finaux
function showFinalResults(data) {
    document.getElementById('game-phase').style.display = 'none';
    document.getElementById('results-phase').style.display = 'none';
    document.getElementById('final-results').style.display = 'block';
    
    const finalScoreboardEl = document.getElementById('final-scoreboard');
    finalScoreboardEl.innerHTML = '<h3>Classement final</h3>';
    
    const sortedScores = Object.entries(data.scores).sort((a, b) => b[1] - a[1]);
    sortedScores.forEach(([username, score], index) => {
        const rank = index + 1;
        let medal = '';
        if (rank === 1) medal = '🥇';
        else if (rank === 2) medal = '🥈';
        else if (rank === 3) medal = '🥉';
        
        finalScoreboardEl.innerHTML += `<p>${medal} ${rank}. ${username}: ${score} points</p>`;
    });
}

// Handler pour les messages WebSocket (appelé depuis websocket.js)
function handleBlindTestMessage(data) {
    switch (data.type) {
        case 'round_start':
            startRound(data);
            break;
        case 'round_end':
            showRoundResults(data);
            break;
        case 'scoreboard_update':
            updateScoreboard(data.scores);
            break;
        case 'game_end':
            showFinalResults(data);
            break;
        default:
            console.log('Unknown message type:', data.type);
    }
}
