# Setup build env

Run these commands:

'''bash
scripts/setup_docker.sh
sudo usermod -a -G docker centos
# relogin to use docker without sudo
scripts/setup_build_softfare.sh
python3 -m venv ~/env
source ~/env/bin/activate
'''

# Build containers
'''bash
scripts/build_containers_bazel.sh
'''

