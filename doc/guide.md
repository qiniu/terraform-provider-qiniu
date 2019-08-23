# Qiniu Cloud Provider 用户指南

## 使用案例

```hcl
# Configure Qiniu Account
provider "qiniu" {
  access_key = "<Qiniu Access Key>"
  secret_key = "<Qiniu Secret Key>"
}

# Create Qiniu Bucket
resource "qiniu_bucket" "basic_bucket" {
  name      = "basic-test-terraform-1"  # Bucket Name
  region_id = "z0"                      # Bucket Region, "z0" means East China
  private   = false                     # Public bucket
}

# Create Qiniu Object
resource "qiniu_bucket_object" "basic_object" {
  bucket    = "basic-test-terraform-1"  # Bucket Name
  key       = "qiniu-key"               # File Key
  source    = "/path/to/file"           # File Path to upload
}

# Qiniu Buckets Data Source
data "qiniu_buckets" "z1" {
  name_regex = "^bucket-"
  region_id = "z1"
}

# Qiniu Buckets Objects Data Source
data "qiniu_buckets_objects" "all" {
  bucket = "basic-test-terraform-1"
}
```

## 认证

支持静态认证和环境变量认证两种方式

### 静态认证

```hcl
provider "qiniu" {
  access_key = "${var.access_key}"
  secret_key = "${var.secret_key}"
}
```

### 环境变量认证

```hcl
provider "qiniu" {
}
```

```bash
export QINIU_ACCESS_KEY="anaccesskey"
export QINIU_SECRET_KEY="asecretkey"
terraform plan
```

## 私有云配置

在私有云情况下，需要手工配置各个私有云服务器的地址

```hcl
provider "qiniu" {
  use_https      = true
  central_rs_url = "https://<中心化 RS 域名>"
  rs_url         = "https://<RS 域名>"
  rsf_url        = "https://<RSF 域名>"
  up_url         = "https://<UP 域名>"
  api_url        = "https://<API 域名>"
  io_url         = "https://<IO 域名>"
}
```

## 资源创建

### 创建 Bucket

创建 Bucket 时，`name` 和 `region_id` 是必须指定的

```hcl
resource "qiniu_bucket" "basic_bucket" {
  name      = "basic-test-terraform-1"  # Bucket Name
  region_id = "z0"                      # Bucket Region, "z0" means East China
}
```

可以额外指定 `private` 参数创建私有 Bucket

```hcl
resource "qiniu_bucket" "basic_bucket" {
  name      = "basic-test-terraform-1"
  region_id = "z0"
  private   = true
}
```

可以额外指定 `image_url` 和 `image_host` 参数指定镜像回源地址

```hcl
resource "qiniu_bucket" "basic_bucket" {
  name      = "basic-test-terraform-1"
  region_id = "z0"
  image_url = "http://target_url"
}
```

```hcl
resource "qiniu_bucket" "basic_bucket" {
  name       = "basic-test-terraform-1"
  region_id  = "z0"
  image_url  = "http://target_url"
  image_host = "www.qiniu.com"          # 该参数可选，仅在给出了 image_url 参数之后才生效
}
```

### 上传文件

可以将指定路径的文件上传至指定 Bucket

```hcl
resource "qiniu_bucket" "basic_object" {
  bucket    = "basic-test-terraform-1"  # Bucket Name
  key       = "keyname"                 # Key Name
  source    = "/path/to/file"
}
```

## 数据源

### 列出 Bucket

列出当前账户所有 Bucket

```hcl
data "qiniu_buckets" "all" {
}
```

列出匹配指定正则表达式的 Bucket

```hcl
data "qiniu_buckets" "all" {
  name_regex = "^data-"
}
```

列出指定区域的 Bucket

```hcl
data "qiniu_buckets" "all" {
  region_id = "z0"
}
```

可以获取到所有 Bucket 的名称

```hcl
${data.qiniu_buckets.all.names}
```

也可以获取到所有 Bucket 的详细信息

```hcl
${data.qiniu_buckets.all.buckets}
```

### 列出文件

列出指定 Bucket 的文件

```hcl
data "qiniu_buckets_objects" "all" {
  bucket = "basic-test-terraform-1"
}
```

列出指定 Bucket 中所有文件名以指定字符串开头的文件

```hcl
data "qiniu_buckets_objects" "all" {
  bucket = "basic-test-terraform-1"
  prefix = "data-"
}
```

可以获取到所有文件的名称

```hcl
${data.qiniu_buckets_objects.all.keys}
```

也可以获取到所有文件的详细信息

```hcl
${data.qiniu_buckets_objects.all.key_infos}
```
