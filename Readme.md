# Requirements
- CentOS 7
- K8s >= 1.16 installed

# Install kubernetes prepared for tf  using kubespray
```bash
git clone https://github.com/tungstenfabric/tf-devstack.git
./tf-devstack/k8s-manifests/run.sh platform
```
# Build tf-operator and CRDs container

```bash
cd tf-operator
./scripts/setup_build_software.sh
# source profile or relogin for add /usr/local/go/bin to the PATH
./scripts/build.sh
```

# Run tf-operator and AIO Tingsten fabric cluster
```bash
./scripts/run_operator.sh
```
