# Forgejo SDK v16.0 Breaking Changes

## ⚠️ Breaking Changes für SDK-Nutzer

Dieses Dokument beschreibt die Breaking Changes, die bei der Aktualisierung auf Forgejo SDK v16.0 zu beachten sind.

## 1. EXIF Stripping Removed

### Betroffene Methoden

Die folgenden Methoden haben **EXIF-Stripping-Optionen entfernt**:

#### `Client.UpdateAvatar()`
```go
// ALT (v15.x):
client.UpdateAvatar(user, UpdateAvatarOption{
    Avatar: imageData,
    StripEXIF: true,  // ⚠️ ENTFALLEN
})

// NEU (v16.0):
client.UpdateAvatar(user, UpdateAvatarOption{
    Avatar: imageData,
})
```

#### `Client.UploadAvatar()`
```go
// ALT (v15.x):
client.UploadAvatar(user, avatarFile, true)  // ⚠️ stripEXIF Parameter ENTFALLEN

// NEU (v16.0):
client.UploadAvatar(user, avatarFile)
```

### Migration Guide

1. Entfernen Sie alle `StripEXIF` oder `stripEXIF` Parameter
2. Wenn Sie EXIF-Stripping benötigen, implementieren Sie es client-seitig
3. Die Methoden werfen jetzt einen Fehler wenn EXIF-Optionen übergeben werden

### Grund

Forgejo v16.0 hat das EXIF-Stripping aufgrund von Lizenzproblemen (AGPL-Bibliothek) entfernt. Siehe [Forgejo PR #13105](https://codeberg.org/forgejo/forgejo/pulls/13105).

---

## 2. Mirror SSRF Hardening

### Betroffene Methoden

Die Mirror-API-Methoden **verhalten sich anders bei Redirects**:

#### `Client.CreateMirror()`
```go
// ALT (v15.x):
// Redirects wurden automatisch gefolgt
mirror, err := client.CreateMirror(repo, CreateMirrorOption{
    CloneAddr: "https://github.com/old/repo",  // Wenn redirect -> automatisch folgend
})

// NEU (v16.0):
// Redirects werden als Error behandelt
mirror, err := client.CreateMirror(repo, CreateMirrorOption{
    CloneAddr: "https://github.com/old/repo",  // Wenn redirect -> ERROR!
})

// Lösung: Verwende die finale URL:
mirror, err := client.CreateMirror(repo, CreateMirrorOption{
    CloneAddr: "https://github.com/new/repo",  // Direkte URL
})
```

#### `Client.SyncMirror()`
```go
// ALT (v15.x):
err := client.SyncMirror(mirrorID)  // Redirects -> automatisch folgend

// NEU (v16.0):
err := client.SyncMirror(mirrorID)  // Redirects -> ERROR
```

### Migration Guide

1. **Aktualisieren Sie alle Mirror-URLs** auf die finale Ziel-URL
2. **Testen Sie Sync-Mirror** auf renamed/transferred Repos
3. **Behandeln Sie Redirect-Errors** explizit:
   ```go
   mirror, err := client.CreateMirror(repo, opts)
   if err != nil {
       if strings.Contains(err.Error(), "redirect") {
           // URL aktualisieren
       }
   }
   ```

### Grund

SSRF-Security-Hardening: Git-Redirects konnten Security-Settings umgehen. Siehe [Forgejo PR #13129](https://codeberg.org/forgejo/forgejo/pulls/13129).

---

## 3. Trusted Proxies Setting (Container-Umgebungen)

### Betroffene Config-Settings

Wenn Sie **containerisierte Forgejo-Instanzen** betreiben:

```go
// ALT (v15.x - Container Default):
config := &ForgejoConfig{
    ReverseProxyTrustedProxies: "*",  // ⚠️ NICHT MEHR DEFAULT
}

// NEU (v16.0 - MUSS explizit gesetzt werden):
config := &ForgejoConfig{
    ReverseProxyTrustedProxies: "127.0.0.1/8,::1/128",  // Beispiel
}
```

### Migration Guide

1. **Überprüfen Sie Ihre Reverse-Proxy-Konfiguration**
2. **Setzen Sie `REVERSE_PROXY_TRUSTED_PROXIES` explizit** wenn:
   - `[service].ENABLE_REVERSE_PROXY_AUTHENTICATION = true`
   - Container-Umgebung (Docker/Podman)
3. **Verwenden Sie sichere Werte:**
   - Docker mit iptables: `"127.0.0.1/8,::1/128"`
   - Hinter bekannter Proxy: Proxy-IP-Range

### Grund

Security-Hardening: Verhindert Header-Injection wenn kein Reverse-Proxy verwendet wird.

---

## 4. Entfernte Endpunkte (Deprecated)

### Keine Endpunkte entfernt

Forgejo v16.0 hat **keine API-Endpunkte entfernt**, aber:
- EXIF-Optionen entfernt (siehe oben)
- Mirror-Redirect-Verhalten geändert (siehe oben)

---

## 🧪 Test-Checkliste

Bevor Sie auf v16.0 upgraden, testen Sie:

- [ ] Avatar-Upload **ohne** EXIF-Optionen
- [ ] Mirror-Create/Sync mit **direkten URLs** (keine Redirects)
- [ ] Mirror-Create mit **Redirects** (sollte Error werfen)
- [ ] Reverse-Proxy-Auth mit **korrektem `REVERSE_PROXY_TRUSTED_PROXIES`**
- [ ] Alle Actions-Features (Logs, Artifacts, Cancel)

---

## 📋 Upgrade-Steps

1. **Code-Audit:** Suchen Sie nach `StripEXIF`, `stripEXIF` Parameter
2. **Mirror-Audit:** Prüfen Sie alle Mirror-URLs auf Redirect-Ketten
3. **Config-Update:** Setzen Sie `REVERSE_PROXY_TRUSTED_PROXIES` wenn nötig
4. **Tests:** Führen Sie alle Tests gegen v16.0-Server durch
5. **Deploy:** Rollout mit Monitoring

---

## 🆘 Hilfe & Support

- **Forgejo Docs:** [v16.0 Upgrade Guide](https://forgejo.org/docs/v16.0/admin/upgrade/)
- **Forgejo Chat:** [Matrix #forgejo-chat](https://matrix.to/#/#forgejo-chat:matrix.org)
- **Issues:** [Forgejo Issues](https://codeberg.org/forgejo/forgejo/issues)

---

*SDK Version: forgejo/v2.4.0-forgejo-16.0.2*  
*Forgejo Version: 16.0.2*  
*Stand: 2026-08-03*
