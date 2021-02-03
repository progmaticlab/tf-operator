package config

import (
	"github.com/ghodss/yaml"
	appsv1 "k8s.io/api/apps/v1"
)

var yamlDataConfigSts = `
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: config
spec:
  selector:
    matchLabels:
      app: config
  serviceName: "config"
  replicas: 1
  template:
    metadata:
      labels:
        app: config
        contrail_manager: config
    spec:
      dnsPolicy: ClusterFirstWithHostNet
      hostNetwork: true
      restartPolicy: Always
      nodeSelector:
        node-role.kubernetes.io/master: ""
      tolerations:
        - effect: NoSchedule
          operator: Exists
        - effect: NoExecute
          operator: Exists
      initContainers:
        - name: init
          image: busybox:latest
          command:
            - sh
            - -c
            - until grep ready /tmp/podinfo/pod_labels > /dev/null 2>&1; do sleep 1; done
          volumeMounts:
            - mountPath: /tmp/podinfo
              name: status
        - name: init2
          image: busybox:latest
          command:
            - sh
            - -c
            - until grep true /tmp/podinfo/peers_ready > /dev/null 2>&1; do sleep 1; done
          volumeMounts:
            - mountPath: /tmp/podinfo
              name: status
      containers:
        - name: api
          image: tungstenfabric/contrail-controller-config-api:latest
          env:
            - name: POD_IP
              valueFrom:
                fieldRef:
                  fieldPath: status.podIP
          #startupProbe:
          #  failureThreshold: 30
          #  periodSeconds: 5
          #  httpGet:
          #    scheme: HTTPS
          #    path: /
          #    port: 8082
          #readinessProbe:
          #  failureThreshold: 3
          #  periodSeconds: 3
          #  httpGet:
          #    scheme: HTTPS
          #    path: /
          #    port: 8082
          volumeMounts:
            - mountPath: /var/log/contrail
              name: config-logs
        - name: devicemanager
          image: tungstenfabric/contrail-controller-config-devicemgr:latest
          env:
            - name: VENDOR_DOMAIN
              value: tungsten.io
            - name: POD_IP
              valueFrom:
                fieldRef:
                  fieldPath: status.podIP
          volumeMounts:
            - mountPath: /var/log/contrail
              name: config-logs
        - name: dnsmasq
          image: tungstenfabric/contrail-external-dnsmasq:latest
          env:
            - name: VENDOR_DOMAIN
              value: tungsten.io
            - name: POD_IP
              valueFrom:
                fieldRef:
                  fieldPath: status.podIP
          volumeMounts:
            - mountPath: /var/log/contrail
              name: config-logs
        - name: schematransformer
          image: tungstenfabric/contrail-controller-config-schema:latest
          env:
            - name: POD_IP
              valueFrom:
                fieldRef:
                  fieldPath: status.podIP
          volumeMounts:
            - mountPath: /var/log/contrail
              name: config-logs
        - name: servicemonitor
          image: tungstenfabric/contrail-controller-config-svcmonitor:latest
          env:
            - name: POD_IP
              valueFrom:
                fieldRef:
                  fieldPath: status.podIP
          volumeMounts:
            - mountPath: /var/log/contrail
              name: config-logs
        - name: analyticsapi
          image: tungstenfabric/contrail-analytics-api:latest
          env:
            - name: POD_IP
              valueFrom:
                fieldRef:
                  fieldPath: status.podIP
            - name: ANALYTICSDB_ENABLE
              value: "true"
            - name: ANALYTICS_ALARM_ENABLE
              value: "true"
          volumeMounts:
            - mountPath: /var/log/contrail
              name: config-logs
        - name: queryengine
          image: tungstenfabric/contrail-analytics-query-engine:latest
          env:
            - name: POD_IP
              valueFrom:
                fieldRef:
                  fieldPath: status.podIP
          volumeMounts:
            - mountPath: /var/log/contrail
              name: config-logs
        - name: collector
          image: tungstenfabric/contrail-analytics-collector:latest
          env:
            - name: POD_IP
              valueFrom:
                fieldRef:
                  fieldPath: status.podIP
          volumeMounts:
            - mountPath: /var/log/contrail
              name: config-logs
        - name: redis
          image: redis:4.0.14
          env:
            - name: POD_IP
              valueFrom:
                fieldRef:
                  fieldPath: status.podIP
          volumeMounts:
            - mountPath: /var/log/contrail
              name: config-logs
            - mountPath: /var/lib/redis
              name: config-data
        - name: nodemanagerconfig
          image: tungstenfabric/contrail-nodemgr:latest
          env:
            - name: VENDOR_DOMAIN
              value: tungsten.io
            - name: NODE_TYPE
              value: config
            - name: POD_IP
              valueFrom:
                fieldRef:
                  fieldPath: status.podIP
            - name: PROVISION_HOSTNAME
              valueFrom:
                fieldRef:
                  fieldPath: metadata.annotations['hostname']
          volumeMounts:
            - mountPath: /var/log/contrail
              name: config-logs
            - mountPath: /var/run
              name: var-run
            - mountPath: /var/crashes
              name: crashes
        - name: nodemanageranalytics
          image: tungstenfabric/contrail-nodemgr:latest
          env:
            - name: VENDOR_DOMAIN
              value: tungsten.io
            - name: NODE_TYPE
              value: analytics
            - name: POD_IP
              valueFrom:
                fieldRef:
                  fieldPath: status.podIP
            - name: PROVISION_HOSTNAME
              valueFrom:
                fieldRef:
                  fieldPath: metadata.annotations['hostname']
          volumeMounts:
            - mountPath: /var/log/contrail
              name: config-logs
            - mountPath: /var/run
              name: var-run
            - mountPath: /var/crashes
              name: crashes
        - name: provisioneranalytics
          image: tungstenfabric/contrail-provisioner:latest
          env:
            - name: NODE_TYPE
              value: analytics
            - name: PROVISION_RETRIES
              value: 1000
            - name: PROVISION_DELAY
              value: 5
            - name: POD_IP
              valueFrom:
                fieldRef:
                  fieldPath: status.podIP
            - name: PROVISION_HOSTNAME
              valueFrom:
                fieldRef:
                  fieldPath: metadata.annotations['hostname']
          volumeMounts:
            - mountPath: /var/log/contrail
              name: config-logs
            - mountPath: /var/crashes
              name: crashes
        - name: provisionerconfig
          image: tungstenfabric/contrail-provisioner:latest
          env:
            - name: NODE_TYPE
              value: config
            - name: PROVISION_RETRIES
              value: 1000
            - name: PROVISION_DELAY
              value: 5
            - name: POD_IP
              valueFrom:
                fieldRef:
                  fieldPath: status.podIP
            - name: PROVISION_HOSTNAME
              valueFrom:
                fieldRef:
                  fieldPath: metadata.annotations['hostname']
          volumeMounts:
            - mountPath: /var/log/contrail
              name: config-logs
            - mountPath: /var/crashes
              name: crashes
      volumes:
        - hostPath:
            path: /var/lib/tftp
            type: ""
          name: tftp
        - hostPath:
            path: /var/lib/dnsmasq
            type: ""
          name: dnsmasq
        - hostPath:
            path: /var/log/contrail/config
            type: ""
          name: config-logs
        - hostPath:
            path: /var/crashes/contrail
            type: ""
          name: crashes
        - hostPath:
            path: /var/lib/contrail/config
            type: ""
          name: config-data
        - hostPath:
            path: /var/run
            type: ""
          name: var-run
        - hostPath:
            path: /usr/local/bin
            type: ""
          name: host-usr-local-bin
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
              path: peers_ready
            - fieldRef:
                apiVersion: v1
                fieldPath: metadata.labels
              path: pod_labelsx
          name: status`

// GetSTS returns StatesfulSet object created from yamlDataConfigSts
func GetSTS() *appsv1.StatefulSet {
	sts := appsv1.StatefulSet{}
	err := yaml.Unmarshal([]byte(yamlDataConfigSts), &sts)
	if err != nil {
		panic(err)
	}
	jsonData, err := yaml.YAMLToJSON([]byte(yamlDataConfigSts))
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal([]byte(jsonData), &sts)
	if err != nil {
		panic(err)
	}
	return &sts
}
