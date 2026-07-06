# Forgejo SDK – Upgrade auf v15 (Analyse & Roadmap)

Stand: 2026-07-06 · Zielinstanz: **v15.0.3+gitea-1.22.0** (`https://git.cluster.mirtuell.net`)

## 1. Durchgeführt (Modell-Ebene)

| Schritt | Ergebnis |
|---|---|
| `swagger.v1.json` ausgetauscht | 13.0.3 → **15.0.3** (828 KB) |
| `make swagger-generate` | go-swagger v0.33.1, Modelle regeneriert |
| `go mod tidy` | openapi-Pakete (validate/strfmt) vorhanden |
| `go build ./...` / `go vet ./...` | **grün (Exit 0)** – bestehender Client-Code kompiliert gegen neue Modelle |

> Die Makefile generiert **nur die Modelle** (`models/`), **nicht** die Client-Methoden. Die fehlenden Endpoint-Methoden (s. §3) müssen handgeschrieben werden.

### Neue Modelle (5 – neue v15-Konzepte)
- `action_runner.go` – erweitertes/einheitliches Runner-Modell
- `register_runner_options.go` – Eingabe `POST /{level}/actions/runners` (Name, `ephemeral`, …)
- `register_runner_response.go` – Antwort (`id`, `uuid`, `token`) der interaktiven Runner-Registrierung
- `repo_target_option.go` – **repo-spezifische Access-Tokens** (PR #11696)
- `forge_outbox.go` – ActivityPub Outbox

### Geänderte Modelle (15 – Schema-Additionen, +183 Zeilen netto)
`access_token.go`, `create_access_token_option.go` (Token-Scopes/Repo-Target) · `action_run_job.go` · `comment.go`, `create/edit_issue_comment_option.go` · `commit_status_state.go` (neue States) · `create_or_update_secret_option.go`, `create/update_variable_option.go` (Actions) · `create/edit_team_option.go`, `team.go` · `edit_user_option.go` · `repository.go`

## 2. Zu prüfende Breaking-/Orphan-Routen (3)

Pfad-basierter Diff – **vermutlich False Positives**, manuell verifizieren:
- `/repos/{owner}/{repo}/archive/{a}{b}` und `.../raw/{a}{b}` – SDK baut Pfad via `Sprintf` mit zwei Args; nur Matching-Artefakt.
- `/user/actions/secrets` – Route evtl. verschoben/umbenannt (v15 Runner-API wurde refactored).

## 3. Endpoint-Backlog (201 Operationen / 138 Pfade fehlen im SDK)

Diff: Swagger-Pfade normalisiert (`{x}`→`{}`) gegen SDK-Routen (`%s`→`{}`). **Hinweis:** pfadbasiert, enthält False Positives wo der SDK denselben Pfad mit Query-Params ruft – vor Implementierung gegen bestehende Client-Methoden prüfen.

### `repository` (44 Pfade) – Highlights
Runners: `GET/POST /repos/{o}/{r}/actions/runners` (getRepoRunners, **registerRepoRunner**), `GET/DELETE .../runners/{runner_id}` · `GET .../actions/runners/jobs` · Wiki: `POST .../wiki/new`, `GET/DELETE/PATCH .../wiki/page/{pageName}` · Tag-Protection: `GET/POST .../tag_protections`, `.../tag_protections/{id}` · Sync-Fork: `GET/POST .../sync_fork`, `.../sync_fork/{branch}` · `GET .../activities/feeds`, `.../subscribers`, `.../times/{user}`

### `admin` (24 Pfade)
Runners: `GET/POST /admin/actions/runners` (**registerAdminRunner**), `GET .../runners/jobs`, `GET .../runners/registration-token`, `GET/DELETE .../runners/{runner_id}` · **Quota** (komplett neu): `/admin/quota/groups`, `/admin/quota/groups/{group}/{rules|users}`, `/admin/quota/rules` · `GET /admin/cron`, `/admin/emails`, `/admin/emails/search` · Hooks `GET/POST /admin/hooks`, `.../hooks/{id}` · `/admin/unadopted`, `/admin/unadopted/{o}/{r}` · `/admin/users/{username}/{emails,quota,rename}` · `GET /admin/orgs`

### `organization` (23 Pfade)
Runners: `GET/POST /orgs/{org}/actions/runners` (**registerOrgRunner**) · **Quota**: `/orgs/{org}/quota`, `/quota/{artifacts,attachments,packages,check}` · Block: `PUT /orgs/{org}/block/{username}`, `/unblock/...`, `GET .../list_blocked` · Labels `/orgs/{org}/labels(/{id})` · `POST /orgs/{org}/rename` · Avatar `POST/DELETE .../avatar` · `/orgs/{org}/activities/feeds` · Teams `/teams/{id}/{members,repos,activities/feeds}`, `/orgs/{org}/teams/search` · `/user/orgs`, `/users/{username}/orgs`

### `user` (24 Pfade)
Runners: `GET/POST /user/actions/runners` (**registerUserRunner**), `.../runners/{runner_id}`, `.../registration-token` · **Quota**: `/user/quota`, `/quota/{artifacts,attachments,packages,check}` · Block: `/user/block/{username}`, `/unblock/...`, `GET .../list_blocked` · Avatar `POST/DELETE /user/avatar` · GPG: `GET /user/gpg_key_token`, `POST /user/gpg_key_verify` · `/user/teams`, `/user/followers`, `/user/following` · `/users/{username}/{activities/feeds,followers,following,gpg_keys,heatmap,keys,repos}`

### `issue` (17 Pfade) – zu ergänzen
(issue-Stopwatch/Subscriptions/Reactions z.T. vorhanden; Diff verifizieren)

### `miscellaneous` (12 Pfade)
`GET /signing-key.gpg`, `/signing-key.ssh` · `GET /nodeinfo` · `POST /markdown`, `/markdown/raw`, `/markup` · Templates: `/gitignore/templates(/{name})`, `/label/templates(/{name})`, `/licenses(/{name})`

### `activitypub` (11 Pfade)
`/activitypub/actor(/inbox/outbox)` · `/activitypub/repository-id/{id}(/inbox/outbox)` · `/activitypub/user-id/{id}(/inbox/outbox/activities/...)`

### `package` (3 Pfade)
`/packages/{owner}?...`, `/packages/{owner}/{type}/{name}/{version}(/files)`

## 4. Empfohlene Reihenfolge

1. **Runner v15** (interaktiv + ephemeral) – Admin/Org/Repo/User `register*Runner`, `get/delete*Runner`, Jobs. Nutzt neue Modelle `register_runner_options/response`, `action_runner`.
2. **Quota** – vollständiger CRUD (Admin Groups/Rules + User/Org Check & Artifacts).
3. **Block/Unblock** + `list_blocked` (User/Org).
4. **Repo-spezifische Tokens** – `CreateAccessTokenOption` um Repo-Target ergänzen (Modell `repo_target_option.go` vorhanden).
5. **Wiki, Tag-Protections, Sync-Fork, Activity Feeds**.
6. **Misc** (signing-key, nodeinfo, markdown/markup, Templates) + **ActivityPub** + **Packages**.
7. Nach Bedarf: `FORGEJO_VERSION` in `Makefile` (13.0.3 → 15.0.3) und README „Version Requirements" anpassen (betrifft `make test-instance`).

## Quellen
- Release-Notes: `codeberg.org/forgejo/forgejo/src/branch/forgejo/release-notes-published/{14.0.0,15.0.0}.md`
- Ankündigungen: `forgejo.org/2026-01-release-v14-0/`, `forgejo.org/2026-04-release-v15-0/`
- Instanz-Spec: `swagger.v1.json` (v15.0.3)
