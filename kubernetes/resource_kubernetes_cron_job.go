package kubernetes

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"k8s.io/api/batch/v1beta1"
	"k8s.io/api/batch/v2alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const cronJobResourceGroupName = "cronjobs"

var cronJobAPIGroups = []APIGroup{batchV1beta1, batchV2alpha1}

var cronJobNotSupportedError = fmt.Errorf("could not find Kubernetes API group that supports CronJob resources")

func resourceKubernetesCronJob() *schema.Resource {
	return &schema.Resource{
		Create: resourceKubernetesCronJobCreate,
		Read:   resourceKubernetesCronJobRead,
		Update: resourceKubernetesCronJobUpdate,
		Delete: resourceKubernetesCronJobDelete,
		Exists: resourceKubernetesCronJobExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"metadata": namespacedMetadataSchema("cronjob", true),
			"spec": {
				Type:        schema.TypeList,
				Description: "Spec of the cron job owned by the cluster",
				Required:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: cronJobSpecFields(),
				},
			},
		},
	}
}

func resourceKubernetesCronJobCreate(d *schema.ResourceData, meta interface{}) error {
	kp := meta.(*kubernetesProvider)
	conn := kp.conn

	metadata := expandMetadata(d.Get("metadata").([]interface{}))
	spec, err := expandCronJobSpec(d.Get("spec").([]interface{}))
	if err != nil {
		return err
	}
	spec.JobTemplate.ObjectMeta.Annotations = metadata.Annotations

	job := v1beta1.CronJob{
		ObjectMeta: metadata,
		Spec:       spec,
	}

	created := &v1beta1.CronJob{}

	log.Printf("[INFO] Creating new cron job: %#v", job)
	apiGroup, err := kp.highestSupportedAPIGroup(cronJobResourceGroupName, cronJobAPIGroups...)
	if err != nil {
		return err
	}
	switch apiGroup {
	case batchV1beta1:
		created, err = conn.BatchV1beta1().CronJobs(metadata.Namespace).Create(&job)

	case batchV2alpha1:
		beta := &v2alpha1.CronJob{}
		Convert(job, beta)
		out, err2 := conn.BatchV2alpha1().CronJobs(metadata.Namespace).Create(beta)
		if err2 != nil {
			err = err2
			break
		}
		Convert(out, created)

	default:
		err = cronJobNotSupportedError
	}
	if err != nil {
		return err
	}

	log.Printf("[INFO] Submitted new cron job: %#v", created)

	d.SetId(buildId(created.ObjectMeta))

	return resourceKubernetesCronJobRead(d, meta)
}

func resourceKubernetesCronJobUpdate(d *schema.ResourceData, meta interface{}) error {
	kp := meta.(*kubernetesProvider)
	conn := kp.conn

	namespace, _, err := idParts(d.Id())
	if err != nil {
		return err
	}

	metadata := expandMetadata(d.Get("metadata").([]interface{}))
	spec, err := expandCronJobSpec(d.Get("spec").([]interface{}))
	if err != nil {
		return err
	}
	spec.JobTemplate.ObjectMeta.Annotations = metadata.Annotations

	cronjob := &v1beta1.CronJob{
		ObjectMeta: metadata,
		Spec:       spec,
	}

	log.Printf("[INFO] Updating cron job %s: %s", d.Id(), cronjob)

	out := &v1beta1.CronJob{}
	apiGroup, err := kp.highestSupportedAPIGroup(cronJobResourceGroupName, cronJobAPIGroups...)
	if err != nil {
		return err
	}
	switch apiGroup {
	case batchV1beta1:
		out, err = conn.BatchV1beta1().CronJobs(namespace).Update(cronjob)

	case batchV2alpha1:
		alpha := &v2alpha1.CronJob{}
		Convert(cronjob, alpha)
		alphaOut, err2 := conn.BatchV2alpha1().CronJobs(namespace).Update(alpha)
		if err2 != nil {
			err = err2
			break
		}
		Convert(alphaOut, out)

	default:
		err = cronJobNotSupportedError
	}
	if err != nil {
		return err
	}
	log.Printf("[INFO] Submitted updated cron job: %#v", out)

	d.SetId(buildId(out.ObjectMeta))
	return resourceKubernetesCronJobRead(d, meta)
}

func resourceKubernetesCronJobRead(d *schema.ResourceData, meta interface{}) error {
	kp := meta.(*kubernetesProvider)

	namespace, name, err := idParts(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Reading cron job %s", name)
	job, err := readCronJob(kp, namespace, name)
	if err != nil {
		log.Printf("[DEBUG] Received error: %#v", err)
		return err
	}
	log.Printf("[INFO] Received cron job: %#v", job)

	// Remove server-generated labels unless using manual selector
	if _, ok := d.GetOk("spec.0.manual_selector"); !ok {
		labels := job.ObjectMeta.Labels

		if _, ok := labels["controller-uid"]; ok {
			delete(labels, "controller-uid")
		}

		if _, ok := labels["cron-job-name"]; ok {
			delete(labels, "cron-job-name")
		}

		if job.Spec.JobTemplate.Spec.Selector != nil &&
			job.Spec.JobTemplate.Spec.Selector.MatchLabels != nil {
			labels = job.Spec.JobTemplate.Spec.Selector.MatchLabels
		}

		if _, ok := labels["controller-uid"]; ok {
			delete(labels, "controller-uid")
		}
	}

	err = d.Set("metadata", flattenMetadata(job.ObjectMeta, d))
	if err != nil {
		return err
	}

	jobSpec, err := flattenCronJobSpec(job.Spec, d)
	if err != nil {
		return err
	}

	err = d.Set("spec", jobSpec)
	if err != nil {
		return err
	}

	return nil
}

func resourceKubernetesCronJobDelete(d *schema.ResourceData, meta interface{}) error {
	kp := meta.(*kubernetesProvider)
	conn := kp.conn

	namespace, name, err := idParts(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleting cron job: %#v", name)
	apiGroup, err := kp.highestSupportedAPIGroup(cronJobResourceGroupName, cronJobAPIGroups...)
	if err != nil {
		return err
	}
	switch apiGroup {
	case batchV1beta1:
		err = conn.BatchV1beta1().CronJobs(namespace).Delete(name, nil)

	case batchV2alpha1:
		err = conn.BatchV2alpha1().CronJobs(namespace).Delete(name, nil)

	default:
		err = cronJobNotSupportedError
	}
	if err != nil {
		return err
	}

	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		_, err := readCronJob(kp, namespace, name)
		if err != nil {
			if statusErr, ok := err.(*errors.StatusError); ok && statusErr.ErrStatus.Code == 404 {
				return nil
			}
			return resource.NonRetryableError(err)
		}

		e := fmt.Errorf("Cron Job %s still exists", name)
		return resource.RetryableError(e)
	})
	if err != nil {
		return err
	}

	log.Printf("[INFO] Cron Job %s deleted", name)

	d.SetId("")
	return nil
}

func resourceKubernetesCronJobExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	kp := meta.(*kubernetesProvider)

	namespace, name, err := idParts(d.Id())
	if err != nil {
		return false, err
	}

	log.Printf("[INFO] Checking cron job %s", name)
	_, err = readCronJob(kp, namespace, name)
	if err != nil {
		if statusErr, ok := err.(*errors.StatusError); ok && statusErr.ErrStatus.Code == 404 {
			return false, nil
		}
		log.Printf("[DEBUG] Received error: %#v", err)
	}
	return true, err
}

func readCronJob(kp *kubernetesProvider, namespace, name string) (cj *v1beta1.CronJob, err error) {
	conn := kp.conn

	log.Printf("[INFO] Reading CronJob %s", name)
	cj = &v1beta1.CronJob{}

	apiGroup, err := kp.highestSupportedAPIGroup(cronJobResourceGroupName, cronJobAPIGroups...)
	if err != nil {
		return nil, err
	}
	log.Printf("[INFO] Reading CronJob using %s API Group", apiGroup)

	switch apiGroup {
	case batchV1beta1:
		cj, err = conn.BatchV1beta1().CronJobs(namespace).Get(name, metav1.GetOptions{})
		return cj, err

	case batchV2alpha1:
		out, err := conn.BatchV2alpha1().CronJobs(namespace).Get(name, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}
		err = Convert(out, cj)

	default:
		return nil, cronJobNotSupportedError
	}

	return cj, err
}
