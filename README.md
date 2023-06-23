# clustermanager
// TODO(user): Add simple overview of use/purpose

## Description
// TODO(user): An in-depth paragraph about your project and overview of use

## Getting Started
You’ll need a Kubernetes cluster to run against. You can use [KIND](https://sigs.k8s.io/kind) to get a local cluster for testing, or run against a remote cluster.
**Note:** Your controller will automatically use the current context in your kubeconfig file (i.e. whatever cluster `kubectl cluster-info` shows).

### Running on the cluster
1. Install Instances of Custom Resources:

```sh
kubectl apply -f config/samples/
```

2. Build and push your image to the location specified by `IMG`:

```sh
make docker-build docker-push IMG=<some-registry>/clustermanager:tag
```

3. Deploy the controller to the cluster with the image specified by `IMG`:

```sh
make deploy IMG=<some-registry>/clustermanager:tag
```

### Uninstall CRDs
To delete the CRDs from the cluster:

```sh
make uninstall
```

### Undeploy controller
UnDeploy the controller from the cluster:

```sh
make undeploy
```

## Contributing
// TODO(user): Add detailed information on how you would like others to contribute to this project

### How it works
This project aims to follow the Kubernetes [Operator pattern](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/).

It uses [Controllers](https://kubernetes.io/docs/concepts/architecture/controller/),
which provide a reconcile function responsible for synchronizing resources until the desired state is reached on the cluster.

### Test It Out
1. Install the CRDs into the cluster:

```sh
make install
```

2. Run your controller (this will run in the foreground, so switch to a new terminal if you want to leave it running):

```sh
make run
```

**NOTE:** You can also run this in one step by running: `make install run`

### Modifying the API definitions
If you are editing the API definitions, generate the manifests such as CRs or CRDs using:

```sh
make manifests
```

**NOTE:** Run `make --help` for more information on all potential `make` targets

More information can be found via the [Kubebuilder Documentation](https://book.kubebuilder.io/introduction.html)


kubebuilder creation command
--------------------------------------------------------
kubebuilder create api --group cluster.cdx.foc --version v1 --kind ClusterWatchNamespace


mockgen -source=C:\Work\cdx\macer\clustermanager\internal\controller\messenger.go -destination=C:\Work\cdx\macer\clustermanager\mocks\messenger_mock.go -package=mocks

--------------------------------------------------------

 To build 
 -------------------------------------------------------------

 go build -o bin/manager cmd/main.go

 // only run from the docker image
 make install manifest

 // install the crds 
 kustomize build config/crd | kubectl apply -f -

 // create instance of the kind
 kubectl apply -f .\config\samples\cluster_v1_clusterwatchnamespace.yaml -n cdx-system



please use the following lanch.json 
{
    // Use IntelliSense to learn about possible attributes.
    // Hover to view descriptions of existing attributes.
    // For more information, visit: https://go.microsoft.com/fwlink/?linkid=830387
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Launch Package",
            "type": "go",
            "request": "launch",
            "mode": "auto",
            "program": "./clustermanager/cmd/main.go"
        }
    ]
}

kubectl auth can-i get services --as=system:serviceaccount:cdx-system:cmserviceaccount -n cdx-system


kubectl auth can-i list clusterwatchnamespaces --as=system:serviceaccount:cdx-system:cmserviceaccount -n cdx-system

