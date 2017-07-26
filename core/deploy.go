package core

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
)

// execute ansible-playbook command && return the execute result
func ExecuteDeploy(name string) ([]byte, error) {
	var out []byte
	deploys, err := ioutil.ReadDir(DefaultCacheDir)
	if err != nil {
		return out, err
	}

	found := false
	for _, f := range deploys {
		if !f.IsDir() {
			continue
		}
		if f.Name() == name {
			found = true
			break
		}
	}

	if found {
		var outBuf bytes.Buffer
		targetFile := fmt.Sprintf("%s/%s/%s", DefaultCacheDir, name, DefaultYmlFile)

		f, err := os.OpenFile(targetFile+".status", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0655)
		if err != nil {
			return outBuf.Bytes(), err
		}
		mWriter := io.MultiWriter(f, os.Stdout, &outBuf)
		cmd := exec.Command("ansible-playbook", targetFile)
		cmd.Stdout = mWriter
		cmd.Stderr = mWriter

		err = cmd.Run()
		if err != nil {
			return outBuf.Bytes(), err
		}

		return outBuf.Bytes(), nil
	}

	return out, errors.New("no found target deploy")
}
