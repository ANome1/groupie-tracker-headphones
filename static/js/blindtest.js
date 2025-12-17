let timerInterval = null;
let currentRound = 0;
let totalRounds = 5;
let wsReady = false;
let pendingPlaylistSelect = null;
let currentRoundData = null;

window.handleBlindTestMessage = handleBlindTestMessage;

let playlistSelected = false;

window.onWSReady = function() {
    wsReady = true;
    if (typeof sendWebSocketMessage === 'function') {
        sendWebSocketMessage({ type: 'join_game' });
    }
};

document.addEventListener('DOMContentLoaded', () => {
    const playlistButtons = document.querySelectorAll('.playlist-btn');
    playlistButtons.forEach(btn => {
        btn.addEventListener('click', function(e) {
            if (!playlistSelected) {
                const playlistId = this.getAttribute('data-playlist-id');
                const playlistName = this.getAttribute('data-playlist-name');
                selectPlaylist(playlistId, playlistName);
                playlistSelected = true;
                this.style.opacity = '0.6';
                playlistButtons.forEach(b => {
                    if (b !== this) b.style.pointerEvents = 'none';
                });
            }
        });
    });

    const submitBtn = document.getElementById('submit-answer-btn');
    if (submitBtn) {
        submitBtn.addEventListener('click', submitAnswer);
    }

    const nextRoundBtn = document.getElementById('next-round-btn');
    if (nextRoundBtn) {
        nextRoundBtn.addEventListener('click', () => {
            sendWebSocketMessage({ type: 'next_round' });
        });
    }

    const newGameBtn = document.getElementById('new-game-btn');
    if (newGameBtn) {
        newGameBtn.addEventListener('click', () => {
            window.location.href = `/room/lobby?code=${window.roomCode}`;
        });
    }
});

// Envoi de la sélection de playlist au serveur
function selectPlaylist(playlistId, playlistName) {
    sendWebSocketMessage({
        type: 'select_playlist',
        playlist_id: playlistId,
        playlist_name: playlistName
    });
    
    document.getElementById('playlist-selection').style.display = 'none';
    document.getElementById('game-phase').style.display = 'block';
}

// Soumission de la réponse du joueur
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
    currentRoundData = data; // Stocker pour utilisation dans le timer
    currentRound = data.round_number;
    totalRounds = data.total_rounds || totalRounds;
    document.getElementById('current-round').textContent = currentRound;
    document.getElementById('total-rounds').textContent = totalRounds;
    
    // Cacher la sélection de playlist (au cas où c'est un client qui rejoint)
    document.getElementById('playlist-selection').style.display = 'none';
    
    const audioPlayer = document.getElementById('track-audio');
    if (data.preview_url) {
        const proxyUrl = `/audio/proxy?url=${encodeURIComponent(data.preview_url)}`;
        audioPlayer.pause();
        audioPlayer.currentTime = 0;
        audioPlayer.src = proxyUrl;
        
        const playAudio = () => {
            audioPlayer.play().catch(err => {
                console.error('Audio play error:', err);
            });
            audioPlayer.removeEventListener('canplay', playAudio);
        };
        
        audioPlayer.addEventListener('canplay', playAudio, { once: true });
        audioPlayer.load();
    }
    
    const coverDiv = document.getElementById('track-cover');
    const coverImg = document.getElementById('cover-image');
    if (data.cover_image) {
        coverImg.style.filter = 'blur(20px)';
        coverImg.src = data.cover_image;
        coverDiv.style.display = 'block';
    } else {
        coverDiv.style.display = 'none';
    }
    
    document.getElementById('answer-form').style.display = 'block';
    document.getElementById('answer-sent').style.display = 'none';
    document.getElementById('track-answer').value = '';
    document.getElementById('artist-answer').value = '';
    
    document.getElementById('results-phase').style.display = 'none';
    document.getElementById('game-phase').style.display = 'block';
    
    const timerDuration = data.remaining_time !== undefined ? Math.ceil(data.remaining_time) : (data.duration || 30);
    startTimer(timerDuration);
}

// Synchronisation client WebSocket: Timer avec révélation progressive du blur
function startTimer(duration) {
    let timeLeft = duration;
    const timerEl = document.getElementById('timer');
    timerEl.textContent = timeLeft;
    
    const revealTime = Math.ceil(duration / 3);
    
    if (timerInterval) clearInterval(timerInterval);
    
    timerInterval = setInterval(() => {
        timeLeft--;
        timerEl.textContent = timeLeft;
        
        if (timeLeft === revealTime) {
            const coverImg = document.getElementById('cover-image');
            if (coverImg) {
                coverImg.style.filter = 'none';
            }
        }
        
        if (timeLeft <= 0) {
            clearInterval(timerInterval);
            const answerForm = document.getElementById('answer-form');
            if (answerForm.style.display !== 'none') {
                revealAnswer();
            }
        }
    }, 1000);
}

// Révélation de la réponse: affiche le titre et l'artiste quand le timer expire
function revealAnswer() {
    if (!currentRoundData) return;
        
    const coverImg = document.getElementById('cover-image');
    if (coverImg) {
        coverImg.style.filter = 'none';
    }
    
    const trackNameDisplay = document.createElement('div');
    trackNameDisplay.id = 'revealed-track-info';
    trackNameDisplay.style.cssText = `
        text-align: center;
        margin-top: 1rem;
        padding: 1rem;
        background: rgba(166, 142, 224, 0.1);
        border-radius: 8px;
        border: 1px solid var(--primary);
    `;
    trackNameDisplay.innerHTML = `
        <p style="margin: 0.5rem 0;"><strong>Titre:</strong> ${currentRoundData.track_name || 'N/A'}</p>
        <p style="margin: 0.5rem 0;"><strong>Artiste:</strong> ${currentRoundData.artist_name || 'N/A'}</p>
    `;
    
    const trackCover = document.getElementById('track-cover');
    if (trackCover && !document.getElementById('revealed-track-info')) {
        trackCover.parentNode.insertBefore(trackNameDisplay, trackCover.nextSibling);
    }
    
    document.getElementById('answer-form').style.display = 'none';
    
    const answerSent = document.getElementById('answer-sent');
    if (answerSent) {
        answerSent.style.display = 'block';
        answerSent.querySelector('p').textContent = '⏰ Temps écoulé! La réponse a été révélée. En attente des autres joueurs...';
    }
}

// Affichage des résultats: scores de la manche et classement actualisé
function showRoundResults(data) {
    clearInterval(timerInterval);
    
    const audioPlayer = document.getElementById('track-audio');
    audioPlayer.pause();
    
    document.getElementById('correct-track').textContent = data.correct_track;
    document.getElementById('correct-artist').textContent = data.correct_artist;
    
    const roundScoresEl = document.getElementById('round-scores');
    roundScoresEl.innerHTML = '<h4>Scores de cette manche:</h4>';
    
    if (data.round_scores && Object.keys(data.round_scores).length > 0) {
        const medals = ['🥇', '🥈', '🥉', '4️⃣', '5️⃣'];
        const sortedEntries = Object.entries(data.round_scores)
            .sort((a, b) => b[1] - a[1])
            .slice(0, 5);
        
        sortedEntries.forEach((entry, index) => {
            const username = entry[0];
            const points = entry[1];
            const medal = medals[index] || '';
            roundScoresEl.innerHTML += `<p>${medal} <strong>${username}</strong>: +${points} points</p>`;
        });
    } else {
        roundScoresEl.innerHTML += '<p>Aucun joueur n\'a trouvé la bonne réponse.</p>';
    }
    
    document.getElementById('game-phase').style.display = 'none';
    document.getElementById('results-phase').style.display = 'block';
    
    if (parseInt(window.currentUserID) === parseInt(data.host_id)) {
        document.getElementById('next-round-btn').style.display = 'block';
    } else {
        document.getElementById('next-round-btn').style.display = 'none';
    }
}

// Mise à jour du scoreboard en direct
function updateScoreboard(scores) {
    const scoreboardEl = document.getElementById('scoreboard-content');
    
    if (!scoreboardEl) {
        return;
    }
    
    scoreboardEl.innerHTML = '';
    
    const sortedScores = Object.entries(scores).sort((a, b) => b[1] - a[1]);
    
    sortedScores.forEach(([username, score], index) => {
        const rank = index + 1;
        const playerDiv = document.createElement('div');
        playerDiv.className = 'scoreboard-item';
        playerDiv.innerHTML = `<span>${rank}. ${username}</span><span>${score} pts</span>`;
        scoreboardEl.appendChild(playerDiv);
    });
}

// Affichage du classement final avec palmarès
function showFinalResults(data) {
    document.getElementById('game-phase').style.display = 'none';
    document.getElementById('results-phase').style.display = 'none';
    document.getElementById('final-results').style.display = 'block';
    
    const finalScoreboardEl = document.getElementById('final-scoreboard');
    finalScoreboardEl.innerHTML = '<h2>🎵 Partie terminée! 🎵</h2><h3 style="margin-top: 30px; font-size: 1.8em; color: white;">Classement final</h3>';
    
    const sortedScores = Object.entries(data.scores).sort((a, b) => b[1] - a[1]);
    const medals = ['🥇', '🥈', '🥉'];
    
    sortedScores.forEach(([username, score], index) => {
        const rank = index + 1;
        const medal = medals[index] || '⭐';
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
                    <span class="rank-username">${username}</span>
                </div>
            </div>
            <div class="rank-points">${score} pts</div>
        `;
        finalScoreboardEl.appendChild(rankDiv);
    });
}

// Routage des messages WebSocket vers les handlers de phase
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
    }
}

// Gestion des playlists personnalisées
document.addEventListener('DOMContentLoaded', () => {
    const addPlaylistBtn = document.getElementById('add-playlist-btn');
    if (addPlaylistBtn) {
        addPlaylistBtn.addEventListener('click', () => {
            const idInput = document.getElementById('custom-playlist-id');
            const nameInput = document.getElementById('custom-playlist-name');
            
            const id = idInput.value.trim();
            const name = nameInput.value.trim();
            
            if (id && name) {
                addCustomPlaylist(id, name);
                idInput.value = '';
                nameInput.value = '';
            } else {
                alert("Veuillez entrer un ID et un Nom pour la playlist.");
            }
        });
    }
});

function addCustomPlaylist(id, name) {
    const grid = document.querySelector('.playlists-grid');
    if (!grid) return;

    const btn = document.createElement('button');
    btn.className = 'playlist-btn';
    btn.setAttribute('data-playlist-id', id);
    btn.setAttribute('data-playlist-name', name);
    btn.textContent = name;
    
    btn.addEventListener('click', () => {
        if (!playlistSelected) {
            selectPlaylist(id, name);
            playlistSelected = true;
        }
    });
    
    btn.style.animation = 'fadeInUp 0.5s ease';
    
    grid.appendChild(btn);
}
