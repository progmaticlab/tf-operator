# Setup build env

```bash
cd contrail-operator
scripts/setup_docker.sh
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

