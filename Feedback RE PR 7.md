- [x] In examples/dot-export/README.md around line 5, the heading "## Files" is
missing a preceding blank line; insert a single blank line above that heading so
there is an empty line separating it from the previous paragraph or content.

---

- [x] In examples/dot-export/tasks.json around lines 1 to 49 the example data is
unrealistic and repetitive: all three tasks share the same feature_id, identical
duration estimates, and empty acceptance_checks, source_evidence, and resources.
Update the JSON to present varied, plausible tasks: assign distinct feature_id
values (or group related tasks under the same feature only where reasonable),
vary optimistic/mostLikely/pessimistic durations and durationUnits to reflect
realistic effort, and add representative acceptance_checks (clear pass/fail
criteria), sample source_evidence entries (links, file names, or commit IDs),
and non-empty resources where applicable (roles or tooling). Also add at least
one dependency between tasks to demonstrate ordering and a small
resource_conflicts example to show conflict handling.

---

- [x] Makefile: Please fix
[Makefile](https://github.com/flyingrobots/TASKS/pull/7/files/e72c0b5a01ad83df929503ce62c0f94fce9fba0e#diff-76ed074a9305c04054cdebb9e9aad2d818052b07091de1f20cad0bbac34ffb52)

Comment on lines +1 to +14

|   |
|---|
|.PHONY: setup setup-dry-run setup-win example-dot|
||
|setup:|
|@bash scripts/setup-deps.sh --install --yes|
||
|setup-dry-run:|
|@bash scripts/setup-deps.sh|
||
|setup-win:|
|@pwsh -File scripts/setup-deps.ps1 -Install -Yes|
||
|example-dot:|
|@./examples/dot-export/render.sh|
||

### 

![@coderabbitai](https://avatars.githubusercontent.com/in/347564?s=48&v=4) **[coderabbitai](https://github.com/apps/coderabbitai) bot** [5 hours ago](https://github.com/flyingrobots/TASKS/pull/7#discussion_r2415090076)

_🧹 Nitpick_ | _🔵 Trivial_

**Where are your standard targets?**

You have a Makefile with custom targets but no `all`, `clean`, or `test` targets. These are CONVENTIONS for a reason. People expect them to exist.

At minimum add:

- `all`: default target that probably runs tests
- `clean`: removes build artifacts
- `test`: runs your test suite

The fact that `checkmake` is screaming about this should have been a clue.

Add standard targets:

```diff
 .PHONY: setup setup-dry-run setup-win example-dot
 
+all: test
+
+test:
+	@cd planner && go test ./...
+
+clean:
+	@cd planner && go clean
+
 setup:
 	@bash scripts/setup-deps.sh --install --yes
```

📝 Committable suggestion

> ‼️ **IMPORTANT**  
> Carefully review the code before committing. Ensure that it accurately replaces the highlighted code, contains no missing lines, and has no issues with indentation. Thoroughly test & benchmark the code to ensure it meets the requirements.

Suggested change

|   |
|---|
|.PHONY: setup setup-dry-run setup-win example-dot|
||
|setup:|
|@bash scripts/setup-deps.sh --install --yes|
||
|setup-dry-run:|
|@bash scripts/setup-deps.sh|
||
|setup-win:|
|@pwsh -File scripts/setup-deps.ps1 -Install -Yes|
||
|example-dot:|
|@./examples/dot-export/render.sh|
|.PHONY: setup setup-dry-run setup-win example-dot|
||
|all: test|
||
|test:|
|@cd planner && go test ./...|
||
|clean:|
|@cd planner && go clean|
||
|setup:|
|@bash scripts/setup-deps.sh --install --yes|
||
|setup-dry-run:|
|@bash scripts/setup-deps.sh|
||
|setup-win:|
|@pwsh -File scripts/setup-deps.ps1 -Install -Yes|
||
|example-dot:|
|@./examples/dot-export/render.sh|

---

- [x] In planner/internal/analysis/census.go around lines 48 to 57, remove the
os.ReadDir workaround that attempts to read each directory to surface permission
errors (the if d.IsDir() { _, statErr := os.ReadDir(currentPath) ... } block)
because filepath.WalkDir already returns permission errors and this
per-directory read adds unnecessary overhead; replace it by relying on WalkDir's
error propagation (or, if you must keep it, replace it with a short
justification comment that includes the platform/Go version and a link to the Go
issue or bug report that necessitates the workaround).

---

- [x] In planner/internal/analysis/census.go around line 67, remove the leftover
commented-out debug log line (`// t.Logf("DEBUG: RunCensus returning error from
WalkDir: %v", err) // DEBUG PRINT`) so no test/debug logging remains in
production code; delete the entire commented statement and ensure no trailing
whitespace or extraneous comments remain on that line.

---

- [x] In planner/internal/analysis/census.go at line 45 there's a commented-out test
log call referencing a test variable `t` (// t.Logf("DEBUG: WalkDir callback
error at %s: %v", currentPath, err)), which is test/debug code left in
production; remove that commented line entirely. If runtime logging is required
replace it with a real logger call (e.g., logger.Debugf or logger.Errorf) using
the existing logging infrastructure, otherwise just delete the commented debug
statement.

---

- [x] [planner/internal/emitter/writer.go](https://github.com/flyingrobots/TASKS/pull/7/files/e72c0b5a01ad83df929503ce62c0f94fce9fba0e#diff-b0a6ed47a0ffae0a6b8805ba05edf5f94bffc93167d069c2a29363ea2e1ceb43)

Comment on lines +12 to +14

|   |
|---|
|// WriteWithArtifactHash writes v to path as canonical JSON, computing meta.artifact_hash over the|
|// canonical bytes with the hash field present but empty. The setHash callback should embed the computed|
|// hash value into v (e.g., set v.Meta.ArtifactHash = h) before the final write.|

### 

![@coderabbitai](https://avatars.githubusercontent.com/in/347564?s=48&v=4) **[coderabbitai](https://github.com/apps/coderabbitai) bot** [5 hours ago](https://github.com/flyingrobots/TASKS/pull/7#discussion_r2415090087)

_⚠️ Potential issue_ | _🟠 Major_

**Document the thread safety (or lack thereof)!**

Is `v` safe to mutate concurrently? Can multiple goroutines call this with the same `v`? You don't document any of this. If `v` is shared and someone else is reading it while you're mutating it via setHash, you've got a race.

---

- [x] [planner/internal/emitter/writer.go](https://github.com/flyingrobots/TASKS/pull/7/files/e72c0b5a01ad83df929503ce62c0f94fce9fba0e#diff-b0a6ed47a0ffae0a6b8805ba05edf5f94bffc93167d069c2a29363ea2e1ceb43)

Comment on lines +15 to +28

|   |
|---|
|func WriteWithArtifactHash(path string, v any, setHash func(h string)) (string, error) {|
|raw, err := json.Marshal(v)|
|if err != nil { return "", fmt.Errorf("marshal %s: %w", path, err) }|
|can, err := canonjson.ToCanonicalJSON(raw)|
|if err != nil { return "", fmt.Errorf("canon %s: %w", path, err) }|
|h := hash.HashCanonicalBytes(can)|
|if setHash != nil { setHash(h) }|
|raw2, err := json.Marshal(v)|
|if err != nil { return "", fmt.Errorf("marshal2 %s: %w", path, err) }|
|can2, err := canonjson.ToCanonicalJSON(raw2)|
|if err != nil { return "", fmt.Errorf("canon2 %s: %w", path, err) }|
|if err := os.WriteFile(path, can2, 0o644); err != nil { return "", fmt.Errorf("write %s: %w", path, err) }|
|return h, nil|
|}|

### 

![@coderabbitai](https://avatars.githubusercontent.com/in/347564?s=48&v=4) **[coderabbitai](https://github.com/apps/coderabbitai) bot** [5 hours ago](https://github.com/flyingrobots/TASKS/pull/7#discussion_r2415090089)

_⚠️ Potential issue_ | _🔴 Critical_

**This function is a goddamn mess!**

You compute the hash over the FIRST canonicalization (line 20), then let some random callback mutate the object (line 21), then serialize AGAIN (line 22-24) and write THAT. So the hash you computed doesn't match what you wrote!

This is completely broken logic. The hash should match the bytes you write to disk, not some earlier version.

Also, if `setHash` is nil, you're doing marshal/canonicalize TWICE for no reason (lines 22-25 are completely wasted).

And where's the goddamn panic recovery around the callback? What if setHash panics? Your whole program dies.

Fix this disaster:

```diff
 func WriteWithArtifactHash(path string, v any, setHash func(h string)) (string, error) {
     raw, err := json.Marshal(v)
     if err != nil { return "", fmt.Errorf("marshal %s: %w", path, err) }
     can, err := canonjson.ToCanonicalJSON(raw)
     if err != nil { return "", fmt.Errorf("canon %s: %w", path, err) }
     h := hash.HashCanonicalBytes(can)
-    if setHash != nil { setHash(h) }
-    raw2, err := json.Marshal(v)
-    if err != nil { return "", fmt.Errorf("marshal2 %s: %w", path, err) }
-    can2, err := canonjson.ToCanonicalJSON(raw2)
-    if err != nil { return "", fmt.Errorf("canon2 %s: %w", path, err) }
-    if err := os.WriteFile(path, can2, 0o644); err != nil { return "", fmt.Errorf("write %s: %w", path, err) }
+    
+    if setHash != nil {
+        defer func() {
+            if r := recover(); r != nil {
+                err = fmt.Errorf("setHash panicked: %v", r)
+            }
+        }()
+        setHash(h)
+        // Re-marshal after mutation
+        raw, err = json.Marshal(v)
+        if err != nil { return "", fmt.Errorf("marshal after setHash %s: %w", path, err) }
+        can, err = canonjson.ToCanonicalJSON(raw)
+        if err != nil { return "", fmt.Errorf("canon after setHash %s: %w", path, err) }
+    }
+    
+    if err := os.WriteFile(path, can, 0o644); err != nil { return "", fmt.Errorf("write %s: %w", path, err) }
     return h, nil
 }
```

Actually, WAIT. Looking at the function comment on lines 12-14, it says you're INTENTIONALLY computing the hash with the field empty, then setting it. But that means the hash in the written file doesn't include itself, which makes sense for avoiding circular references.

So the logic is: hash(v with empty hash field) → set hash field → write(v with populated hash field)

That makes more sense, but it's STILL broken because:

1. You're not verifying the hash field is actually empty in the first marshal
2. No panic recovery
3. Wasteful double-marshal when setHash is nil

Document this weird workflow properly!

📝 Committable suggestion

> ‼️ **IMPORTANT**  
> Carefully review the code before committing. Ensure that it accurately replaces the highlighted code, contains no missing lines, and has no issues with indentation. Thoroughly test & benchmark the code to ensure it meets the requirements.

Suggested change

|   |
|---|
|func WriteWithArtifactHash(path string, v any, setHash func(h string)) (string, error) {|
|raw, err := json.Marshal(v)|
|if err != nil { return "", fmt.Errorf("marshal %s: %w", path, err) }|
|can, err := canonjson.ToCanonicalJSON(raw)|
|if err != nil { return "", fmt.Errorf("canon %s: %w", path, err) }|
|h := hash.HashCanonicalBytes(can)|
|if setHash != nil { setHash(h) }|
|raw2, err := json.Marshal(v)|
|if err != nil { return "", fmt.Errorf("marshal2 %s: %w", path, err) }|
|can2, err := canonjson.ToCanonicalJSON(raw2)|
|if err != nil { return "", fmt.Errorf("canon2 %s: %w", path, err) }|
|if err := os.WriteFile(path, can2, 0o644); err != nil { return "", fmt.Errorf("write %s: %w", path, err) }|
|return h, nil|
|}|
|func WriteWithArtifactHash(path string, v any, setHash func(h string)) (string, error) {|
|raw, err := json.Marshal(v)|
|if err != nil {|
|return "", fmt.Errorf("marshal %s: %w", path, err)|
|}|
|can, err := canonjson.ToCanonicalJSON(raw)|
|if err != nil {|
|return "", fmt.Errorf("canon %s: %w", path, err)|
|}|
|h := hash.HashCanonicalBytes(can)|
||
|if setHash != nil {|
|// Catch panics in the user-supplied callback|
|defer func() {|
|if r := recover(); r != nil {|
|err = fmt.Errorf("setHash panicked: %v", r)|
|}|
|}()|
|setHash(h)|
||
|// Re-marshal after the callback may have mutated v|
|raw, err = json.Marshal(v)|
|if err != nil {|
|return "", fmt.Errorf("marshal after setHash %s: %w", path, err)|
|}|
|can, err = canonjson.ToCanonicalJSON(raw)|
|if err != nil {|
|return "", fmt.Errorf("canon after setHash %s: %w", path, err)|
|}|
|}|
||
|if err := os.WriteFile(path, can, 0o644); err != nil {|
|return "", fmt.Errorf("write %s: %w", path, err)|
|}|
|return h, nil|
|}|

---

- [ ] In planner/internal/planner/dag/dag.go around lines 22 to 39, the loop that
builds idx and order currently overwrites earlier entries when task IDs
duplicate; detect duplicate IDs and fail loudly instead of silently overwriting:
before assigning idx[t.ID] = i, check if t.ID already exists in idx and return
an error (use fmt.Errorf with the duplicate ID and index details) or otherwise
stop processing so the caller knows the DAG is invalid; add the necessary import
for fmt.

---

- [ ] [planner/internal/planner/dag/dag.go](https://github.com/flyingrobots/TASKS/pull/7/files/e72c0b5a01ad83df929503ce62c0f94fce9fba0e#diff-d4fa045af807dc1b19fbf522ae0dfdb2526e581079fda2552787b8225247d13b)

Comment on lines +94 to +116

|   |
|---|
|// Longest path (by edges count) and predecessor|
|dist := make([]int, n)|
|pred := make([]int, n)|
|for i := range pred { pred[i] = -1 }|
|for _, u := range topo {|
|for _, v := range adj[u] {|
|if dist[v] < dist[u]+1 {|
|dist[v] = dist[u] + 1|
|pred[v] = u|
|}|
|}|
|}|
|// find sink on critical path|
|sink := 0|
|for i := 1; i < n; i++ { if dist[i] > dist[sink] { sink = i } }|
|critPath := []string{}|
|for x := sink; x != -1; x = pred[x] {|
|critPath = append(critPath, tasks[x].ID)|
|}|
|// reverse critPath|
|for i, j := 0, len(critPath)-1; i < j; i, j = i+1, j-1 { critPath[i], critPath[j] = critPath[j], critPath[i] }|
||
|// Transitive reduction: remove (u->v) if there exists u->w->...->v|

### 

![@coderabbitai](https://avatars.githubusercontent.com/in/347564?s=48&v=4) **[coderabbitai](https://github.com/apps/coderabbitai) bot** [5 hours ago](https://github.com/flyingrobots/TASKS/pull/7#discussion_r2415090114)

_⚠️ Potential issue_ | _🔴 Critical_

**HANDLING EMPTY INPUT BY CRASHING? SERIOUSLY?**

Feed this thing an empty task list and it detonates on `tasks[sink]` / `pred[sink]`. That’s a straight-up panic path. Guard it and fail with a real error instead of nuking the process.

```diff
@@
-   // build adjacency
-   n := len(tasks)
+   // build adjacency
+   n := len(tasks)
+   if n == 0 {
+       df.Analysis.OK = false
+       df.Analysis.Errors = append(df.Analysis.Errors, "no tasks to build DAG")
+       return df, errors.New("no tasks to build DAG")
+   }
```

Stop leaving crash grenades in public APIs.

📝 Committable suggestion

> ‼️ **IMPORTANT**  
> Carefully review the code before committing. Ensure that it accurately replaces the highlighted code, contains no missing lines, and has no issues with indentation. Thoroughly test & benchmark the code to ensure it meets the requirements.

Suggested change

|                                                                                                                |
| -------------------------------------------------------------------------------------------------------------- |
| // Longest path (by edges count) and predecessor                                                               |
| dist := make([]int, n)                                                                                         |
| pred := make([]int, n)                                                                                         |
| for i := range pred { pred[i] = -1 }                                                                           |
| for _, u := range topo {                                                                                       |
| for _, v := range adj[u] {                                                                                     |
| if dist[v] < dist[u]+1 {                                                                                       |
| dist[v] = dist[u] + 1                                                                                          |
| pred[v] = u                                                                                                    |
| }                                                                                                              |
| }                                                                                                              |
| }                                                                                                              |
| // find sink on critical path                                                                                  |
| sink := 0                                                                                                      |
| for i := 1; i < n; i++ { if dist[i] > dist[sink] { sink = i } }                                                |
| critPath := []string{}                                                                                         |
| for x := sink; x != -1; x = pred[x] {                                                                          |
| critPath = append(critPath, tasks[x].ID)                                                                       |
| }                                                                                                              |
| // reverse critPath                                                                                            |
| for i, j := 0, len(critPath)-1; i < j; i, j = i+1, j-1 { critPath[i], critPath[j] = critPath[j], critPath[i] } |
|                                                                                                                |
| // Transitive reduction: remove (u->v) if there exists u->w->...->v                                            |
| // build adjacency                                                                                             |
| n := len(tasks)                                                                                                |
| if n == 0 {                                                                                                    |
| df.Analysis.OK = false                                                                                         |
| df.Analysis.Errors = append(df.Analysis.Errors, "no tasks to build DAG")                                       |
| return df, errors.New("no tasks to build DAG")                                                                 |

---

- [x] In planner/internal/planner/docparse/docparse.go around lines 61-70, the code
currently ignores json.Unmarshal errors when parsing an `accept` block; instead,
check and handle the error: capture the error returned by json.Unmarshal for
both array and single-object paths, and either return that error from the
parsing function (preferred) or attach a parse-error to the containing task
(e.g., tasks[lastTaskIdx].ParseErrors = append(...)) so the failure is surfaced;
include context (payload string and task index or source location) in the error
message to make debugging possible.

---

- [x] In planner/internal/validate/schemas/coordinator.schema.json around lines 8 to
15, the "graph" object schema lacks an additionalProperties constraint so
arbitrary extra keys are permitted; update the schema to set
"additionalProperties": false on the "graph" object and (optionally) tighten
"nodes" and "edges" by specifying item schemas (e.g., arrays of objects with
required fields) to prevent random garbage from being accepted by the validator.

---

- [x] In planner/internal/validate/schemas/coordinator.schema.json around line 7, the
"version" property currently allows empty strings; update the schema to require
non-empty strings by adding a minLength constraint (e.g., "version":
{"type":"string","minLength":1}) so empty versions are rejected.

---

- [x] In planner/internal/validate/schemas/coordinator.schema.json around lines 12–13
the "nodes" and "edges" properties are declared as arrays with no item schema,
which permits any garbage; replace each with a proper item schema: for "nodes"
use "type":"array","items":{ "type":"object", "required":["id","type"],
"properties":{ "id":{"type":"string"}, "type":{"type":"string"},
"metadata":{"type":"object"} }, "additionalProperties":false } (add
uniqueItems:true if node ids must be unique); for "edges" use
"type":"array","items":{ "type":"object", "required":["source","target"],
"properties":{ "source":{"type":"string"}, "target":{"type":"string"},
"label":{"type":"string"}, "weight":{"type":"number"} },
"additionalProperties":false } and adjust required/fields to match your domain
model so validation enforces structure and types rather than allowing arbitrary
values.

---

- [x] In planner/internal/validate/schemas/dag.schema.json around lines 20 to 28, the
node object schema omits an additionalProperties control so arbitrary fields are
currently allowed; add "additionalProperties": false to the node schema to
forbid unknown fields (or if allowing extras is intentional, add a brief
comment/doc entry explaining that behavior), and ensure the schema still
validates required properties and constraints after this change.

---

- [x] In planner/internal/validate/schemas/dag.schema.json around lines 32 to 41, the
edge item schema lacks an additionalProperties policy; add
"additionalProperties": false to the "items" object (or to the inner object
schema) to prevent unexpected fields (or set to true if your project policy is
to allow extras) and ensure this same policy is applied consistently across all
schemas.

---

- [x] In planner/internal/validate/schemas/dag.schema.json around line 43, the
"metrics": {"type":"object"} entry is too permissive (accepts any object,
including {}). Replace it with a concrete schema: either reference an existing
metrics schema (e.g., "$ref": "#/definitions/metrics") or define explicit
properties, types and required keys for metrics, and set "additionalProperties":
false to prevent arbitrary keys; if metrics may be optional allow null or an
empty array with a proper oneOf/anyOf. Ensure the new schema enforces the
expected metric fields and types rather than accepting any object.

---

- [x] [planner/internal/validate/schemas/features.schema.json](https://github.com/flyingrobots/TASKS/pull/7/files/e72c0b5a01ad83df929503ce62c0f94fce9fba0e#diff-d1d5e061cb54b9df1cbab334cabb6656b369ab24e20c3732185f60fa28e2bca9)

Comment on lines +1 to +29

|   |
|---|
|{|
|"$schema": "https://json-schema.org/draft/2020-12/schema",|
|"title": "features.json",|
|"type": "object",|
|"required": ["meta", "features"],|
|"properties": {|
|"meta": {|
|"type": "object",|
|"required": ["version", "artifact_hash"],|
|"properties": {|
|"version": {"type": "string"},|
|"artifact_hash": {"type": "string"}|
|}|
|},|
|"features": {|
|"type": "array",|
|"minItems": 1,|
|"items": {|
|"type": "object",|
|"required": ["id", "title"],|
|"properties": {|
|"id": {"type": "string"},|
|"title": {"type": "string"}|
|}|
|}|
|}|
|},|
|"additionalProperties": true|
|}|

### 

![@coderabbitai](https://avatars.githubusercontent.com/in/347564?s=48&v=4) **[coderabbitai](https://github.com/apps/coderabbitai) bot** [5 hours ago](https://github.com/flyingrobots/TASKS/pull/7#discussion_r2415090141)

_🛠️ Refactor suggestion_ | _🟠 Major_

**Schema is too permissive. Did you even think about validation?**

Your schema has `"additionalProperties": true` at the top level (line 28), which means any random garbage can be added to the JSON and it'll pass validation. Is that what you want?

Also, your `meta` and `features` objects have no validation on `additionalProperties` for their nested objects either. The feature items can contain anything beyond `id` and `title`. If that's intentional, fine, but it smells like lazy schema design.

At minimum, document WHY you're allowing arbitrary properties. Otherwise, tighten this up.

Consider this diff if you want to restrict additional properties:

```diff
       }
     }
   },
-  "additionalProperties": true
+  "additionalProperties": false
 }
```

And for feature items:

```diff
         "properties": {
           "id": {"type": "string"},
           "title": {"type": "string"}
-        }
+        },
+        "additionalProperties": false
       }
```

📝 Committable suggestion

> ‼️ **IMPORTANT**  
> Carefully review the code before committing. Ensure that it accurately replaces the highlighted code, contains no missing lines, and has no issues with indentation. Thoroughly test & benchmark the code to ensure it meets the requirements.

Suggested change

|   |
|---|
|{|
|"$schema": "https://json-schema.org/draft/2020-12/schema",|
|"title": "features.json",|
|"type": "object",|
|"required": ["meta", "features"],|
|"properties": {|
|"meta": {|
|"type": "object",|
|"required": ["version", "artifact_hash"],|
|"properties": {|
|"version": {"type": "string"},|
|"artifact_hash": {"type": "string"}|
|}|
|},|
|"features": {|
|"type": "array",|
|"minItems": 1,|
|"items": {|
|"type": "object",|
|"required": ["id", "title"],|
|"properties": {|
|"id": {"type": "string"},|
|"title": {"type": "string"}|
|}|
|}|
|}|
|},|
|"additionalProperties": true|
|}|
|{|
|"$schema": "https://json-schema.org/draft/2020-12/schema",|
|"title": "features.json",|
|"type": "object",|
|"required": ["meta", "features"],|
|"properties": {|
|"meta": {|
|"type": "object",|
|"required": ["version", "artifact_hash"],|
|"properties": {|
|"version": { "type": "string" },|
|"artifact_hash": { "type": "string" }|
|}|
|},|
|"features": {|
|"type": "array",|
|"minItems": 1,|
|"items": {|
|"type": "object",|
|"required": ["id", "title"],|
|"properties": {|
|"id": { "type": "string" },|
|"title": { "type": "string" }|
|},|
|"additionalProperties": false|
|}|
|}|
|},|
|"additionalProperties": false|
|}|

---

- [ ] In planner/internal/validate/schemas/waves.schema.json around lines 11 to 13,
the string fields "version", "planId", and "artifact_hash" currently accept
empty strings; update their schemas to enforce non-empty values by adding
"minLength": 1 for "version" and "planId", and for "artifact_hash" require a
SHA-256 hex pattern (64 hex chars) plus "minLength": 64 and "maxLength": 64 or a
"pattern": "^[A-Fa-f0-9]{64}$" to strictly validate the hash.

---

- [ ] In planner/internal/validate/schemas/waves.schema.json around lines 16 to 22,
the current schema allows empty inner wave arrays (e.g. "waves": [[]]) which is
likely unintended; update the inner array schema to require at least one item by
adding "minItems": 1 to the inner "items" array definition (optionally also
consider adding "minItems": 1 to the outer "waves" array if an empty waves list
should be disallowed).

---

- [ ] In scripts/setup-deps.ps1 around lines 21 to 39, the Have() function currently
pipes Get-Command to Out-Null so it returns nothing, causing all if (Have '...')
checks to evaluate as $null; change Have() to return a proper boolean by testing
Get-Command's result and returning $true when the command exists and $false when
it does not (e.g., call Get-Command with -ErrorAction SilentlyContinue, check if
the result is non-$null, and return the corresponding boolean) so the downstream
if checks behave correctly.

---

- [ ] In planner/internal/validate/schemas/waves.schema.json around lines 7 to 15, the
"meta" object is missing an additionalProperties constraint which allows
arbitrary properties; update the "meta" schema to disallow unknown fields by
adding "additionalProperties": false alongside the existing "type", "required",
and "properties" so only version, planId, and artifact_hash are permitted.

---

- [ ] In package.json lines 1-10, the package is missing basic metadata (repository,
author, license); add these fields at the top level of package.json: a
"repository" object with "type" and "url" (pointing to the VCS remote), an
"author" string (name and optional email), and a "license" string (choose an
SPDX identifier such as MIT or a private/license file note). Optionally add
"homepage" and "bugs" fields for contact/issue tracking; ensure the added values
are valid JSON strings/objects and commit the updated package.json.

---

- [ ] In scripts/todo/task.js lines 1-143, replace the custom
parseFrontmatter/stringifyFrontmatter with gray-matter: install gray-matter (npm
i gray-matter and add to package.json), add const matter =
require('gray-matter') at the top, remove the parseFrontmatter and
stringifyFrontmatter functions, and in moveTask use const parsed =
matter(srcContent); update parsed.data.status = targetLane; const newContent =
matter.stringify(parsed.content, parsed.data); then write newContent to dst as
before; ensure any other direct calls to the removed functions are updated
similarly.

---

- [ ] scripts/todo/task.js around lines 29-46: the custom regex-based frontmatter
parser is fragile, swallows parse errors and tries ad-hoc JSON parsing; replace
it with a proper frontmatter/YAML parser (e.g. add gray-matter as a dependency
and use gray-matter(s) to extract data and content), remove the manual
line-splitting/regex/JSON.parse logic, and ensure parse failures are surfaced
(throw or log and return a failure) instead of silently ignoring errors; update
package.json to include the new dependency and adjust any callers to expect
parsed data from the library.

---

- [ ] In todo/README.md around lines 26 to 30, the command examples only show the
happy path; update the doc to include concrete examples for error cases and
idempotent operations: for each command (set-active, set-finished, set-merged)
add one example showing successful output, one showing the error when the task
ID does not exist (include exact error message and exit code), and one showing
the behavior when the task is already in the target state (idempotent case and
expected exit code or message). Also document any other common failure modes
(permission/config issues, invalid ID format) with sample outputs so readers can
reproduce and understand the CLI behavior. Ensure examples are verbatim console
output and keep formatting consistent with other README examples.

---

- [ ] In todo/README.md around line 38, the phrase "Keep titles verb-first and bounded
(2–8 hours)" asserts a 2–8 hour bound with no justification; either remove the
specific "2–8 hours" range or replace it with a justified guideline and
citation: add a brief rationale for the chosen timebox (e.g., tie it to our
sprint cadence, INVEST/Small-ticket guidance, or throughput goals) or give a
more flexible phrasing like "keep tasks small and completable within a workday"
plus a link or short note explaining how teams should choose the exact bound for
their context.

---

- [ ] In todo/milestones/M3-validators.md around lines 11 to 13, there's a heading
placed immediately after a list item without a separating blank line which
violates markdownlint MD022; fix it by inserting a single blank line between the
list and the heading (i.e., add an empty line after the list item before the
next heading).

---

- [ ] [scripts/todo/task.js](https://github.com/flyingrobots/TASKS/pull/7/files/a67fdf2f1572d06fdee3526dfb72bb0fadb6747e#diff-1b41df9efafa843a9016571bcf505af1e3c5cfd7fccd050caa49c013dd7edaa8)

Comment on lines +13 to +27

|   |
|---|
|function findTaskFile(taskId) {|
|const ms = fs.readdirSync(path.join(TODO, 'tasks'));|
|for (const m of ms) {|
|const base = path.join(TODO, 'tasks', m);|
|for (const lane of ['backlog','active','finished','merged']) {|
|const dir = path.join(base, lane);|
|if (!fs.existsSync(dir)) continue;|
|const files = fs.readdirSync(dir);|
|for (const f of files) {|
|if (f.startsWith(taskId + '-')) return { milestone: m, lane, file: path.join(dir, f) };|
|}|
|}|
|}|
|return null;|
|}|

### 

![@coderabbitai](https://avatars.githubusercontent.com/in/347564?s=48&v=4) **[coderabbitai](https://github.com/apps/coderabbitai) bot** [4 hours ago](https://github.com/flyingrobots/TASKS/pull/7#discussion_r2415127015)

_⚠️ Potential issue_ | _🟠 Major_

**This nested loop disaster will scale terribly.**

You're scanning every directory, every lane, every file on every single operation. This is O(n_m_k) complexity for something that should be O(1) with proper indexing. What happens when you have 100 milestones with 500 tasks each?

And you're doing synchronous readdirSync in nested loops! This will block the entire process.

Consider either:

1. Maintaining an index file that maps task IDs to file paths
2. Using a consistent naming convention that allows direct path construction
3. At minimum, making this async and adding early exits

---

- [ ] [scripts/todo/task.js](https://github.com/flyingrobots/TASKS/pull/7/files/a67fdf2f1572d06fdee3526dfb72bb0fadb6747e#diff-1b41df9efafa843a9016571bcf505af1e3c5cfd7fccd050caa49c013dd7edaa8)

Comment on lines +8 to +9

|   |
|---|
|const ROOT = process.cwd();|
|const TODO = path.join(ROOT, 'todo');|

### 

![@coderabbitai](https://avatars.githubusercontent.com/in/347564?s=48&v=4) **[coderabbitai](https://github.com/apps/coderabbitai) bot** [4 hours ago](https://github.com/flyingrobots/TASKS/pull/7#discussion_r2415127010)

_⚠️ Potential issue_ | _🟠 Major_

**You're using process.cwd() as the root without any validation.**

What happens when someone runs this from a subdirectory? Or from a completely different project? This will silently fail or create files in the wrong places. At minimum, validate you're in the right repository by checking for a marker file.

Apply this diff to add basic validation:

```diff
 const ROOT = process.cwd();
+if (!fs.existsSync(path.join(ROOT, 'package.json'))) {
+  die('Must run from repository root');
+}
+if (!fs.existsSync(path.join(ROOT, 'todo'))) {
+  die('todo/ directory not found - are you in the right repository?');
+}
 const TODO = path.join(ROOT, 'todo');
```

📝 Committable suggestion

> ‼️ **IMPORTANT**  
> Carefully review the code before committing. Ensure that it accurately replaces the highlighted code, contains no missing lines, and has no issues with indentation. Thoroughly test & benchmark the code to ensure it meets the requirements.

Suggested change

|   |
|---|
|const ROOT = process.cwd();|
|const TODO = path.join(ROOT, 'todo');|
|const ROOT = process.cwd();|
|if (!fs.existsSync(path.join(ROOT, 'package.json'))) {|
|die('Must run from repository root');|
|}|
|if (!fs.existsSync(path.join(ROOT, 'todo'))) {|
|die('todo/ directory not found - are you in the right repository?');|
|}|
|const TODO = path.join(ROOT, 'todo');|

---

- [ ] In docs/todo-workflow.md lines 1-56 you have duplicated content already present
in todo/README.md (specific overlaps around directory structure, scripts, and
frontmatter), so remove duplication by either: A) Replace the duplicated
sections with a short pointer to todo/README.md (e.g., a single sentence saying
“see todo/README.md for directory structure and commands”) and keep only the
unique "Hubless Integration" content here; or B) Delete docs/todo-workflow.md
entirely and move the unique "Hubless Integration" paragraph into todo/README.md
under a new “Hubless Integration” heading; or C) Consolidate both files into a
single source (use your docs generator/include mechanism) so that the canonical
content lives in one place. After choosing one option, remove the duplicated
blocks, update any internal links or cross-references, and run the docs
build/check to ensure there are no broken references.

---

- [ ] [scripts/todo/task.js](https://github.com/flyingrobots/TASKS/pull/7/files/9869b3d20d9d82507931dbccab6b2e339ba7495b#diff-1b41df9efafa843a9016571bcf505af1e3c5cfd7fccd050caa49c013dd7edaa8)

Comment on lines +48 to +51

|   |
|---|
|function stringifyFrontmatter(data, body) {|
|const lines = Object.keys(data).map(k => `${k}: ${Array.isArray(data[k]) ? JSON.stringify(data[k]) : data[k]}`);|
|return `---\n${lines.join('\n')}\n---\n\n${body}`;|
|}|

### 

![@coderabbitai](https://avatars.githubusercontent.com/in/347564?s=48&v=4) **[coderabbitai](https://github.com/apps/coderabbitai) bot** [2 hours ago](https://github.com/flyingrobots/TASKS/pull/7#discussion_r2415282130)

_⚠️ Potential issue_ | _🔴 Critical_

**stringifyFrontmatter is broken for any non-trivial data.**

Line 49: You only handle arrays with `Array.isArray()`, but what about:

- Nested objects? They'll be stringified as `[object Object]`
- Strings with colons, newlines, or special YAML characters? No escaping
- Numbers that look like strings? No type preservation
- Boolean values? They'll become strings

This will corrupt frontmatter the moment someone tries to store anything complex. Combined with your broken parser that was already flagged, you have a data corruption pipeline.

This is why the previous reviewer told you to use `gray-matter`. It handles serialization correctly. But if you insist on rolling your own:

```diff
 function stringifyFrontmatter(data, body) {
-  const lines = Object.keys(data).map(k => `${k}: ${Array.isArray(data[k]) ? JSON.stringify(data[k]) : data[k]}`);
+  const lines = Object.keys(data).map(k => {
+    const val = data[k];
+    if (Array.isArray(val) || typeof val === 'object') {
+      return `${k}: ${JSON.stringify(val)}`;
+    }
+    if (typeof val === 'string' && (val.includes(':') || val.includes('\n'))) {
+      return `${k}: "${val.replace(/"/g, '\\"')}"`;
+    }
+    return `${k}: ${val}`;
+  });
   return `---\n${lines.join('\n')}\n---\n\n${body}`;
 }
```

But seriously, just use `gray-matter` like you were already told.

📝 Committable suggestion

> ‼️ **IMPORTANT**  
> Carefully review the code before committing. Ensure that it accurately replaces the highlighted code, contains no missing lines, and has no issues with indentation. Thoroughly test & benchmark the code to ensure it meets the requirements.

Suggested change

|   |
|---|
|function stringifyFrontmatter(data, body) {|
|const lines = Object.keys(data).map(k => `${k}: ${Array.isArray(data[k]) ? JSON.stringify(data[k]) : data[k]}`);|
|return `---\n${lines.join('\n')}\n---\n\n${body}`;|
|}|
|function stringifyFrontmatter(data, body) {|
|const lines = Object.keys(data).map(k => {|
|const val = data[k];|
|if (Array.isArray(val) \| typeof val === 'object') {|
|return `${k}: ${JSON.stringify(val)}`;|
|}|
|if (typeof val === 'string' && (val.includes(':') \| val.includes('\n'))) {|
|return `${k}: "${val.replace(/"/g, '\\"')}"`;|
|}|
|return `${k}: ${val}`;|
|});|
|return `---\n${lines.join('\n')}\n---\n\n${body}`;|
|}|

---

- [x] In scripts/todo/task.js around lines 134-136, taskId is taken from user input
and used directly which allows path traversal; validate and whitelist the
expected taskId format (e.g. only alphanumeric, dashes, underscores) and reject
anything containing path separators or unexpected characters, returning an
error; additionally sanitize/normalize before use (or use path.basename) and
ensure moveTask only operates within the todo directory (resolve and verify the
resulting path is inside the allowed base) so arbitrary file moves outside the
todo tree are impossible.

---

- [x] In todo/milestones/M6-resilience.md around lines 11 to 13, the Markdown heading
"## Links" is immediately followed by a bullet list which violates markdownlint;
insert a single blank line between the "## Links" heading and the subsequent
bullet list so the heading is separated from the list and the linter error is
resolved.

--- 

- [x] In todo/milestones/M7-audit-admin.md around lines 11 to 13, the "## Links"
heading is immediately followed by a bullet list with no blank line; insert a
single empty line after the "## Links" heading so the markdown has a blank line
between the heading and the list (i.e., add one newline after line 11).

---

- [x] In todo/milestones/M8-security-visualization.md around lines 11 to 13 the
unordered list is placed immediately after the "## Links" heading which violates
Markdown convention/MD022; insert a single blank line between the "## Links"
heading and the list (i.e., add an empty line before the "- Features..." line)
so the heading and list are separated and the linter warning is resolved.

---

- [x] In todo/README.md lines 1–46, the README documents the CLI but omits
prerequisites; add a "## Prerequisites" section immediately before the "Usage"
header that lists "Node.js 14+ and npm" and instructs to "Run `npm install` from
repository root before using these commands" so users know to install Node/npm
and dependencies prior to running the npm scripts.

---

- [x] In planner/internal/canonjson/canonjson.go around lines 96-113 the code only
clears the sign for zero but leaves the exponent (so "0e10" canonicalizes
differently than "0"); change the zero-detection to treat magnitude-zero when
intPart == "0" and either there is no fractional part or the fractional part is
all zeros, and when that condition holds clear both sign and exp (i.e., set exp
= ""), before recomposing the string.

---

- [x] In planner/cmd/tasksd/main.go around lines 92–203, runExportDot mixes three
modes (directory, coordinator single-file, dag+tasks) and duplicates DOT
generation/writing logic; extract the DOT creation and file-write logic into
small helpers to reduce duplication and make testing easier. Create helpers such
as writeCoordinatorDot(coordPath, outPath, nodeLabel, edgeLabel) and
writeDagDot(dagPath, tasksPath, outPath, nodeLabel, edgeLabel) that load JSON,
build the title map (for tasks), call dot.FromCoordinatorWithOptions or
dot.FromDagWithOptions, and either return the dot string or perform the write
and return an error; add an emitDirArtifacts(dir, nodeLabel, edgeLabel) helper
that checks for coordinator.json and dag+tasks.json and calls the above helpers
for each file found. Refactor runExportDot to only parse flags, call the
appropriate helper for the chosen mode, and handle returned errors (avoid
duplicating loadJSON/writeFile/error messages across modes) so behavior remains
the same but code is concise and testable.

---

- [ ] In planner/cmd/tasksd/main.go around lines 225-478, runPlan is a huge function
mixing parsing, task construction, dependency inference, DAG building, and
artifact emission; it also uses unnecessary immediately-invoked anonymous
functions to initialize tasks (lines ~280-305). Split responsibilities by
extracting: (1) buildTasksFromDoc(doc, repo) ([]m.Task, []m.Edge, error) to
handle reading --doc, parsing features/tasks, creating tasks (avoid IIFE —
create structs then set fields), and return any doc-specified dependency edges;
(2) inferDependencies(tasks, existingEdges) []m.Edge to encapsulate linear
fallback and resource-conflict edge generation; (3) writeArtifacts(outDir
string, tf m.TasksFile, df m.DagFile, coord m.Coordinator, waves map[string]any)
error to write JSON/DOT files and compute hashes; wire runPlan to call these
helpers, propagate errors instead of os.Exit inside helpers (return errors to
runPlan), and remove all immediate anonymous-function initializations in favor
of plain struct literals followed by explicit field assignments for clarity and
testability.

---

- [ ] In todo/tasks/m2/finished/T060-generate-waves-preview.md around line 4, the
front matter uses the plural key "features: [F008]" but the repository schema
expects the singular "feature". Replace the plural key with "feature: F008" (or
"feature: \"F008\"" if quoting is used elsewhere) so the file matches the other
task files and tooling can parse it consistently.

---

- [ ] In docs/TEMPLATES.md around lines 1 to 90: several fenced code blocks (lines 7,
37, 41, 68, 72, 89) lack language specifiers and many fenced blocks and headings
lack required blank lines; add appropriate language tags (use markdown for
prose/code snippets, yaml for front-matter blocks, bash for shell/command
examples) to each fence and ensure there is an empty line before and after every
fenced code block and before each heading listed (lines 27, 32, 35, 43, 48, 51,
54, 57, 60, 63, 66, 73, 83, 86); keep front-matter fences as ```yaml and
acceptance/command fences as ```bash (or ```markdown if they are pure markdown),
and verify rendering after making these spacing and language-specifier fixes.

---

- [ ] In docs/todo-workflow.md around lines 31 to 59, the fenced code blocks lack
language specifiers and are missing surrounding blank lines which breaks
Markdown rendering; update the first code fence to ```yaml (for the front-matter
block) and add a blank line before and after it, then for each command example
replace plain ``` with ```bash and ensure there is a blank line before each
opening fence and after each closing fence, and also normalize heading spacing
(single blank line above headings) so the document follows consistent Markdown
hygiene.

---

- [ ] In scripts/todo/branch.js around lines 16 to 19, the script currently runs `git
checkout -B` which force-resets an existing branch; change this to use `git
checkout -b` and handle the case where the branch already exists: run the
checkout command and detect a non-zero exit or specific error output indicating
the branch exists, then print a clear error message and exit (or prompt for
confirmation) instead of force-overwriting; alternatively, before checkout, call
a git command to check if the branch exists and if so abort with a descriptive
message asking the user to delete/rename or confirm a force reset.

---

- [ ] In planner/internal/planner/wavesim/wavesim.go around line 30, the local type
definition "type set map[string]struct{}" is declared inside a function which
reduces reusability and is unusual; move this type declaration to package scope
(top of the file) so it can be reused across the package, or if it's only used
once replace occurrences with the inline type "map[string]struct{}"; update any
references accordingly and run `go vet`/`go test` to ensure no name collisions
or unused declarations.

---

- [ ] In planner/internal/planner/wavesim/wavesim.go around lines 12 to 60, the
function assumes every DAG node ID exists in the tasks slice and directly
indexes idToTask without presence checks; this yields zero-value Tasks when a
lookup fails and produces silent incorrect waves. Change the function to return
an error (e.g., Generate(df m.DagFile, tasks []m.Task) ([][]string, error)),
validate each lookup with the comma-ok idiom when doing t, ok := idToTask[id];
if !ok return nil, fmt.Errorf("missing task for node ID %q", id), and propagate
any other validation errors similarly; update the caller at main.go line 428 to
handle the returned error (check error, log/return it) and adjust tests
accordingly.

---

- [ ] planner/cmd/tasksd/main.go around lines 546-551: the function currently returns
any despite defining a concrete fileCensus type; change the signature to return
(fileCensus, error) (or *fileCensus if pointer preferred), return the concrete
type instead of any, and remove the local type if you promote it to package
scope. Then update the caller at line 319 to accept the concrete fileCensus
return type (adjust variable types and error handling accordingly) and replace
any usage of a generic analysis field in Meta with the fileCensus type (embed or
add a field of type fileCensus) so the code is fully typed end-to-end.

---

- [ ] In planner/cmd/tasksd/main.go around lines 432-436, the mustWriteJSONWithHash
closure aborts the process on the first write error (fmt.Fprintf + os.Exit),
which prevents attempting the remaining artifact writes and kills proper error
handling; change the pattern so you collect and return errors instead of
exiting: make the writer function return an error (and only set hashes[name] on
success), call it for each artifact, aggregate any errors (e.g., append to a
slice or use a multi-error), and after all writes handle/return a single
combined error to the caller for proper reporting and cleanup rather than
exiting inside the helper.

---

- [ ] In planner/cmd/tasksd/main.go around lines 308 to 314, the code silently ignores
errors from analysisPkg when a --repo is provided; change it to log the error
(including the repo value and the error) instead of swallowing it so the user
knows the flag had no effect; call the project’s logger (or log.Printf if no
logger exists) to record a concise message like "analysisPkg failed for repo %s:
%v" and then continue as before.

---

- [ ] In scripts/todo/task.js around line 63, the mkdirSync call with { recursive:
true } isn't error-handled so failures (path exists as a file, no write
permission, disk full, etc.) will cause the later write to throw an unhelpful
error; wrap the mkdirSync in a try-catch, and before creating the directory
check fs.existsSync(dstDir) and if true use fs.statSync to ensure it's a
directory (throw or log a clear error if it's a file), otherwise call
fs.mkdirSync inside try-catch and handle common errno cases (ENOTDIR/EEXIST ->
treat as fatal with a clear message, EACCES -> permission denied, ENOSPC -> disk
full), log or rethrow a descriptive error so the process fails with useful
diagnostics instead of crashing later.
