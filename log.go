package flamingo

import(
    "log"
    "os"
    "io"
)

var Log *log.Logger = log.New(os.Stdout, "", log.LstdFlags)

func SetLogWriter(w io.Writer) {
    Log = log.New(w, "", log.LstdFlags)
}