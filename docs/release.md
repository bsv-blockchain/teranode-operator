# Catalog bundle release workflow
```bash
$ IMG=434394763103.dkr.ecr.eu-north-1.amazonaws.com/teranode-operator:<version> make bundle-push
$ aws ecr get-login-password --region eu-north-1 | docker login --username AWS --password-stdin 434394763103.dkr.ecr.eu-north-1.amazonaws.com
$ BUNDLE_IMG=434394763103.dkr.ecr.eu-north-1.amazonaws.com/teranode-operator-bundle:v0.5.1 make bundle-build bundle-push
# Update catalog template
$ opm alpha render-template semver -o yaml < catalog-templates/teranode-operator.yaml > bsva-catalog/catalog.yaml
```
