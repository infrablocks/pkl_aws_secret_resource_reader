# Pkl AWS Secret Resource Reader

An implementation of the Pkl Resource Reader SPI that allows you to pull secrets from AWS Secrets Manager.

This tool assumes you have your AWS credentials configured in the environment and uses AWS Go SDK under the hood.

To use with Pkl:

`pkl eval example.pkl --external-resource-reader awssecret=pkl_aws_secret_reader`