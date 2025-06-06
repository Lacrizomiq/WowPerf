# Google OAuth Setup Guide

## 📋 Vue d'ensemble

Ce guide détaille la configuration complète de l'authentification Google OAuth pour WowPerf, permettant aux utilisateurs de se connecter avec leur compte Google.

## 🎯 Prérequis

- Compte Google Cloud Console
- Projet WowPerf déployé avec HTTPS
- Accès aux variables d'environnement

## 1. Configuration Google Cloud Console

### 1.1 Créer/Configurer le Projet

1. **Accéder à Google Cloud Console**

   - Aller sur [Google Cloud Console](https://console.cloud.google.com/)
   - Sélectionner ou créer un projet

2. **Activer les APIs nécessaires**
   ```
   APIs à activer :
   - Google+ API (ou People API)
   - Google Identity and Access Management (IAM) API
   ```

### 1.2 Créer les Identifiants OAuth 2.0

1. **Navigation**

   ```
   APIs & Services > Credentials > Create Credentials > OAuth 2.0 Client IDs
   ```

2. **Configuration de l'Application**

   ```
   Application type: Web application
   Name: WowPerf OAuth Client
   ```

3. **Origines JavaScript Autorisées**

   ```
   Development:
   - https://localhost

   Production:
   - https://wowperf.com
   - https://www.wowperf.com
   ```

4. **URIs de Redirection Autorisées**

   ```
   Development:
   - https://localhost/api/auth/google/callback

   Production:
   - https://api.wowperf.com/auth/google/callback
   - https://wowperf.com/api/auth/google/callback
   ```

### 1.3 Récupérer les Identifiants

Après création, noter :

- **Client ID** : `xxx.apps.googleusercontent.com`
- **Client Secret** : `GOCSPX-xxx`

## 2. Configuration Variables d'Environnement

### 2.1 Fichier `.env` (Development)

```env
# Google OAuth Configuration
GOOGLE_CLIENT_ID=votre-client-id.googleusercontent.com
GOOGLE_CLIENT_SECRET=votre-client-secret
GOOGLE_REDIRECT_URL=https://localhost/api/auth/google/callback

# Frontend URLs pour redirection après auth
FRONTEND_URL=https://localhost
FRONTEND_DASHBOARD_PATH=/dashboard
FRONTEND_AUTH_ERROR_PATH=/login

# Autres variables requises
DOMAIN=localhost
JWT_SECRET=votre-jwt-secret-securise
```

### 2.2 Variables Production

```env
# Google OAuth Configuration
GOOGLE_CLIENT_ID=prod-client-id.googleusercontent.com
GOOGLE_CLIENT_SECRET=prod-client-secret
GOOGLE_REDIRECT_URL=https://api.wowperf.com/auth/google/callback

# Frontend URLs pour redirection après auth
FRONTEND_URL=https://wowperf.com
FRONTEND_DASHBOARD_PATH=/dashboard
FRONTEND_AUTH_ERROR_PATH=/login

# Production settings
DOMAIN=wowperf.com
JWT_SECRET=production-jwt-secret-tres-securise
```

### 2.3 Docker Compose

```yaml
environment:
  # ... autres variables ...
  - GOOGLE_CLIENT_ID=${GOOGLE_CLIENT_ID}
  - GOOGLE_CLIENT_SECRET=${GOOGLE_CLIENT_SECRET}
  - GOOGLE_REDIRECT_URL=${GOOGLE_REDIRECT_URL}
  - FRONTEND_URL=${FRONTEND_URL}
  - FRONTEND_DASHBOARD_PATH=${FRONTEND_DASHBOARD_PATH}
  - FRONTEND_AUTH_ERROR_PATH=${FRONTEND_AUTH_ERROR_PATH}
```

## 3. Configuration Frontend

### 3.1 Bouton de Connexion Google

```html
<!-- Exemple de bouton -->
<a href="/api/auth/google/login" class="google-login-btn">
  <img src="google-icon.svg" alt="Google" />
  Se connecter avec Google
</a>
```

### 3.2 Pages de Redirection

**Page de Dashboard** (`/dashboard`)

- Page où arrivent les utilisateurs existants après connexion

**Page d'Erreur** (`/login`)

- Page d'erreur avec paramètres : `?error=code&message=description`

## 4. Validation de la Configuration

### 4.1 Checklist de Vérification

- [ ] Projet Google Cloud créé
- [ ] APIs activées (Google+ ou People API)
- [ ] Identifiants OAuth 2.0 créés
- [ ] URIs de redirection configurées
- [ ] Variables d'environnement définies
- [ ] Frontend configuré avec bouton Google

### 4.2 Test de Configuration

1. **Démarrer l'application**

   ```bash
   go run cmd/server/main.go
   ```

2. **Tester l'initiation**

   ```
   GET https://localhost/api/auth/google/login
   → Doit rediriger vers Google
   ```

3. **Vérifier les logs**
   ```
   ✅ "Google OAuth routes registered"
   ✅ "Initiating Google OAuth flow"
   ✅ "Google OAuth state cookie set"
   ```

## 5. Dépannage

### 5.1 Erreurs Communes

**❌ "redirect_uri_mismatch"**

```
Solution : Vérifier que l'URI de redirection est exactement
la même dans Google Cloud Console et GOOGLE_REDIRECT_URL
```

**❌ "invalid_client"**

```
Solution : Vérifier GOOGLE_CLIENT_ID et GOOGLE_CLIENT_SECRET
```

**❌ "access_denied"**

```
Normal : L'utilisateur a refusé l'autorisation
```

### 5.2 Logs de Debug

Activer les logs détaillés :

```go
log.Printf("🔍 Google API response: %s", responseBody)
log.Printf("🔍 Decoded - ID:'%s', Email:'%s'", userInfo.ID, userInfo.Email)
```

## 6. Sécurité

### 6.1 Bonnes Pratiques

- ✅ **HTTPS obligatoire** en production
- ✅ **Client Secret sécurisé** (variables d'environnement)
- ✅ **URIs de redirection restreintes**
- ✅ **Validation email_verified**
- ✅ **Protection CSRF** via state

### 6.2 Variables Sensibles

**⚠️ Ne JAMAIS commiter :**

- `GOOGLE_CLIENT_SECRET`
- `JWT_SECRET`
- Fichiers `.env`

## 7. Ressources

- [Google OAuth 2.0 Documentation](https://developers.google.com/identity/protocols/oauth2)
- [Google Cloud Console](https://console.cloud.google.com/)
- [OAuth 2.0 Security Best Practices](https://tools.ietf.org/html/draft-ietf-oauth-security-topics)

---

✅ **Configuration terminée !** Votre application supporte maintenant l'authentification Google OAuth.
