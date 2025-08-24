#//skip-test

# Construct a callable value for the external program `$program`. Example:
#
# ```elvish-transcript
# ~> var x = (external man)
# ~> $x ls # opens the manpage for ls
# ```
#
# See also [`has-external`]() and [`search-external`]().
fn external {|program| }

# Test whether `$command` names a valid external command. Examples (your output
# might differ):
#
# ```elvish-transcript
# ~> has-external cat
# ▶ $true
# ~> has-external lalala
# ▶ $false
# ```
#
# See also [`external`]() and [`search-external`]().
fn has-external {|command| }

# Output the full path of the external `$command`. Throws an exception when not
# found. Example (your output might vary):
#
# ```elvish-transcript
# ~> search-external cat
# ▶ /bin/cat
# ```
#
# See also [`external`]() and [`has-external`]().
fn search-external {|command| }

# Replace the Elvish process with an external `$command`, defaulting to
# `elvish`, passing the given arguments. This decrements `$E:SHLVL` before
# starting the new process.
#
# This command always raises an exception on Windows with the message "not
# supported on Windows".
fn exec {|command? @args| }

# Bring one or more background processes to the foreground by their process IDs.
# All PIDs must belong to the same process group.
#
# This function performs the following operations:
# - Verifies all processes are in the same process group
# - Sets the process group as the foreground process group of the terminal
# - Sends SIGCONT to resume the processes (in case they were stopped)
# - Waits for each process to complete or be stopped
#
# Examples:
#
# ```elvish-transcript
# ~> # Start a background job and note its PID
# ~> sleep 10 &
# ▶ 12345
# ~> fg 12345  # Bring the sleep process to foreground
# ```
#
# ```elvish-transcript
# ~> # Multiple processes in the same group can be brought to foreground together
# ~> bash -c 'sleep 20 & sleep 25 & wait' &
# ▶ 12346
# ~> fg 12346  # Brings the entire process group to foreground
# ```
#
# **Note**: This command always raises an exception on Windows with the message
# "not supported on Windows" as Windows doesn't have the same job control
# concepts as Unix-like systems.
#
# See also [`exec`]().
fn fg {|@pids| }

# Exit the Elvish process with `$status` (defaulting to 0).
fn exit {|status?| }
