// +build !ignore_autogenerated

// Copyright (c) 2021 Red Hat, Inc.
// Copyright Contributors to the Open Cluster Management project
/*


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ImageOverrides) DeepCopyInto(out *ImageOverrides) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ImageOverrides.
func (in *ImageOverrides) DeepCopy() *ImageOverrides {
	if in == nil {
		return nil
	}
	out := new(ImageOverrides)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PodResource) DeepCopyInto(out *PodResource) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PodResource.
func (in *PodResource) DeepCopy() *PodResource {
	if in == nil {
		return nil
	}
	out := new(PodResource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SearchCustomization) DeepCopyInto(out *SearchCustomization) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SearchCustomization.
func (in *SearchCustomization) DeepCopy() *SearchCustomization {
	if in == nil {
		return nil
	}
	out := new(SearchCustomization)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *SearchCustomization) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SearchCustomizationList) DeepCopyInto(out *SearchCustomizationList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]SearchCustomization, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SearchCustomizationList.
func (in *SearchCustomizationList) DeepCopy() *SearchCustomizationList {
	if in == nil {
		return nil
	}
	out := new(SearchCustomizationList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *SearchCustomizationList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SearchCustomizationSpec) DeepCopyInto(out *SearchCustomizationSpec) {
	*out = *in
	if in.Persistence != nil {
		in, out := &in.Persistence, &out.Persistence
		*out = new(bool)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SearchCustomizationSpec.
func (in *SearchCustomizationSpec) DeepCopy() *SearchCustomizationSpec {
	if in == nil {
		return nil
	}
	out := new(SearchCustomizationSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SearchCustomizationStatus) DeepCopyInto(out *SearchCustomizationStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SearchCustomizationStatus.
func (in *SearchCustomizationStatus) DeepCopy() *SearchCustomizationStatus {
	if in == nil {
		return nil
	}
	out := new(SearchCustomizationStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SearchOperator) DeepCopyInto(out *SearchOperator) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SearchOperator.
func (in *SearchOperator) DeepCopy() *SearchOperator {
	if in == nil {
		return nil
	}
	out := new(SearchOperator)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *SearchOperator) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SearchOperatorList) DeepCopyInto(out *SearchOperatorList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]SearchOperator, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SearchOperatorList.
func (in *SearchOperatorList) DeepCopy() *SearchOperatorList {
	if in == nil {
		return nil
	}
	out := new(SearchOperatorList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *SearchOperatorList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SearchOperatorSpec) DeepCopyInto(out *SearchOperatorSpec) {
	*out = *in
	out.SearchImageOverrides = in.SearchImageOverrides
	out.Redisgraph_Resource = in.Redisgraph_Resource
	if in.NodeSelector != nil {
		in, out := &in.NodeSelector, &out.NodeSelector
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SearchOperatorSpec.
func (in *SearchOperatorSpec) DeepCopy() *SearchOperatorSpec {
	if in == nil {
		return nil
	}
	out := new(SearchOperatorSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SearchOperatorStatus) DeepCopyInto(out *SearchOperatorStatus) {
	*out = *in
	if in.DeployRedisgraph != nil {
		in, out := &in.DeployRedisgraph, &out.DeployRedisgraph
		*out = new(bool)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SearchOperatorStatus.
func (in *SearchOperatorStatus) DeepCopy() *SearchOperatorStatus {
	if in == nil {
		return nil
	}
	out := new(SearchOperatorStatus)
	in.DeepCopyInto(out)
	return out
}
