# Setup kubernetes
```bash
git clone https://github.com/tungstenfabric/tf-devstack.git
./tf-devstack/k8s_manifests/run.sh platform
```

# Setup build env

```bash
cd contrail-operator
sudo usermod -a -G docker centos
# relogin here to use docker without sudo
scripts/setup_build_sofware.sh
python3 -m venv ~/env
source ~/env/bin/activate
```

# Build containers
```bash
scripts/build_containers_bazel.sh
```
# Run contrail AIO
```bash
kubectl apply -f deploy/kustomize/operator/latest
kubectl apply -k deploy/kustomize/contrail/1node/latest
```
