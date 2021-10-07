package aws

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/service/s3control"
	"github.com/hashicorp/aws-sdk-go-base/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
)

func ResourceBucketPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceBucketPolicyCreate,
		Read:   resourceBucketPolicyRead,
		Update: resourceBucketPolicyUpdate,
		Delete: resourceBucketPolicyDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateArn,
			},
			"policy": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateFunc:     validation.StringIsJSON,
				DiffSuppressFunc: suppressEquivalentAwsPolicyDiffs,
			},
		},
	}
}

func resourceBucketPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).S3ControlConn

	bucket := d.Get("bucket").(string)

	input := &s3control.PutBucketPolicyInput{
		Bucket: aws.String(bucket),
		Policy: aws.String(d.Get("policy").(string)),
	}

	_, err := conn.PutBucketPolicy(input)

	if err != nil {
		return fmt.Errorf("error creating S3 Control Bucket Policy (%s): %w", bucket, err)
	}

	d.SetId(bucket)

	return resourceBucketPolicyRead(d, meta)
}

func resourceBucketPolicyRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).S3ControlConn

	parsedArn, err := arn.Parse(d.Id())

	if err != nil {
		return fmt.Errorf("error parsing S3 Control Bucket ARN (%s): %w", d.Id(), err)
	}

	if parsedArn.AccountID == "" {
		return fmt.Errorf("error parsing S3 Control Bucket ARN (%s): unknown format", d.Id())
	}

	input := &s3control.GetBucketPolicyInput{
		AccountId: aws.String(parsedArn.AccountID),
		Bucket:    aws.String(d.Id()),
	}

	output, err := conn.GetBucketPolicy(input)

	if !d.IsNewResource() && tfawserr.ErrCodeEquals(err, "NoSuchBucket") {
		log.Printf("[WARN] S3 Control Bucket Policy (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if !d.IsNewResource() && tfawserr.ErrCodeEquals(err, "NoSuchBucketPolicy") {
		log.Printf("[WARN] S3 Control Bucket Policy (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if !d.IsNewResource() && tfawserr.ErrCodeEquals(err, "NoSuchOutpost") {
		log.Printf("[WARN] S3 Control Bucket Policy (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("error reading S3 Control Bucket Policy (%s): %w", d.Id(), err)
	}

	if output == nil {
		return fmt.Errorf("error reading S3 Control Bucket Policy (%s): empty response", d.Id())
	}

	d.Set("bucket", d.Id())
	d.Set("policy", output.Policy)

	return nil
}

func resourceBucketPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).S3ControlConn

	input := &s3control.PutBucketPolicyInput{
		Bucket: aws.String(d.Id()),
		Policy: aws.String(d.Get("policy").(string)),
	}

	_, err := conn.PutBucketPolicy(input)

	if err != nil {
		return fmt.Errorf("error updating S3 Control Bucket Policy (%s): %w", d.Id(), err)
	}

	return resourceBucketPolicyRead(d, meta)
}

func resourceBucketPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).S3ControlConn

	parsedArn, err := arn.Parse(d.Id())

	if err != nil {
		return fmt.Errorf("error parsing S3 Control Bucket ARN (%s): %w", d.Id(), err)
	}

	input := &s3control.DeleteBucketPolicyInput{
		AccountId: aws.String(parsedArn.AccountID),
		Bucket:    aws.String(d.Id()),
	}

	_, err = conn.DeleteBucketPolicy(input)

	if tfawserr.ErrCodeEquals(err, "NoSuchBucket") {
		return nil
	}

	if tfawserr.ErrCodeEquals(err, "NoSuchBucketPolicy") {
		return nil
	}

	if tfawserr.ErrCodeEquals(err, "NoSuchOutpost") {
		return nil
	}

	if err != nil {
		return fmt.Errorf("error deleting S3 Control Bucket Policy (%s): %w", d.Id(), err)
	}

	return nil
}
