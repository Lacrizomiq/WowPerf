# Google OAuth Deployment Guide

## 📋 Vue d'ensemble

Ce guide détaille le déploiement en production du système Google OAuth, incluant la configuration infrastructure, le monitoring, et les procédures de maintenance.

## 🚀 Déploiement Production

### 1. Prérequis Infrastructure

**Composants requis :**

- ✅ **Backend** : Go application avec HTTPS
- ✅ **Frontend** : React/Next.js avec routing OAuth
- ✅ **Base de données** : PostgreSQL avec contraintes uniques
- ✅ **Cache** : Redis pour tokens/sessions
- ✅ **Load Balancer** : Nginx/Traefik avec SSL termination
- ✅ **Monitoring** : Logs centralisés + métriques

### 2. Configuration Environment Production

**Variables d'environnement :**

```bash
# Google OAuth Production
GOOGLE_CLIENT_ID=prod-client-id.googleusercontent.com
GOOGLE_CLIENT_SECRET=GOCSPX-production-secret
GOOGLE_REDIRECT_URL=https://api.wowperf.com/auth/google/callback

# Frontend URLs
FRONTEND_URL=https://wowperf.com
FRONTEND_DASHBOARD_PATH=/dashboard
FRONTEND_AUTH_ERROR_PATH=/login

# Security
JWT_SECRET=production-jwt-secret-very-long-and-secure
ENCRYPTION_KEY=32-bytes-encryption-key-for-tokens

# Infrastructure
DOMAIN=wowperf.com
ENVIRONMENT=production
DATABASE_URL=postgresql://user:pass@db:5432/wowperf
REDIS_URL=redis://redis:6379/0
```

## 📋 Checklist Déploiement

### Pre-Deployment

- [ ] **Variables d'environnement** configurées et validées
- [ ] **Google Cloud Console** configuré avec URLs production
- [ ] **Certificats SSL** valides et renouvelables
- [ ] **Base de données** migrée avec contraintes uniques
- [ ] **Rate limiting** configuré sur Nginx / Traefik - Optionnel
- [ ] **Monitoring** et alertes configurés - Optionnel
- [ ] **Backup base de données** effectué

## 🚨 Troubleshooting Guide

### Erreurs Fréquentes

| Erreur                          | Symptôme                  | Solution                           |
| ------------------------------- | ------------------------- | ---------------------------------- |
| `redirect_uri_mismatch`         | Erreur Google Console     | Vérifier GOOGLE_REDIRECT_URL exact |
| `invalid_client`                | Erreur d'authentification | Vérifier GOOGLE_CLIENT_SECRET      |
| `email_not_verified`            | Utilisateurs rejetés      | Normal - emails non vérifiés       |
| `CSRF state mismatch`           | Échecs de callback        | Vérifier cookies sécurisés + HTTPS |
| `Database constraint violation` | Erreur création user      | Vérifier contraintes uniques       |

### Commandes de Debug

```bash


# Test manuel endpoint
curl -v https://api.wowperf.com/auth/google/login


```

---

🚀 **Déploiement production robuste et monitoring complet pour Google OAuth !**
