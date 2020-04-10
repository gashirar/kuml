package resource

import (
	"bytes"
	"io/ioutil"
	appsv1 "k8s.io/api/apps/v1"
	autoscalingv1 "k8s.io/api/autoscaling/v1"
	batchv1 "k8s.io/api/batch/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	extenshionsv1beta1 "k8s.io/api/extensions/v1beta1"
	policyv1beta1 "k8s.io/api/policy/v1beta1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"log"
	"os"
	"path/filepath"
	"sigs.k8s.io/yaml"
)

type APIResource interface {
	GetName() string
	GetNamespace() string
	GetLabels() map[string]string
	GroupVersionKind() schema.GroupVersionKind
}

type APIResourceList struct {
	Items []APIResource
}

func NewAPIResourceList(yamlByteSlice [][]byte) APIResourceList {
	var res APIResourceList
	for _, yamlByte := range yamlByteSlice {
		kind, _ := checkResourceKind(yamlByte)
		switch kind {
		case "CronJob":
			r := batchv1beta1.CronJob{}
			yaml.Unmarshal(yamlByte, &r)
			res.Items = append(res.Items, &r)

			pod := corev1.Pod{}
			pod.Name = r.Name
			pod.Spec = r.Spec.JobTemplate.Spec.Template.Spec
			pod.Labels = r.Spec.JobTemplate.Spec.Template.Labels
			res.Items = append(res.Items, &pod)
		case "Deployment":
			r := appsv1.Deployment{}
			yaml.Unmarshal(yamlByte, &r)
			res.Items = append(res.Items, &r)

			rs := appsv1.ReplicaSet{}
			rs.Kind = "ReplicaSet"
			rs.Name = r.Name
			rs.Labels = r.Spec.Template.Labels
			rs.Spec.Selector = r.Spec.Selector
			res.Items = append(res.Items, &rs)

			pod := corev1.Pod{}
			pod.Kind = "Pod"
			pod.Name = r.Name
			pod.Spec = r.Spec.Template.Spec
			pod.Labels = r.Spec.Template.Labels
			res.Items = append(res.Items, &pod)
		case "Job":
			r := batchv1.Job{}
			yaml.Unmarshal(yamlByte, &r)
			res.Items = append(res.Items, &r)

			pod := corev1.Pod{}
			pod.Kind = "Pod"
			pod.Name = r.Name
			pod.Spec = r.Spec.Template.Spec
			pod.Labels = r.Spec.Template.Labels
			res.Items = append(res.Items, &pod)
		case "Pod":
			r := corev1.Pod{}
			yaml.Unmarshal(yamlByte, &r)
			res.Items = append(res.Items, &r)
		case "ReplicaSet":
			r := appsv1.ReplicaSet{}
			yaml.Unmarshal(yamlByte, &r)
			res.Items = append(res.Items, &r)

			pod := corev1.Pod{}
			pod.Kind = "Pod"
			pod.Name = r.Name
			pod.Spec = r.Spec.Template.Spec
			pod.Labels = r.Spec.Template.Labels
			res.Items = append(res.Items, &pod)
		case "StatefulSet":
			r := appsv1.StatefulSet{}
			yaml.Unmarshal(yamlByte, &r)
			res.Items = append(res.Items, &r)

			pod := corev1.Pod{}
			pod.Kind = "Pod"
			pod.Name = r.Name
			pod.Spec = r.Spec.Template.Spec
			pod.Labels = r.Spec.Template.Labels
			res.Items = append(res.Items, &pod)
		case "Ingress":
			r := extenshionsv1beta1.Ingress{}
			yaml.Unmarshal(yamlByte, &r)
			res.Items = append(res.Items, &r)
		case "Service":
			r := corev1.Service{}
			yaml.Unmarshal(yamlByte, &r)
			res.Items = append(res.Items, &r)
		case "ConfigMap":
			r := corev1.ConfigMap{}
			yaml.Unmarshal(yamlByte, &r)
			res.Items = append(res.Items, &r)
		case "Secret":
			r := corev1.Secret{}
			yaml.Unmarshal(yamlByte, &r)
			res.Items = append(res.Items, &r)
		case "HorizontalPodAutoscaler":
			r := autoscalingv1.HorizontalPodAutoscaler{}
			yaml.Unmarshal(yamlByte, &r)
			res.Items = append(res.Items, &r)
		case "PodDisruptionBudget":
			r := policyv1beta1.PodDisruptionBudget{}
			yaml.Unmarshal(yamlByte, &r)
			res.Items = append(res.Items, &r)
		default:
		}
	}

	return res
}

func checkResourceKind(yamlByte []byte) (string, error) {
	yamlMap := make(map[string]interface{})
	err := yaml.Unmarshal(yamlByte, &yamlMap)

	if err != nil {
		return "", err
	}
	if len(yamlMap) == 0 {
		return "", nil
	}

	kind := yamlMap["kind"].(string)

	return kind, nil
}

func IsDirectory(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		log.Fatal(err)
	}
	mode := fi.Mode()
	return mode.IsDir()
}

func ReadYamlFile(path string) [][]byte {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	return bytes.Split(buf, []byte("\n---"))
}

func ReadYaml(recursive bool, paths ...string) [][]byte {
	var yamlByteSlice [][]byte
	for _, path := range paths {
		if IsDirectory(path) {
			err := filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
				if !info.IsDir() {
					yamlByteSlice = append(yamlByteSlice, ReadYamlFile(p)...)
				} else {
					if recursive {
						yamlByteSlice = append(yamlByteSlice, ReadYaml(recursive, p)...)
					}
				}
				return nil
			})
			if err != nil {
				log.Fatal(err)
			}
		} else {
			yamlByteSlice = append(yamlByteSlice, ReadYamlFile(path)...)
		}
	}
	return yamlByteSlice
}
