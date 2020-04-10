package plantuml

import (
	"fmt"
	"github.com/gashirar/kuml/pkg/resource"
	appsv1 "k8s.io/api/apps/v1"
	autoscalingv1 "k8s.io/api/autoscaling/v1"
	corev1 "k8s.io/api/core/v1"
	extenshionsv1beta1 "k8s.io/api/extensions/v1beta1"
	policyv1beta1 "k8s.io/api/policy/v1beta1"
	"sort"
	"strings"
)

type RenderOption struct {
	showLinkLabel bool
}

type PlantUML struct {
	elementList  ElementList
	linkList     LinkList
	renderOption RenderOption
}

func (u *PlantUML) Render() {
	fmt.Println("@startuml")

	e := u.elementList.Items
	sort.Slice(e, func(i, j int) bool { return e[i].UniqueId < e[j].UniqueId })
	for _, elem := range e {
		elem.Render()
	}

	for _, link := range u.linkList.Items {
		link.Render()
	}

	fmt.Println("@enduml")
}

type Element struct {
	UniqueId    string
	Description string
}

func (e *Element) Render() {
	fmt.Printf("rectangle \"%s\" as %s\n", e.Description, e.UniqueId)
}

type ElementList struct {
	Items []Element
}

type Link struct {
	From      string
	To        string
	Connector string
	Label     string
}

func (l Link) Render() {
	fmt.Printf("%s %s %s : \"%s\"\n", l.From, l.Connector, l.To, "")
}

type LinkList struct {
	Items []Link
}

func NewPlantUML(resource resource.APIResourceList) PlantUML {
	elementList := NewElementList(resource)
	linkList := NewLinkList(resource)

	return PlantUML{elementList: elementList, linkList: linkList}
}

func NewElement(name string, description string) Element {
	return Element{
		UniqueId:    name,
		Description: description,
	}
}
func NewElementList(list resource.APIResourceList) ElementList {
	var elementList ElementList

	for _, apiRes := range list.Items {
		kind := apiRes.GroupVersionKind().Kind
		namespace := apiRes.GetNamespace()
		name := apiRes.GetName()
		uniqueId := createUniqueId(namespace, kind, name)
		description := fmt.Sprintf("kind: %s\\nname: %s", kind, name)

		elementList.Items = append(elementList.Items, NewElement(uniqueId, description))
	}

	return elementList
}

func NewLink(from string, to string, connector string, label string) Link {
	return Link{
		From:      from,
		To:        to,
		Connector: connector,
		Label:     label,
	}
}

func NewLinkList(resource resource.APIResourceList) LinkList {
	linkList := LinkList{}

	linkList.Items = append(linkList.Items, DeploymentToReplicaSet(resource).Items...)
	linkList.Items = append(linkList.Items, ReplicaSetToPod(resource).Items...)
	linkList.Items = append(linkList.Items, PodToConfigMap(resource).Items...)
	linkList.Items = append(linkList.Items, PodToSecret(resource).Items...)
	linkList.Items = append(linkList.Items, ServiceToPod(resource).Items...)
	linkList.Items = append(linkList.Items, IngressToService(resource).Items...)
	linkList.Items = append(linkList.Items, PodDisruptionBudgetToPod(resource).Items...)
	linkList.Items = append(linkList.Items, HorizontalPodAutoscalerToDeployment(resource).Items...)
	linkList.Items = append(linkList.Items, CronJobToJob(resource).Items...)
	linkList.Items = append(linkList.Items, JobToPod(resource).Items...)
	linkList.Items = append(linkList.Items, StatefulSetToPod(resource).Items...)
	return linkList
}

func DeploymentToReplicaSet(apiList resource.APIResourceList) LinkList {
	linkList := LinkList{}

	for _, res := range apiList.Items {
		if res.GroupVersionKind().Kind == "Deployment" {
			for _, targetRes := range apiList.Items {
				if targetRes.GroupVersionKind().Kind == "ReplicaSet" {
					matchLabels := res.(*appsv1.Deployment).Spec.Selector.MatchLabels
					if IsMapContainsMap(targetRes.GetLabels(), matchLabels) {
						from := createUniqueId(res.GetNamespace(), res.GroupVersionKind().Kind, res.GetName())
						to := createUniqueId(targetRes.GetNamespace(), targetRes.GroupVersionKind().Kind, targetRes.GetName())
						linkList.Items = append(linkList.Items, NewLink(from, to, "-DOWN->", labelMapToString(matchLabels)))
					}
				}
			}
		}
	}
	return linkList
}

func ReplicaSetToPod(apiList resource.APIResourceList) LinkList {
	linkList := LinkList{}

	for _, res := range apiList.Items {
		if res.GroupVersionKind().Kind == "ReplicaSet" {
			for _, targetRes := range apiList.Items {
				if targetRes.GroupVersionKind().Kind == "Pod" {
					matchLabels := res.(*appsv1.ReplicaSet).Spec.Selector.MatchLabels
					if IsMapContainsMap(targetRes.GetLabels(), matchLabels) {
						from := createUniqueId(res.GetNamespace(), res.GroupVersionKind().Kind, res.GetName())
						to := createUniqueId(targetRes.GetNamespace(), targetRes.GroupVersionKind().Kind, targetRes.GetName())
						linkList.Items = append(linkList.Items, NewLink(from, to, "-DOWN->", labelMapToString(matchLabels)))
					}
				}
			}
		}
	}
	return linkList
}

func StatefulSetToPod(apiList resource.APIResourceList) LinkList {
	linkList := LinkList{}

	for _, res := range apiList.Items {
		if res.GroupVersionKind().Kind == "StatefulSet" {
			for _, targetRes := range apiList.Items {
				if targetRes.GroupVersionKind().Kind == "Pod" {
					matchLabels := res.(*appsv1.StatefulSet).Spec.Selector.MatchLabels
					if IsMapContainsMap(targetRes.GetLabels(), matchLabels) {
						from := createUniqueId(res.GetNamespace(), res.GroupVersionKind().Kind, res.GetName())
						to := createUniqueId(targetRes.GetNamespace(), targetRes.GroupVersionKind().Kind, targetRes.GetName())
						linkList.Items = append(linkList.Items, NewLink(from, to, "-DOWN->", labelMapToString(matchLabels)))
					}
				}
			}
		}
	}
	return linkList
}

func CronJobToJob(apiList resource.APIResourceList) LinkList {
	linkList := LinkList{}

	for _, res := range apiList.Items {
		if res.GroupVersionKind().Kind == "CronJob" {
			from := createUniqueId(res.GetNamespace(), res.GroupVersionKind().Kind, res.GetName())
			to := createUniqueId(res.GetNamespace(), "Job", res.GetName())
			linkList.Items = append(linkList.Items, NewLink(from, to, "-DOWN->", ""))
		}
	}
	return linkList
}

func JobToPod(apiList resource.APIResourceList) LinkList {
	linkList := LinkList{}

	for _, res := range apiList.Items {
		if res.GroupVersionKind().Kind == "Job" {
			from := createUniqueId(res.GetNamespace(), res.GroupVersionKind().Kind, res.GetName())
			to := createUniqueId(res.GetNamespace(), "Pod", res.GetName())
			linkList.Items = append(linkList.Items, NewLink(from, to, "-DOWN->", ""))
		}
	}
	return linkList
}

func PodToConfigMap(apiList resource.APIResourceList) LinkList {
	linkList := LinkList{}

	for _, res := range apiList.Items {
		if res.GroupVersionKind().Kind == "Pod" {
			volumes := res.(*corev1.Pod).Spec.Volumes
			for _, volume := range volumes {
				if volume.ConfigMap != nil {
					from := createUniqueId(res.GetNamespace(), res.GroupVersionKind().Kind, res.GetName())
					to := createUniqueId(res.GetNamespace(), "ConfigMap", volume.ConfigMap.Name)
					linkList.Items = append(linkList.Items, NewLink(from, to, "-DOWN->", ".spec.volume.configMap"))
				}
				if volume.Projected != nil {
					for _, projected := range volume.Projected.Sources {
						if projected.ConfigMap != nil {
							from := createUniqueId(res.GetNamespace(), res.GroupVersionKind().Kind, res.GetName())
							to := createUniqueId(res.GetNamespace(), "ConfigMap", projected.ConfigMap.Name)
							linkList.Items = append(linkList.Items, NewLink(from, to, "-DOWN->", ".spec.volume.projected.sources.configMap"))
						}
					}
				}
			}
		}
	}
	return linkList
}

func PodToSecret(apiList resource.APIResourceList) LinkList {
	linkList := LinkList{}

	for _, res := range apiList.Items {
		if res.GroupVersionKind().Kind == "Pod" {
			volumes := res.(*corev1.Pod).Spec.Volumes
			for _, volume := range volumes {
				if volume.Secret != nil {
					from := createUniqueId(res.GetNamespace(), res.GroupVersionKind().Kind, res.GetName())
					to := createUniqueId(res.GetNamespace(), "Secret", volume.Secret.SecretName)
					linkList.Items = append(linkList.Items, NewLink(from, to, "-DOWN->", ".spec.volume.secret"))
				}
				if volume.Projected != nil {
					for _, projected := range volume.Projected.Sources {
						if projected.Secret != nil {
							from := createUniqueId(res.GetNamespace(), res.GroupVersionKind().Kind, res.GetName())
							to := createUniqueId(res.GetNamespace(), "Secret", projected.Secret.Name)
							linkList.Items = append(linkList.Items, NewLink(from, to, "-DOWN->", ".spec.volume.projected.sources.secret"))
						}
					}
				}
			}
		}
	}
	return linkList
}

func ServiceToPod(apiList resource.APIResourceList) LinkList {
	linkList := LinkList{}

	for _, res := range apiList.Items {
		if res.GroupVersionKind().Kind == "Service" {
			for _, targetRes := range apiList.Items {
				if targetRes.GroupVersionKind().Kind == "Pod" {
					matchLabels := res.(*corev1.Service).Spec.Selector
					if IsMapContainsMap(targetRes.GetLabels(), matchLabels) {
						from := createUniqueId(res.GetNamespace(), res.GroupVersionKind().Kind, res.GetName())
						to := createUniqueId(targetRes.GetNamespace(), targetRes.GroupVersionKind().Kind, targetRes.GetName())
						linkList.Items = append(linkList.Items, NewLink(from, to, "-RIGHT->", labelMapToString(matchLabels)))
					}
				}
			}
		}
	}
	return linkList
}

func PodDisruptionBudgetToPod(apiList resource.APIResourceList) LinkList {
	linkList := LinkList{}

	for _, res := range apiList.Items {
		if res.GroupVersionKind().Kind == "PodDisruptionBudget" {
			matched := false
			matchLabels := res.(*policyv1beta1.PodDisruptionBudget).Spec.Selector.MatchLabels
			for _, targetRes := range apiList.Items {
				if targetRes.GroupVersionKind().Kind == "Pod" {
					if IsMapContainsMap(targetRes.GetLabels(), matchLabels) {
						matched = true
						from := createUniqueId(res.GetNamespace(), res.GroupVersionKind().Kind, res.GetName())
						to := createUniqueId(targetRes.GetNamespace(), targetRes.GroupVersionKind().Kind, targetRes.GetName())
						linkList.Items = append(linkList.Items, NewLink(from, to, "-LEFT->", labelMapToString(matchLabels)))
					}
				}
			}
			if !matched {
				from := createUniqueId(res.GetNamespace(), res.GroupVersionKind().Kind, res.GetName())
				to := "(No Target Pod)"
				linkList.Items = append(linkList.Items, NewLink(from, to, "-LEFT->", labelMapToString(matchLabels)))
			}
		}
	}
	return linkList
}

func HorizontalPodAutoscalerToDeployment(apiList resource.APIResourceList) LinkList {
	linkList := LinkList{}

	for _, res := range apiList.Items {
		if res.GroupVersionKind().Kind == "HorizontalPodAutoscaler" {
			matched := false
			scaleTargetRef := res.(*autoscalingv1.HorizontalPodAutoscaler).Spec.ScaleTargetRef
			for _, targetRes := range apiList.Items {
				if targetRes.GroupVersionKind().Kind == "Deployment" {
					if scaleTargetRef.Name == targetRes.GetName() {
						matched = true
						from := createUniqueId(res.GetNamespace(), res.GroupVersionKind().Kind, res.GetName())
						to := createUniqueId(targetRes.GetNamespace(), targetRes.GroupVersionKind().Kind, targetRes.GetName())
						label := fmt.Sprintf(".spec.scaleTargetRef.kind: %s\\n.spec.scaleTargetRef.name: %s", targetRes.GroupVersionKind().Kind, targetRes.GetName())
						linkList.Items = append(linkList.Items, NewLink(from, to, "-LEFT->", label))
					}
				}
			}
			if !matched {
				from := createUniqueId(res.GetNamespace(), res.GroupVersionKind().Kind, res.GetName())
				to := "(No Target Deployment)"
				linkList.Items = append(linkList.Items, NewLink(from, to, "-LEFT->", ""))
			}
		}
	}
	return linkList
}

func IngressToService(apiList resource.APIResourceList) LinkList {
	linkList := LinkList{}

	for _, res := range apiList.Items {
		if res.GroupVersionKind().Kind == "Ingress" {
			for _, targetRes := range apiList.Items {
				if targetRes.GroupVersionKind().Kind == "Service" {
					ing := res.(*extenshionsv1beta1.Ingress)
					svc := targetRes.(*corev1.Service)
					if ing.Spec.Backend != nil && ing.Spec.Backend.ServiceName == svc.Name {
						matched := false
						for _, port := range svc.Spec.Ports {
							if ing.Spec.Backend.ServicePort.StrVal == port.Name ||
								ing.Spec.Backend.ServicePort.IntVal == port.Port {
								matched = true
								from := createUniqueId(res.GetNamespace(), res.GroupVersionKind().Kind, res.GetName())
								to := createUniqueId(targetRes.GetNamespace(), targetRes.GroupVersionKind().Kind, targetRes.GetName())
								label := fmt.Sprintf(".spec.backend.serviceName: %s\\n.spec.backend.servicePort: %s", ing.Spec.Backend.ServiceName, ing.Spec.Backend.ServicePort.String())
								linkList.Items = append(linkList.Items, NewLink(from, to, "-RIGHT->", label))
							}
						}
						if !matched {
							from := createUniqueId(res.GetNamespace(), res.GroupVersionKind().Kind, res.GetName())
							to := "(No backend Service)"
							label := fmt.Sprintf(".spec.backend.serviceName: %s\\n.spec.backend.servicePort: %s", ing.Spec.Backend.ServiceName, ing.Spec.Backend.ServicePort.String())
							linkList.Items = append(linkList.Items, NewLink(from, to, "-RIGHT->", label))
						}
					}
					if ing.Spec.Rules != nil {
						for _, rule := range ing.Spec.Rules {
							for _, path := range rule.HTTP.Paths {
								if path.Backend.ServiceName == svc.Name {
									matched := false
									for _, port := range svc.Spec.Ports {
										if path.Backend.ServicePort.StrVal == port.Name ||
											path.Backend.ServicePort.IntVal == port.Port {
											matched = true
											from := createUniqueId(res.GetNamespace(), res.GroupVersionKind().Kind, res.GetName())
											to := createUniqueId(targetRes.GetNamespace(), targetRes.GroupVersionKind().Kind, targetRes.GetName())
											label := fmt.Sprintf(".spec.rules.http.paths.backend.serviceName: %s\\n.spec.rules.http.paths.backend.servicePort: %s", path.Backend.ServiceName, path.Backend.ServicePort.String())
											linkList.Items = append(linkList.Items, NewLink(from, to, "-RIGHT->", label))
										}
									}
									if !matched {
										from := createUniqueId(res.GetNamespace(), res.GroupVersionKind().Kind, res.GetName())
										to := "(No backend Service)"
										label := fmt.Sprintf(".spec.backend.serviceName: %s\\n.spec.backend.servicePort: %s", path.Backend.ServiceName, path.Backend.ServicePort.String())
										linkList.Items = append(linkList.Items, NewLink(from, to, "-RIGHT->", label))
									}
								}
							}
						}
					}
				}
			}
		}
	}
	return linkList
}

func IsMapContainsMap(mainMap map[string]string, subMap map[string]string) bool {
	for sk, sv := range subMap {
		isContains := false
		for mk, mv := range mainMap {
			if mk == sk && mv == sv {
				isContains = true
			}
		}
		if !isContains {
			return false
		}
	}
	return true
}

func createUniqueId(namespace string, kind string, name string) string {
	if namespace == "" {
		namespace = "default"
	}
	return namespace + "_" + kind + "_" + strings.ReplaceAll(name, "-", "_")
}

func labelMapToString(label map[string]string) string {
	labelString := ""
	for k, v := range label {
		labelString += string(k) + " : " + string(v) + "\\n"
	}
	return strings.TrimRight(labelString, "\\n")
}
