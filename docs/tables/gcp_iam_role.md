# Table: gcp_iam_role

An IAM role is an IAM entity that defines a set of permissions for making GCP service requests.

## Examples

### IAM role basic info

```sql
select
  name,
  role_id,
  deleted,
  description,
  title
from
  gcp_iam_role;
```

### List of IAM roles which are in BETA stage

```sql
select
  name,
  description,
  stage
from
  gcp_iam_role
where
  stage = 'BETA';
```

### List of IAM customer managed roles

```sql
select
  name,
  description,
  stage
from
  gcp_iam_role
where
  is_gcp_managed = false;
```
