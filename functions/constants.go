package functions

// Message Constants
const (
	// GreetingMessage - displayed when a user connects
	GreetingMessage = "Welcome to TCP-Chat!\n         _nnnn_\n        dGGGGMMb\n       @p~qp~~qMb\n       M|@||@) M|\n       @,----.JM|\n      JS^\\__/  qKL\n     dZP        qKRb\n    dZP          qKKb\n   fZP            SMMb\n   HZM            MMMM\n   FqM            MMMM\n __| \".        |\\dS\"qML\n |    `.       | `' \\Zq\n_)      \\.___.,|     .'\n\\____   )MMMMMP|   .'\n     `-'       `--'\n[ENTER YOUR NAME]: "
	// [time][username]:
	UserPrompt = "[%v][%v]: "
	// [time][username]: message
	UserMessage  = "[%v][%s]: %s"
	UserJoined   = "%s has joined our chat...\n"
	UserLeft     = "%s has left our chat...\n"
)

// MessageModes - Enum for different chat modes (join, send, leave)
const (
	ModeJoined = iota // 0
	ModeSend // 1
	ModeLeft // 2
)

// TimeFormat - Format for displaying time in messages
const (
	TimeFormat = "2006-01-02 15:04:05"
)

// Color Patterns
const (
	ColorReset  = "\u001b[0m"
	ColorRed = "\u001b[31m"
	ColorYellow = "\u001b[33m"
)
