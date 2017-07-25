package core

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os/exec"
)

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
		targetFile := fmt.Sprintf("%s/%s", name, DefaultYmlFile)
		cmd := exec.Command("ansible-playbook", targetFile)
		var outBuf bytes.Buffer
		cmd.Stdout = &outBuf
		cmd.Stderr = &outBuf
		err = cmd.Run()
		if err != nil {
			return outBuf.Bytes(), err
		}
		return outBuf.Bytes(), nil
	}

	return out, errors.New("no found target deploy")
}
