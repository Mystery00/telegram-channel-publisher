{
  "apiVersion": "moment.halo.run/v1alpha1",
  "kind": "Moment",
  "metadata": {
    "generateName": "telegram-moment-"
  },
  "spec": {
    "content": {
      "html": "{{.Html}}",
      "medium": {{.Medium}},
      "raw": "{{.Content}}"
    },
    "owner": "",
    "releaseTime": "{{.ReleaseTime.Format "2006-01-02T15:04:05.000Z"}}",
    "tags": {{.Tags}},
    "visible": "PUBLIC"
  }
}