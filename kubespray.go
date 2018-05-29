package kubespray_go

import (
	"os"
	"fmt"
	"errors"
	"strings"
	"io/ioutil"
	"encoding/json"
	"os/user"
	"os/exec"
	"io"
)

type Kubespray struct {
	Verbose bool
	Addresses string
	AnsibleUser string
	AnsibleBecome bool
}

func NewKubespray(addresses string) (*Kubespray, error) {
	usr, err := user.Current()
	if err != nil {
		return nil, err
	}
	return &Kubespray{Verbose: false, Addresses: addresses, AnsibleUser: usr.Name , AnsibleBecome: false}, nil
}

func (kubespray *Kubespray) Provision() error {
	fmt.Println("Provisioning bare metal cluster...")

	if _, err := os.Stat("kubespray"); os.IsNotExist(err) {
		fmt.Println("Cannot find Kubespray. Downloading...")
		err := runCommand(true,"git", "clone", "git@github.com:kubernetes-incubator/kubespray.git")
		if err != nil {
			return err
		}
		fmt.Println(true, "Downloading Kubespray succeeded.")
	} else {
		fmt.Println(true, "Existing Kubespray installation found. Reusing it.")
	}

	var kubeApiserverAddress = ""
	if len(kubespray.Addresses) == 0 {
		return errors.New("Addresses list cannot be empty.")
	}
	publicAddresses := []string{}
	for i, addressPair := range strings.Split(kubespray.Addresses, " ") {
		pairSplit := strings.Split(addressPair, ":")
		publicAddress := pairSplit[0]
		if i == 0 {
			kubeApiserverAddress = publicAddress
		}
		publicAddresses = append(publicAddresses, publicAddress)
	}

	if _, err := os.Stat("kubespray/inventory/mycluster/hosts.ini"); os.IsNotExist(err) {
		fmt.Println(true, "Cluster inventory file not found - generating it...\n")

		err := runCommand(true,"cp", "-rfp", "kubespray/inventory/sample", "kubespray/inventory/mycluster")
		if err != nil {
			return err
		}

		args := []string{"kubespray/contrib/inventory_builder/inventory.py"}
		args = append(args, publicAddresses...)
		err = runCommandWithEnv(true, "python3", []string{"CONFIG_FILE=kubespray/inventory/mycluster/hosts.ini"}, args...)
		if err != nil {
			return err
		}

		input, err := ioutil.ReadFile("kubespray/inventory/mycluster/hosts.ini")
		if err != nil {
			return err
		}
		for _, addressPair := range strings.Split(kubespray.Addresses, " ") {
			pairSplit := strings.Split(addressPair, ":")
			publicAddress := pairSplit[0]
			privateAddress := pairSplit[1]
			input = []byte(strings.Replace(string(input), "ip="+publicAddress, "ip="+privateAddress, 1))
		}
		err = ioutil.WriteFile("kubespray/inventory/mycluster/hosts.ini", []byte(input), 0644)
		if err != nil {
			return err
		}
	} else {
		fmt.Println(true, "Existing cluster inventory file found. Reusing it.\n")
	}

	vars := map[string]interface{} {
		"kubeconfig_localhost": "True",
		"kube_apiserver_ip": kubeApiserverAddress,
		"kube_apiserver_address": kubeApiserverAddress,
		"supplementary_addresses_in_ssl_keys": []string{"10.233.0.1"},
	}
	varsJson, err := json.Marshal(vars)
	if err != nil {
		return err
	}
	args := []string{"-i", "kubespray/inventory/mycluster/hosts.ini", "kubespray/cluster.yml", "--user=" + kubespray.AnsibleUser,
		"-e", string(varsJson)}
	if kubespray.AnsibleBecome {
		args = append(args, "--become")
	}
	err = runCommand(true,"ansible-playbook", args...)
	if err != nil {
		return err
	}

	usr, err := user.Current()
	if err != nil {
		return err
	}
	err = CopyFile("kubespray/inventory/mycluster/artifacts/admin.conf", usr.HomeDir + "/.kube/config")
	if err != nil {
		return err
	}

	return nil
}

func runCommand(verbose bool, name string, args ...string) error {
	return runCommandWithEnv(verbose, name, nil, args...)
}

func runCommandWithEnv(verbose bool, name string, env []string, args ...string) error {
	e := exec.Command(name, args...)
	if verbose {
		e.Stdout = os.Stdout
		e.Stderr = os.Stderr
	}
	if env != nil {
		e.Env = env
	}
	err := e.Run()
	if err != nil {
		fmt.Println("Error: Command failed  %s %s", name, strings.Join(args, " "))
	}
	return err
}

// credit https://gist.github.com/r0l1/92462b38df26839a3ca324697c8cba04
func CopyFile(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		if e := out.Close(); e != nil {
			err = e
		}
	}()

	_, err = io.Copy(out, in)
	if err != nil {
		return
	}

	err = out.Sync()
	if err != nil {
		return
	}

	si, err := os.Stat(src)
	if err != nil {
		return
	}
	err = os.Chmod(dst, si.Mode())
	if err != nil {
		return
	}

	return
}