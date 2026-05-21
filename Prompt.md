# Development

## Do I continue

* Pull changes.
* Check the project status in the `status.yaml`
    * If status is `blocked`, `error`, `stopped`, or `paused` - exit, you cannot continue working at the moment.
    * If status is `active` move on to `Establish Context`.
    * If the status is `working`, evaluate the `updated` field - if the updated timestamp is greater than 10 minutes, the task got stuck. Log the fact you will attempt to recover the stuck task. `Establish Context` and resume work.
    * If the status is `working` but less than 30 minutes, another process might be working - exit. Do not proceed since this might overwrite any changes.
* Quote all values in the status file.
* If the model chosen has problem - log the problem.

## Establish Context

* Read the `AGENTS.md`
* Read the `MEMORY.md`
* Read the `TODO.md`

## Update Status

* Change status to `working`, record agent name, and current project phase
* Choose the next available task, mark the item with `W` to note the task currently being worked.
* Confirm you can work the task, if not, choose another task - if possible.
* If no tasks are available to work because human intervention is required, mark status as `blocked`
* If unable to continue because there are no more tasks, mark status as `paused`
* Add, commit, and push changes.

## Perform Work

* Do the work as prescribed by `AGENTS.md` and `AGENT_WORKFLOW.md`
* Mark task complete in the `TODO.md` with an `X`
* Update status back to `active`
* Update `MEMORY.md` of any decisions, milestones reached, or other important details for next pass.
* Sort the tasks in the `TODO.md` in order of priority
* Commit and push changes.

## Exit

Only complete one task at a time. Exit.