// Copyright (c) 2020 Red Hat, Inc.
// Copyright Contributors to the Open Cluster Management project

package controllers

import (
	"context"
	"os"
	"strconv"
	"testing"

	searchv1alpha1 "github.com/stolostron/search-operator/api/v1alpha1"
	"github.com/stretchr/testify/assert"
	appv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type testSetup struct {
	scheme                *runtime.Scheme
	request               reconcile.Request
	srchOperator          *searchv1alpha1.SearchOperator
	secret                *corev1.Secret
	statefulsetWithPVC    *appv1.StatefulSet
	statefulsetWithOutPVC *appv1.StatefulSet
	pvc                   *corev1.PersistentVolumeClaim
	podWithPVC            *corev1.Pod
	podWithOutPVC         *corev1.Pod
	unSchedulablePod      *corev1.Pod
	customizationCR       *searchv1alpha1.SearchCustomization
	context               context.Context
}

func commonSetup() testSetup {
	testScheme := scheme.Scheme

	namespace = "test-cluster"
	searchv1alpha1.AddToScheme(testScheme)
	testScheme.AddKnownTypes(corev1.SchemeGroupVersion, &corev1.Secret{})
	waitSecondsForPodChk = 2
	redisPodResource := searchv1alpha1.PodResource{
		RequestMemory: "64Mi",
		RequestCPU:    "25m",
		LimitMemory:   "1Gi",
		LimitCPU:      "250m",
	}
	testSearchOperator := &searchv1alpha1.SearchOperator{
		TypeMeta: metav1.TypeMeta{
			APIVersion: searchv1alpha1.GroupVersion.String(),
			Kind:       "SearchOperator",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "searchoperator",
			Namespace: namespace,
		},
		Spec: searchv1alpha1.SearchOperatorSpec{
			Redisgraph_Resource: redisPodResource,
		},
	}
	testSecret := newRedisSecret(testSearchOperator, testScheme)
	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      "searchoperator",
			Namespace: namespace,
		},
	}
	client := fake.NewFakeClientWithScheme(testScheme)
	testSearchOperatorReconciler := SearchOperatorReconciler{client, log, testScheme}

	testStatefulsetWithPVC := testSearchOperatorReconciler.executeDeployment(client, testSearchOperator, true, true)
	testStatefulsetWithOutPVC := testSearchOperatorReconciler.executeDeployment(client, testSearchOperator, false, true)
	// Set PVC Size to 10Gi
	fakePVC := createFakeNamedPVC("10Gi", testSearchOperator.Namespace, nil)
	fakePodWithPVC := createFakeRedisGraphPod(namespace, true, true)
	fakePodWithOutPVC := createFakeRedisGraphPod(namespace, false, true)
	fakeUnschedulablePod := createFakeRedisGraphPod(namespace, false, false)
	fakeSearchCustCR := createFakeSearchCustomizationCR(namespace, false)
	testSetup := testSetup{scheme: testScheme,
		request:               req,
		srchOperator:          testSearchOperator,
		secret:                testSecret,
		statefulsetWithPVC:    testStatefulsetWithPVC,
		statefulsetWithOutPVC: testStatefulsetWithOutPVC,
		pvc:                   fakePVC,
		podWithPVC:            fakePodWithPVC,
		podWithOutPVC:         fakePodWithOutPVC,
		unSchedulablePod:      fakeUnschedulablePod,
		customizationCR:       fakeSearchCustCR,
		context:               context.TODO()}
	return testSetup
}

func Test_searchOperatorNotFound(t *testing.T) {
	testSetup := commonSetup()
	req := testSetup.request
	client := fake.NewFakeClientWithScheme(testSetup.scheme)
	nilSearchOperator := SearchOperatorReconciler{client, log, testSetup.scheme}

	_, err := nilSearchOperator.Reconcile(testSetup.context, req)
	assert.Nil(t, err, "Expected Reconcile Error to be Nil. Got error: %v", err)

	instance := &searchv1alpha1.SearchOperator{}
	err = client.Get(context.TODO(), req.NamespacedName, instance)
	assert.True(t, errors.IsNotFound(err), "Expected error: SearchOperator to be Not Found. Got %v", err.Error())
}

func Test_secretCreatedWithOwnerRef(t *testing.T) {
	testSetup := commonSetup()
	testSecret := testSetup.secret
	client := fake.NewFakeClientWithScheme(testSetup.scheme, testSetup.srchOperator)
	nilSearchOperator := SearchOperatorReconciler{client, log, testSetup.scheme}

	_, err := nilSearchOperator.Reconcile(testSetup.context, testSetup.request)

	found := &corev1.Secret{}
	err = client.Get(context.TODO(), types.NamespacedName{Name: testSecret.Name, Namespace: testSecret.Namespace}, found)
	assert.Nil(t, err, "Expected secret to be created. Got error: %v", err)
	assert.Equal(t, testSecret.Name, found.Name, "Secret is created with expected name.")
	assert.Equal(t, testSecret.Namespace, found.Namespace, "Secret is created in expected namespace.")
	assert.EqualValues(t, testSecret.GetLabels(), found.GetLabels(), "Secret is created with expected labels.")
	ownerRefArray := found.GetOwnerReferences()
	assert.NotNil(t, ownerRefArray, "Created secret should have an ownerReference.")
	assert.Len(t, ownerRefArray, 1, "Created secret should have an ownerReference.")

	ownerRef := ownerRefArray[0]
	assert.Equal(t, testSetup.srchOperator.APIVersion, ownerRef.APIVersion, "secret's ownerRef has expected APIVersion.")
	assert.Equal(t, testSetup.srchOperator.Kind, ownerRef.Kind, "secret's ownerRef has expected Kind.")
	assert.Equal(t, testSetup.srchOperator.Name, ownerRef.Name, "secret's ownerRef has expected Name.")

}

func Test_secretAlreadyExists(t *testing.T) {
	testSetup := commonSetup()
	testSecret := testSetup.secret

	client := fake.NewFakeClientWithScheme(testSetup.scheme, testSetup.srchOperator, testSecret)
	nilSearchOperator := SearchOperatorReconciler{client, log, testSetup.scheme}

	_, err := nilSearchOperator.Reconcile(testSetup.context, testSetup.request)

	found := &corev1.Secret{}
	err = client.Get(context.TODO(), types.NamespacedName{Name: testSecret.Name, Namespace: testSecret.Namespace}, found)

	assert.Nil(t, err, "Expected secret to be created. Got error: %v", err)
	assert.EqualValues(t, testSecret.GetObjectMeta(), found.GetObjectMeta(), "Secret is created with expected labels.")
	assert.EqualValues(t, testSecret.Data, found.Data, "Secret is created with expected data.")

}

func Test_EmptyDirStatefulsetCreatedWithOwnerRef(t *testing.T) {
	testSetup := commonSetup()

	client := fake.NewFakeClientWithScheme(testSetup.scheme, testSetup.srchOperator, testSetup.secret, testSetup.podWithOutPVC)
	nilSearchOperator := SearchOperatorReconciler{client, log, testSetup.scheme}
	var err error

	_, err = nilSearchOperator.Reconcile(testSetup.context, testSetup.request)

	foundStatefulset := &appv1.StatefulSet{}
	err = client.Get(context.TODO(), types.NamespacedName{Name: testSetup.statefulsetWithOutPVC.Name, Namespace: testSetup.statefulsetWithOutPVC.Namespace}, foundStatefulset)
	assert.Nil(t, err, "Expected statefulset to be created. Got error: %v", err)

	assert.Equal(t, testSetup.statefulsetWithOutPVC.Name, foundStatefulset.Name, "Statefulset is created with expected name.")
	assert.Equal(t, testSetup.statefulsetWithOutPVC.Namespace, foundStatefulset.Namespace, "Statefulset is created in expected namespace.")

	ownerRefArray := foundStatefulset.GetOwnerReferences()

	assert.NotNil(t, ownerRefArray, "Created Statefulset should have an ownerReference.")
	assert.Len(t, ownerRefArray, 1, "Created Statefulset should have an ownerReference.")

	ownerRef := ownerRefArray[0]
	assert.Equal(t, testSetup.srchOperator.APIVersion, ownerRef.APIVersion, "redisgraph statefulset's ownerRef has expected APIVersion.")
	assert.Equal(t, testSetup.srchOperator.Kind, ownerRef.Kind, "redisgraph statefulset's ownerRef has expected Kind.")
	assert.Equal(t, testSetup.srchOperator.Name, ownerRef.Name, "redisgraph statefulset's ownerRef has expected Name.")
}

func Test_StatefulsetWithNoPersistenceStatus(t *testing.T) {
	testSetup := commonSetup()

	client := fake.NewFakeClientWithScheme(testSetup.scheme, testSetup.srchOperator, testSetup.secret, testSetup.podWithOutPVC, testSetup.customizationCR)
	nilSearchOperator := SearchOperatorReconciler{client, log, testSetup.scheme}
	var err error
	instance := &searchv1alpha1.SearchOperator{}
	//Turn persistence to false in customizationCR
	persistence := false
	testSetup.customizationCR.Spec.Persistence = &persistence
	err = client.Update(context.TODO(), testSetup.customizationCR)

	_, err = nilSearchOperator.Reconcile(testSetup.context, testSetup.request)
	//Calling reconcile again to check if shorter path with 1 sec wait time is used the second time
	_, err = nilSearchOperator.Reconcile(testSetup.context, testSetup.request)
	err = client.Get(context.TODO(), testSetup.request.NamespacedName, instance)

	assert.Nil(t, err, "Expected searchoperator to be found with no error. Got error: %v", err)
	assert.Equal(t, statusNoPersistence, instance.Status.PersistenceStatus, "Search Operator status updated with statusNoPersistence as expected.")
}

func Test_StatefulsetWithPVC(t *testing.T) {
	testSetup := commonSetup()
	req := testSetup.request
	testStatefulset := testSetup.statefulsetWithPVC

	client := fake.NewFakeClientWithScheme(testSetup.scheme, testSetup.srchOperator, testSetup.secret, testSetup.pvc, testSetup.podWithPVC)
	nilSearchOperator := SearchOperatorReconciler{client, log, testSetup.scheme}
	var err error

	instance := &searchv1alpha1.SearchOperator{}
	err = client.Get(context.TODO(), req.NamespacedName, instance)
	assert.Nil(t, err, "Expected search Operator to be created. Got error: %v", err)

	//Persistence is enabled by default in search operator
	_, err = nilSearchOperator.Reconcile(testSetup.context, req)
	//Calling reconcile again to check if shorter path with 1 sec wait time is used the second time
	_, err = nilSearchOperator.Reconcile(testSetup.context, req)
	assert.Nil(t, err, "Expected search Operator reconcile to complete successfully. Got error: %v", err)

	err = client.Get(context.TODO(), req.NamespacedName, instance)
	assert.Nil(t, err, "Expected search Operator to be created. Got error: %v", err)

	foundStatefulset := &appv1.StatefulSet{}
	err = client.Get(context.TODO(), types.NamespacedName{Name: testStatefulset.Name, Namespace: testStatefulset.Namespace}, foundStatefulset)

	assert.Nil(t, err, "Expected Statefulset to be created. Got error: %v", err)

	assert.Equal(t, testStatefulset.Name, foundStatefulset.Name, "Statefulset is created with expected name.")
	assert.Equal(t, testStatefulset.Namespace, foundStatefulset.Namespace, "Statefulset is created in expected namespace.")
	assert.EqualValues(t, testStatefulset.Spec.Template.Spec, foundStatefulset.Spec.Template.Spec, "Statefulset is created with expected template spec.")
	assert.Equal(t, statusUsingPVC, instance.Status.PersistenceStatus, "Search Operator status updated with statusUsingPVC as expected.")

}

func Test_FailSettingupWithPVC(t *testing.T) {
	testSetup := commonSetup()

	req := testSetup.request
	testStatefulset := testSetup.statefulsetWithPVC

	//TODO: Passing already existing secret doesn't set ownerRef - testSecret
	client := fake.NewFakeClientWithScheme(testSetup.scheme, testSetup.srchOperator, testSetup.secret, testSetup.unSchedulablePod)
	nilSearchOperator := SearchOperatorReconciler{client, log, testSetup.scheme}
	var err error

	instance := &searchv1alpha1.SearchOperator{}
	err = client.Get(context.TODO(), req.NamespacedName, instance)
	assert.Nil(t, err, "Expected search Operator to be created. Got error: %v", err)

	//Set persistence to true in customizationCR - this should cause statefulset to fall back to empty dir since we don't have PVC
	persistence := true
	testSetup.customizationCR.Spec.Persistence = &persistence
	err = client.Create(context.TODO(), testSetup.customizationCR)

	_, err = nilSearchOperator.Reconcile(testSetup.context, req)
	assert.NotNil(t, err, "Expected Reconcile error to be not nil. Got nil.")

	err = client.Get(context.TODO(), req.NamespacedName, instance)
	assert.Nil(t, err, "Expected search Operator to be created. Got error: %v", err)

	foundStatefulset := &appv1.StatefulSet{}
	err = client.Get(context.TODO(), types.NamespacedName{Name: testStatefulset.Name, Namespace: testStatefulset.Namespace}, foundStatefulset)

	assert.True(t, errors.IsNotFound(err), "Expected statefulset Not Found error. Got %v", err.Error())
	assert.Equal(t, statusFailedUsingPVC, instance.Status.PersistenceStatus, "Search Operator status updated with statusFailedUsingPVC as expected.")
}

func Test_UnschedulablePodWithPersistence(t *testing.T) {
	testSetup := commonSetup()

	req := testSetup.request
	testStatefulset := testSetup.statefulsetWithOutPVC

	client := fake.NewFakeClientWithScheme(testSetup.scheme, testSetup.srchOperator, testSetup.secret, testStatefulset, testSetup.unSchedulablePod)
	nilSearchOperator := SearchOperatorReconciler{client, log, testSetup.scheme}
	var err error

	instance := &searchv1alpha1.SearchOperator{}
	err = client.Get(context.TODO(), req.NamespacedName, instance)
	assert.Nil(t, err, "Expected search Operator to be created. Got error: %v", err)

	//Persistence is enabled by default in operator - this should cause statefulset to fall back to empty dir since we don't have PVC
	_, err = nilSearchOperator.Reconcile(testSetup.context, req)
	assert.NotNil(t, err, "Expected Reconcile error to be not nil. Got nil.")
	assert.Equal(t, "Redisgraph Pod not running", err.Error(), "Expected Redisgraph Pod not running error. Got %v", err)

	foundStatefulset := &appv1.StatefulSet{}
	err = client.Get(context.TODO(), types.NamespacedName{Name: testStatefulset.Name, Namespace: testStatefulset.Namespace}, foundStatefulset)
	assert.True(t, errors.IsNotFound(err), "Expected Statefulset Not Found error. Got %v", err.Error())

	err = client.Get(context.TODO(), req.NamespacedName, instance)
	assert.Equal(t, statusFailedDegraded, instance.Status.PersistenceStatus, "Search Operator status updated with statusFailedDegraded as expected.")
}

func Test_UnschedulablePodWithOutPersistence(t *testing.T) {
	testSetup := commonSetup()

	req := testSetup.request
	testStatefulset := testSetup.statefulsetWithOutPVC

	client := fake.NewFakeClientWithScheme(testSetup.scheme, testSetup.srchOperator, testSetup.customizationCR, testSetup.secret, testStatefulset, testSetup.unSchedulablePod)
	nilSearchOperator := SearchOperatorReconciler{client, log, testSetup.scheme}
	var err error

	instance := &searchv1alpha1.SearchOperator{}
	err = client.Get(context.TODO(), req.NamespacedName, instance)
	assert.Nil(t, err, "Expected search Operator to be created. Got error: %v", err)

	//Set persistence to false in customizationCR - check if status is updated correctly if pod fails to get scheduled
	persistence := false
	testSetup.customizationCR.Spec.Persistence = &persistence
	err = client.Update(context.TODO(), testSetup.customizationCR)
	_, err = nilSearchOperator.Reconcile(testSetup.context, req)

	assert.NotNil(t, err, "Expected Reconcile error to be not nil. Got nil.")
	assert.Equal(t, "Redisgraph Pod not running", err.Error(), "Expected Redisgraph Pod not running error. Got %v", err)
	err = client.Get(context.TODO(), req.NamespacedName, instance)
	assert.Equal(t, statusFailedNoPersistence, instance.Status.PersistenceStatus, "Search Operator status updated with statusFailedNoPersistence as expected.")
}

func Test_DoNotDeployRedisPod(t *testing.T) {
	testSetup := commonSetup()
	//Set DEPLOY_REDISGRAPH env to false to stop redisgraph pod from being deployed
	os.Setenv("DEPLOY_REDISGRAPH", "false")
	deployRedisgraphPod, deployVarPresent = os.LookupEnv("DEPLOY_REDISGRAPH")
	deploy, deployVarErr = strconv.ParseBool(deployRedisgraphPod)
	req := testSetup.request
	testStatefulset := testSetup.statefulsetWithPVC

	client := fake.NewFakeClientWithScheme(testSetup.scheme, testSetup.srchOperator, testSetup.secret, testSetup.pvc, testSetup.statefulsetWithPVC)

	nilSearchOperator := SearchOperatorReconciler{client, log, testSetup.scheme}
	var err error
	instance := &searchv1alpha1.SearchOperator{}
	err = client.Get(context.TODO(), req.NamespacedName, instance)
	assert.Nil(t, err, "Expected search Operator to be created. Got error: %v", err)

	_, err = nilSearchOperator.Reconcile(testSetup.context, req)
	assert.Nil(t, err, "Expected search Operator reconcile to complete successfully. Got error: %v", err)

	err = client.Get(context.TODO(), req.NamespacedName, instance)
	assert.Nil(t, err, "Expected search Operator to be created. Got error: %v", err)

	foundStatefulset := &appv1.StatefulSet{}
	err = client.Get(context.TODO(), types.NamespacedName{Name: testStatefulset.Name, Namespace: testStatefulset.Namespace}, foundStatefulset)
	assert.True(t, errors.IsNotFound(err), "Expected error: redisgraph statefulset to be Not Found. Got %v", err.Error())
	assert.Equal(t, redisNotRunning, instance.Status.PersistenceStatus, "Search Operator status not set as expected.")
	//Resetting the variables
	os.Unsetenv("DEPLOY_REDISGRAPH")
	deployRedisgraphPod, deployVarPresent = os.LookupEnv("DEPLOY_REDISGRAPH")
	deploy, deployVarErr = strconv.ParseBool(deployRedisgraphPod)
}

func Test_UnschedulablePodWithDisAllowDegradedMode(t *testing.T) {
	testSetup := commonSetup()

	req := testSetup.request
	testStatefulset := testSetup.statefulsetWithOutPVC

	client := fake.NewFakeClientWithScheme(testSetup.scheme, testSetup.srchOperator, testSetup.customizationCR, testSetup.secret, testStatefulset, testSetup.unSchedulablePod)
	nilSearchOperator := SearchOperatorReconciler{client, log, testSetup.scheme}
	var err error

	instance := &searchv1alpha1.SearchOperator{}
	err = client.Get(context.TODO(), req.NamespacedName, instance)
	assert.Nil(t, err, "Expected search Operator to be created. Got error: %v", err)

	//Set persistence to true in customizationCR - pod should fail to start as we don't have PVC
	persistence := true
	testSetup.customizationCR.Spec.Persistence = &persistence
	err = client.Update(context.TODO(), testSetup.customizationCR)
	_, err = nilSearchOperator.Reconcile(testSetup.context, req)

	assert.NotNil(t, err, "Expected Reconcile error to be not nil. Got nil.")
	assert.Equal(t, "Redisgraph Pod not running", err.Error(), "Expected Redisgraph Pod not running error. Got %v", err)
	err = client.Get(context.TODO(), req.NamespacedName, instance)
	assert.Equal(t, statusFailedUsingPVC, instance.Status.PersistenceStatus, "Search Operator status updated with statusFailedUsingPVC as expected.")
}

func TestUpdateCR(t *testing.T) {
	testSetup := commonSetup()
	client := fake.NewFakeClientWithScheme(testSetup.scheme)
	var err error
	err = updateCRs(client, testSetup.srchOperator, "status", testSetup.customizationCR, false, "", "10G", true)
	assert.True(t, errors.IsNotFound(err), "Expected searchOperator Not Found error. Got %v", err.Error())

	client = fake.NewFakeClientWithScheme(testSetup.scheme, testSetup.srchOperator)
	err = updateCRs(client, testSetup.srchOperator, "status", testSetup.customizationCR, false, "", "10G", true)
	assert.True(t, errors.IsNotFound(err), "Expected customizationCR Not Found error. Got %v", err.Error())

	client = fake.NewFakeClientWithScheme(testSetup.scheme, testSetup.srchOperator, testSetup.customizationCR)
	err = updateCRs(client, testSetup.srchOperator, "status", testSetup.customizationCR, false, "", "10G", true)
	assert.Nil(t, err, "Expected CR statuses to be updated successfully. Got error: %v", err)
}

func TestGetPVC(t *testing.T) {
	storageClass = "test"
	pvc := getPVC()
	assert.NotNil(t, pvc.Spec.StorageClassName, "Expected StorageClassName to be not nil.")
	assert.Equal(t, "test", *pvc.Spec.StorageClassName, "Expected StorageClassName not found. Got %v", *pvc.Spec.StorageClassName)

	storageClass = ""
	pvc = getPVC()
	assert.Nil(t, pvc.Spec.StorageClassName, "Expected empty StorageClassName. Got: %s", pvc.Spec.StorageClassName)
}

func TestRestartCollector(t *testing.T) {
	//create fake Collector Pod
	labels := map[string]string{}
	labels["app"] = "search-prod"
	labels["component"] = "search-collector"
	testSetup := commonSetup()
	collectorPod := &corev1.Pod{ObjectMeta: metav1.ObjectMeta{Namespace: namespace, Name: "search-collector-pod", Labels: labels}}

	client := fake.NewFakeClientWithScheme(testSetup.scheme, testSetup.srchOperator, collectorPod)
	nilSearchOperator := SearchOperatorReconciler{client, log, testSetup.scheme}

	nilSearchOperator.restartSearchComponents()
	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Namespace: namespace,
			Name:      "search-collector-pod",
		},
	}
	collectorPod1 := &corev1.Pod{}
	err := client.Get(context.TODO(), req.NamespacedName, collectorPod1)
	assert.True(t, errors.IsNotFound(err), "Expected error: SearchCollector pod to be Not Found. Got %v", err.Error())
}

func createFakeNamedPVC(requestBytes string, namespace string, userAnnotations map[string]string) *corev1.PersistentVolumeClaim {
	annotations := map[string]string{}
	for k, v := range userAnnotations {
		annotations[k] = v
	}

	return &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			UID:         "testid",
			Name:        pvcName,
			Namespace:   namespace,
			Annotations: annotations,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			Selector: nil, // Provisioner doesn't support selector
			Resources: corev1.ResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceName(corev1.ResourceStorage): resource.MustParse(requestBytes),
				},
			},
		},
	}
}

func createFakeRedisGraphPod(namespace string, persistence, schedulable bool) *corev1.Pod {
	labels := map[string]string{}
	labels["app"] = appName
	labels["component"] = component
	image := "quay.io/stolostron/search-operator:latest"
	containerStatuses := []corev1.ContainerStatus{}
	containerStatus := corev1.ContainerStatus{Ready: true}
	containerStatuses = append(containerStatuses, containerStatus)
	status := corev1.PodStatus{ContainerStatuses: containerStatuses}

	if !schedulable {
		podStatusConditions := []corev1.PodCondition{}
		podUnschedulableStatusCondition := corev1.PodCondition{Reason: "Unschedulable"} //ContainerStatus{Reason:"Unschedulable"}
		podStatusConditions = append(podStatusConditions, podUnschedulableStatusCondition)
		status = corev1.PodStatus{Conditions: podStatusConditions}
	}
	// unschedulableStatus := corev1.PodStatus{ContainerStatuses: containerStatuses}

	persistentVolSource := corev1.PersistentVolumeClaimVolumeSource{ClaimName: pvcName}
	emptyDirVolSource := corev1.EmptyDirVolumeSource{}
	var volSource corev1.VolumeSource

	if !schedulable {
		return &corev1.Pod{ObjectMeta: metav1.ObjectMeta{Namespace: namespace, Labels: labels}, Spec: corev1.PodSpec{Volumes: []corev1.Volume{{VolumeSource: volSource}}, Containers: []corev1.Container{{Image: image}}}, Status: status}
	}
	if persistence {
		volSource = corev1.VolumeSource{PersistentVolumeClaim: &persistentVolSource}
		return &corev1.Pod{ObjectMeta: metav1.ObjectMeta{Namespace: namespace, Labels: labels}, Spec: corev1.PodSpec{Volumes: []corev1.Volume{{Name: "persist", VolumeSource: volSource}}, Containers: []corev1.Container{{Image: image}}}, Status: status}
	}
	volSource = corev1.VolumeSource{EmptyDir: &emptyDirVolSource}
	return &corev1.Pod{ObjectMeta: metav1.ObjectMeta{Namespace: namespace, Labels: labels}, Spec: corev1.PodSpec{Volumes: []corev1.Volume{{Name: "persist", VolumeSource: volSource}}, Containers: []corev1.Container{{Image: image}}}, Status: status}

}

func createFakeSearchCustomizationCR(namespace string, persistence bool) *searchv1alpha1.SearchCustomization {
	return &searchv1alpha1.SearchCustomization{TypeMeta: metav1.TypeMeta{
		APIVersion: searchv1alpha1.GroupVersion.String(),
		Kind:       "SearchCustomization"},
		ObjectMeta: metav1.ObjectMeta{Namespace: namespace, Name: "searchcustomization"},
		Spec:       searchv1alpha1.SearchCustomizationSpec{Persistence: &persistence, StorageSize: "1Gi"}}
}
