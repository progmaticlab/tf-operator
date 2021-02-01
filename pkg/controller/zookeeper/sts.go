package zookeeper

import (
	"github.com/ghodss/yaml"
	appsv1 "k8s.io/api/apps/v1"
)

var yamlDatazookeeper_sts = `
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: zookeeper
spec:
  selector:
    matchLabels:
      app: zookeeper
  serviceName: "zookeeper"
  replicas: 1
  template:
    metadata:
      labels:
        app: zookeeper
        zookeeper_cr: zookeeper
        contrail_manager: zookeeper
    spec:
      nodeSelector:
        node-role.kubernetes.io/master: ''
      tolerations:
      - operator: Exists
        effect: NoSchedule
      - operator: Exists
        effect: NoExecute
      dnsPolicy: ClusterFirstWithHostNet
      hostNetwork: true
      initContainers:
      - command:
        - sh
        - -c
        - until grep ready /tmp/podinfo/pod_labels > /dev/null 2>&1; do sleep 1; done
        image: busybox
        name: init
        resources: {}
        securityContext:
          privileged: false
          procMount: Default
        terminationMessagePath: /dev/termination-log
        terminationMessagePolicy: File
        volumeMounts:
        - mountPath: /tmp/podinfo
          name: status
      - image: busybox
        command:
          - sh
          - -c
          - cp /zookeeper-conf/* /mnt/zookeeper/ && cp /mnt/zookeeper/zoo.cfg.$POD_IP /mnt/zookeeper/zoo.cfg && sleep 120
        env:
        - name: POD_IP
          valueFrom:
            fieldRef:
              fieldPath: status.podIP
        name: conf-init
        resources: {}
        securityContext:
          privileged: false
          procMount: Default
        terminationMessagePath: /dev/termination-log
        terminationMessagePolicy: File
      containers:
      - name: zookeeper
        env:
        - name: POD_IP
          valueFrom:
            fieldRef:
              fieldPath: status.podIP
        - name: POD_NAME
          valueFrom:
            fieldRef:
              fieldPath: metadata.name
        # TODO: dont provide till 2 DBs be supported
        #- name: NODE_TYPE
        #  value: config-database
        image: tungstenfabric/contrail-external-zookeeper:latest
        resources:
          requests:
            memory: "1Gi"
        startupProbe:
          periodSeconds: 3
          failureThreshold: 30
          exec:
            command:
            - /bin/bash
            - -c
            - "OK=$(echo ruok | nc ${POD_IP} 2181); if [[ ${OK} == \"imok\" ]]; then exit 0; else exit 1;fi"
        readinessProbe:
          initialDelaySeconds: 30
          timeoutSeconds: 3
          failureThreshold: 3
          exec:
            command:
            - /bin/bash
            - -c
            - "OK=$(echo ruok | nc ${POD_IP} 2181); if [[ ${OK} == \"imok\" ]]; then exit 0; else exit 1;fi"
        volumeMounts:
        - mountPath: /tmp/conf
          name: conf
        - mountPath: /var/lib/zookeeper
          name: zookeeper-data
        - mountPath: /var/log/zookeeper
          name: zookeeper-logs
      volumes:
      - name: zookeeper-data
        hostPath:
          path: /var/lib/contrail/zookeeper
      - name: zookeeper-logs
        hostPath:
          path: /var/log/contrail/zookeeper
      - downwardAPI:
          defaultMode: 420
          items:
          - fieldRef:
              apiVersion: v1
              fieldPath: metadata.labels
            path: pod_labels
        name: status
      - downwardAPI:
          defaultMode: 420
          items:
          - fieldRef:
              apiVersion: v1
              fieldPath: metadata.labels
            path: zoo.cfg
        name: conf`

func GetSTS() *appsv1.StatefulSet {
	sts := appsv1.StatefulSet{}
	err := yaml.Unmarshal([]byte(yamlDatazookeeper_sts), &sts)
	if err != nil {
		panic(err)
	}
	jsonData, err := yaml.YAMLToJSON([]byte(yamlDatazookeeper_sts))
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal([]byte(jsonData), &sts)
	if err != nil {
		panic(err)
	}
	return &sts
}
