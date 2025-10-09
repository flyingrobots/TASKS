#!/usr/bin/env node
// Simple task workflow utility: move a task file between status folders and update frontmatter/status
// Usage: node scripts/todo/task.js set-active T001

const fs = require('fs');
const path = require('path');

const ROOT = process.cwd();
const TODO = path.join(ROOT, 'todo');

function die(msg) { console.error(msg); process.exit(1); }

function findTaskFile(taskId) {
  const ms = fs.readdirSync(path.join(TODO, 'tasks'));
  for (const m of ms) {
    const base = path.join(TODO, 'tasks', m);
    for (const lane of ['backlog','active','finished','merged']) {
      const dir = path.join(base, lane);
      if (!fs.existsSync(dir)) continue;
      const files = fs.readdirSync(dir);
      for (const f of files) {
        if (f.startsWith(taskId + '-')) return { milestone: m, lane, file: path.join(dir, f) };
      }
    }
  }
  return null;
}

function parseFrontmatter(s) {
  if (!s.startsWith('---')) return { data: {}, body: s };
  const end = s.indexOf('\n---', 3);
  if (end === -1) return { data: {}, body: s };
  const fm = s.slice(3, end).trim();
  const body = s.slice(end + 4).replace(/^\s*\n/, '');
  const data = {};
  fm.split(/\r?\n/).forEach(line => {
    const m = line.match(/^([A-Za-z0-9_]+):\s*(.*)$/);
    if (m) {
      const key = m[1];
      let val = m[2].trim();
      if (val.startsWith('[')) { try { val = JSON.parse(val); } catch { /* noop */ } }
      data[key] = val;
    }
  });
  return { data, body };
}

function stringifyFrontmatter(data, body) {
  const lines = Object.keys(data).map(k => {
    const val = data[k];
    if (Array.isArray(val) || (val && typeof val === 'object')) {
      return `${k}: ${JSON.stringify(val)}`;
    }
    if (typeof val === 'string' && (val.includes(':') || val.includes('\n'))) {
      return `${k}: "${val.replace(/"/g, '\\"')}"`;
    }
    return `${k}: ${val}`;
  });
  return `---\n${lines.join('\n')}\n---\n\n${body}`;
}

function moveTask(taskId, targetLane) {
  const rec = findTaskFile(taskId);
  if (!rec) die(`Task ${taskId} not found`);
  if (rec.lane === targetLane) { console.log(`Task ${taskId} already in ${targetLane}`); return; }
  const srcContent = fs.readFileSync(rec.file, 'utf8');
  const { data, body } = parseFrontmatter(srcContent);
  data.status = targetLane;
  const newContent = stringifyFrontmatter(data, body);
  const basename = path.basename(rec.file);
  const dstDir = path.join(TODO, 'tasks', rec.milestone, targetLane);
  if (!fs.existsSync(dstDir)) fs.mkdirSync(dstDir, { recursive: true });
  const dst = path.join(dstDir, basename);
  fs.writeFileSync(dst, newContent);
  fs.unlinkSync(rec.file);
  console.log(`Moved ${taskId} -> ${targetLane}`);
}

function scanTasks() {
  const result = {};
  const ms = fs.readdirSync(path.join(TODO, 'tasks'));
  for (const m of ms) {
    const base = path.join(TODO, 'tasks', m);
    const lanes = {};
    for (const lane of ['backlog','active','finished','merged']) {
      const dir = path.join(base, lane);
      let count = 0;
      if (fs.existsSync(dir)) {
        for (const f of fs.readdirSync(dir)) if (f.endsWith('.md')) count++;
      }
      lanes[lane] = count;
    }
    result[m] = lanes;
  }
  return result;
}

function updateProgress() {
  const stats = scanTasks();
  // Update todo/README.md and each milestone file between markers
  function replaceBlock(file, marker, text) {
    let s = fs.readFileSync(file, 'utf8');
    const start = s.indexOf(`<!-- PROGRESS:START ${marker} -->`);
    const end = s.indexOf(`<!-- PROGRESS:END ${marker} -->`);
    if (start === -1 || end === -1) return;
    const head = s.slice(0, start + (`<!-- PROGRESS:START ${marker} -->`).length);
    const tail = s.slice(end);
    const mid = `\n${text}\n`;
    fs.writeFileSync(file, head + mid + tail);
  }
  // Root roadmap
  const lines = [];
  for (const m of Object.keys(stats).sort()) {
    const st = stats[m];
    const total = st.backlog + st.active + st.finished + st.merged;
    const done = st.finished + st.merged;
    const pct = total ? Math.round((done/total)*100) : 0;
    lines.push(`- ${m.toUpperCase()} – ${done}/${total} (${pct}%) [backlog:${st.backlog} active:${st.active} finished:${st.finished} merged:${st.merged}]`);
    // Milestone file
    const map = {
      m1: 'M1-planner-foundation.md',
      m2: 'M2-plan-compiler.md',
      m3: 'M3-validators.md',
      m4: 'M4-contract-runtime-stub.md',
      m5: 'M5-resources-throughput.md',
      m6: 'M6-resilience.md',
      m7: 'M7-audit-admin.md',
      m8: 'M8-security-visualization.md',
      m9: 'M9-ci-docs.md'
    };
    const file = map[m] ? path.join(TODO, 'milestones', map[m]) : null;
    if (file && fs.existsSync(file) && fs.statSync(file).isFile()) {
      replaceBlock(file, m.toUpperCase(), `${done}/${total} done (${pct}%)`);
    }
  }
  replaceBlock(path.join(TODO, 'README.md'), 'ROADMAP', lines.join('\n'));
}

function validTaskId(id) { return /^T\d{3}$/.test(id); }

function main() {
  const [cmd, taskId] = process.argv.slice(2);
  if (!cmd) die('Usage: todo task <set-active|set-finished|set-merged> <TASK_ID>');
  if (cmd.startsWith('set-')) {
    if (!taskId) die('Task ID required');
    if (!validTaskId(taskId)) die('Invalid task ID format. Expected T### (e.g., T070)');
    const lane = cmd.replace('set-','');
    moveTask(taskId, lane);
    updateProgress();
    return;
  }
  if (cmd === 'update-progress') { updateProgress(); return; }
  die('Unknown command');
}

main();
