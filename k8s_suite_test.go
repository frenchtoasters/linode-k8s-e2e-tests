package e2e_test

import (
	"flag"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/onsi/ginkgo/reporters"
	"github.com/linode/linode-k8s-e2e-tests/framework"
	"github.com/linode/linode-k8s-e2e-tests/rand"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	useExisting   = false
	kubecofigFile = filepath.Join(homedir.HomeDir(), ".kube/config")
	ClusterName string
)

func init() {
	flag.StringVar(&framework.Image, "image", framework.Image, "registry/repository:tag")
	flag.StringVar(&framework.ApiToken, "api-token", os.Getenv("LINODE_API_TOKEN"), "linode api token")

	flag.BoolVar(&useExisting, "use-existing", useExisting, "Use existing kubernetes cluster")
	flag.StringVar(&kubecofigFile, "kubeconfig", kubecofigFile, "To use existing cluster provide kubeconfig file")

	var errRandom error

	ClusterName, errRandom = rand.WithRandomSuffix("ccm-linode")
	if errRandom != nil {
		panic(errRandom)
	}
}

const (
	TIMEOUT = 20 * time.Minute
)

var (
	root *framework.Framework
)

func TestE2e(t *testing.T) {
	RegisterFailHandler(Fail)
	SetDefaultEventuallyTimeout(TIMEOUT)

	junitReporter := reporters.NewJUnitReporter("junit.xml")
	RunSpecsWithDefaultAndCustomReporters(t, "e2e Suite", []Reporter{junitReporter})
}

var _ = BeforeSuite(func() {

	if !useExisting {
		err := framework.CreateCluster(ClusterName)
		Expect(err).NotTo(HaveOccurred())
		dir, err := os.Getwd()
		Expect(err).NotTo(HaveOccurred())
		kubecofigFile = filepath.Join(dir, ClusterName+".conf")
	}

	By("Using kubeconfig from " + kubecofigFile)
	config, err := clientcmd.BuildConfigFromFlags("", kubecofigFile)
	Expect(err).NotTo(HaveOccurred())

	// Clients
	kubeClient := kubernetes.NewForConfigOrDie(config)

	// Framework
	root, err = framework.New(config, kubeClient, kubecofigFile)
	Expect(err).NotTo(HaveOccurred())

	By("Using namespace " + root.Namespace())

	// Create namespace
	err = root.CreateNamespace()
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	if !useExisting {
		err := framework.DeleteCluster()
		Expect(err).NotTo(HaveOccurred())
	}
})
