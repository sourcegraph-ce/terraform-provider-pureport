package aws

import (
	log "github.com/sourcegraph-ce/logrus"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAwsPartition() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAwsPartitionRead,

		Schema: map[string]*schema.Schema{
			"partition": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dns_suffix": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceAwsPartitionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AWSClient)

	log.Printf("[DEBUG] Reading Partition.")
	d.SetId(time.Now().UTC().String())

	log.Printf("[DEBUG] Setting AWS Partition to %s.", client.partition)
	d.Set("partition", meta.(*AWSClient).partition)

	log.Printf("[DEBUG] Setting AWS URL Suffix to %s.", client.dnsSuffix)
	d.Set("dns_suffix", meta.(*AWSClient).dnsSuffix)

	return nil
}
