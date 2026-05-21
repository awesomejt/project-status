# Development

## Do I continue

* Pull changes.
* If `status` field is `error`, `paused`, `stopped`, or `blocked` - log the status and exit.
* If the model chosen has problem - log the problem.

## Establish Context

* Read the `AGENTS.md`
* Read the `MEMORY.md`
* Read the `PROJECT_BRIEF.md` to understand the project context
* Read the `AGENT_WORKFLOW.md` to understand the approved workflow
* Read the `TODO.md`
* Read all documents in `./docs` folder

## Resume Work

* If current status is `working`, check the `TODO.md` for the task marked by `W` indicating the currently working task. Resume working on that task.
    * You can skip the `Save Status Change` since status is already in `working` state.
    * Move to the `Perform Work` section of this document.
* If no task is currently marked with a `W` or if the current status is `active`:
    * Move on to `Task Selection` section of this document.

## Task Selection

* Choose the next available task, mark the item with `W` to note the task currently being worked.
* Confirm you can work the task, if not, choose another task - if possible.
* If no tasks are available to work because human intervention is required, mark status as `blocked`
* If unable to continue because there are no more tasks, mark status as `paused`
* Add, commit, and push changes.

## Save Status Change

Update the `status.yaml`:

* Change status to `working`
* Record agent name
* current project phase
* update the `updated` field with the current system time

With Git: stage changed files, commit changes, and push to GitHub:

```bash
git add TODO.md  # if file was modified
git add status.yaml
git commit -m "Change status to working"
git push
```

## Perform Work

* Do the work as prescribed by `AGENTS.md` and `AGENT_WORKFLOW.md`
* Mark task complete in the `TODO.md` with an `X`
* Update status back to `active`
* Update `MEMORY.md` of any decisions, milestones reached, or other important details for next pass.
* Sort the tasks in the `TODO.md` in order of priority
* Git stage all changed files.
* Commit changes and push to GitHub.

## Exit

Only complete one task at a time. Exit.