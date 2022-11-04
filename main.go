package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/hcl2/gohcl"
	"github.com/hashicorp/hcl2/hclwrite"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type AwsIAMPolicyResource struct {
	Name   string `hcl:"name"`
	Path   string `hcl:"path"`
	Desc   string `hcl:"description"`
	Policy string `hcl:"policy"`
}

type IAMpolicy struct {
	Statement IAMStatement
	Version   string
}

type IAMStatement struct {
	Action   []string
	Effect   string
	Resource string
	Sid      string
}

var m = make(map[string]bool)
var a = []string{}

type ClientSet struct {
	Client *kubernetes.Clientset
}

func NewConfig() *rest.Config {
	_, err := os.Stat(os.Getenv("HOME") + "/.kube/config")
	if err == nil {
		log.Print("Config File is found \n")
		path := os.Getenv("HOME") + "/.kube/config"
		config, _ := clientcmd.BuildConfigFromFlags("", path)
		return config
	} else {
		log.Print("Config File is not found \n")
		config, err := rest.InClusterConfig()
		if err != nil {
			log.Fatal(err, "failed to get InClusterConfig")
		}
		return config
	}
}

func NewK8sClientSet() (*ClientSet, error) {
	var err error
	config := NewConfig()
	k := &ClientSet{}
	k.Client, err = kubernetes.NewForConfig(config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get k8s k8sClientset")
	}
	return k, nil
}

func main() {
	// k8s, err := newK8s()
	k8s, err := NewK8sClientSet()
	if err != nil {
		fmt.Println("Can't create ClientSet", err)
		return
	}

	pods, err := k8s.Client.CoreV1().Pods("default").List(context.TODO(),
		v1.ListOptions{LabelSelector: "app.kubernetes.io/name=localstack"})

	podName := pods.Items[0].Name

	log.Printf("Localstack podName is : ", string(podName))

	podLogs, err := k8s.Client.CoreV1().Pods("default").GetLogs(podName, &corev1.PodLogOptions{}).Stream(context.TODO())
	if err != nil {
		log.Fatalf("error in opening stream", err)
	}
	defer podLogs.Close()

	re := regexp.MustCompile(`(AWS)\s\S+`)

	for i := 0; i < 10; i++ {
		reader := bufio.NewScanner(podLogs)
		for reader.Scan() {
			line := reader.Text()
			match, _ := regexp.MatchString("localstack.request.aws", line)
			if err != nil {
				log.Fatal(err)
			} else if match {
				request := strings.Replace(re.FindString(line), ".", ":", 1)
				final := strings.Replace(request, "AWS ", "", 1)
				add(final)
			}

		}
		time.Sleep(1 * time.Millisecond)
	}

	iamString, _ := prettyStruct(NewIAMpolicy(a))

	ins := AwsIAMPolicyResource{
		Name: "S3 IAM policy",
		Path: "/",
		Desc: "IAM policy generated on top of results provided by sniffer",
		// Policy: PrintIAM(a),
		Policy: iamString,
	}
	f := hclwrite.NewEmptyFile()
	gohcl.EncodeIntoBody(&ins, f.Body())

	fmt.Println("====== Policy ===============")
	fmt.Println(iamString)
	fmt.Println("====== ============== =======\n")
	fmt.Println("====== Terraform plan =======")
	fmt.Println("resource \"aws_iam_policy\" \"s3-policy\" {")
	// fmt.Printf(ins.Desc)
	fmt.Printf("%s", f.Bytes())
	fmt.Printf("}\n")
	fmt.Println("====== ============== =======")

}

func add(s string) {
	if m[s] {
		return // Already in the map
	}
	a = append(a, s)
	m[s] = true
}

func NewIAMpolicy(actions []string) IAMpolicy {
	return IAMpolicy{
		Statement: IAMStatement{
			Action:   actions,
			Effect:   "Allow",
			Resource: "*",
			Sid:      "RANDOM",
		},
		Version: "",
	}
}

func PrintIAM(actions []string) string {
	iam := NewIAMpolicy(actions)

	b, err := json.Marshal(iam)
	if err != nil {
		fmt.Println(err)
		return "error"
	}
	return string(b)
}

func prettyStruct(data interface{}) (string, error) {
	val, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return "", err
	}
	return string(val), nil
}
