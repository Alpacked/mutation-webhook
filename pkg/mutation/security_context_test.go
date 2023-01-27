package mutation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestInjectSecurityContext(t *testing.T) {
	want := &corev1.Pod{
		ObjectMeta: v1.ObjectMeta{
			Name: "test",
		},
		Spec: corev1.PodSpec{
			SecurityContext: &corev1.PodSecurityContext{
				RunAsUser:    &[]int64{1000}[0],
				RunAsGroup:   &[]int64{1000}[0],
				RunAsNonRoot: &[]bool{true}[0],
			},
			Containers: []corev1.Container{{
				Name: "test",
				SecurityContext: &corev1.SecurityContext{
					ReadOnlyRootFilesystem: &[]bool{true}[0],
					AllowPrivilegeEscalation: &[]bool{false}[0],
					RunAsNonRoot: &[]bool{true}[0],
					Capabilities: &corev1.Capabilities{
						Drop: []corev1.Capability{"ALL"},
					},
				},
			}},
			InitContainers: []corev1.Container{{
				Name: "inittest",
				SecurityContext: &corev1.SecurityContext{
					ReadOnlyRootFilesystem: &[]bool{true}[0],
					AllowPrivilegeEscalation: &[]bool{false}[0],
					RunAsNonRoot: &[]bool{true}[0],
					Capabilities: &corev1.Capabilities{
						Drop: []corev1.Capability{"ALL"},
					},
				},
			}},
		},
	}

	pod := &corev1.Pod{
		ObjectMeta: v1.ObjectMeta{
			Name: "test",
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{{
				Name: "test",
			}},
			InitContainers: []corev1.Container{{
				Name: "inittest",
			}},
		},
	}

	got, err := injectSecurityContext{Logger: logger()}.Mutate(pod)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, want, got)
}

func TestSkipInjectSecurityContext(t *testing.T) {
	want := &corev1.Pod{
		ObjectMeta: v1.ObjectMeta{
			Name: "test",
		},
		Spec: corev1.PodSpec{
			SecurityContext: &corev1.PodSecurityContext{
				RunAsGroup: &[]int64{1000}[0],
			},
			Containers: []corev1.Container{{
				Name: "test",
				SecurityContext: &corev1.SecurityContext{
					ReadOnlyRootFilesystem: &[]bool{true}[0],
					AllowPrivilegeEscalation: &[]bool{false}[0],
				},
			}},
			InitContainers: []corev1.Container{{
				Name: "inittest",
				SecurityContext: &corev1.SecurityContext{
					Capabilities: &corev1.Capabilities{
						Drop: []corev1.Capability{"ALL"},
					},
				},
			}},
		},
	}

	pod := &corev1.Pod{
		ObjectMeta: v1.ObjectMeta{
			Name: "test",
		},
		Spec: corev1.PodSpec{
			SecurityContext: &corev1.PodSecurityContext{
				RunAsGroup: &[]int64{1000}[0],
			},
			Containers: []corev1.Container{{
				Name: "test",
				SecurityContext: &corev1.SecurityContext{
					ReadOnlyRootFilesystem: &[]bool{true}[0],
					AllowPrivilegeEscalation: &[]bool{false}[0],
				},
			}},
			InitContainers: []corev1.Container{{
				Name: "inittest",
				SecurityContext: &corev1.SecurityContext{
					Capabilities: &corev1.Capabilities{
						Drop: []corev1.Capability{"ALL"},
					},
				},
			}},
		},
	}

	got, err := injectSecurityContext{Logger: logger()}.Mutate(pod)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, want, got)
}
