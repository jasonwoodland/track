package view

const (
	AddedProject                           = "Added project <magenta>%s</>\n"
	AddedProjectTaskDurationTotal          = "Added: <magenta>%s</> <blue>%s</> (%s, %s total)\n"
	AddedTask                              = "Added task <blue>%s</>\n"
	AlreadyRunningProjectTaskElapsedTotal  = "Already running: <magenta>%s</> <blue>%s</> (%s, %s total)\033[J\n"
	CancelledProjectTaskDurationTotal      = "Cancelled: <magenta>%s</> <blue>%s</> (%s, %s total)\n"
	ConfirmDeleteFrameTimeProjectTask      = "Delete frame <green>%s - %s</> on <magenta>%s</> <blue>%s</>?"
	ConfirmDeleteProject                   = "Delete project <magenta>%s</>?"
	ConfirmDeleteTaskFramesOnProject       = "Delete task <blue>%s</> and %d frame%s on project <magenta>%s</>?"
	ConfirmStopRunningTask                 = "Stop running task?"
	Deleted                                = "Delete"
	DeletedProject                         = "Deleted project <magenta>%s</>\n"
	FinishedAtTimeElapsed                  = "Finished at <green>%s</> (%s)\n"
	FrameDoesNotExistForProjectTask        = "Frame <gray>[%v]</> doesn't exist on <magenta>%s</> <blue>%s</>\n"
	FrameTimesDuration                     = "  <gray>[%v]</> <green>%s - %s</> %6s\n"
	FrameTimesDurationLog                  = "  <gray>[%v]</> <green>%s - %s</> %6s\n"
	FrameTimesDurationTask                 = "  <green>%s - %s</> %6s <blue>%-*s</>\n"
	DailyDateHours                         = "<green>%s</> %6s\n"
	DailyHoursProject                      = "  %5s <magenta>%s</>\n"
	DailyHoursTask                         = "  %5s   <blue>%-*s</>\n"
	ConfirmMoveFrameTimesFromToProjectTask = "Move frame <green>%s - %s</> from <magenta>%s</> <blue>%s</> to <magenta>%s</> <blue>%s</>?"
	Moved                                  = "Moved"
	NotRunning                             = "Not running"
	Project                                = "<magenta>%s</>\n"
	ProjectAlreadyExists                   = "Project <magenta>%s</> already exists\n"
	ProjectDoesNotExist                    = "Project <magenta>%s</> doesn't exist\n"
	ProjectHours                           = "<magenta>%s</> %.2fh\n"
	TotalHours                             = "Total: %.2fh\n"
	RenamedProject                         = "Renamed project <magenta>%s</> to <magenta>%s</>\n"
	RenamedTaskOnProject                   = "Renamed task <blue>%s</> to <blue>%s</> on project <magenta>%s</>\n"
	RunningProjectTaskElapsedTotal         = "Running: <magenta>%s</> <blue>%s</> (%s, %s total)\033[J\n"
	RunningProjectTaskPrevElapsedTotal     = "Running: <magenta>%s</> <blue>%s</> (%s -> %s, %s -> %s total)\033[J\n"
	RunningProjectTaskTotal                = "Running: <magenta>%s</> <blue>%s</> (%s)\033[J\n"
	StartedAtTime                          = "Started at <green>%s</>\n"
	StartedAtTimeElapsed                   = "Started at <green>%s</> (%s ago)\033[J\n"
	StartedAtPrevTimeElapsed               = "Started at <green>%s -> %s</> (%s -> %s ago)\033[J\n"
	StoppedProjectTaskElapsedTotal         = "Stopped: <magenta>%s</> <blue>%s</> (%s, %s total)\033[J\n"
	Task                                   = "  <blue>%s</>\n"
	TaskAlreadyExistsForProject            = "Task <blue>%s</> already exists on <magenta>%s</>\n"
	TaskDoesNotExistForProject             = "Task <blue>%s</> doesn't exist on <magenta>%s</>\n"
	ConfirmMergeFramesFromToProjectTask    = "Merge %d frame%s from <magenta>%s</> <blue>%s</> into <magenta>%s</> <blue>%s</>?"
	Merged                                 = "Merged"
)
