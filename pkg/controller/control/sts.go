package control

import (
	"github.com/ghodss/yaml"
	appsv1 "k8s.io/api/apps/v1"
)

var yamlDatacontrol_sts = `
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: control
spec:
  selector:
    matchLabels:
      app: control
  serviceName: "control"
  replicas: 1
  template:
    metadata:
      labels:
        app: control
        contrail_manager: control
    spec:
      securityContext:
        fsGroup: 1999
      initContainers:
        - name: init
          image: busybox:latest
          command:
            - sh
            - -c
            - until grep ready /tmp/podinfo/pod_labels > /dev/null 2>&1; do sleep 1; done
          env:
            - name: POD_IP
              valueFrom:
                fieldRef:
                  fieldPath: status.podIP
          imagePullPolicy: Always
          volumeMounts:
            - mountPath: /tmp/podinfo
              name: status
      containers:
        - name: control
          image: tungstenfabric/contrail-controller-control-control:latest
          env:
            - name: POD_IP
              valueFrom:
                fieldRef:
                  fieldPath: status.podIP
          imagePullPolicy: Always
          volumeMounts:
            - mountPath: /var/log/contrail
              name: control-logs
        - name: dns
          image: tungstenfabric/contrail-controller-control-dns:latest
          env:
            - name: POD_IP
              valueFrom:
                fieldRef:
                  fieldPath: status.podIP
          imagePullPolicy: Always
          volumeMounts:
            - mountPath: /var/log/contrail
              name: control-logs
            - mountPath: /etc/contrail
              name: etc-contrail
            - mountPath: /etc/contrail/dns
              name: etc-contrail-dns
        - name: named
          image: tungstenfabric/contrail-controller-control-named:latest
          env:
            - name: POD_IP
              valueFrom:
                fieldRef:
                  fieldPath: status.podIP
          imagePullPolicy: Always
          securityContext:
            privileged: true
            runAsGroup: 1999
          volumeMounts:
            - mountPath: /var/log/contrail
              name: control-logs
            - mountPath: /etc/contrail
              name: etc-contrail
            - mountPath: /etc/contrail/dns
              name: etc-contrail-dns
        - name: nodemanager
          image: tungstenfabric/contrail-nodemgr:latest
          env:
            - name: VENDOR_DOMAIN
              value: tungsten.io
            - name: NODE_TYPE
              value: control
            - name: POD_IP
              valueFrom:
                fieldRef:
                  fieldPath: status.podIP
            - name: CONTROL_HOSTNAME
              valueFrom:
                fieldRef:
                  fieldPath: metadata.annotations['hostname']
          imagePullPolicy: Always
          volumeMounts:
            - mountPath: /var/log/contrail
              name: control-logs
            - mountPath: /var/crashes
              name: crashes
            - mountPath: /mnt
              name: var-run
        - name: provisioner
          image: tungstenfabric/contrail-provisioner:latest
          env:
            - name: NODE_TYPE
              value: control
            - name: POD_IP
              valueFrom:
                fieldRef:
                  fieldPath: status.podIP
            - name: CONTROL_HOSTNAME
              valueFrom:
                fieldRef:
                  fieldPath: metadata.annotations['hostname']
          imagePullPolicy: Always
          lifecycle:
            preStop:
              exec:
                command:
                  - python /etc/contrailconfigmaps/deprovision.py.${POD_IP}
          volumeMounts:
            - mountPath: /var/log/contrail
              name: control-logs
            - mountPath: /var/crashes
              name: crashes
            - mountPath: /var/run
              name: var-run
      dnsPolicy: ClusterFirst
      hostNetwork: true
      nodeSelector:
        node-role.kubernetes.io/master: ""
      restartPolicy: Always
      tolerations:
        - effect: NoSchedule
          operator: Exists
        - effect: NoExecute
          operator: Exists
      volumes:
        - hostPath:
            path: /var/log/contrail/control
            type: ""
          name: control-logs
        - hostPath:
            path: /var/contrail/crashes
            type: ""
          name: crashes
        - hostPath:
            path: /var/run
            type: ""
          name: var-run
        - hostPath:
            path: /usr/local/bin
            type: ""
          name: host-usr-local-bin
        - emptyDir: {}
          name: etc-contrail
        - emptyDir: {}
          name: etc-contrail-dns
        - downwardAPI:
            defaultMode: 420
            items:
              - fieldRef:
                  apiVersion: v1
                  fieldPath: metadata.labels
                path: pod_labels
              - fieldRef:
                  apiVersion: v1
                  fieldPath: metadata.labels
                path: pod_labelsx
          name: status`

// GetSTS returns StatesfulSet object created from yamlDatacontrol_sts
func GetSTS() *appsv1.StatefulSet {
	sts := appsv1.StatefulSet{}
	err := yaml.Unmarshal([]byte(yamlDatacontrol_sts), &sts)
	if err != nil {
		panic(err)
	}
	jsonData, err := yaml.YAMLToJSON([]byte(yamlDatacontrol_sts))
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal([]byte(jsonData), &sts)
	if err != nil {
		panic(err)
	}
	return &sts
}
