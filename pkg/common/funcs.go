package common

import (
	"encoding/base64"
	"regexp"
	"strings"

	sprigcrypto "github.com/crowdstrike/falcon-operator/pkg/sprig"
	"k8s.io/apimachinery/pkg/version"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func InitContainerArgs() []string {
	return []string{
		"-c",
		"mkdir -p " + FalconDataDir +
			" && " +
			"touch " + FalconStoreFile,
	}
}

func FCAdmissionReviewVersions() []string {
	kubeVersion := GetKubernetesVersion()
	fcArv := []string{"v1"}

	if strings.Compare(kubeVersion.Minor, "22") < 0 {
		fcArv = []string{"v1", "v1beta"}
	}

	return fcArv
}

func GenCert(cn string, ips []interface{}, alternateDNS []interface{}, validity int, ca sprigcrypto.Certificate) sprigcrypto.Certificate {
	certs, err := sprigcrypto.GenerateSignedCertificate(cn, ips, alternateDNS, validity, ca)
	if err != nil {
		panic(err.Error())
	}

	return certs
}

func GenCA(cn string, validity int) sprigcrypto.Certificate {
	ca, err := sprigcrypto.GenerateCertificateAuthority(cn, validity)
	if err != nil {
		panic(err.Error())
	}
	return ca
}

func EncodedBase64String(data string) []byte {
	base64EncodedData := make([]byte, base64.StdEncoding.EncodedLen(len(data)))
	base64.StdEncoding.Encode(base64EncodedData, []byte(data))
	return base64EncodedData
}

func CleanDecodedBase64(s []byte) []byte {
	re := regexp.MustCompile(`[\t|\n]*`)
	b64byte, err := base64.StdEncoding.DecodeString(string(s))
	if err != nil {
		return []byte(re.ReplaceAllString(string(s), ""))
	}
	return []byte(re.ReplaceAllString(string(b64byte), ""))
}

func GetKubernetesVersion() *version.Info {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	version, err := clientset.DiscoveryClient.ServerVersion()
	if err != nil {
		panic(err.Error())
	}

	return version
}
