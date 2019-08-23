# Terraform Provider for Qiniu Cloud
==================

- Website: https://www.terraform.io
- [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
- [![LICENSE](https://img.shields.io/badge/license-Mozilla--2.0-yellowgreen)](https://www.mozilla.org/en-US/MPL/2.0/)
- [![Build Status](https://api.travis-ci.org/bachue/terraform-provider-qiniu.svg?branch=master)](https://travis-ci.org/bachue/terraform-provider-qiniu)
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)

<img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" width="400px">
<img src="https://mars-assets.qnssl.com/qiniulog/img-slogan-blue-en.png" alt="Qiniu Cloud">

Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) 0.12.x
-	[Go](https://golang.org/doc/install) 1.12 (to build the provider plugin)

Building The Provider
---------------------

Clone repository

```sh
$ git clone git@github.com:qiniu/terraform-provider-qiniu.git --recurse-submodules
```

Enter the provider directory and build the provider

```sh
$ cd terraform-provider-qiniu
$ make
```

Using the Provider
----------------------

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

Developing the Provider
---------------------------

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.12+ is *required*).

To compile the provider, run `make`. This will build the provider and put the provider binary in the `bin/` directory.

```sh
$ make
...
$ bin/terraform-provider-qiniu
...
```

In order to test the provider, you can copy `.env.example` to `.env`, and edit the `.env` file

```sh
QINIU_ACCESS_KEY=<Qiniu Access Key>
QINIU_SECRET_KEY=<Qiniu Secret Key>
TF_ACC=1
```

And then simply run

```sh
make test
```
