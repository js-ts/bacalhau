bacalhau_version       = "v0.1.39"
bacalhau_port          = "1235"
bacalhau_connect_node0 = "QmP6RVpStuEoShqTTTiS2e3PYazcd54sj2RaZTeJP9VCeh"
ipfs_version           = "v0.12.2"
gcp_project            = "bacalhau-staging"
instance_count         = 3
region                 = "southamerica-east1"
zone                   = "southamerica-east1-b"
volume_size_gb         = 10
machine_type           = "e2-standard-4"
protect_resources      = true
auto_subnets           = true
ingress_cidrs          = ["0.0.0.0/0"]
ssh_access_cidrs       = ["0.0.0.0/0"]
num_gpu_machines       = 0
