# Live DAG State Updates

This DAG viewer now supports real-time state updates via WebSocket connection.

## Architecture

```
Task Updates → Unix Socket → State Service → WebSocket → Browser
     ↓             ↓             ↓              ↓           ↓
dag-update    /tmp/socket   dag-state.json   port 8080   Live DAG
```

## Quick Start

### 1. Start the State Service

```bash
node dag-state-service.js
```

This starts:
- Unix socket listener at `/tmp/dag-state.sock`
- WebSocket server on port 8080
- HTTP API on port 8080

### 2. Generate Live DAG Visualization

```bash
node src/cli.js --dir ../../../ --renderer live -o output/live-dag.html
open output/live-dag.html
```

### 3. Update Task States

Now you can update task states in real-time using ANY of these methods:

#### Method A: Command Line Tool
```bash
# Update single task
./dag-update task1 completed
./dag-update task2 failed
./dag-update task3 in_progress

# Batch update
./dag-update batch '{"task1":"completed","task2":"started","task3":"blocked"}'

# Query current state
./dag-update query
```

#### Method B: Simple curl (HTTP API)
```bash
# YES! You can just curl it!
curl -X POST http://localhost:8080/update/task1/completed
curl -X POST http://localhost:8080/update/task2/failed
curl -X POST http://localhost:8080/update/task3/in_progress
```

#### Method C: Unix Socket (from your scripts)
```bash
echo '{"type":"update","taskId":"task1","state":"completed"}' | nc -U /tmp/dag-state.sock
```

## Task States

The following states are supported:
- `pending` - Gray (task not started)
- `started` - Light blue (task just started)
- `in_progress` - Blue (task actively running)
- `failed` - Red (task failed)
- `blocked` - Orange (task blocked by dependencies)
- `completed` - Green (task successfully completed)

## Features

- **Real-time Updates**: Changes appear instantly in the browser
- **No Polling**: Uses WebSocket for efficient real-time communication
- **Persistent State**: State is saved to `dag-state.json`
- **Reconnection**: Browser automatically reconnects if connection drops
- **Visual Feedback**: Nodes animate when state changes
- **Connection Status**: Live indicator shows connection state
- **State Statistics**: Header shows count of tasks in each state

## Environment Variables

```bash
# State service
WS_PORT=8080                    # WebSocket/HTTP port
UNIX_SOCKET=/tmp/dag-state.sock # Unix socket path
STATE_FILE=./dag-state.json     # State persistence file

# HTML generator
dag-viewer --ws-port 8080 --ws-host localhost
```

## Integration with T.A.S.K.S

Your task runner can update the DAG visualization by sending updates to the Unix socket or HTTP endpoint whenever a task changes state:

```javascript
// Example: Update from Node.js
import { createConnection } from 'net';

function updateTaskState(taskId, state) {
  const client = createConnection('/tmp/dag-state.sock');
  client.write(JSON.stringify({
    type: 'update',
    taskId: taskId,
    state: state
  }) + '\n');
  client.end();
}

// Or use HTTP
fetch('http://localhost:8080/update/task1/completed', { method: 'POST' });
```

## macOS Compatibility

Everything works perfectly on macOS:
- Unix domain sockets are fully supported
- No Docker required
- Native Node.js performance
- File watching works reliably

## Troubleshooting

If the WebSocket connection fails:
1. Check that `dag-state-service.js` is running
2. Verify port 8080 is not in use: `lsof -i :8080`
3. Check browser console for connection errors
4. Ensure no firewall is blocking local connections