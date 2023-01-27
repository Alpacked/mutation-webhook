package mutation

import (
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
)

// injectSecurityContext is a container for the mutation injecting security context
type injectSecurityContext struct {
	Logger logrus.FieldLogger
}

// injectSecurityContext implements the podMutator interface
var _ podMutator = (*injectSecurityContext)(nil)

// Name returns the injectSecurityContext short name
func (mpl injectSecurityContext) Name() string {
	return "inj_sec_context"
}

// Mutate returns a new mutated pod according to preconfigured security context
func (mpl injectSecurityContext) Mutate(pod *corev1.Pod) (*corev1.Pod, error) {

	mpl.Logger = mpl.Logger.WithField("mutation", mpl.Name())
	mpod := pod.DeepCopy()

	mpl.Logger.WithField("p_sc_before", mpod.Spec.SecurityContext).Printf("pod sc before mutation")
	mpl.Logger.WithField("c_sc_before", mpod.Spec.Containers[0].SecurityContext).Printf("con sc before mutation")

	if IsPodSecurityContextEmpty(mpod.Spec.SecurityContext) {
		mpod.Spec.SecurityContext = &corev1.PodSecurityContext{
			RunAsUser:    &[]int64{1000}[0],
			RunAsGroup:   &[]int64{1000}[0],
			RunAsNonRoot: &[]bool{true}[0],
		}
	}

	for i := range mpod.Spec.Containers {
		if mpod.Spec.Containers[i].SecurityContext == nil {
			mpod.Spec.Containers[i].SecurityContext = &corev1.SecurityContext{
				ReadOnlyRootFilesystem: &[]bool{true}[0],
				AllowPrivilegeEscalation: &[]bool{false}[0],
				RunAsNonRoot: &[]bool{true}[0],
				Capabilities: &corev1.Capabilities{
					Drop: []corev1.Capability{"ALL"},
				},
			}
		}
	}

	for i := range mpod.Spec.InitContainers {
		if mpod.Spec.InitContainers[i].SecurityContext == nil {
			mpod.Spec.InitContainers[i].SecurityContext = &corev1.SecurityContext{
				ReadOnlyRootFilesystem: &[]bool{true}[0],
				AllowPrivilegeEscalation: &[]bool{false}[0],
				RunAsNonRoot: &[]bool{true}[0],
				Capabilities: &corev1.Capabilities{
					Drop: []corev1.Capability{"ALL"},
				},
			}
		}
	}

	mpl.Logger.WithField("inj_sec_context", mpod.Spec.SecurityContext).Printf("setting security context")

	return mpod, nil
}

func IsPodSecurityContextEmpty(psc *corev1.PodSecurityContext) bool {
	if psc == nil {
		return true
	}
	if psc.SELinuxOptions != nil {
		return false
	}
	if psc.RunAsUser != nil {
		return false
	}
	if psc.RunAsNonRoot != nil {
		return false
	}
	if len(psc.SupplementalGroups) > 0 {
		return false
	}
	if psc.FSGroup != nil {
		return false
	}
	if psc.RunAsGroup != nil {
		return false
	}
	if len(psc.Sysctls) > 0 {
		return false
	}
	if psc.WindowsOptions != nil {
		return false
	}
	if psc.FSGroupChangePolicy != nil {
		return false
	}
	if psc.SeccompProfile != nil {
		return false
	}
	return true
}