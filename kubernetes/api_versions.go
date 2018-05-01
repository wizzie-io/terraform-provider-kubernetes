package kubernetes

import (
	"encoding/json"
	"log"
	"strconv"

	"time"

	"fmt"

	"k8s.io/client-go/kubernetes"
)

type APIGroup int

const (
	none APIGroup = iota
	appsV1
	appsV1beta1
	appsV1beta2
	extensionsV1beta1
)

func (g APIGroup) String() string {
	switch g {
	case appsV1:
		return "apps/v1"
	case appsV1beta1:
		return "apps/v1beta1"
	case appsV1beta2:
		return "apps/v1beta2"
	case extensionsV1beta1:
		return "extensions/v1beta1"
	default:
		return "none"
	}
}

func (kp *kubernetesProvider) lowestSupportedAPIGroup(rtype string, groups ...APIGroup) (APIGroup, error) {
	for i := len(groups) - 1; i > -1; i-- {
		g := groups[i]
		match, err := kp.serverSupportsResourceAPIVersion(rtype, g.String())
		if err != nil {
			return none, err
		} else if match {
			return g, nil
		}
	}
	return none, nil
}

func (kp *kubernetesProvider) highestSupportedAPIGroup(rtype string, groups ...APIGroup) (APIGroup, error) {
	for _, g := range groups {
		match, err := kp.serverSupportsResourceAPIVersion(rtype, g.String())
		if err != nil {
			return none, err
		} else if match {
			return g, nil
		}
	}
	return none, nil
}

func (kp *kubernetesProvider) serverSupportsResourceAPIVersion(rname string, groupVersion string) (bool, error) {
	start := time.Now()
	resList, err := kp.discoClient.ServerResources()
	//resList, err := providerInstance.conn.DiscoveryClient.ServerResources()
	if err != nil {
		log.Printf("[WARN] discovery client could not resource list: %v\n", err)
		return false, err
	}
	log.Printf("[DEBUG] retrieved resource list in %v\n", time.Now().Sub(start))

	for _, v := range resList {
		if v.GroupVersion == groupVersion {
			for _, v2 := range v.APIResources {
				if v2.Name == rname {
					log.Printf("[DEBUG] api group [%s] supports %s resource type\n", groupVersion, rname)
					return true, nil
				}
			}
		}
	}
	log.Printf("[DEBUG] api group [%s] does not supports %s resource type on Kubernetes server\n", groupVersion, rname)

	return false, nil
}

// Convert between two types by converting to/from JSON. Intended to switch
// between multiple API versions, as they are strict supersets of one another.
// item and out are pointers to structs
func Convert(item, out interface{}) error {
	bytes, err := json.Marshal(item)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, out)
	if err != nil {
		return err
	}

	return nil
}

// ServerVersionPre1_9 reads the Kubernetes API verions and returns true if less
// than v1.9
func (kp *kubernetesProvider) ServerVersionPre1_9(conn *kubernetes.Clientset) bool {
	ver, _ := kp.discoClient.ServerVersion()
	minor, _ := strconv.Atoi(string(ver.Minor[0]))
	log.Printf("[INFO] Kubernetes Server version: %#v", ver)

	if ver.Major == "1" && minor < 9 {
		return true
	}

	return false
}

func printObjectJSON(item interface{}) error {
	bytes, err := json.MarshalIndent(item, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(bytes))

	return nil
}
