# Google OAuth Security Guide

## 📋 Vue d'ensemble

Ce document détaille les mesures de sécurité implémentées dans le système Google OAuth de WowPerf et les bonnes pratiques à suivre.

## 🛡️ Mesures de Sécurité Implémentées

### 1. Protection CSRF (Cross-Site Request Forgery)

**Mécanisme :**

```go
// Génération state aléatoire (32 bytes)
state := base64.URLEncoding.EncodeToString(randomBytes(32))

// Stockage sécurisé en cookie temporaire
c.SetCookie("google_oauth_state", state, 600, "/", "", true, true)
                                        // 10min  HttpOnly Secure

// Validation au callback
if receivedState != storedState {
    return fmt.Errorf("CSRF attack detected")
}
```

**Protection contre :**

- Attaques CSRF sur initiation OAuth
- Replay attacks avec anciens states
- Manipulation des paramètres OAuth

### 2. Validation Stricte des Données Google

**Email vérifié obligatoire :**

```go
if !userInfo.VerifiedEmail {
    return &OAuthError{
        Code: "email_not_verified",
        Message: "Email not verified by Google",
    }
}
```

**Validation complète :**

```go
func validateGoogleUserInfo(userInfo *GoogleUserInfo) error {
    // 1. Email vérifié (recommandation Google)
    if !userInfo.VerifiedEmail {
        return errors.New("email not verified")
    }

    // 2. Champs obligatoires présents
    if userInfo.Email == "" || userInfo.ID == "" {
        return errors.New("missing required fields")
    }

    // 3. Format email valide (sanity check)
    if !isValidEmail(userInfo.Email) {
        return errors.New("invalid email format")
    }

    return nil
}
```

### 3. Cookies Sécurisés

**Configuration JWT cookies :**

```go
c.SetCookie(
    "access_token",
    accessToken,
    int(7*24*time.Hour.Seconds()), // 7 jours
    "/",                           // Path
    domain,                        // Domain
    true,                          // Secure (HTTPS uniquement)
    true,                          // HttpOnly (pas accessible JS)
)
```

**Configuration state cookie :**

```go
c.SetCookie(
    "google_oauth_state",
    state,
    600,                           // 10 minutes uniquement
    "/",                           // Path
    "",                            // Domain (auto)
    true,                          // Secure (HTTPS uniquement)
    true,                          // HttpOnly (pas accessible JS)
)
```

### 4. Gestion Sécurisée des Tokens

**Échange code → token avec timeout :**

```go
func ExchangeCodeForToken(ctx context.Context, code string) (*oauth2.Token, error) {
    // Timeout de 30 secondes
    ctxWithTimeout, cancel := context.WithTimeout(ctx, 30*time.Second)
    defer cancel()

    token, err := oauthConfig.Exchange(ctxWithTimeout, code)
    if err != nil {
        return nil, fmt.Errorf("token exchange failed: %w", err)
    }

    return token, nil
}
```

**Retry sécurisé avec backoff :**

```go
func GetUserInfoWithRetry(ctx context.Context, token *oauth2.Token) (*GoogleUserInfo, error) {
    const maxRetries = 3

    for i := 0; i < maxRetries; i++ {
        // Timeout par tentative
        ctxWithTimeout, cancel := context.WithTimeout(ctx, 30*time.Second)
        userInfo, err := GetUserInfo(ctxWithTimeout, token)
        cancel()

        if err == nil {
            return userInfo, nil
        }

        // Backoff exponentiel pour éviter le spam
        time.Sleep(time.Duration(1<<uint(i)) * time.Second)
    }

    return nil, fmt.Errorf("failed after %d retries", maxRetries)
}
```

## 🔐 Sécurité Base de Données

### Contraintes d'Intégrité

**Contraintes uniques :**

```go
type User struct {
    Username    string  `gorm:"uniqueIndex;not null"`
    Email       string  `gorm:"uniqueIndex;not null"`
    GoogleID    *string `gorm:"uniqueIndex"`       // NULL ou unique
    BattleNetID *string `gorm:"uniqueIndex"`       // NULL ou unique
}
```

### Transactions Sécurisées

**Création utilisateur atomique :**

```go
func CreateUser(user *models.User) error {
    tx := db.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()

    if err := tx.Create(user).Error; err != nil {
        tx.Rollback()
        return err
    }

    return tx.Commit().Error
}
```

### Gestion Race Conditions

**Username unique garanti :**

```go
func CreateUserFromGoogle(userInfo *GoogleUserInfo) (*models.User, error) {
    // Retry automatique en cas de conflit username
    for attempt := 1; attempt <= 100; attempt++ {
        username := generateUsername(userInfo, attempt)

        if user, err := tryCreateUser(username, userInfo); err == nil {
            return user, nil // ✅ Succès
        } else if !isUniqueConstraintError(err) {
            return nil, err // Erreur autre que contrainte → fail
        }
        // Contrainte unique → retry avec nouveau username
    }
}
```

## 🚨 Gestion des Erreurs et Attaques

### Types d'Erreurs Gérées

| Erreur               | Cause                           | Action                             |
| -------------------- | ------------------------------- | ---------------------------------- |
| `access_denied`      | Utilisateur refuse autorisation | Redirection avec message explicite |
| `invalid_request`    | Paramètres OAuth malformés      | Log + redirection erreur           |
| `server_error`       | Erreur côté Google              | Retry automatique                  |
| `email_not_verified` | Email non vérifié par Google    | Rejet avec message                 |
| `invalid_state`      | Attaque CSRF potentielle        | Rejet immédiat + log sécurité      |

### Gestion Centralisée des Erreurs

```go
func handleOAuthError(c *gin.Context, errorType, description string) {
    log.Printf("🚨 OAuth security event: type=%s, ip=%s, user_agent=%s",
        errorType, c.ClientIP(), c.GetHeader("User-Agent"))

    switch errorType {
    case "access_denied":
        redirectWithError(c, "auth_cancelled", "Authentication cancelled")
    case "invalid_state":
        // Possible CSRF attack
        log.Printf("🔴 POTENTIAL CSRF ATTACK: ip=%s", c.ClientIP())
        redirectWithError(c, "security_error", "Security validation failed")
    default:
        redirectWithError(c, "auth_failed", "Authentication failed")
    }
}
```

## 🔍 Logging et Monitoring de Sécurité

### Events de Sécurité Loggés

```go
// Tentatives d'authentification
log.Printf("🚀 OAuth initiation: ip=%s, user_agent=%s", c.ClientIP(), c.GetHeader("User-Agent"))

// Succès d'authentification
log.Printf("✅ OAuth success: method=%s, user_id=%d, ip=%s",
    result.Method, result.User.ID, c.ClientIP())

// Échecs de sécurité
log.Printf("🔴 OAuth security failure: error=%s, ip=%s", errorType, c.ClientIP())

// Nouvelles liaisons de compte
log.Printf("🔗 Account linking: user_id=%d, google_id=%s, ip=%s",
    user.ID, userInfo.ID, c.ClientIP())
```

### Métriques de Sécurité

**À monitorer :**

- Taux d'échec CSRF (invalid_state)
- Tentatives avec emails non vérifiés
- Patterns d'IPs suspectes
- Tentatives de liaison de comptes multiples
- Fréquence des retry username

## 🌐 Sécurité Réseau et Infrastructure

### HTTPS Obligatoire

**Configuration Production :**

```go
// Headers de sécurité
func securityHeaders() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Header("X-Content-Type-Options", "nosniff")
        c.Header("X-Frame-Options", "DENY")
        c.Header("X-XSS-Protection", "1; mode=block")

        if os.Getenv("ENVIRONMENT") == "production" {
            c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
        }

        c.Next()
    }
}
```

### CORS Sécurisé

```go
cors.Config{
    AllowOrigins: []string{
        "https://wowperf.com",        // Production uniquement
        "https://www.wowperf.com",    // Pas de wildcards
    },
    AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
    AllowHeaders: []string{
        "Content-Type", "Authorization", "X-CSRF-Token",
    },
    AllowCredentials: true,           // Pour les cookies
    MaxAge: 12 * time.Hour,
}
```

## 🔒 Secrets et Configuration

### Variables Sensibles

**⚠️ NE JAMAIS exposer :**

```bash
# Secrets absolus
GOOGLE_CLIENT_SECRET=xxx
JWT_SECRET=xxx
ENCRYPTION_KEY=xxx

# Informations système
DATABASE_URL=xxx
REDIS_PASSWORD=xxx
```

**✅ Gestion sécurisée :**

```bash
# Variables d'environnement
export GOOGLE_CLIENT_SECRET=$(vault kv get -field=secret auth/google)

# Kubernetes secrets
kubectl create secret generic google-oauth \
  --from-literal=client-secret="xxx"

# Docker secrets
echo "xxx" | docker secret create google_client_secret -
```

### Rotation des Secrets

**Stratégie de rotation :**

1. **JWT_SECRET** : Rotation mensuelle
2. **GOOGLE_CLIENT_SECRET** : Rotation trimestrielle
3. **State cookies** : Auto-expiration (10 minutes)

## 🎯 Checklist de Sécurité

### Pre-Deployment

- [ ] **HTTPS configuré** en production
- [ ] **Variables sensibles** dans vault/secrets
- [ ] **URIs de redirection** restreintes aux domaines légitimes
- [ ] **Client Secret** non committé dans le code
- [ ] **Logs de sécurité** configurés
- [ ] **Rate limiting** activé sur les endpoints OAuth
- [ ] **Monitoring** des métriques de sécurité

### Runtime

- [ ] **Validation email_verified** active
- [ ] **Protection CSRF** fonctionnelle
- [ ] **Cookies sécurisés** (Secure, HttpOnly)
- [ ] **Timeouts** appropriés (30s max)
- [ ] **Retry limité** (3 tentatives max)
- [ ] **Logs d'audit** complets
- [ ] **Alertes** sur events suspects

### Monitoring Continue

- [ ] **Taux d'échec CSRF** < 0.1%
- [ ] **Temps de réponse Google** < 5s
- [ ] **Rate d'emails non vérifiés** monitored
- [ ] **Patterns d'IPs** analysés
- [ ] **Tentatives multiples** détectées

## ⚠️ Threats et Mitigations

### Menaces Potentielles

| Menace                | Impact                       | Mitigation                        |
| --------------------- | ---------------------------- | --------------------------------- |
| **CSRF Attack**       | Liaison compte non autorisée | State validation + cookies secure |
| **Session Hijacking** | Vol de session               | HTTPS + HttpOnly + expiration     |
| **Replay Attack**     | Réutilisation anciens tokens | State unique + expiration courte  |
| **Email Spoofing**    | Faux comptes Google          | Validation email_verified         |
| **Race Condition**    | Doublon username             | Contraintes DB + retry logic      |

### Plan de Réponse Incident

**En cas de suspicion d'attaque :**

1. **Isolation**

   ```bash
   # Désactiver temporairement OAuth
   export GOOGLE_CLIENT_ID=""
   ```

2. **Investigation**

   ```bash
   # Analyser les logs
   grep "CSRF\|security" /var/log/app.log
   ```

3. **Remediation**
   - Rotation des secrets si nécessaire
   - Blacklist IPs suspectes
   - Notification utilisateurs si compromission

---

🛡️ **Sécurité robuste et multicouche pour l'authentification Google OAuth !**
