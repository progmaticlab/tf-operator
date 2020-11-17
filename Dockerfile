ARG LINUX_DISTR=centos
ARG LINUX_DISTR_VER=7
FROM $LINUX_DISTR:$LINUX_DISTR_VER

COPY . /tf-operator
RUN yum install -y dnf wget patch gcc gcc-c++ python3 && \
    wget https://dl.google.com/go/go1.14.2.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.14.2.linux-amd64.tar.gz && \
    rm -f go1.14.2.linux-amd64.tar.gz && \
cat <<EOF >> ~/.bash_profile \
export PATH=\$PATH:/usr/local/go/bin \
EOF && \
curl -LO https://github.com/operator-framework/operator-sdk/releases/download/v0.13.0/operator-sdk-v0.13.0-x86_64-linux-gnu && \
chmod u+x ./operator-sdk-v0.13.0-x86_64-linux-gnu && \
mv ./operator-sdk-v0.13.0-x86_64-linux-gnu /usr/local/bin/operator-sdk && \
dnf install -y dnf-plugins-core && \
dnf copr enable -y vbatts/bazel && \
dnf install -y bazel3