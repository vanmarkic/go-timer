package main

import (
	"context"
	"embed"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os/exec"
	"runtime"
	"time"
)

//go:embed web/*
var webFS embed.FS

func openBrowser(url string) {
	var cmd string
	var args []string
	switch runtime.GOOS {
	case "darwin":
		cmd = "open"
	case "windows":
		cmd = "rundll32"
		args = []string{"url.dll,FileProtocolHandler"}
	default:
		cmd = "xdg-open"
	}
	args = append(args, url)
	// #nosec G204 -- cmd is one of three hardcoded OS-specific binaries;
	// exec.Command does not invoke a shell, so url cannot be interpreted as a command.
	_ = exec.Command(cmd, args...).Start()
}

func handler() (http.Handler, error) {
	sub, err := fs.Sub(webFS, "web")
	if err != nil {
		return nil, fmt.Errorf("embedded web assets missing: %w", err)
	}
	return http.FileServer(http.FS(sub)), nil
}

func run(ctx context.Context, addr string, openUI bool) error {
	h, err := handler()
	if err != nil {
		return err
	}
	srv := &http.Server{
		Handler:           h,
		ReadHeaderTimeout: 5 * time.Second,
	}
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("listen %s: %w", addr, err)
	}
	url := fmt.Sprintf("http://%s", ln.Addr().String())
	fmt.Printf("go-timer running on %s (Ctrl+C to quit)\n", url)
	if openUI {
		openBrowser(url)
	}
	errc := make(chan error, 1)
	go func() {
		if err := srv.Serve(ln); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errc <- err
			return
		}
		errc <- nil
	}()
	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		_ = srv.Shutdown(shutdownCtx)
		return nil
	case err := <-errc:
		return err
	}
}

func main() {
	addr := flag.String("addr", "127.0.0.1:8080", "address to listen on")
	noOpen := flag.Bool("no-open", false, "do not open a browser window on startup")
	flag.Parse()
	if err := run(context.Background(), *addr, !*noOpen); err != nil {
		log.Fatal(err)
	}
}
