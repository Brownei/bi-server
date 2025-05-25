package cmd

import (
	"fmt"
	"net/http"
	"os"
	"syscall"
)

const workerCount = 4

func runWorker() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello from Go Worker PID %d\n", os.Getpid())
	})
	fmt.Println("Go Worker PID", os.Getpid(), "listening on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error:", err)
	}
}

func something() {
	if os.Getenv("WORKER") == "1" {
		runWorker()
		return
	}

	for i := range workerCount {
		fmt.Print(i)
		env := append(os.Environ(), "WORKER=1")
		_, err := syscall.ForkExec(os.Args[0], os.Args, &syscall.ProcAttr{
			Env: env,
			Files: []uintptr{
				os.Stdin.Fd(),
				os.Stdout.Fd(),
				os.Stderr.Fd(),
			},
		})
		if err != nil {
			fmt.Printf("Failed to fork: %v and for range %d", err, i)
		}
	}

	select {} // master process blocks forever
}
