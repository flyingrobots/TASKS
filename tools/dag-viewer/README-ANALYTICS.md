# Full Analytics System - Tasks + Git

## What Gets Tracked

### 📊 **Task Metrics**
- Tasks completed/failed/started per agent
- Median & average completion times
- Success rates
- Task paths through DAG

### 🔧 **Git Metrics**
- Commits per agent
- Lines added/removed (code churn)
- Files changed
- Commit sizes
- Commit frequency patterns

### 🔥 **Hot Files**
- Most frequently changed files
- File type distribution
- Code ownership by agent

### ⏰ **Activity Patterns**
- Hourly commit distribution (when are agents most active?)
- Daily activity trends
- Time since last activity

### 🏆 **Productivity Score** (Composite metric)
- 40 points: Task completion
- 30 points: Code contributions
- 20 points: Success rate
- 10 points: Recent activity

### 📈 **Cool Insights You Get**

1. **Tasks per Commit Ratio**
   - Are agents completing tasks without committing code?
   - Or committing without completing tasks?

2. **Lines per Task**
   - How much code does each task typically require?
   - Who writes the most efficient code?

3. **Code Churn**
   - Total lines added + removed
   - Net contribution (added - removed)
   - Who's refactoring vs adding new code?

4. **Hot Files & Bottlenecks**
   - Which files change most often? (potential design issues)
   - File type distribution (frontend vs backend work)

5. **Agent Patterns**
   - Who works on what parts of the codebase?
   - Commit message patterns
   - Work session duration

## Usage

### Start the server:
```bash
node dag-state-server-full.js
```

### Send task updates:
```bash
curl -X POST http://localhost:8080/event \
  -H "Content-Type: application/json" \
  -d '{
    "event": "task_update",
    "data": {
      "agent": "mr_clean",
      "task": "P1.T003",
      "status": "completed",
      "timestamp": 1755633601
    }
  }'
```

### Send git commits:
```bash
curl -X POST http://localhost:8080/event \
  -H "Content-Type: application/json" \
  -d '{
    "event": "git_commit",
    "agent": "mr_clean",
    "data": {
      "sha": "abc123def",
      "message": "Fix authentication bug in OAuth flow",
      "branch": "feature/auth-fix",
      "timestamp": 1755633701,
      "files": [
        {
          "path": "src/auth/oauth.js",
          "additions": 45,
          "deletions": 12
        },
        {
          "path": "tests/auth.test.js",
          "additions": 78,
          "deletions": 5
        }
      ]
    }
  }'
```

### Query endpoints:
```bash
# Full leaderboard with all metrics
curl http://localhost:8080/leaderboard

# Git insights (hot files, activity patterns, etc)
curl http://localhost:8080/git-insights

# Recent events (last 100)
curl http://localhost:8080/events

# Complete state dump
curl http://localhost:8080/state
```

## Leaderboard Response Example:
```json
{
  "name": "mr_clean",
  "productivityScore": 95,
  "tasksCompleted": 42,
  "commits": 18,
  "successRate": "95.2",
  "avgTaskTime": 234000,
  "medianTaskTime": 180000,
  "linesAdded": 3421,
  "linesRemoved": 892,
  "netLines": 2529,
  "avgCommitSize": "234",
  "codeChurn": 4313,
  "filesChangedCount": 67,
  "tasksPerCommit": "2.33",
  "linesPerTask": "81",
  "lastCommit": {
    "sha": "abc123",
    "timestamp": 1755633701,
    "message": "Fix auth bug"
  }
}
```

## Git Insights Response:
```json
{
  "summary": {
    "totalCommits": 342,
    "totalLinesAdded": 45632,
    "totalLinesRemoved": 12893,
    "netLines": 32739,
    "avgCommitSize": "178"
  },
  "hotFiles": [
    {"path": "src/api/handler.js", "changes": 47},
    {"path": "src/auth/oauth.js", "changes": 31}
  ],
  "hourlyActivity": [0,0,0,0,1,3,8,15,22,28,35,42,38,31,25,18,12,8,5,2,1,0,0,0],
  "topContributors": [
    {"agent": "mr_clean", "commits": 89},
    {"agent": "freddy", "commits": 67}
  ],
  "fileTypes": {
    "js": 234,
    "test": 89,
    "json": 45,
    "md": 23
  }
}
```

## Why This Is Cool

1. **Identify Patterns**: See when agents are most productive
2. **Spot Bottlenecks**: Hot files might need refactoring
3. **Measure Real Output**: Not just task completion, but actual code
4. **Team Dynamics**: Who collaborates on what files
5. **Quality Metrics**: Success rate + test coverage changes
6. **Refactoring Detection**: High churn with low net lines = refactoring
7. **Work Session Analysis**: Commit clusters show work sessions

## Integration Ideas

Your task orchestrator could:
- Auto-capture git commits when tasks complete
- Track which commits relate to which tasks
- Build a complete audit trail
- Identify "ghost work" (commits without tasks)
- Find "incomplete work" (tasks without commits)
- Calculate velocity trends over time
- Predict task completion based on historical data