# Manual backend testing

This guide uses PowerShell and `curl.exe` to exercise the current API from end to end. No frontend is required.

## What curl cannot do yet

`curl` can only send HTTP requests to routes the server exposes. LeeTKOB does not yet have endpoints for:

- Creating the first admin or student account. Use the existing `create-user` setup command.
- Creating a lab. Insert a local test lab into PostgreSQL with the SQL command below.
- Applying database migrations or starting the API. Run the setup commands below.
- Judging a submission or checking its later status. The worker, sandbox, and status endpoint are not implemented yet.

The login, logout, student enrollment, and submission actions **are** HTTP endpoints, so test those with curl. The setup commands exist to prepare the server and data those requests need.

## 1. Start PostgreSQL, migrate, and run the API

Configure `.env` as in the [README](../README.md). For local HTTP use `APP_ORIGIN=http://localhost:5173` and `COOKIE_SECURE=false`.

From the repository root in PowerShell, start the database and apply migrations:

```powershell
docker compose up -d
Get-Content .env | Where-Object { $_ -match '^[A-Za-z_][A-Za-z0-9_]*=' } | ForEach-Object { $name, $value = $_ -split '=', 2; [Environment]::SetEnvironmentVariable($name, $value, 'Process') }; Set-Location backend; go run ./cmd/migrate
```

Start the API in that same terminal:

```powershell
go run ./cmd/server
```

Leave the API running. Open a second PowerShell terminal at the repository root and load `.env` there too:

```powershell
Get-Content .env | Where-Object { $_ -match '^[A-Za-z_][A-Za-z0-9_]*=' } | ForEach-Object { $name, $value = $_ -split '=', 2; [Environment]::SetEnvironmentVariable($name, $value, 'Process') }
```

The Go server needs database, port, and origin settings before it starts; curl cannot configure or launch that server. Those are process settings, not HTTP request fields.

## 2. Create local accounts and a lab

These are bootstrap actions without corresponding HTTP endpoints. Run from the second terminal. Replace the sample passwords with local test passwords of at least 12 characters. The command prints each new account's UUID; save the admin and student IDs.

```powershell
Set-Location backend; $env:USER_EMAIL='admin@example.test'; $env:USER_PASSWORD='local-admin-password-123'; $env:USER_ROLE='admin'; go run ./cmd/create-user
```

```powershell
$env:USER_EMAIL='student@example.test'; $env:USER_PASSWORD='local-student-password-123'; $env:USER_ROLE='student'; go run ./cmd/create-user; Remove-Item Env:USER_EMAIL,Env:USER_PASSWORD,Env:USER_ROLE
```

Go to the repository root, create a lab UUID, and insert the lab using the local PostgreSQL container. This SQL creates only a new test lab:

```powershell
Set-Location ..; $labId=[guid]::NewGuid().ToString(); docker compose exec -T postgres psql -U $env:POSTGRES_USER -d $env:POSTGRES_DB -c "INSERT INTO labs (id,title) VALUES ('$labId','Manual test lab')"; $labId
```

Save the displayed lab ID. Use your actual admin/student UUIDs and passwords in the curl requests below. For a repeat run, use new email addresses and a new lab ID, or remove old test records using PostgreSQL.

## 3. Set request values

In PowerShell, set values for convenient copy/paste. These are client-side variables used to form curl arguments; credentials still go in the login request body.

```powershell
$api='http://localhost:8080'; $origin='http://localhost:5173'; $labId='<paste-lab-uuid>'; $studentId='<paste-student-uuid>'; $adminEmail='admin@example.test'; $adminPassword='local-admin-password-123'; $studentEmail='student@example.test'; $studentPassword='local-student-password-123'
```

If `API_PORT` or `APP_ORIGIN` differ, update `$api` or `$origin` to match. Use `curl.exe` explicitly because PowerShell may define `curl` as an alias.

## 4. Test login, enrollment, and logout

Log in as the admin. `-c admin.cookies` saves the session cookie for later requests:

```powershell
curl.exe -i -c admin.cookies -H "Origin: $origin" -H "Content-Type: application/json" --data-raw ('{"email":"'+$adminEmail+'","password":"'+$adminPassword+'"}') "$api/auth/login"
```

**Expected:** `204 No Content` and `Set-Cookie`.

Enroll the student. Admin authorization comes from the saved session cookie:

```powershell
curl.exe -i -b admin.cookies -c admin.cookies -H "Origin: $origin" -H "Content-Type: application/json" --data-raw ('{"studentId":"'+$studentId+'"}') "$api/labs/$labId/enrollments"
```

**Expected:** `204 No Content`. Repeating enrollment also returns `204`. Trying enrollment without the admin cookie returns `401`; trying with a student cookie returns `403`.

Log in as the student, saving a separate session:

```powershell
curl.exe -i -c student.cookies -H "Origin: $origin" -H "Content-Type: application/json" --data-raw ('{"email":"'+$studentEmail+'","password":"'+$studentPassword+'"}') "$api/auth/login"
```

**Expected:** `204 No Content` and a student session cookie. Incorrect credentials return `401`. An unapproved or missing `Origin` returns `403`.

## 5. Test submission and idempotency

Submit one C++ source file. Use a fresh key for each new intended submission:

```powershell
$key=[guid]::NewGuid().ToString(); $body='{"source":"int main() { return 0; }"}'; curl.exe -i -b student.cookies -H "Origin: $origin" -H 'Content-Type: application/json' -H "Idempotency-Key: $key" --data-raw $body "$api/labs/$labId/submissions"
```

**Expected:** `201 Created`, a submission UUID, and `"status":"queued"`. The endpoint saves the submission and pending job together. It does not compile or judge the code yet.

Retry the exact request with the same key and body:

```powershell
curl.exe -i -b student.cookies -H "Origin: $origin" -H 'Content-Type: application/json' -H "Idempotency-Key: $key" --data-raw $body "$api/labs/$labId/submissions"
```

**Expected:** `201`, the same submission response, and `Idempotent-Replayed: true`.

Reuse the key with a changed body:

```powershell
curl.exe -i -b student.cookies -H "Origin: $origin" -H 'Content-Type: application/json' -H "Idempotency-Key: $key" --data-raw '{"source":"int main() { return 1; }"}' "$api/labs/$labId/submissions"
```

**Expected:** `422` problem details with type `https://leetkob.example/problems/idempotency-key-reused`.

Other useful failure checks:

| curl variation | Expected |
|---|---|
| Remove `-b student.cookies` | `401` unauthenticated |
| Use an un-enrolled student's session | `404` not-found |
| Omit `Idempotency-Key` | `400` invalid-idempotency-key |
| Change content type to `text/plain` | `415` unsupported-media-type |
| Send malformed JSON | `400` invalid-json |
| Send `{"source":""}` | `400` invalid-request |

Errors use RFC 9457 problem details. See the [submission contract](contract/submissions.md) for the complete rules and error list.

## 6. Test logout

```powershell
curl.exe -i -b student.cookies -c student.cookies -H "Origin: $origin" -X POST "$api/auth/logout"
```

**Expected:** `204 No Content`. Repeat the submission request with `student.cookies`; it should return `401` because logout destroyed that session.
