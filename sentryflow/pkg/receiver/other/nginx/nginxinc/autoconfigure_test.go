// SPDX-License-Identifier: Apache-2.0
// Copyright 2024 Authors of SentryFlow

package nginxinc

import (
	"context"
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/5GSEC/SentryFlow/pkg/config"
)

func Test_autoConfigure(t *testing.T) {
	type args struct {
		cfg *config.Config
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "when `autoconfigure` is true for `nginx-inc-ingress-controller` then return true",
			args: args{
				cfg: &config.Config{
					Filters: &config.Filters{
						NginxIngress: &config.NginxIngressConfig{
							DeploymentName:             "nginx-inc-ingress-controller",
							ConfigMapName:              "nginx-ingress",
							SentryFlowNjsConfigMapName: "sentryflow-nginx-inc",
						},
						Server: &config.Server{Port: 9999},
					},
					Receivers: &config.Receivers{
						Others: []*config.NameAndNamespace{
							{
								Name:          "nginx-inc-ingress-controller",
								Namespace:     "nginx-ingress",
								AutoConfigure: true,
							},
						},
					},
					Exporter: &config.ExporterConfig{
						Grpc: &config.Server{Port: 8888},
					},
				},
			},
			want: true,
		},
		{
			name: "when `autoconfigure` is not set for `nginx-inc-ingress-controller` then return false",
			args: args{
				cfg: &config.Config{
					Filters: &config.Filters{
						NginxIngress: &config.NginxIngressConfig{
							DeploymentName:             "nginx-inc-ingress-controller",
							ConfigMapName:              "nginx-ingress",
							SentryFlowNjsConfigMapName: "sentryflow-nginx-inc",
						},
						Server: &config.Server{Port: 9999},
					},
					Receivers: &config.Receivers{
						Others: []*config.NameAndNamespace{
							{
								Name:      "nginx-inc-ingress-controller",
								Namespace: "nginx-ingress",
							},
						},
					},
					Exporter: &config.ExporterConfig{
						Grpc: &config.Server{Port: 8888},
					},
				},
			},
			want: false,
		},
		{
			name: "when `autoconfigure` is true for other receiver then return false",
			args: args{
				cfg: &config.Config{
					Filters: &config.Filters{
						Server: &config.Server{Port: 9999},
					},
					Receivers: &config.Receivers{
						Others: []*config.NameAndNamespace{
							{
								Name:          "nginx-webserver",
								AutoConfigure: true,
							},
						},
					},
					Exporter: &config.ExporterConfig{
						Grpc: &config.Server{Port: 8888},
					},
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := autoConfigure(tt.args.cfg); got != tt.want {
				t.Errorf("autoConfigure() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_containsVolumeAndVolumeMount(t *testing.T) {
	type args struct {
		spec v1.PodSpec
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "when both volume and volumeMount path and subPath are present then return true",
			args: args{
				spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							VolumeMounts: []v1.VolumeMount{
								{
									Name:      volumeName,
									MountPath: "/etc/nginx/njs/sentryflow.js",
									SubPath:   "sentryflow.js",
								},
							},
						},
					},
					Volumes: []v1.Volume{
						{
							Name: volumeName,
							VolumeSource: v1.VolumeSource{
								ConfigMap: &v1.ConfigMapVolumeSource{
									LocalObjectReference: v1.LocalObjectReference{
										Name: volumeName,
									},
								},
							},
						},
					},
				},
			},
			want: true,
		},
		{
			name: "when only volume is present then return false",
			args: args{
				spec: v1.PodSpec{
					Volumes: []v1.Volume{
						{
							Name: volumeName,
							VolumeSource: v1.VolumeSource{
								ConfigMap: &v1.ConfigMapVolumeSource{
									LocalObjectReference: v1.LocalObjectReference{
										Name: volumeName,
									},
								},
							},
						},
					},
				},
			},
			want: false,
		},
		{
			name: "when only volumeMount is present then return false",
			args: args{
				spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							VolumeMounts: []v1.VolumeMount{
								{
									Name:      volumeName,
									MountPath: "/etc/nginx/njs/sentryflow.js",
									SubPath:   "sentryflow.js",
								},
							},
						},
					},
				},
			},
			want: false,
		},
		{
			name: "when both volume and volumeMount are not present then return false",
			args: args{
				spec: v1.PodSpec{
					Containers: []v1.Container{},
					Volumes:    []v1.Volume{},
				},
			},
			want: false,
		},
		{
			name: "when incorrect volumeName is present then return false",
			args: args{
				spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							VolumeMounts: []v1.VolumeMount{
								{
									Name:      volumeName,
									MountPath: "/etc/nginx/njs/sentryflow.js",
									SubPath:   "sentryflow.js",
								},
							},
						},
					},
					Volumes: []v1.Volume{
						{
							Name: "other-volume",
							VolumeSource: v1.VolumeSource{
								ConfigMap: &v1.ConfigMapVolumeSource{
									LocalObjectReference: v1.LocalObjectReference{
										Name: volumeName,
									},
								},
							},
						},
					},
				},
			},
			want: false,
		},
		{
			name: "when incorrect volumeMount path is present then return false",
			args: args{
				spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							VolumeMounts: []v1.VolumeMount{
								{
									Name:      volumeName,
									MountPath: "/etc/nginx/njs/other.js",
									SubPath:   "sentryflow.js",
								},
							},
						},
					},
					Volumes: []v1.Volume{
						{
							Name: volumeName,
							VolumeSource: v1.VolumeSource{
								ConfigMap: &v1.ConfigMapVolumeSource{
									LocalObjectReference: v1.LocalObjectReference{
										Name: volumeName,
									},
								},
							},
						},
					},
				},
			},
			want: false,
		},
		{
			name: "when no volumeMount subPath is present then return false",
			args: args{
				spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							VolumeMounts: []v1.VolumeMount{
								{
									Name:      volumeName,
									MountPath: "/etc/nginx/njs/other.js",
								},
							},
						},
					},
					Volumes: []v1.Volume{
						{
							Name: volumeName,
							VolumeSource: v1.VolumeSource{
								ConfigMap: &v1.ConfigMapVolumeSource{
									LocalObjectReference: v1.LocalObjectReference{
										Name: volumeName,
									},
								},
							},
						},
					},
				},
			},
			want: false,
		},
		{
			name: "when incorrect volumeMount subPath is present then return false",
			args: args{
				spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							VolumeMounts: []v1.VolumeMount{
								{
									Name:      volumeName,
									MountPath: "/etc/nginx/njs/other.js",
									SubPath:   "other.js",
								},
							},
						},
					},
					Volumes: []v1.Volume{
						{
							Name: volumeName,
							VolumeSource: v1.VolumeSource{
								ConfigMap: &v1.ConfigMapVolumeSource{
									LocalObjectReference: v1.LocalObjectReference{
										Name: volumeName,
									},
								},
							},
						},
					},
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := containsVolumeAndVolumeMount(tt.args.spec); got != tt.want {
				t.Errorf("containsVolumeAndVolumeMount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_deployResources(t *testing.T) {
	ctx := context.Background()
	k8sClient := getFakeClient()
	sentryFlowSvc := sentryFlowService()
	cfg := configWithAutoConfigureTrue()

	if err := k8sClient.Create(ctx, sentryFlowSvc); err != nil {
		t.Errorf("failed to create sentryflow service: %v", err)
	}

	t.Run("when ingress deployment and configmap don't contain patches then patch them and return nil", func(t *testing.T) {
		ingDeploy := ingressDeployment()
		ingCm := ingressConfigMap()
		wantErr := false

		if err := k8sClient.Create(ctx, ingDeploy); err != nil {
			t.Errorf("failed to create ingress deployment: %v", err)
		}
		if err := k8sClient.Create(ctx, ingCm); err != nil {
			t.Errorf("failed to create ingress configmap: %v", err)
		}
		defer cleanup(t, k8sClient, ingDeploy, ingCm, nil)

		if err := deployResources(ctx, k8sClient, cfg); (err != nil) != wantErr {
			t.Errorf("deployResources() error = %v, wantErr %v", err, wantErr)
		}
	})

	t.Run("when ingress deployment and configmap contain patches then don't patch them and return nil", func(t *testing.T) {
		ingDeploy := ingressDeployment()
		ingDeploy.Spec.Template.Spec.Containers[0].VolumeMounts = append(ingDeploy.Spec.Template.Spec.Containers[0].VolumeMounts, v1.VolumeMount{
			Name:      volumeName,
			MountPath: "/etc/nginx/njs/sentryflow.js",
			SubPath:   "sentryflow.js",
		})
		ingDeploy.Spec.Template.Spec.Volumes = append(ingDeploy.Spec.Template.Spec.Volumes, v1.Volume{
			Name: volumeName,
			VolumeSource: v1.VolumeSource{
				ConfigMap: &v1.ConfigMapVolumeSource{
					LocalObjectReference: v1.LocalObjectReference{
						Name: volumeName,
					},
				},
			},
		})
		ingCm := ingressConfigMap()
		patchConfigMap(ingCm)

		wantErr := false

		if err := k8sClient.Create(ctx, ingDeploy); err != nil {
			t.Errorf("failed to create ingress deployment: %v", err)
		}
		if err := k8sClient.Create(ctx, ingCm); err != nil {
			t.Errorf("failed to create ingress configmap: %v", err)
		}
		defer cleanup(t, k8sClient, ingDeploy, ingCm, nil)

		if err := deployResources(ctx, k8sClient, cfg); (err != nil) != wantErr {
			t.Errorf("deployResources() error = %v, wantErr %v", err, wantErr)
		}
	})

	t.Run("when ingress deployment doesn't exist then return doesn't exist error", func(t *testing.T) {
		ingCm := ingressConfigMap()
		wantErr := true
		errMessage := `deployments.apps "nginx-ingress" not found`

		if err := k8sClient.Create(ctx, ingCm); err != nil {
			t.Errorf("failed to create ingress configmap: %v", err)
		}
		defer cleanup(t, k8sClient, nil, ingCm, nil)

		err := deployResources(ctx, k8sClient, cfg)
		if (err != nil) != wantErr {
			t.Errorf("deployResources() error = %v, wantErr %v", err, wantErr)
		}
		if err.Error() != errMessage {
			t.Errorf("deployResources() errorMessage = %v, wantErrMessage %v", err.Error(), errMessage)
		}
	})

	t.Run("when ingress configMap doesn't exist then return doesn't exist error", func(t *testing.T) {
		ingressDeploy := ingressDeployment()
		wantErr := true
		errMessage := `configmaps "nginx-ingress" not found`

		if err := k8sClient.Create(ctx, ingressDeploy); err != nil {
			t.Errorf("failed to create ingress configmap: %v", err)
		}
		defer cleanup(t, k8sClient, ingressDeploy, nil, nil)

		err := deployResources(ctx, k8sClient, cfg)
		if (err != nil) != wantErr {
			t.Errorf("deployResources() error = %v, wantErr %v", err, wantErr)
		}
		if err.Error() != errMessage {
			t.Errorf("deployResources() errorMessage = %v, wantErrMessage %v", err.Error(), errMessage)
		}
	})
}

func Test_removePatchFromIngressConfigMap(t *testing.T) {
	ctx := context.Background()
	k8sClient := getFakeClient()

	t.Run("when ingress configmap doesn't contain any sentryflow patches then return nil", func(t *testing.T) {
		cm := ingressConfigMap()
		cfg := configWithAutoConfigureTrue()
		wantErr := false

		if err := k8sClient.Create(ctx, cm); err != nil {
			t.Errorf("failed to create ingress configmap: %v", err)
		}
		defer cleanup(t, k8sClient, nil, cm, nil)

		if err := removePatchFromIngressConfigMap(ctx, k8sClient, cfg, cm.Namespace); (err != nil) != wantErr {
			t.Errorf("removePatchFromIngressConfigMap() error = %v, wantErr %v", err, wantErr)
		}
	})

	t.Run("when ingress configmap contains sentryflow patches then remove then and return nil", func(t *testing.T) {
		cm := ingressConfigMap()
		patchConfigMap(cm)
		cfg := configWithAutoConfigureTrue()
		wantErr := false
		if err := k8sClient.Create(ctx, cm); err != nil {
			t.Errorf("failed to create ingress configmap: %v", err)
		}
		defer cleanup(t, k8sClient, nil, cm, nil)

		if err := removePatchFromIngressConfigMap(ctx, k8sClient, cfg, cm.Namespace); (err != nil) != wantErr {
			t.Errorf("removePatchFromIngressConfigMap() error = %v, wantErr %v", err, wantErr)
		}

		var updatedCm corev1.ConfigMap
		if err := k8sClient.Get(ctx, client.ObjectKey{Namespace: cm.Namespace, Name: cm.Name}, &updatedCm); err != nil {
			t.Errorf("failed to get updated ingress configmap: %v", err)
		}

		if snippets, exists := updatedCm.Data["http-snippets"]; exists {
			t.Errorf("removePatchFromIngressConfigMap() got http-snippets = %v,\nwant = %v", snippets, "")
		}
		if snippets, exists := updatedCm.Data["location-snippets"]; exists {
			t.Errorf("removePatchFromIngressConfigMap() got location-snippets = %v,\nwant = %v", snippets, "")
		}
		if snippets, exists := updatedCm.Data["server-snippets"]; exists {
			t.Errorf("removePatchFromIngressConfigMap() got server-snippets = %v,\nwant =  %v", snippets, "")
		}
	})

}

func Test_removePatchFromIngressDeploy(t *testing.T) {
	ctx := context.Background()
	k8sClient := getFakeClient()

	t.Run("when ingress deployment doesn't contain patches then return nil", func(t *testing.T) {
		deploy := ingressDeployment()
		if err := k8sClient.Create(ctx, deploy); err != nil {
			t.Errorf("failed to create ingress deployment: %v", err)
		}
		defer cleanup(t, k8sClient, deploy, nil, nil)

		wantErr := false

		if err := removePatchFromIngressDeploy(ctx, k8sClient, deploy.Name, deploy.Namespace); (err != nil) != wantErr {
			t.Errorf("removePatchFromIngressDeploy() error = %v, wantErr %v", err, wantErr)
		}
	})

	t.Run("when ingress deployment contains patches then remove them and return nil", func(t *testing.T) {
		ingDeploy := ingressDeployment()
		ingDeploy.Spec.Template.Spec.Containers[0].VolumeMounts = append(ingDeploy.Spec.Template.Spec.Containers[0].VolumeMounts, v1.VolumeMount{
			Name:      volumeName,
			MountPath: "/etc/nginx/njs/sentryflow.js",
			SubPath:   "sentryflow.js",
		})
		ingDeploy.Spec.Template.Spec.Volumes = append(ingDeploy.Spec.Template.Spec.Volumes, v1.Volume{
			Name: volumeName,
			VolumeSource: v1.VolumeSource{
				ConfigMap: &v1.ConfigMapVolumeSource{
					LocalObjectReference: v1.LocalObjectReference{
						Name: volumeName,
					},
				},
			},
		})
		if err := k8sClient.Create(ctx, ingDeploy); err != nil {
			t.Errorf("failed to create ingress deployment: %v", err)
		}
		defer cleanup(t, k8sClient, ingDeploy, nil, nil)
		wantErr := false

		if err := removePatchFromIngressDeploy(ctx, k8sClient, ingDeploy.Name, ingDeploy.Namespace); (err != nil) != wantErr {
			t.Errorf("removePatchFromIngressDeploy() error = %v, wantErr %v", err, wantErr)
		}

		var updatedDeploy appsv1.Deployment
		if err := k8sClient.Get(ctx, client.ObjectKey{Namespace: ingDeploy.Namespace, Name: ingDeploy.Name}, &updatedDeploy); err != nil {
			t.Errorf("failed to get updated deployment: %v", err)
		}

		if len(updatedDeploy.Spec.Template.Spec.Containers[0].VolumeMounts) != 0 {
			t.Errorf("removePatchFromIngressDeploy() volumeMounts = %v, want 0", len(updatedDeploy.Spec.Template.Spec.Containers[0].VolumeMounts))
		}
		if len(updatedDeploy.Spec.Template.Spec.Volumes) != 0 {
			t.Errorf("removePatchFromIngressDeploy() volumes = %v, want 0", len(updatedDeploy.Spec.Template.Spec.Volumes))
		}
	})

	t.Run("when ingress deployment contains multiple volumes and volumeMounts then remove patched ones and return nil", func(t *testing.T) {
		ingDeploy := ingressDeployment()
		ingDeploy.Spec.Template.Spec.Containers[0].VolumeMounts = append(ingDeploy.Spec.Template.Spec.Containers[0].VolumeMounts,
			v1.VolumeMount{
				Name:      volumeName,
				MountPath: "/etc/nginx/njs/sentryflow.js",
				SubPath:   "sentryflow.js",
			},
			v1.VolumeMount{
				Name:      "other-volume",
				MountPath: "/some/mount/path",
				ReadOnly:  true,
			},
		)
		ingDeploy.Spec.Template.Spec.Volumes = append(ingDeploy.Spec.Template.Spec.Volumes,
			v1.Volume{
				Name: volumeName,
				VolumeSource: v1.VolumeSource{
					ConfigMap: &v1.ConfigMapVolumeSource{
						LocalObjectReference: v1.LocalObjectReference{
							Name: volumeName,
						},
					},
				},
			},
			v1.Volume{
				Name: "other-volume",
				VolumeSource: v1.VolumeSource{
					EmptyDir: &v1.EmptyDirVolumeSource{},
				},
			},
		)
		if err := k8sClient.Create(ctx, ingDeploy); err != nil {
			t.Errorf("failed to create ingress deployment: %v", err)
		}
		defer cleanup(t, k8sClient, ingDeploy, nil, nil)
		wantErr := false

		if err := removePatchFromIngressDeploy(ctx, k8sClient, ingDeploy.Name, ingDeploy.Namespace); (err != nil) != wantErr {
			t.Errorf("removePatchFromIngressDeploy() error = %v, wantErr %v", err, wantErr)
		}

		var updatedDeploy appsv1.Deployment
		if err := k8sClient.Get(ctx, client.ObjectKey{Namespace: ingDeploy.Namespace, Name: ingDeploy.Name}, &updatedDeploy); err != nil {
			t.Errorf("failed to get updated deployment: %v", err)
		}

		if len(updatedDeploy.Spec.Template.Spec.Containers[0].VolumeMounts) != 1 {
			t.Errorf("removePatchFromIngressDeploy() volumeMounts = %v, want 1", len(updatedDeploy.Spec.Template.Spec.Containers[0].VolumeMounts))
		}
		if len(updatedDeploy.Spec.Template.Spec.Volumes) != 1 {
			t.Errorf("removePatchFromIngressDeploy() volumes = %v, want 1", len(updatedDeploy.Spec.Template.Spec.Volumes))
		}
	})
}

func Test_sentryFlowSvcNameAndNs(t *testing.T) {
	ctx := context.Background()
	k8sClient := getFakeClient()

	t.Run("when sentrflow doesn't exist then return empty name and namespace and doesn't exist error", func(t *testing.T) {
		wantErr := true
		svcName, namespace := "", ""
		errMessage := "sentryFlow svc was not found"

		gotSvcName, gotNsName, err := sentryFlowSvcNameAndNs(ctx, k8sClient)
		if (err != nil) != wantErr {
			t.Errorf("sentryFlowSvcNameAndNs() error = %v, wantErr %v", err, wantErr)
		}
		if err.Error() != errMessage {
			t.Errorf("sentryFlowSvcNameAndNs() errorMessage = %v, wantErrMessage %v", err.Error(), errMessage)
		}
		if gotSvcName != svcName {
			t.Errorf("sentryFlowSvcNameAndNs() got = %v, want %v", gotSvcName, svcName)
		}
		if gotNsName != namespace {
			t.Errorf("sentryFlowSvcNameAndNs() got1 = %v, want %v", gotNsName, namespace)
		}
	})

	t.Run("when sentrflow exist then return its name and namespace and no error", func(t *testing.T) {
		wantErr := false
		sentryFlowSvc := sentryFlowService()

		if err := k8sClient.Create(ctx, sentryFlowSvc); err != nil {
			t.Errorf("failed to create sentry flow service: %v", err)
		}
		defer cleanup(t, k8sClient, nil, nil, sentryFlowSvc)

		gotSvcName, gotNsName, err := sentryFlowSvcNameAndNs(ctx, k8sClient)
		if (err != nil) != wantErr {
			t.Errorf("sentryFlowSvcNameAndNs() error = %v, wantErr %v", err, wantErr)
		}
		if gotSvcName != sentryFlowSvc.Name {
			t.Errorf("sentryFlowSvcNameAndNs() got = %v, want %v", gotSvcName, sentryFlowSvc.Name)
		}
		if gotNsName != sentryFlowSvc.Namespace {
			t.Errorf("sentryFlowSvcNameAndNs() got1 = %v, want %v", gotNsName, sentryFlowSvc.Namespace)
		}
	})

	t.Run("when sentrflow exist with different labels then return empty name and namespace and doesn't exist error ", func(t *testing.T) {
		wantErr := true
		sentryFlowSvc := sentryFlowService()
		sentryFlowSvc.Labels = map[string]string{
			"app": "sentryflow",
		}
		errMessage := "sentryFlow svc was not found"

		if err := k8sClient.Create(ctx, sentryFlowSvc); err != nil {
			t.Errorf("failed to create sentry flow service: %v", err)
		}
		defer cleanup(t, k8sClient, nil, nil, sentryFlowSvc)

		gotSvcName, gotNsName, err := sentryFlowSvcNameAndNs(ctx, k8sClient)
		if (err != nil) != wantErr {
			t.Errorf("sentryFlowSvcNameAndNs() error = %v, wantErr %v", err, wantErr)
		}
		if err.Error() != errMessage {
			t.Errorf("sentryFlowSvcNameAndNs() errorMessage = %v, wantErrMessage %v", err, errMessage)
		}
		if gotSvcName != "" {
			t.Errorf("sentryFlowSvcNameAndNs() got = %v, want %v", gotSvcName, "")
		}
		if gotNsName != "" {
			t.Errorf("sentryFlowSvcNameAndNs() got1 = %v, want %v", gotNsName, "")
		}
	})
}

func getFakeClient() client.WithWatch {
	scheme := runtime.NewScheme()
	utilruntime.Must(corev1.AddToScheme(scheme))
	utilruntime.Must(appsv1.AddToScheme(scheme))
	return fake.
		NewClientBuilder().
		WithScheme(scheme).
		Build()
}

func cleanup(t *testing.T, k8sClient client.WithWatch, deployment *appsv1.Deployment, configMap *v1.ConfigMap, svc *v1.Service) {
	t.Helper()

	ctx := context.Background()
	if deployment != nil {
		if err := k8sClient.Delete(ctx, deployment); err != nil {
			t.Errorf("failed to delete deployment: %v", err)
		}
	}
	if configMap != nil {
		if err := k8sClient.Delete(ctx, configMap); err != nil {
			t.Errorf("failed to delete configmap: %v", err)
		}
	}
	if svc != nil {
		if err := k8sClient.Delete(ctx, svc); err != nil {
			t.Errorf("failed to delete service: %v", err)
		}
	}
}

func sentryFlowService() *corev1.Service {
	return &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "sentryflow",
			Namespace: "sentryflow",
			Labels: map[string]string{
				"app.kubernetes.io/name": "sentryflow",
			},
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:     "http",
					Port:     9999,
					Protocol: corev1.ProtocolTCP,
				},
				{
					Name:     "grpc",
					Port:     8888,
					Protocol: corev1.ProtocolTCP,
				},
			},
			Selector: map[string]string{
				"app.kubernetes.io/name": "sentryflow",
			},
			Type: corev1.ServiceTypeClusterIP,
		},
	}
}

func ingressConfigMap() *corev1.ConfigMap {
	return &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "nginx-ingress",
			Namespace: "nginx-ingress",
		},
	}
}

func configWithAutoConfigureTrue() *config.Config {
	return &config.Config{
		Filters: &config.Filters{
			NginxIngress: &config.NginxIngressConfig{
				DeploymentName:             "nginx-ingress",
				ConfigMapName:              "nginx-ingress",
				SentryFlowNjsConfigMapName: volumeName,
			},
			Server: &config.Server{
				Port: 9999,
			},
		},
		Receivers: &config.Receivers{
			Others: []*config.NameAndNamespace{
				{
					Name:          "nginx-inc-ingress-controller",
					Namespace:     "nginx-ingress",
					AutoConfigure: true,
				},
			},
		},
		Exporter: &config.ExporterConfig{
			Grpc: &config.Server{Port: 8888},
		},
	}
}

func ingressDeployment() *appsv1.Deployment {
	return &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "nginx-ingress",
			Namespace: "nginx-ingress",
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "nginx-ingress",
				},
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "nginx-ingress",
					},
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:  "nginx-ingress",
							Image: "nginx/nginx-ingress:4.0.0",
						},
					},
				},
			},
		},
	}
}

func patchConfigMap(cm *v1.ConfigMap) {
	cm.Data = map[string]string{
		"http-snippets": `js_path "/etc/nginx/njs/";
subrequest_output_buffer_size 8k;
js_shared_dict_zone zone=apievents:1M timeout=300s evict;
js_import main from sentryflow.js;`,

		"location-snippets": `js_body_filter main.requestHandler buffer_type=buffer;
mirror      /mirror_request;
mirror_request_body on;`,

		"server-snippets": `
location /mirror_request {
  internal;
  js_content main.dispatchHttpCall;
}

location /sentryflow {
  internal;
  # Update SentryFlow URL with path to ingest access logs if required.
  proxy_pass http://{{ .SentryFlowSvcName }}.{{ .SentryFlowSvcNamespace }}:{{ .SentryFlowFilterServerPort }}/api/v1/events;
  proxy_method      POST;
  proxy_set_header accept "application/json";
  proxy_set_header Content-Type "application/json";
}`,
	}
}
