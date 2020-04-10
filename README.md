# kuml

![Process](./docs/images/process.png "Process")

Kuml is a kubernetes manifest visualization tool.

## Installation

### Source

```bash
go get -u github.com/gashirar/kuml/cmd
```

### krew plugin manager.

```bash
Sorry!
UnderConstruction!
```

## Usage

### Run
```bash
kuml example/application
```

### Output
```bash
@startuml
rectangle "kind: ConfigMap\nname: adapter-app-properties" as default_ConfigMap_adapter_app_properties
rectangle "kind: ConfigMap\nname: adapter-infra-properties" as default_ConfigMap_adapter_infra_properties
rectangle "kind: ConfigMap\nname: application-app-properties" as default_ConfigMap_application_app_properties
rectangle "kind: ConfigMap\nname: application-infra-properties" as default_ConfigMap_application_infra_properties
rectangle "kind: Deployment\nname: sample-deployment" as default_Deployment_sample_deployment
rectangle "kind: HorizontalPodAutoscaler\nname: sample-horizontalpodautoscaler" as default_HorizontalPodAutoscaler_sample_horizontalpodautoscaler
rectangle "kind: Ingress\nname: sample-ingress" as default_Ingress_sample_ingress
rectangle "kind: PodDisruptionBudget\nname: sample-poddisruptionbudget" as default_PodDisruptionBudget_sample_poddisruptionbudget
rectangle "kind: Pod\nname: sample-deployment" as default_Pod_sample_deployment
rectangle "kind: ReplicaSet\nname: sample-deployment" as default_ReplicaSet_sample_deployment
rectangle "kind: Service\nname: sample-service" as default_Service_sample_service
default_Deployment_sample_deployment -DOWN-> default_ReplicaSet_sample_deployment : ""
default_ReplicaSet_sample_deployment -DOWN-> default_Pod_sample_deployment : ""
default_Pod_sample_deployment -DOWN-> default_ConfigMap_adapter_app_properties : ""
default_Pod_sample_deployment -DOWN-> default_ConfigMap_adapter_infra_properties : ""
default_Pod_sample_deployment -DOWN-> default_ConfigMap_application_app_properties : ""
default_Pod_sample_deployment -DOWN-> default_ConfigMap_application_infra_properties : ""
default_Service_sample_service -RIGHT-> default_Pod_sample_deployment : ""
default_Ingress_sample_ingress -RIGHT-> default_Service_sample_service : ""
default_Ingress_sample_ingress -RIGHT-> default_Service_sample_service : ""
default_PodDisruptionBudget_sample_poddisruptionbudget -LEFT-> default_Pod_sample_deployment : ""
default_HorizontalPodAutoscaler_sample_horizontalpodautoscaler -LEFT-> default_Deployment_sample_deployment : ""
@enduml
```

### Generate UML diagram
In your favorite way.

-> [call it from your script using command line - PlantUML](https://plantuml.com/en/command-line/)  
-> [PlantUML Web Server](http://www.plantuml.com/plantuml/uml/)

![Process](./docs/images/uml.png "Process")
