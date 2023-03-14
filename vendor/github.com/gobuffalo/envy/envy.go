package envy

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
)

var gil = &sync.Mutex{}
var env = map[string]string{}

func init() {
	fmt.Println("ACA")
	// set the GOPATH if using >= 1.8 and the GOPATH isn't set
	if os.Getenv("GOPATH") == "" {
		out, err := exec.Command("go", "env", "GOPATH").Output()
		fmt.Println(" ACA 2")
		fmt.Println(out)
		if err == nil {
			fmt.Println(" ACA 3")
			gp := strings.TrimSpace(string(out))
			fmt.Println(gp)
			os.Setenv("GOPATH", gp)
		}
	}


	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		env[pair[0]] = os.Getenv(pair[0])
	}
}

// Get a value from the ENV. If it doesn't exist the
// default value will be returned.
func Get(key string, value string) string {
	gil.Lock()
	defer gil.Unlock()
	if v, ok := env[key]; ok {
		return v
	}
	return value
}

// Get a value from the ENV. If it doesn't exist
// an error will be returned
func MustGet(key string) (string, error) {
	gil.Lock()
	defer gil.Unlock()
	if v, ok := env[key]; ok {
		return v, nil
	}
	return "", fmt.Errorf("could not file ENV var with %s", key)
}

// Set a value into the ENV. This is NOT permanent. It will
// only affect values accessed through envy.
func Set(key string, value string) {
	gil.Lock()
	defer gil.Unlock()
	env[key] = value
}

// MustSet the value into the underlying ENV, as well as envy.
// This may return an error if there is a problem setting the
// underlying ENV value.
func MustSet(key string, value string) error {
	gil.Lock()
	defer gil.Unlock()
	err := os.Setenv(key, value)
	if err != nil {
		return err
	}
	env[key] = value
	return nil
}

// Map all of the keys/values set in envy.
func Map() map[string]string {
	gil.Lock()
	defer gil.Unlock()
	return env
}

// Temp makes a copy of the values and allows operation on
// those values temporarily during the run of the function.
// At the end of the function run the copy is discarded and
// the original values are replaced. This is useful for testing.
func Temp(f func()) {
	oenv := env
	env = map[string]string{}
	for k, v := range oenv {
		env[k] = v
	}
	defer func() { env = oenv }()
	f()
}

// GoPath returns the first GOPATH that is set
func GoPath() string {
	paths := GoPaths()
	if len(paths) > 0 {
		return paths[0]
	}
	return ""
}

// GoPaths returns all possible GOPATHS that are set.
func GoPaths() []string {
	gp := Get("GOPATH", "")
	if runtime.GOOS == "windows" {
		return strings.Split(gp, ";") // Windows uses a different separator
	}
	return strings.Split(gp, ":")
}
