package puppet

import (
	"encoding/json"
	"github.com/pkg/errors"
	"os/exec"
)

type Puppet struct {
	CaCert      string `json:"cacert"`
	CaCrl       string `json:"cacrl"`
	CaKey       string `json:"cakey"`
	HostCert    string `json:"hostcert"`
	HostPrivKey string `json:"hostprivkey"`
	MasterPort  int    `json:"masterport"`
	RequestDir  string `json:"requestdir"`
	SignedDir   string `json:"signeddir"`
}

func New() (*Puppet, error) {
	var puppet Puppet
    return &puppet, loadConfig(&puppet)
}

func (puppet *Puppet) GenerateCertificate(certname string) error {
	cmd := exec.Command("openssl", "x509", "-req",
		"-in", puppet.RequestDir + "/" + certname,
		"-out", puppet.SignedDir + "/" + certname,
		"-CA", puppet.CaCert,
		"-CAkey", puppet.CaKey,
		"-CAcreateserial", "-sha256")

	if err := cmd.Run() ; err != nil {
		return errors.Wrapf(err, "could not sign certificate for %s", certname)
	}

	return nil
}

func loadConfig(puppet *Puppet) error {
	cmd := exec.Command("puppet", "config", "print", "--render-as", "json")

	out, err := cmd.Output()
	if err != nil {
		return errors.Wrap(err, "could not load Puppet config")
	}

	err = json.Unmarshal(out, puppet)
	if err != nil {
		return errors.Wrapf(err, "could not parse Puppet config")
	}

	return nil
}
