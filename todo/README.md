# Project Roadmap (Hubless-Style, Git-Backed)

This `todo/` directory is a markdown-first, Git-native planner for T.A.S.K.S. + S.L.A.P.S. It complements the v8 spec by tracking milestones, features, and tasks as files you can move through a simple CLI (`npm run todo:*`).

- Milestones live in `todo/milestones/`
- Features live in `todo/features/`
- Tasks live in `todo/tasks/{milestone}/{backlog|active|finished|merged}/`

Use the scripts to move tasks across states; progress bars and indexes update automatically.

## Milestones

- M1 – Planner Foundation (canonical JSON, schemas, census)
- M2 – Plan Compiler (parse → tasks → deps → DAG → waves)
- M3 – Validators (acceptance, evidence, interface)
- M4 – Contract + Runtime Stub (coordinator + basic frontier)

<!-- PROGRESS:START ROADMAP -->
- M1 – 4/4 (100%) [backlog:0 active:0 finished:4 merged:0]
- M2 – 5/5 (100%) [backlog:0 active:0 finished:5 merged:0]
- M3 – 0/1 (0%) [backlog:1 active:0 finished:0 merged:0]
- M4 – 0/2 (0%) [backlog:2 active:0 finished:0 merged:0]
- M5 – 0/2 (0%) [backlog:2 active:0 finished:0 merged:0]
- M6 – 0/2 (0%) [backlog:2 active:0 finished:0 merged:0]
- M7 – 0/2 (0%) [backlog:2 active:0 finished:0 merged:0]
- M8 – 1/2 (50%) [backlog:1 active:0 finished:1 merged:0]
- M9 – 0/2 (0%) [backlog:2 active:0 finished:0 merged:0]
<!-- PROGRESS:END ROADMAP -->

## Usage

## Prerequisites
- Node.js 14+ and npm
- Run `npm install` from repository root before using the commands below

- Set a task active:
  - `npm run todo:task:set-active -- T001`
- Finish a task:
  - `npm run todo:task:set-finished -- T001`
- Merge/land a task:
  - `npm run todo:task:set-merged -- T001`

The scripts move the markdown file between status folders, update its frontmatter `status`, and refresh progress bars here and in each milestone README.

## Conventions

- Task file frontmatter:
  - `id`, `milestone`, `features`, `title`, `status`, `deps`
- Keep titles verb-first and bounded (2–8 hours).
- Link tasks to features where applicable.
