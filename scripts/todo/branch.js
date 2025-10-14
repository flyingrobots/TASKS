#!/usr/bin/env node
const { spawnSync } = require('node:child_process');

function run(cmd, args, opts = {}) {
  const stdio = opts.stdio ?? 'pipe';
  const result = spawnSync(cmd, args, {
    stdio,
    encoding: 'utf8',
    ...opts,
  });
  if (result.status !== 0) {
    const err = new Error(`command failed: ${cmd} ${args.join(' ')} (exit ${result.status ?? 'unknown'})`);
    err.exitStatus = result.status ?? 1;
    if (result.error) {
      err.cause = result.error;
      err.message += `: ${result.error.message}`;
      err.code = result.error.code ?? err.code;
    }
    if (result.signal) {
      err.signal = result.signal;
    }
    if (stdio === 'pipe') {
      err.stdout = result.stdout?.toString();
      err.stderr = result.stderr?.toString();
    }
    throw err;
  }
  return result;
}

function main() {
  const [feature, taskId] = process.argv.slice(2);
  if (!feature || !taskId) {
    console.error('Usage: node scripts/todo/branch.js <feature-slug> <TASK_ID>');
    console.error('Example: node scripts/todo/branch.js validators T070');
    process.exit(1);
  }
  const branch = `feat/${feature}-task-${taskId}`;
  try {
    run('git', ['fetch', 'origin']);
  } catch (err) {
    console.error('git fetch failed:', err.stderr?.trim() || err.message);
    process.exit(err.exitStatus || 1);
  }

  try {
    run('git', ['rev-parse', '--verify', branch], { stdio: 'pipe' });
    console.error(`Branch ${branch} already exists. Aborting.`);
    process.exit(1);
  } catch (err) {
    if ((err.exitStatus ?? 0) !== 128) {
      console.error('Failed to check existing branches:', err.stderr?.trim() || err.message);
      process.exit(err.exitStatus || 1);
    }
  }

  try {
    run('git', ['checkout', '-b', branch, 'origin/main']);
  } catch (err) {
    console.error('git checkout failed:', err.stderr?.trim() || err.message);
    process.exit(err.exitStatus || 1);
  }
  console.log(`Switched to branch ${branch}`);
}

main();
