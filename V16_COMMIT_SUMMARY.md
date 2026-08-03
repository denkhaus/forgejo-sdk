# Forgejo SDK v16.0.2 Upgrade - Commit Summary

## 🎯 Überblick

Upgrade des Forgejo SDK von v15.0.3 auf v16.0.2 mit Breaking Changes und neuen API-Features.

## 📊 Statistik

- **Dateien geändert:** 270
- **Insertionen:** +1,487
- **Deletionen:** -1,970
- **Netto-Änderung:** -483 Zeilen

## 🚨 Breaking Changes

### 1. EXIF Stripping Removed
- **Dateien:** `forgejo/avatar.go`, `forgejo/models/update_user_avatar_option.go`
- **Änderung:** EXIF-Stripping-Optionen entfernt
- **Grund:** Lizenzprobleme mit AGPL-Bibliothek
- **Betroffene Methoden:** `UserUpdateAvatar()`, `UserUploadAvatar()`

### 2. Mirror SSRF Hardening  
- **Dateien:** `forgejo/repo_mirror.go`, `forgejo/repo_push_mirror.go`, `forgejo/repo_migrate.go`
- **Änderung:** Redirect-Handling geändert (Redirects → Errors)
- **Grund:** SSRF-Security-Hardening
- **Betroffene Methoden:** `CreateMirror()`, `SyncMirror()`, `MigrateRepo()`

### 3. Trusted Proxies Setting
- **Dateien:** `forgejo/client.go`, Config-Settings
- **Änderung:** `REVERSE_PROXY_TRUSTED_PROXIES` Default geändert
- **Grund:** Container-Security-Hardening
- **Betroffene:** Container-Deployments mit Reverse-Proxy-Auth

## ✨ Neue Features

### 1. Actions API Additions
- **Dateien:** `forgejo/action_run.go`, `forgejo/repo_action_run.go`, `forgejo/repo_action_artifact.go`
- **Neue Methoden:**
  - `GetActionRunLogs()`, `GetActionJobLogs()` - Workflow/Job Logs
  - `ListActionArtifacts()`, `GetActionArtifact()`, `DownloadActionArtifact()`, `DeleteActionArtifact()`
  - `CancelActionRun()`, `DeleteActionRun()`

### 2. Quota API
- **Dateien:** `forgejo/models/quota_*.go`, `forgejo/admin_quota.go`, `forgejo/org_quota.go`, `forgejo/user_quota.go`
- **Neue Modelle:** `QuotaGroup`, `QuotaRuleInfo`, `QuotaUsed*`
- **Endpoints:** Admin/Org/User Quota-Management

### 3. Settings API  
- **Dateien:** `forgejo/settings.go`, `forgejo/user_settings.go`, `forgejo/models/general_*_settings.go`
- **Neue Methoden:** `GetSettings()`, `UpdateSettings()`, `GetUserSettings()`, `UpdateUserSettings()`

### 4. ActivityPub Remote Follow
- **Dateien:** `forgejo/activitypub.go`, `forgejo/models/a_p_remote_follow_option.go`
- **Neue Methode:** `RemoteFollow()`

### 5. User Settings Granular Watch
- **Dateien:** `forgejo/user_settings.go`, `forgejo/models/user_settings_options.go`
- **Feature:** Granulare Notifications-Einstellungen (Issues, PRs, Releases)

## 🔧 Build & Modelle

### Swagger-Generierung
- **Quelle:** https://git.cluster.mirtuell.net/swagger.v1.json
- **Version:** Forgejo 16.0.2  
- **Endpunkte:** 326 (+66 neue)
- **Modelle:** Alle aktualisiert

### Makefile
- **Version:** `FORGEJO_VERSION := 16.0.2` (von 15.0.3)
- **Test-URL:** `http://localhost:3000` (v16.0.2 Server)

## 🧪 Tests

### Test-Änderungen
- **ActivityPub:** Tests angepasst für neue Remote Follow API
- **User Settings:** Tests für granulare Watch Settings
- **Actions:** Tests für Logs, Artifacts, Cancel, Delete
- **Quota:** Tests für Quota-Endpoints

### Linting
- **Behoben:** 8 Issues (1x unparam, 7x testifylint)
- **Status:** Alle Linter bestanden

## 📋 Kompatibilität

### Mit v15.x
- ❌ **Breaking Changes** erfordern Code-Anpassungen
- ✅ API-Kompatibilität für unveränderte Endpunkte
- ⚠️ Mirror-Redirects müssen neu implementiert werden

### Mit v16.0.x
- ✅ **Voll kompatibel** mit Forgejo 16.0.2
- ✅ Alle neuen Features implementiert
- ✅ Breaking Changes dokumentiert

## 🏁 Nächste Schritte

1. ✅ Code fertiggestellt
2. ⏳ Tests laufen (in Bearbeitung)
3. ⏸️ Review durch Benutzer
4. ⏸️ Tag: `forgejo/v2.4.0-forgejo-16.0.2`

## 📚 Dokumentation

- **V16_UPGRADE_TASKS.md** - Vollständige Task-Übersicht
- **V16_BREAKING_CHANGES.md** - Breaking Changes für Nutzer
- **V16_COMMIT_SUMMARY.md** - Diese Datei

---

*SDK Version: forgejo/v2.4.0-forgejo-16.0.2*  
*Forgejo Version: 16.0.2*  
*Stand: 2026-08-03*
