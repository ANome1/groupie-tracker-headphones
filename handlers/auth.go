package handlers

// RESPONSABLE: @Nome
// Handlers pour l'authentification (inscription, connexion, déconnexion)

// TODO @Nome: RegisterHandler - Gère l'inscription d'un nouvel utilisateur
// - Méthode GET: Afficher le formulaire d'inscription (templates/auth/register.html)
// - Méthode POST:
//   1. Récupérer username, email, password depuis le formulaire
//   2. Valider les données (utils/validator.go):
//      - Username: 3-20 caractères, alphanumérique
//      - Email: format valide
//      - Password: minimum 8 caractères
//   3. Vérifier que l'username et l'email n'existent pas déjà
//   4. Hasher le mot de passe (utils/hash.go)
//   5. Créer l'utilisateur en base de données
//   6. Créer une session (cookie)
//   7. Rediriger vers la landing page

// TODO @Nome: LoginHandler - Gère la connexion d'un utilisateur
// - Méthode GET: Afficher le formulaire de connexion (templates/auth/login.html)
// - Méthode POST:
//   1. Récupérer username, password depuis le formulaire
//   2. Récupérer l'utilisateur depuis la base de données
//   3. Vérifier le mot de passe (utils/hash.go)
//   4. Créer une session (cookie)
//   5. Rediriger vers la landing page

// TODO @Nome: LogoutHandler - Gère la déconnexion
// - Supprimer la session/cookie
// - Rediriger vers la landing page
