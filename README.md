# bookstore-controller-kubebuilder

A Kubernetes Controller written using Kubebuilder. It reconciles a custom resource named Bookstore. A bookstore
object creates a deployment and service. We provide Replica count, Service Type, ContainerPort and
KindNodePort (if the service is type NodePort) with
other optional fields when creating an object of type Bookstore. The container image is a simple
[bookstore api server](https://github.com/Shaad7/bookstore-api-server). `port` value should match with `hostPort`
value in `cluster-config.yaml` file. If the `serviceType` is `NodePort` the server listens and  serve request on given port. 

## How to Use
```shell
git clone https://github.com/Shaad7/bookstore-controller-kubebuilder
cd bookstore-controller-kubebuilder
```

Create a cluster using  [Kind](https://kind.sigs.k8s.io/docs/user/quick-start/#installation)
```shell
kind create cluster --config=clusterconfig.yaml 
```

Generate CRD yaml and generated codes
```shell
make manifest
```

Register CRD
```shell
make install
```

Run the Controller
```shell
make run 
```
Create an example Custom Resource
```shell
kubectl apply -f config/samples/core_v1_bookstore.yaml 
```

