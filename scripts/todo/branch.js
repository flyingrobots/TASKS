#!/usr/bin/env node
const { spawnSync } = require('node:child_process');

function run(cmd, args, opts = {}) {
  const result = spawnSync(cmd, args, {
    stdio: 'inherit',
    encoding: 'utf8',
    ...opts,
  });
  if (result.status !== 0) {
    const err = new Error(`command failed: ${cmd} ${args.join(' ')}`);
    err.result = result;
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
    console.error('git fetch failed:', err.result?.stderr || err.message);
    process.exit(err.result?.status || 1);
  }

  const existsCheck = spawnSync('git', ['rev-parse', '--verify', branch], {
    stdio: 'pipe',
    encoding: 'utf8',
  });
  if (existsCheck.status === 0) {
    console.error(`Branch ${branch} already exists. Aborting.`);
    process.exit(1);
  }
  if (existsCheck.error) {
    console.error('Failed to check existing branches:', existsCheck.error.message);
    process.exit(1);
  }

  const checkout = spawnSync('git', ['checkout', '-b', branch, 'origin/main'], {
    stdio: 'inherit',
    encoding: 'utf8',
  });
  if (checkout.status !== 0) {
    console.error('git checkout failed:', checkout.stderr || 'unknown error');
    process.exit(checkout.status || 1);
  }
  console.log(`Switched to branch ${branch}`);
}

main();
